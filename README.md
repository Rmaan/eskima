# eskima
Generate Elasticsearch/OpenSearch schema from your Go structs

## Basic usage
```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Rmaan/eskima"
)

type Book struct {
	// 🧑‍🔧 Eskima will automatically generate schema for everything with a json tag
	Title string `json:"title,omitempty"`
	// 📚 All strings will be text by default, you can override type with eskima tag
	ISBN string `json:"isbn,omitempty" eskima:"keyword"`
	// 😌 Structs/Slices/Pointers will be followed automatically
	Authors []Author `json:"author,omitempty"`
}

type Author struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty" eskima:"byte"`
}

func main() {
	schema, err := eskima.Generate(Book{})
	if err != nil {
		log.Fatal(err)
	}
	buf, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))
}

```

Output:
```json
{
  "properties": {
    "author": {
      "properties": {
        "age": {
          "type": "byte"
        },
        "name": {
          "type": "text"
        }
      }
    },
    "isbn": {
      "type": "keyword"
    },
    "title": {
      "type": "text"
    }
  }
}
```

## Excluding fields
To remove a field from schema, just remove it from JSON
```go
type Book struct {
	InternalID int `json:"-"`
}
```

If it's needed in JSON but not in schema, you can use alternative syntax
```go
type Book struct {
	InternalID int `json:"internalID" eskima:"-"`
}
```

To disable a field, i.e. `{"enabled": false}`, use eskima tag "!"
```go
type Book struct {
	Comments []Comment `json:"comments" eskima:"!"`
}
```

Will emit:
```json
{
  "properties": {
    "comments": {
      "type": "object",
      "enabled": false
    }
  }
}
```

## Custom types
You can implement `eskima.ElasticSchemaer` to have custom control over schema
generation.

For example, we can this Keyword to remove the need to specify type in eskima
tag.
```go
type Keyword string

func (k *Keyword) ElasticSchema() (*ElasticSchema, error) {
	return &ElasticSchema{Type: "keyword"}, nil
}

type Book struct {
    Publisher Keyword `json:"publisher,omitempty"`
}
```
