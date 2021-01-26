package pagination

import (
	"context"
	"fmt"

	"github.com/Scalingo/go-utils/mongo/document"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

func NewByCursorService(opts ServiceOpts) ServiceByCursor {
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

func (s ServiceOpts) paramValidationByCursor(collection string, opts PaginateByCursorOpts) error {
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

	// Request parameters validation
	if opts.AmountItems <= 0 {
		badRequestErr.Errors[perPageErr] =
			append(badRequestErr.Errors[perPageErr], "must be strictly positive")
	}
	if opts.AmountItems > s.MaxPerPage {
		badRequestErr.Errors[perPageErr] =
			append(badRequestErr.Errors[perPageErr], fmt.Sprintf("must be lower than %d", s.MaxPerPage))
	}

	if badRequestErr.Errors != nil && len(badRequestErr.Errors) > 0 {
		return badRequestErr
	}

	return nil
}

// PaginateByCursor queries a MongoDB database and fill the result param with a
// paginate result. The pagination is made in reverse order and the field on
// which the comparison is made must be comparable with the `$lt` mongo operator
func (s ServiceOpts) PaginateByCursor(ctx context.Context,
	collection string,
	result interface{},
	opts PaginateByCursorOpts) error {

	var optsQuery bson.M

	// Default values assignation
	if opts.AmountItems == 0 {
		opts.AmountItems = s.PerPageDefault
	}
	if opts.SortOrder == "" {
		opts.SortOrder = "-_id"
	}

	err := s.paramValidationByCursor(collection, opts)
	if err != nil {
		return err
	}

	if opts.Cursor == nil {
		// In case of empty cursor, we will return the first page.
		optsQuery = opts.Query
	} else {
		// Provide an empty bson.M as cursor will have the same behavior as no
		// cursor: so it will return the first page.
		optsQuery = bson.M{"$and": []bson.M{
			opts.Query,
			opts.Cursor,
		}}
	}

	query, session := document.WhereQuery(ctx, collection, optsQuery)
	defer session.Close()

	err = query.Sort(opts.SortOrder).Limit(opts.AmountItems).All(result)
	if err != nil {
		return errors.Wrapf(err, "fail to query database for collection %s", collection)
	}

	return nil
}
