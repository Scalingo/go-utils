# Mongo Pagination

Pagination is a package providing a function (`Paginate`) that queries the mongo
database and return results by pages. It is useful to avoid the processing of
too much data.

The pagination follow the rules explained in the
[Scalingo pagination documentation](https://developers.scalingo.com/index#pagination)

### How to use it

First, set the default value of number of items per page and the maximum of item
by page as the following snippet:

```Go
pageService := pagination.NewPaginationService(pagination.ServiceOpts{
	PerPageDefault: 5,
	MaxPerPage:     15,
})
```

Then call the `Paginate` as follow:
```Go
// resultObject will be filled with the query result.
resultObject := []*ResultObject{}
dbQuery := bson.M{"searched_field": "field_1"}

paginateOpts := PaginateOpts{
    PageNumber:  5,         // Number of the requested page (can be empty, default 1)
    AmountItems: 10,        // Amount of items by page (can be empty)
    Query:       dbQuery,   // Query which will be executed on the database (can be nil)
    SortOrder:   "-_id",    // The field for the sort order (by default "_id")
}

meta, err := pageService.Paginate(
    ctx,                    // A context (required)
    "RequestedCollection",  // Name of the collection (required)
    &resultObject,          // The object that will contain the data (must be an array)
    paginateOpts)
```

The returned meta object contains pagination metadata that could be used in the
request answer.

## Pagination by Cursor

This package also provide a pagination by cursor, could be useful in the case of
frequently updated data.

It follows the mongoDB pagination using range queries, explained on
[their documentation](https://docs.mongodb.com/manual/reference/method/cursor.skip/#using-range-queries).

This pagination return results lower than the cursor value passed at parameter.
Indeed, the `$lt` mongo operator is used to compare with the cursor value. So
cursor value must be comparable with this operator.

### How to use it

First, set the default value of number of items per page and the maximum of item
by page as the following snippet:

```Go
pageService := pagination.NewPaginationService(pagination.ServiceOpts{
	PerPageDefault: 5,
	MaxPerPage:     15,
})
```

Then call the `PaginateByCursor` described bellow:

```Go
// resultObject will be filled with the query result.
resultObject := []*ResultObject{}
dbQuery := bson.M{"searched_field": "field_1"}

paginateOpts := PaginateByCursorOpts{
    CursorKey   ""      // The mongoDB document field of the cursor (can be empty)
    CursorValue nil     // The cursor it self, it must be comparable with `$lt` mongo operator (can be nil)
    AmountItems 5       // Amount of items by page (can be empty)
    Query       dbQuery // Query which will be executed on the database (can be nil)
    SortOrder   "-_id"   // The field for the sort order (by default "-_id")
}

err := pageService.PaginateByCursor(
    ctx,                    // A context (required)
    "RequestedCollection",  // Name of the collection (required)
    &resultObject,          // The object that will contain the data (must be an array)
    paginateOpts)
// Here `resultObject` contains the first page

// To get the next page, you must provide the cursor in the `PaginateByCursorOpts`
// object, in accordance with the order of the data.
// The following is an example to retrieve the second page in a case of reverse ordered data.
resultObjectLen := len(resultObjectLen)
paginateOpts := PaginateByCursorOpts{
    CursorKey   "_id"
    CursorValue resultObject[resultObjectLen].ID
    AmountItems 5
    Query       dbQuery
    SortOrder   "-_id"
}

err := pageService.PaginateByCursor(
    ctx,
    "RequestedCollection",
    &resultObject,
    paginateOpts)
// So here `resultObject` contains the second page

```
