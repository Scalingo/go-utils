package document

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/mongo"
)

type SortField string

type document interface {
	getID() bson.ObjectId
	ensureID()
	ensureCreatedAt()
	setUpdatedAt(time.Time)
	getUpdatedAt() time.Time
	Validable
}

var _ document = &Base{}
var _ document = &Paranoid{}

type scopable interface {
	scope(bson.M) bson.M
}

var _ scopable = Base{}
var _ scopable = Paranoid{}

type destroyable interface {
	destroy(ctx context.Context, collectionName string) error
}

var _ destroyable = &Base{}
var _ destroyable = &Paranoid{}

type Closer interface {
	Close()
}

var ErrValidateNoInternalErrorFunc = stderrors.New("no validation returning an internal error has been implemented")

type Validable interface {
	// Validate will be used if no ValidateWithInternalError is defined on a document
	// It is not useful to have both defined on a document, only ValidationWithInternalError
	// would be used in this case
	Validate(ctx context.Context) *errors.ValidationErrors

	// ValidateWithInternalError will be used in priority if defined on a document
	// It will be called for all modifying operations (Create, Save, Update)
	// If it returns an internal error, the validation error will be nil.
	ValidateWithInternalError(ctx context.Context) (*errors.ValidationErrors, error)
}

var _ Validable = &Base{}

// Create inserts the document in the database, returns an error if document
// already exists and set CreatedAt timestamp
func Create(ctx context.Context, collectionName string, doc document) error {
	return save(ctx, collectionName, doc, func(ctx context.Context, collectionName string, doc document) error {
		log := logger.Get(ctx)
		c := mongo.Session(log).Clone().DB("").C(collectionName)
		defer c.Database.Session.Close()
		log.WithFields(logrus.Fields{
			"collection": collectionName,
			"doc_id":     doc.getID().Hex(),
		}).Debugf("save '%v'", collectionName)
		return c.Insert(doc)
	})
}

func Save(ctx context.Context, collectionName string, doc document) error {
	return save(ctx, collectionName, doc, func(ctx context.Context, collectionName string, doc document) error {
		log := logger.Get(ctx)
		c := mongo.Session(log).Clone().DB("").C(collectionName)
		defer c.Database.Session.Close()
		log.WithFields(logrus.Fields{
			"collection": collectionName,
			"doc_id":     doc.getID().Hex(),
		}).Debugf("save '%v'", collectionName)
		_, err := c.UpsertId(doc.getID(), doc)
		return err
	})
}

func Update(ctx context.Context, collectionName string, update bson.M, doc document) error {
	return save(ctx, collectionName, doc, func(ctx context.Context, collectionName string, doc document) error {
		log := logger.Get(ctx)
		c := mongo.Session(log).Clone().DB("").C(collectionName)
		defer c.Database.Session.Close()

		if _, ok := update["$set"]; ok {
			update["$set"].(bson.M)["updated_at"] = doc.getUpdatedAt()
		}

		log.WithFields(logrus.Fields{
			"collection": collectionName,
			"doc_id":     doc.getID().Hex(),
		}).Debugf("update %v", collectionName)
		return c.UpdateId(doc.getID(), update)
	})
}

func save(ctx context.Context, collectionName string, doc document, saveFunc func(context.Context, string, document) error) error {
	validationErrors, err := doc.ValidateWithInternalError(ctx)
	if err == ErrValidateNoInternalErrorFunc {
		validationErrors = doc.Validate(ctx)
	} else if err != nil {
		return errors.Wrap(ctx, err, "fail to validate document")
	}
	if validationErrors != nil {
		return validationErrors
	}

	doc.ensureID()
	doc.ensureCreatedAt()
	doc.setUpdatedAt(time.Now())

	return saveFunc(ctx, collectionName, doc)
}

// Destroy really deletes
func Destroy(ctx context.Context, collectionName string, doc destroyable) error {
	return doc.destroy(ctx, collectionName)
}

func ReallyDestroy(ctx context.Context, collectionName string, doc document) error {
	log := logger.Get(ctx)
	c := mongo.Session(log).Clone().DB("").C(collectionName)
	defer c.Database.Session.Close()
	log.WithFields(logrus.Fields{
		"collection": collectionName,
		"doc_id":     doc.getID().Hex(),
	}).Debugf("remove '%v'", collectionName)
	return c.RemoveId(doc.getID())
}

// Find is finding the model with objectid id in the collection name, with its
// default scope for paranoid documents, it won't look at documents tagged as
// deleted
func Find(ctx context.Context, collectionName string, id bson.ObjectId, doc scopable, sortFields ...SortField) error {
	query := doc.scope(bson.M{"_id": id})
	return find(ctx, collectionName, query, doc, sortFields...)
}

// FindUnscoped is similar as Find but does not care of the default scope of
// the document.
func FindUnscoped(ctx context.Context, collectionName string, id bson.ObjectId, doc interface{}, sortFields ...SortField) error {
	query := bson.M{"_id": id}
	return find(ctx, collectionName, query, doc, sortFields...)
}

