package pagination

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Scalingo/go-utils/mongo/document"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type BadRequestError struct {
	Errors map[string][]string `json:"errors"`
}

func (err BadRequestError) Error() string {
	errArray := make([]string, 0, len(err.Errors))
	for errTitle, errValues := range err.Errors {
		errArray = append(errArray, fmt.Sprintf("* %s → %s", errTitle, strings.Join(errValues, ", ")))
	}
	return strings.Join(errArray, "\n")
}

func NewBadRequestErrors() *BadRequestError {
	return &BadRequestError{
		Errors: make(map[string][]string),
	}
}

func NewPaginationService(opts ServiceOpts) Service {
	return opts
}

type Meta struct {
	CurrentPage int  `json:"current_page"`
	PrevPage    *int `json:"prev_page"`
	NextPage    *int `json:"next_page"`
	TotalPages  int  `json:"total_pages"`
	TotalCount  int  `json:"total_count"`
	perPageNum  int
}

type Service interface {
	Paginate(ctx context.Context, DBQuery bson.M, collection string, result interface{}, sortField string) (*Meta, error)
}

type ServiceOpts struct {
	PerPageDefault     int
	MaxPerPage         int
	PageQueryParams    string
	PerPageQueryParams string
}

func (opts ServiceOpts) paramValidation(meta *Meta, collection string) error {
	badRequestErr := NewBadRequestErrors()
	pageErr := "Requested page"
	perPageErr := "per_page"
	var err error

	// Options validation
	if opts.PerPageDefault > opts.MaxPerPage ||
		opts.PerPageDefault <= 0 ||
		opts.MaxPerPage <= 0 {
		return errors.New("Usage: MaxPerPage > PerPageDefault > 0")
	}
	// Parameter validation
	if collection == "" {
		return errors.New("collection must be set")
	}

	// Default values assignation
	if opts.PageQueryParams == "" {
		opts.PageQueryParams = "1"
	}
	if opts.PerPageQueryParams == "" {
		opts.PerPageQueryParams = fmt.Sprintf("%d", opts.PerPageDefault)
	}

	// Request parameters validation
	meta.CurrentPage, err = strconv.Atoi(opts.PageQueryParams)
	if err != nil {
		badRequestErr.Errors[pageErr] =
			append(badRequestErr.Errors[pageErr], "fail to parse page parameter")
	}
	if meta.CurrentPage <= 0 {
		badRequestErr.Errors[pageErr] =
			append(badRequestErr.Errors[pageErr], "must be greater then 0")
	}
	meta.perPageNum, err = strconv.Atoi(opts.PerPageQueryParams)
	if err != nil {
		badRequestErr.Errors[perPageErr] =
			append(badRequestErr.Errors[perPageErr], "fail to parse per_page parameter")
	}
	if meta.perPageNum <= 0 {
		badRequestErr.Errors[perPageErr] =
			append(badRequestErr.Errors[perPageErr], "must be greater then 0")
	}
	if meta.perPageNum < 0 || meta.perPageNum > opts.MaxPerPage {
		badRequestErr.Errors[perPageErr] =
			append(badRequestErr.Errors[perPageErr], fmt.Sprintf("must be between 0 and %d", opts.MaxPerPage))
	}

	if badRequestErr.Errors != nil && len(badRequestErr.Errors) > 0 {
		return badRequestErr
	}

	// Total page calculation
	meta.TotalPages = meta.TotalCount / meta.perPageNum
	// We truncate to the higher integer except in the case of a "round" division
	if meta.TotalCount%meta.perPageNum != 0 {
		meta.TotalPages++
	}
	if meta.CurrentPage > meta.TotalPages {
		badRequestErr.Errors[pageErr] =
			append(badRequestErr.Errors[pageErr], fmt.Sprintf("must be between 0 and %d", meta.TotalPages))
		return badRequestErr
	}

	return nil
}

func (opts ServiceOpts) Paginate(ctx context.Context,
	DBQuery bson.M,
	collection string,
	result interface{},
	sortField string) (*Meta, error) {

	var err error
	meta := Meta{}

	query, session := document.WhereQuery(ctx, collection, DBQuery)
	defer session.Close()

	meta.TotalCount, err = query.Count()
	if err != nil {
		return nil, errors.Wrapf(err, "fail to count items for the collection %s", collection)
	}

	if meta.TotalCount == 0 {
		return &meta, nil
	}

	err = opts.paramValidation(&meta, collection)
	if err != nil {
		return nil, err
	}

	// PrevPage must be null when unavailable
	if meta.CurrentPage > 1 {
		meta.PrevPage = new(int)
		*meta.PrevPage = meta.CurrentPage - 1
	}
	// NextPage must be null when unavailable
	if meta.CurrentPage < meta.TotalPages {
		meta.NextPage = new(int)
		*meta.NextPage = meta.CurrentPage + 1
	}

	// Mongo Skip will become slower as the offset increases. Hence this could be
	// improved using a range. But we have to use an index field to refer to it.
	// More information about the range-queries:
	// https://docs.mongodb.com/manual/reference/method/cursor.skip/#using-range-queries
	offset := (meta.CurrentPage - 1) * meta.perPageNum
	err = query.Skip(offset).Sort(sortField).Limit(meta.perPageNum).All(result)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to query database for collection %s", collection)
	}

	return &meta, nil
}
