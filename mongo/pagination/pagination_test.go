package pagination

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"

	"github.com/Scalingo/go-utils/mongo/document"
)

const dummyCollection = "dummy_documents"

type dummyDocument struct {
	document.Paranoid  `bson:",inline"`
	AppID              string `bson:"app_id" json:"app_id"`
	VirtualStorageName string `bson:"virtual_storage_name" json:"virtual_storage_name"`
}

func newDummyDocuments(t *testing.T, virtualStorageName string, amount int) func() {
	return newDummyDocumentsIncludingDeleted(t, virtualStorageName, amount, 0)
}

func newDummyDocumentsIncludingDeleted(t *testing.T, virtualStorageName string, amount int, amountDeleted int) func() {
	d := make([]*dummyDocument, amount)

	for i := 0; i < amount; i++ {
		doc := dummyDocument{
			AppID:              strconv.Itoa(i),
			VirtualStorageName: virtualStorageName,
		}
		require.NoError(t, document.Save(context.Background(), dummyCollection, &doc))
		d[i] = &doc
	}
	for i := 0; i < amountDeleted; i++ {
		require.NoError(t, document.Destroy(context.Background(), dummyCollection, d[i]))
	}
	return func() {
		for _, doc := range d {
			require.NoError(t, document.ReallyDestroy(context.Background(), dummyCollection, doc))
		}
	}
}

