package pagination

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"

	"github.com/Scalingo/go-utils/mongo/document"
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

type PaginateOpts struct {
	PageNumber  int
	AmountItems int
	Query       bson.M
	SortOrder   string
}

type Service interface {
	Paginate(ctx context.Context, collection string, result interface{}, opts PaginateOpts) (Meta, error)
}

type ServiceOpts struct {
	PerPageDefault int
	MaxPerPage     int
}

func (s ServiceOpts) paramValidation(meta *Meta, collection string, opts *PaginateOpts) error {
	badRequestErr := NewBadRequestErrors()
	pageErr := "Requested page"
	perPageErr := "per_page"

	// Options validation
	if s.PerPageDefault > s.MaxPerPage ||
		s.PerPageDefault <= 0 ||
		s.MaxPerPage <= 0 {
		return errors.New("invalid pagination service configuration: MaxPerPage > PerPageDefault > 0")
	}
	// Parameter validation
	if collection == "" {
		return errors.New("invalid pagination service configuration: collection must be set")
	}

	// Default values assignation
	if opts.PageNumber == 0 {
		opts.PageNumber = 1
	}
	if opts.AmountItems == 0 {
		opts.AmountItems = s.PerPageDefault
	}
	if opts.SortOrder == "" {
		opts.SortOrder = "_id"
	}

	// Request parameters validation
	meta.CurrentPage = opts.PageNumber
	if meta.CurrentPage <= 0 {
		badRequestErr.Errors[pageErr] =
			append(badRequestErr.Errors[pageErr], "must be greater then 0")
	}
	meta.perPageNum = opts.AmountItems
	if meta.perPageNum <= 0 {
		badRequestErr.Errors[perPageErr] =
			append(badRequestErr.Errors[perPageErr], "must be greater then 0")
	}
	if meta.perPageNum < 0 || meta.perPageNum > s.MaxPerPage {
		badRequestErr.Errors[perPageErr] =
			append(badRequestErr.Errors[perPageErr], fmt.Sprintf("must be between 0 and %d", s.MaxPerPage))
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

func (s ServiceOpts) Paginate(ctx context.Context,
	collection string,
	result interface{},
	opts PaginateOpts) (Meta, error) {
	var err error
	meta := Meta{}

	query, session := document.WhereQuery(ctx, collection, opts.Query)
	defer session.Close()

	meta.TotalCount, err = query.Count()
	if err != nil {
		return meta, errors.Wrapf(err, "fail to count items for the collection %s", collection)
	}

	if meta.TotalCount == 0 {
		return meta, nil
	}

	err = s.paramValidation(&meta, collection, &opts)
	if err != nil {
		return meta, err
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
	err = query.Skip(offset).Sort(opts.SortOrder).Limit(meta.perPageNum).All(result)
	if err != nil {
		return meta, errors.Wrapf(err, "fail to query database for collection %s", collection)
	}

	return meta, nil
}
