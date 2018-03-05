package document

import (
	"context"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Base struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

func (d Base) getID() bson.ObjectId {
	return d.ID
}

func (d *Base) setCreatedAt(t time.Time) {
	d.CreatedAt = t
}

func (d *Base) setUpdatedAt(t time.Time) {
	d.UpdatedAt = t
}

func (d Base) scope(query bson.M) bson.M {
	return query
}

func (d Base) destroy(ctx context.Context, collection string) error {
	return ReallyDestroy(ctx, collection, d)
}