func FindOne(ctx context.Context, collectionName string, query bson.M, doc scopable, sortFields ...SortField) error {
	return find(ctx, collectionName, doc.scope(query), doc, sortFields...)
}

func FindOneUnscoped(ctx context.Context, collectionName string, query bson.M, doc interface{}) error {
	return find(ctx, collectionName, query, doc)
}

func find(ctx context.Context, collectionName string, query bson.M, doc interface{}, sortFields ...SortField) error {
	log := logger.Get(ctx)
	c := mongo.Session(log).Clone().DB("").C(collectionName)
	defer c.Database.Session.Close()

	fields := make([]string, len(sortFields))
	for i, f := range sortFields {
		fields[i] = string(f)
	}
	return c.Find(query).Sort(fields...).One(doc)
}

func WhereQuery(ctx context.Context, collectionName string, query bson.M, sortFields ...SortField) (*mgo.Query, Closer) {
	if query == nil {
		query = bson.M{}
	}
	if _, ok := query["deleted_at"]; !ok {
		query["deleted_at"] = nil
	}

	return WhereUnscopedQuery(ctx, collectionName, query, sortFields...)
}

func WhereUnscopedQuery(ctx context.Context, collectionName string, query bson.M, sortFields ...SortField) (*mgo.Query, Closer) {
	log := logger.Get(ctx)
	c := mongo.Session(log).Clone().DB("").C(collectionName)

	if query == nil {
		query = bson.M{}
	}

	fields := make([]string, len(sortFields))
	for i, f := range sortFields {
		fields[i] = string(f)
	}
	return c.Find(query).Sort(fields...), c.Database.Session
}

func Where(ctx context.Context, collectionName string, query bson.M, data interface{}, sortFields ...SortField) error {
	mongoQuery, session := WhereQuery(ctx, collectionName, query, sortFields...)
	defer session.Close()
	err := mongoQuery.All(data)
	if err != nil {
		return fmt.Errorf("fail to query mongo %v: %v", query, err)
	}
	return nil
}

func WhereUnscoped(ctx context.Context, collectionName string, query bson.M, data interface{}, sortFields ...SortField) error {
	mongoQuery, session := WhereUnscopedQuery(ctx, collectionName, query, sortFields...)
	defer session.Close()
	err := mongoQuery.All(data)
	if err != nil {
		return fmt.Errorf("fail to query mongo %v: %v", query, err)
	}
	return nil
}

func Count(ctx context.Context, collectionName string, query bson.M) (int, error) {
	if query == nil {
		query = bson.M{}
	}
	if _, ok := query["deleted_at"]; !ok {
		query["deleted_at"] = nil
	}

	return CountUnscoped(ctx, collectionName, query)
}

func CountUnscoped(ctx context.Context, collectionName string, query bson.M) (int, error) {
	log := logger.Get(ctx)
	// nolint: contextcheck
	c := mongo.Session(log).Clone().DB("").C(collectionName)
	defer c.Database.Session.Close()

	if query == nil {
		query = bson.M{}
	}

	return c.Find(query).Count()
}

func WhereIter(ctx context.Context, collectionName string, query bson.M, fun func(*mgo.Iter) error, sortFields ...SortField) error {
	if query == nil {
		query = bson.M{}
	}
	if _, ok := query["deleted_at"]; !ok {
		query["deleted_at"] = nil
	}
	return WhereIterUnscoped(ctx, collectionName, query, fun, sortFields...)
}

func WhereIterUnscoped(ctx context.Context, collectionName string, query bson.M, fun func(*mgo.Iter) error, sortFields ...SortField) error {
	log := logger.Get(ctx)
	c := mongo.Session(log).Clone().DB("").C(collectionName)
	defer c.Database.Session.Close()

	if query == nil {
		query = bson.M{}
	}

	fields := make([]string, len(sortFields))
	for i, f := range sortFields {
		fields[i] = string(f)
	}
	iter := c.Find(query).Sort(fields...).Iter()
	defer iter.Close()

	err := fun(iter)
	if err != nil {
		return fmt.Errorf("error occurred during iteration over collection %v with query %v: %v", collectionName, query, err)
	}
	if iter.Err() != nil {
		return fmt.Errorf("fail to iterate over collection %v with query %v: %v", collectionName, query, iter.Err())
	}
	return nil
}

func EnsureParanoidIndices(ctx context.Context, collectionNames ...string) {
	log := logger.Get(ctx)

	for _, collectionName := range collectionNames {
		log = logger.Get(ctx).WithFields(logrus.Fields{
			"init":       "setup-indices",
			"collection": collectionName,
		})
		ctx = logger.ToCtx(ctx, log)
		log.Info("Setup the MongoDB index")

		c := mongo.Session(log).Clone().DB("").C(collectionName)
		defer c.Database.Session.Close()
		err := c.EnsureIndexKey("deleted_at")
		if err != nil {
			log.WithError(err).Error("fail to setup the deleted_at index")
			continue
		}
	}
}