func TestPaginationPaginate(t *testing.T) {
	runs := []struct {
		Name              string
		DummyDocument     func(t *testing.T) func()
		PaginationOpts    *ServiceOpts
		PageNumber        int
		AmountItems       int
		SortOrder         string
		ExpectedQuery     bson.M
		ExpectedMeta      func() Meta
		ExpectedResult    []dummyDocument
		ExpectedQueryFunc QueryFunc
		Error             string
	}{
		{
			Name: "It should return an empty list when the requested PageNumber is out of range",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 4)
				return clean
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 2,
				MaxPerPage:     2,
			},
			AmountItems:    0,
			PageNumber:     3,
			ExpectedResult: []dummyDocument{},
			ExpectedMeta: func() Meta {
				prevPage := 2
				return Meta{
					CurrentPage: 3,
					PrevPage:    &prevPage,
					NextPage:    nil,
					TotalPages:  2,
					TotalCount:  4,
					perPageNum:  2,
				}
			},
		},
		{
			Name: "It should return an error when the requested AmountItems is greater than MaxPerPage",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 4)
				return clean
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 2,
				MaxPerPage:     2,
			},
			PageNumber:  0,
			AmountItems: 3,
			Error:       "* per_page → must be lower or equal to 2",
		},
		{
			Name: "It should return an error when the requested AmountItems is lower than 0",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 4)
				return clean
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 2,
				MaxPerPage:     2,
			},
			PageNumber:  0,
			AmountItems: -1,
			Error:       "* per_page → must be greater than 0",
		},
		{
			Name: "It should return an error when the requested item per page is greater than max per page",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 4)
				return clean
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 2,
				MaxPerPage:     2,
			},
			PageNumber:  0,
			AmountItems: 3,
			Error:       "* per_page → must be lower or equal to 2",
		},
		{
			Name: "It should return an error with a perPageDefault lower or equal to 0",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 1)
				return clean
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 0,
				MaxPerPage:     2,
			},
			PageNumber:  0,
			AmountItems: 0,
			Error:       "invalid pagination service configuration: MaxPerPage > PerPageDefault > 0",
		},
		{
			Name: "It should return an error with a perPageDefault greater MaxPerPage",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 1)
				return clean
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 10,
				MaxPerPage:     2,
			},
			PageNumber:  0,
			AmountItems: 0,
			Error:       "invalid pagination service configuration: MaxPerPage > PerPageDefault > 0",
		},
		{
			Name: "It should return an error with a MaxPerPage lower or equal to 0",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 1)
				return clean
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 2,
				MaxPerPage:     -1,
			},
			PageNumber:  0,
			AmountItems: 0,
			Error:       "invalid pagination service configuration: MaxPerPage > PerPageDefault > 0",
		},
		{
			Name: "It should return an error with a requested page lower or equal to 0",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 1)
				return clean
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 5,
				MaxPerPage:     15,
			},
			PageNumber:  -1,
			AmountItems: 0,
			Error:       "* Requested page → must be greater than 0",
		},
		{
			Name: "It should return an empty result array with Meta object nil",
			DummyDocument: func(t *testing.T) func() {
				return func() {}
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 5,
				MaxPerPage:     15,
			},
			PageNumber:  0,
			AmountItems: 0,
			ExpectedMeta: func() Meta {
				return Meta{
					CurrentPage: 0,
					PrevPage:    nil,
					NextPage:    nil,
					TotalPages:  0,
					TotalCount:  0,
				}
			},
			ExpectedResult: []dummyDocument{},
		},
		{
			Name: "It should return the only one element",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 1)
				return clean
			},
			ExpectedMeta: func() Meta {
				return Meta{
					CurrentPage: 1,
					PrevPage:    nil,
					NextPage:    nil,
					TotalPages:  1,
					TotalCount:  1,
					perPageNum:  1,
				}
			},
			ExpectedResult: []dummyDocument{
				{AppID: "0", VirtualStorageName: "vs_name_1"},
			},
		},
		{
			Name: "With two document, it should return only one element thanks to the query",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 1)
				clean2 := newDummyDocuments(t, "vs_name_2", 1)
				return func() {
					clean()
					clean2()
				}
			},
			ExpectedMeta: func() Meta {
				return Meta{
					CurrentPage: 1,
					PrevPage:    nil,
					NextPage:    nil,
					TotalPages:  1,
					TotalCount:  1,
					perPageNum:  1,
				}
			},
			ExpectedResult: []dummyDocument{
				{AppID: "0", VirtualStorageName: "vs_name_2"},
			},
			ExpectedQuery: bson.M{"virtual_storage_name": "vs_name_2"},
		},
		{
			Name: "It should return 2 elements",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 1)
				clean2 := newDummyDocuments(t, "vs_name_2", 2)
				return func() {
					clean()
					clean2()
				}
			},
			ExpectedMeta: func() Meta {
				nextPage := 2
				return Meta{
					CurrentPage: 1,
					PrevPage:    nil,
					NextPage:    &nextPage,
					TotalPages:  2,
					TotalCount:  2,
					perPageNum:  1,
				}
			},
			ExpectedResult: []dummyDocument{
				{AppID: "0", VirtualStorageName: "vs_name_2"},
			},
			ExpectedQuery: bson.M{"virtual_storage_name": "vs_name_2"},
		},
		{
			Name: "It should return the second page",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 2)
				clean2 := newDummyDocuments(t, "vs_name_2", 4)
				return func() {
					clean()
					clean2()
				}
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 2,
				MaxPerPage:     2,
			},
			PageNumber:  2,
			AmountItems: 0,
			ExpectedMeta: func() Meta {
				prevPage := 1
				return Meta{
					CurrentPage: 2,
					PrevPage:    &prevPage,
					NextPage:    nil,
					TotalPages:  2,
					TotalCount:  4,
					perPageNum:  2,
				}
			},
			ExpectedResult: []dummyDocument{
				{AppID: "2", VirtualStorageName: "vs_name_2"},
				{AppID: "3", VirtualStorageName: "vs_name_2"},
			},
			ExpectedQuery: bson.M{"virtual_storage_name": "vs_name_2"},
		},
		{
			Name: "It should return the second page in reverse order",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 2)
				clean2 := newDummyDocuments(t, "vs_name_2", 4)
				return func() {
					clean()
					clean2()
				}
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 2,
				MaxPerPage:     2,
			},
			PageNumber:  2,
			AmountItems: 0,
			ExpectedMeta: func() Meta {
				prevPage := 1
				return Meta{
					CurrentPage: 2,
					PrevPage:    &prevPage,
					NextPage:    nil,
					TotalPages:  2,
					TotalCount:  4,
					perPageNum:  2,
				}
			},
			SortOrder: "-_id",
			ExpectedResult: []dummyDocument{
				{AppID: "1", VirtualStorageName: "vs_name_2"},
				{AppID: "0", VirtualStorageName: "vs_name_2"},
			},
			ExpectedQuery: bson.M{"virtual_storage_name": "vs_name_2"},
		},
		{
			Name: "it should not return soft-deleted documents by default",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocumentsIncludingDeleted(t, "vs_name_1", 2, 1)
				return func() {
					clean()
				}
			},
			ExpectedMeta: func() Meta {
				return Meta{
					CurrentPage: 1,
					PrevPage:    nil,
					NextPage:    nil,
					TotalPages:  1,
					TotalCount:  1,
					perPageNum:  1,
				}
			},
			ExpectedResult: []dummyDocument{
				{AppID: "1", VirtualStorageName: "vs_name_1"},
			},
			ExpectedQuery: bson.M{"virtual_storage_name": "vs_name_1"},
		},
		{
			Name: "With custom QueryFunc to document.WhereQueryUnscoped, it should return soft-deleted documents",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocumentsIncludingDeleted(t, "vs_name_1", 2, 1)
				return func() {
					clean()
				}
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 10,
				MaxPerPage:     10,
			},
			AmountItems: 10,
			ExpectedMeta: func() Meta {
				return Meta{
					CurrentPage: 1,
					PrevPage:    nil,
					NextPage:    nil,
					TotalPages:  1,
					TotalCount:  2,
					perPageNum:  10,
				}
			},
			ExpectedResult: []dummyDocument{
				{AppID: "0", VirtualStorageName: "vs_name_1"},
				{AppID: "1", VirtualStorageName: "vs_name_1"},
			},
			ExpectedQuery:     bson.M{"virtual_storage_name": "vs_name_1"},
			ExpectedQueryFunc: document.WhereUnscopedQuery,
		},
		{
			Name: "With custom QueryFunc to document.WhereQueryUnscoped, it should respect custom sort order",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocumentsIncludingDeleted(t, "vs_name_1", 3, 1)
				return func() {
					clean()
				}
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 10,
				MaxPerPage:     10,
			},
			AmountItems: 10,
			SortOrder:   "-app_id",
			ExpectedMeta: func() Meta {
				return Meta{
					CurrentPage: 1,
					PrevPage:    nil,
					NextPage:    nil,
					TotalPages:  1,
					TotalCount:  3,
					perPageNum:  10,
				}
			},
			ExpectedResult: []dummyDocument{
				{AppID: "2", VirtualStorageName: "vs_name_1"},
				{AppID: "1", VirtualStorageName: "vs_name_1"},
				{AppID: "0", VirtualStorageName: "vs_name_1"},
			},
			ExpectedQuery:     bson.M{"virtual_storage_name": "vs_name_1"},
			ExpectedQueryFunc: document.WhereUnscopedQuery,
		},
	}

	for _, run := range runs {
		t.Run(run.Name, func(t *testing.T) {
			clean := run.DummyDocument(t)
			defer clean()

			results := []*dummyDocument{}

			if run.PaginationOpts == nil {
				run.PaginationOpts = &ServiceOpts{
					PerPageDefault: 1,
					MaxPerPage:     1,
				}
			}

			if run.ExpectedQuery == nil {
				run.ExpectedQuery = bson.M{"virtual_storage_name": "vs_name_1"}
			}

			paginateOpts := PaginateOpts{
				PageNumber:  run.PageNumber,
				AmountItems: run.AmountItems,
				Query:       run.ExpectedQuery,
				SortOrder:   run.SortOrder,
				QueryFunc:   run.ExpectedQueryFunc,
			}

			meta, err := run.PaginationOpts.Paginate(context.Background(),
				dummyCollection, &results, paginateOpts)
			if run.Error != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), run.Error)
			} else {
				require.NoError(t, err)

				expectedMeta := run.ExpectedMeta()
				assert.Equal(t, expectedMeta, meta)
				assert.Equal(t, len(run.ExpectedResult), len(results))

				for idx, expectedResult := range run.ExpectedResult {
					require.Equal(t, expectedResult.VirtualStorageName, results[idx].VirtualStorageName)
					require.Equal(t, expectedResult.AppID, results[idx].AppID)
				}
			}
		})
	}
}
