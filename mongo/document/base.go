package document

import (
	"context"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/Scalingo/go-utils/errors/v2"
)

type Base struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

func (d Base) IsPersisted() bool {
	return !d.CreatedAt.IsZero()
}

func (d Base) getID() bson.ObjectId {
	return d.ID
}

func (d *Base) ensureID() {
	if d.ID == "" {
		d.ID = bson.NewObjectId()
	}
}

func (d *Base) ensureCreatedAt() {
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now()
	}
}

func (d *Base) setUpdatedAt(t time.Time) {
	d.UpdatedAt = t
}

func (d *Base) getUpdatedAt() time.Time {
	return d.UpdatedAt
}

func (d Base) scope(query bson.M) bson.M {
	return query
}

func (d *Base) destroy(ctx context.Context, collection string) error {
	return ReallyDestroy(ctx, collection, d)
}

func (d *Base) Validate(_ context.Context) *errors.ValidationErrors {
	return nil
}

func (d *Base) ValidateWithInternalError(_ context.Context) (*errors.ValidationErrors, error) {
	return nil, ErrValidateNoInternalErrorFunc
}
