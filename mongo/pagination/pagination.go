package pagination

import (
	"context"
	"fmt"
	"net/url"
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

type Meta struct {
	CurrentPage int  `json:"current_page"`
	PrevPage    *int `json:"prev_page"`
	NextPage    *int `json:"next_page"`
	TotalPages  int  `json:"total_pages"`
	TotalCount  int  `json:"total_count"`
}

type PaginationService interface {
	Paginate(ctx context.Context, requestQuery url.Values, DBQuery bson.M, collection string, result interface{}, sortField string) (Meta, error)
}

type PaginationServiceOpts struct {
	PerPageDefault     int
	MaxPerPage         int
	PageQueryParams    string
	PerPageQueryParams string
}

func (opts PaginationServiceOpts) Paginate(ctx context.Context,
	DBQuery bson.M,
	collection string,
	result interface{},
	sortField string) (*Meta, error) {

	meta := Meta{}
	badRequestErr := NewBadRequestErrors()
	var err error
	perPageNum := 0

	if opts.PageQueryParams != "" {
		meta.CurrentPage, err = strconv.Atoi(opts.PageQueryParams)
		if err != nil {
			badRequestErr.Errors["Parameter error"] = append([]string{"fail to parse page parameter"})
		}
	}
	if opts.PerPageQueryParams != "" {
		perPageNum, err = strconv.Atoi(opts.PerPageQueryParams)
		if err != nil {
			badRequestErr.Errors["Parameter error"] = append([]string{"fail to parse page parameter"})
			return nil, badRequestErr
		}
	}

	if perPageNum < 0 || perPageNum > opts.MaxPerPage {
		badRequestErr.Errors["Parameter error"] =
			append([]string{fmt.Sprintf("per_page must be between 0 and %d", opts.MaxPerPage)})
		return nil, badRequestErr
	}

	if perPageNum == 0 {
		perPageNum = opts.PerPageDefault
	}
	if meta.CurrentPage == 0 {
		meta.CurrentPage = 1
	}

	query, session := document.WhereQuery(ctx, collection, DBQuery)
	defer session.Close()

	meta.TotalCount, err = query.Count()
	if err != nil {
		return nil, errors.Wrapf(err, "fail to count items for the collection %s", collection)
	}

	meta.TotalPages = meta.TotalCount / perPageNum
	// We truncate to the higher integer except in the case of a "round" division
	if meta.TotalCount%perPageNum != 0 {
		meta.TotalPages++
	}

	if meta.CurrentPage > meta.TotalPages || meta.CurrentPage <= 0 {
		badRequestErr.Errors["Parameter error"] =
			append([]string{fmt.Sprintf("requested page must be between 0 and %d", meta.TotalPages)})
		return nil, badRequestErr
	}

	if meta.CurrentPage > 1 {
		meta.PrevPage = new(int)
		*meta.PrevPage = meta.CurrentPage - 1
	}
	if meta.CurrentPage < meta.TotalPages {
		meta.NextPage = new(int)
		*meta.NextPage = meta.CurrentPage + 1
	}

	offset := (meta.CurrentPage - 1) * perPageNum
	err = query.Skip(offset).Sort(sortField).Limit(perPageNum).All(result)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to query database for collection %s", collection)
	}

	return &meta, nil
}
