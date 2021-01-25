package pagination

import (
	"context"
	"fmt"

	"github.com/Scalingo/go-utils/mongo/document"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

func NewPaginationByCursorService(opts ServiceOpts) ServiceByCursor {
	return opts
}

type PaginateByCursorOpts struct {
	Cursor      bson.M
	AmountItems int
	Query       bson.M
	SortOrder   string
}

type ServiceByCursor interface {
	PaginateByCursor(ctx context.Context, collection string, result interface{}, opts PaginateByCursorOpts) error
}

func (s ServiceOpts) paramValidationByCursor(collection string, opts *PaginateByCursorOpts) error {
	badRequestErr := NewBadRequestErrors()
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
	if opts.AmountItems == 0 {
		opts.AmountItems = s.PerPageDefault
	}
	if opts.SortOrder == "" {
		opts.SortOrder = "-_id"
	}

	// Request parameters validation
	if opts.AmountItems <= 0 {
		badRequestErr.Errors[perPageErr] =
			append(badRequestErr.Errors[perPageErr], "must be greater then 0")
	}
	if opts.AmountItems < 0 || opts.AmountItems > s.MaxPerPage {
		badRequestErr.Errors[perPageErr] =
			append(badRequestErr.Errors[perPageErr], fmt.Sprintf("must be between 0 and %d", s.MaxPerPage))
	}

	if badRequestErr.Errors != nil && len(badRequestErr.Errors) > 0 {
		return badRequestErr
	}

	return nil
}

// PaginateByCursor query a mongo database and fill the restult param with a
// paginate result. The pagination is made in reverse order and the field on
// which the comparison is made must be comparable with the `$lt` mongo operator
func (s ServiceOpts) PaginateByCursor(ctx context.Context,
	collection string,
	result interface{},
	opts PaginateByCursorOpts) error {

	var err error
	var optsQuery bson.M

	if opts.Cursor == nil {
		optsQuery = opts.Query
	} else {
		optsQuery = bson.M{"$and": []bson.M{
			opts.Query,
			opts.Cursor,
		}}
	}

	err = s.paramValidationByCursor(collection, &opts)
	if err != nil {
		return err
	}

	query, session := document.WhereQuery(ctx, collection, optsQuery)
	defer session.Close()

	//{ $and:[{virtual_storage_name: "vs_name_2"}, {app_id: {"$lt": 2}} ]}
	err = query.Sort(opts.SortOrder).Limit(opts.AmountItems).All(result)
	if err != nil {
		return errors.Wrapf(err, "fail to query database for collection %s", collection)
	}

	return nil
}
