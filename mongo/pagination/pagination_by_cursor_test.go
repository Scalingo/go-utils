package pagination

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
)

func TestPaginationPaginateByCursor(t *testing.T) {
	runs := []struct {
		Name           string
		DummyDocument  func(t *testing.T) func()
		PaginationOpts *ServiceOpts
		CursorKey      string
		CursorValue    int
		AmountItems    int
		SortOrder      string
		ExpectedQuery  bson.M
		ExpectedResult []dummyDocument
		Error          string
	}{
		{
			Name: "It should return an error with a request item per page superior then max per page",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 4)
				return clean
			},
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 2,
				MaxPerPage:     2,
			},
			AmountItems: 3,
			Error:       "* per_page → must be between 0 and 2",
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
			AmountItems: 0,
			Error:       "invalid pagination service configuration: MaxPerPage > PerPageDefault > 0",
		},
		{
			Name: "It should return an empty result with an invalid cursor",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 1)
				return clean
			},
			CursorValue: -1,
			PaginationOpts: &ServiceOpts{
				PerPageDefault: 5,
				MaxPerPage:     15,
			},
			AmountItems:    0,
			ExpectedResult: []dummyDocument{},
		},
		{
			Name: "It should return the only one element",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 1)
				return clean
			},
			ExpectedResult: []dummyDocument{
				{AppID: 0, VirtualStorageName: "vs_name_1"},
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
			ExpectedResult: []dummyDocument{
				{AppID: 0, VirtualStorageName: "vs_name_2"},
			},
			ExpectedQuery: bson.M{"virtual_storage_name": "vs_name_2"},
		},
		{
			Name: "Without cursor parameter, it should return the first page",
			DummyDocument: func(t *testing.T) func() {
				clean := newDummyDocuments(t, "vs_name_1", 1)
				clean2 := newDummyDocuments(t, "vs_name_2", 2)
				return func() {
					clean()
					clean2()
				}
			},
			SortOrder: "_id",
			ExpectedResult: []dummyDocument{
				{AppID: 0, VirtualStorageName: "vs_name_2"},
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
			CursorKey:   "app_id",
			CursorValue: 2,
			AmountItems: 0,
			ExpectedResult: []dummyDocument{
				{AppID: 1, VirtualStorageName: "vs_name_2"},
				{AppID: 0, VirtualStorageName: "vs_name_2"},
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
			CursorKey:   "app_id",
			CursorValue: 2,
			AmountItems: 0,
			SortOrder:   "-_id",
			ExpectedResult: []dummyDocument{
				{AppID: 1, VirtualStorageName: "vs_name_2"},
				{AppID: 0, VirtualStorageName: "vs_name_2"},
			},
			ExpectedQuery: bson.M{"virtual_storage_name": "vs_name_2"},
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

			paginateByCursorOpts := PaginateByCursorOpts{
				CursorKey:   run.CursorKey,
				AmountItems: run.AmountItems,
				Query:       run.ExpectedQuery,
				SortOrder:   run.SortOrder,
			}
			if run.CursorValue != 0 {
				paginateByCursorOpts.CursorValue = run.CursorValue
			}

			err := run.PaginationOpts.PaginateByCursor(context.Background(),
				dummyCollection, &results, paginateByCursorOpts)
			if run.Error != "" {
				assert.NotNil(t, err)
				if err == nil {
					return
				}
				assert.Contains(t, err.Error(), run.Error)
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(run.ExpectedResult), len(results))

				if len(run.ExpectedResult) == len(results) {
					for idx, expectedResult := range run.ExpectedResult {
						require.Equal(t, expectedResult.VirtualStorageName, results[idx].VirtualStorageName)
						require.Equal(t, expectedResult.AppID, results[idx].AppID)
					}
				}
			}
		})
	}
}
