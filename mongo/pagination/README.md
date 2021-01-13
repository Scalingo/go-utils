# Mongo Pagination

Pagination is a package providing a function (`Paginate`) that queries the mongo
database and return results by pages. It is useful to avoid the processing of
too much data.

The pagination follow the rules explained in the
[Scalingo pagination documentation](https://developers.scalingo.com/index#pagination)

## How to use it

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
resultObject := []*ResultObject{}
dbQuery := bson.M{"searched_field": "field_1"}

meta, err := pageService.Paginate(
    ctx,                    // A context (required)
    "5",                    // Number of the requested page (can be empty, default 1)
    "10",                   // Amount of items by page (can be empty)
    dbQuery,                // Query which will be executed on the database (can be nil)
    "RequestedCollection",  // Name of the collection (required)
    &resultObject,          // The object that will contain the data (must be an array)
    "-_id")                 // The field for the sort order (required)
```

The returned meta object contains pagination metadata that could be used in the
request answer.
