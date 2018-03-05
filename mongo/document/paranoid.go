package document

import (
	"context"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Paranoid struct {
	Base      `bson:",inline"`
	DeletedAt time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func (d *Paranoid) SetDeletedAt(t time.Time) {
	d.DeletedAt = t
}

func (d Paranoid) scope(query bson.M) bson.M {
	if _, ok := query["deleted_at"]; !ok {
		query["deleted_at"] = nil
	}
	return query
}

func (d *Paranoid) destroy(ctx context.Context, collectionName string) error {
	now := time.Now()
	d.SetDeletedAt(now)
	return Update(ctx, collectionName, bson.M{"$set": bson.M{"deleted_at": now}}, d)
}
