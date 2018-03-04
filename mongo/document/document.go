package document

import (
	"context"
	"time"

	"gopkg.in/errgo.v1"
	"gopkg.in/mgo.v2/bson"

	"github.com/Scalingo/appsdeck-database/controller/mongo"
	"github.com/Scalingo/go-utils/logger"
)

type Document interface {
	GetID() bson.ObjectId
}

type ParanoiaDeletable interface {
	Updatable
	SetDeletedAt(time.Time)
}

type Updatable interface {
	Document
	SetUpdatedAt(time.Time)
}

func Save(ctx context.Context, collectionName string, doc Document) error {
	log := logger.Get(ctx)
	c := mongo.Collection(collectionName)
	defer c.Database.Session.Close()
	log.WithField(collectionName, doc).Debugf("save '%v'", collectionName)
	_, err := c.UpsertId(doc.GetID(), &doc)
	return err
}

// Remove the volume from the database.
// Handle with care...
func Destroy(ctx context.Context, collectionName string, doc Document) error {
	log := logger.Get(ctx)
	c := mongo.Collection(collectionName)
	defer c.Database.Session.Close()
	log.WithField(collectionName, doc).Debugf("remove '%v'", collectionName)
	return c.RemoveId(doc.GetID())
}

func ParanoiaDelete(ctx context.Context, collectionName string, d ParanoiaDeletable) error {
	now := time.Now()
	d.SetDeletedAt(now)
	err := Update(ctx, collectionName, bson.M{"$set": bson.M{"deleted_at": now}}, d)
	if err != nil {
		return errgo.Notef(err, "fail to run mongo update")
	}
	return nil
}

func Find(ctx context.Context, collectionName string, id bson.ObjectId, doc Document) error {
	c := mongo.Collection(collectionName)
	defer c.Database.Session.Close()
	query := bson.M{"_id": id, "deleted_at": nil}
	return c.Find(query).One(doc)
}

func WhereParanoia(ctx context.Context, collectionName string, query bson.M, data interface{}) error {
	if query == nil {
		query = bson.M{}
	}
	if _, ok := query["deleted_at"]; !ok {
		query["deleted_at"] = nil
	}
	return Where(ctx, collectionName, query, data)
}

func Where(ctx context.Context, collectionName string, query bson.M, data interface{}) error {
	return WhereSort(ctx, collectionName, query, data)
}

func WhereSort(ctx context.Context, collectionName string, query bson.M, data interface{}, sortFields ...string) error {
	c := mongo.Collection(collectionName)
	defer c.Database.Session.Close()

	if query == nil {
		query = bson.M{}
	}

	err := c.Find(query).Sort(sortFields...).All(data)
	if err != nil {
		return errgo.Notef(err, "fail to query mongo %v", query)
	}
	return nil
}

func Update(ctx context.Context, collectionName string, update bson.M, doc Updatable) error {
	log := logger.Get(ctx)
	c := mongo.Collection(collectionName)
	defer c.Database.Session.Close()

	now := time.Now()
	doc.SetUpdatedAt(now)
	if _, ok := update["$set"]; ok {
		update["$set"].(bson.M)["updated_at"] = now
	}
	log.WithField("query", update).Debugf("update %v", collectionName)
	return c.UpdateId(doc.GetID(), update)
}
