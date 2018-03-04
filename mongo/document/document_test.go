package document

import (
	"context"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const DocsCollection = "docs"

type Doc struct {
	ID        bson.ObjectId `bson:"_id"`
	DeletedAt time.Time     `bson:"deleted_at,omitempty"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

func (d Doc) GetID() bson.ObjectId {
	return d.ID
}

func (d *Doc) SetDeletedAt(t time.Time) {
	d.DeletedAt = t
}

func (d *Doc) SetUpdatedAt(t time.Time) {
	d.UpdatedAt = t
}

func NewTestDoc(t *testing.T) (*Doc, func()) {
	d := Doc{ID: bson.NewObjectId()}
	require.NoError(t, Save(context.Background(), DocsCollection, &d))
	return &d, func() {
		require.NoError(t, Destroy(context.Background(), DocsCollection, &d))
	}
}

func TestFind(t *testing.T) {
	examples := []struct {
		Name  string
		Doc   func(t *testing.T) (*Doc, func())
		Error string
	}{
		{
			Name: "it should find existing doc",
			Doc: func(t *testing.T) (*Doc, func()) {
				d, clean := NewTestDoc(t)
				return d, clean
			},
		}, {
			Name: "it should not find unsaved doc",
			Doc: func(t *testing.T) (*Doc, func()) {
				return &Doc{ID: bson.NewObjectId()}, func() {}
			},
			Error: "not found",
		}, {
			Name: "it should not find deleted doc",
			Doc: func(t *testing.T) (*Doc, func()) {
				d, clean := NewTestDoc(t)
				err := ParanoiaDelete(context.Background(), DocsCollection, d)
				if err != nil {
					clean()
					require.NoError(t, err)
				}
				return d, clean
			},
			Error: "not found",
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			fixtureDoc, clean := example.Doc(t)
			defer clean()

			var d Doc
			err := Find(context.Background(), DocsCollection, fixtureDoc.ID, &d)
			if example.Error != "" {
				assert.Contains(t, err.Error(), example.Error)
			} else {
				require.NoError(t, err)
				require.Equal(t, fixtureDoc.ID, d.ID)
			}
		})
	}
}

func TestWhereParanoia(t *testing.T) {
	examples := []struct {
		Name  string
		Query bson.M
		Docs  func(t *testing.T) ([]*Doc, func())
		Count int
	}{
		{
			Name: "it should find existing documents",
			Docs: func(t *testing.T) ([]*Doc, func()) {
				d1, clean1 := NewTestDoc(t)
				d2, clean2 := NewTestDoc(t)
				return []*Doc{d1, d2}, func() {
					clean1()
					clean2()
				}
			},
			Count: 2,
		}, {
			Name: "it should not find paranoia-deleted documents",
			Docs: func(t *testing.T) ([]*Doc, func()) {
				d1, clean1 := NewTestDoc(t)
				err := ParanoiaDelete(context.Background(), DocsCollection, d1)
				require.NoError(t, err)
				d2, clean2 := NewTestDoc(t)
				err = ParanoiaDelete(context.Background(), DocsCollection, d2)
				require.NoError(t, err)
				return []*Doc{d1, d2}, func() {
					clean1()
					clean2()
				}
			},
			Count: 0,
		}, {
			Name:  "it should find deleted document, if queried specifically",
			Query: bson.M{"deleted_at": bson.M{"$exists": true}},
			Docs: func(t *testing.T) ([]*Doc, func()) {
				d1, clean1 := NewTestDoc(t)
				err := ParanoiaDelete(context.Background(), DocsCollection, d1)
				require.NoError(t, err)
				d2, clean2 := NewTestDoc(t)
				err = ParanoiaDelete(context.Background(), DocsCollection, d2)
				return []*Doc{d1, d2}, func() {
					clean1()
					clean2()
				}
			},
			Count: 2,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			_, clean := example.Docs(t)
			defer clean()

			query := bson.M{}
			if example.Query != nil {
				query = example.Query
			}
			var docs []*Doc
			err := WhereParanoia(context.Background(), DocsCollection, query, &docs)
			require.NoError(t, err)
			require.Len(t, docs, example.Count)
		})
	}
}
