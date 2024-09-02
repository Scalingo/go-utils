## Package `pagination` v0

This is a pagination library for Go.

The `PageRequest` struct is used to define the page number and the page size when requesting a page of items.

The `Paginated` struct is used to return a page of items along with the total number of items.

// Example data returned
```json 
{
  "data": [],
  "meta": {
    "current_page": 1,
    "next_page": 1,
    "per_page": 20,
    "prev_page": 1,
    "total_count": 0,
    "total_pages": 1
  }
}
```


### Usage

```go
package main

import (
    "fmt"
    "github.com/Scalingo/go-utils/pagination"
)

type Item struct {
    ID   int
    Name string
}

func main() {
	items := []Item{
		{ID: 1, Name: "Item 1"},
		{ID: 2, Name: "Item 2"},
		{ID: 3, Name: "Item 3"},
		{ID: 4, Name: "Item 4"},
	}
	
    p := pagination.NewPaginated[[]Item](items, pagination.NewPageRequest(1, 2), 4)
    fmt.Println(p)
}
```