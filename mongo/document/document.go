package document

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/mongo"
)

type document interface {
	getID() bson.ObjectId
	ensureID()
	ensureCreatedAt()
	setUpdatedAt(time.Time)
	Validable
}

type scopable interface {
	scope(bson.M) bson.M
}

type destroyable interface {
	destroy(ctx context.Context, collection string) error
}

type Validable interface {
	Validate(ctx context.Context) *ValidationErrors
}

// Create inser the document in the database, returns an error if document already exists and set CreatedAt timestamp
func Create(ctx context.Context, collectionName string, doc document) error {
	log := logger.Get(ctx)
	doc.ensureID()
	doc.ensureCreatedAt()
	doc.setUpdatedAt(time.Now())

	if err := doc.Validate(ctx); err != nil {
		return err
	}

	c := mongo.Session(log).Clone().DB("").C(collectionName)
	defer c.Database.Session.Close()
	log.WithField(collectionName, doc).Debugf("save '%v'", collectionName)
	return c.Insert(&doc)

}

func Save(ctx context.Context, collectionName string, doc document) error {
	log := logger.Get(ctx)
	doc.ensureID()
	doc.ensureCreatedAt()
	doc.setUpdatedAt(time.Now())

	if err := doc.Validate(ctx); err != nil {
		return err
	}

	c := mongo.Session(log).Clone().DB("").C(collectionName)
	defer c.Database.Session.Close()
	log.WithField(collectionName, doc).Debugf("save '%v'", collectionName)
	_, err := c.UpsertId(doc.getID(), &doc)
	return err
}

// Destroy really deletes
func Destroy(ctx context.Context, collectionName string, doc destroyable) error {
	return doc.destroy(ctx, collectionName)
}

func ReallyDestroy(ctx context.Context, collectionName string, doc document) error {
	log := logger.Get(ctx)
	c := mongo.Session(log).Clone().DB("").C(collectionName)
	defer c.Database.Session.Close()
	log.WithField(collectionName, doc).Debugf("remove '%v'", collectionName)
	return c.RemoveId(doc.getID())
}

// Find is finding the model with objectid id in the collection name, with its
// default scope for paranoid documents, it won't look at documents tagged as
// deleted
func Find(ctx context.Context, collectionName string, id bson.ObjectId, doc scopable) error {
	query := doc.scope(bson.M{"_id": id})
	return FindOne(ctx, collectionName, query, doc)
}

// FindUnscoped is similar as Find but does not care of the default scope of
// the document default scope.
func FindUnscoped(ctx context.Context, collectionName string, id bson.ObjectId, doc interface{}) error {
	query := bson.M{"_id": id}
	return FindOne(ctx, collectionName, query, doc)
}

func FindSort(ctx context.Context, collectionName string, query bson.M, doc interface{}, sortFields ...string) error {
	log := logger.Get(ctx)
	c := mongo.Session(log).Clone().DB("").C(collectionName)
	defer c.Database.Session.Close()

	return c.Find(query).Sort(sortFields...).One(doc)
}

func FindOne(ctx context.Context, collectionName string, query bson.M, doc interface{}) error {
	return FindSort(ctx, collectionName, query, doc)
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
	log := logger.Get(ctx)
	c := mongo.Session(log).Clone().DB("").C(collectionName)
	defer c.Database.Session.Close()

	if query == nil {
		query = bson.M{}
	}

	err := c.Find(query).Sort(sortFields...).All(data)
	if err != nil {
		return fmt.Errorf("fail to query mongo %v: %v", query, err)
	}
	return nil
}

func WhereParanoiaIter(ctx context.Context, collectionName string, query bson.M, fun func(*mgo.Iter) error, sortFields ...string) error {
	if query == nil {
		query = bson.M{}
	}
	if _, ok := query["deleted_at"]; !ok {
		query["deleted_at"] = nil
	}
	return WhereIter(ctx, collectionName, query, fun, sortFields...)
}

func WhereIter(ctx context.Context, collectionName string, query bson.M, fun func(*mgo.Iter) error, sortFields ...string) error {
	log := logger.Get(ctx)
	c := mongo.Session(log).Clone().DB("").C(collectionName)
	defer c.Database.Session.Close()

	if query == nil {
		query = bson.M{}
	}

	iter := c.Find(query).Sort(sortFields...).Iter()
	defer iter.Close()

	err := fun(iter)
	if err != nil {
		return fmt.Errorf("error occured during iteration over collection %v with query %v: %v", collectionName, query, err)
	}
	if iter.Err() != nil {
		return fmt.Errorf("fail to iterate over collection %v with query %v: %v", collectionName, query, err)
	}
	return nil
}

func Update(ctx context.Context, collectionName string, update bson.M, doc document) error {
	log := logger.Get(ctx)
	c := mongo.Session(log).Clone().DB("").C(collectionName)
	defer c.Database.Session.Close()

	now := time.Now()
	doc.setUpdatedAt(now)
	if _, ok := update["$set"]; ok {
		update["$set"].(bson.M)["updated_at"] = now
	}

	if err := doc.Validate(ctx); err != nil {
		return err
	}

	log.WithField("query", update).Debugf("update %v", collectionName)
	return c.UpdateId(doc.getID(), update)
}
