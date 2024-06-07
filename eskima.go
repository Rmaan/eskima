package eskima

import (
	"fmt"
	"reflect"
	"strings"
)

type KNNMethod struct {
	Name          string `json:"name,omitempty"`
	SpaceType     string `json:"space_type,omitempty"`
	Engine        string `json:"engine,omitempty"`
	KNNParameters `json:"parameters,omitempty"`
}

type KNNParameters struct {
	EfConstruction int `json:"ef_construction,omitempty"`
	M              int `json:"m,omitempty"`
}

type ElasticSchema struct {
	Type       string                    `json:"type,omitempty"`
	Properties map[string]*ElasticSchema `json:"properties,omitempty"`

	Analyzer       string     `json:"analyzer,omitempty"`
	Coerce         *bool      `json:"coerce,omitempty"`
	CopyTo         []string   `json:"copy_to,omitempty"`
	DocValues      *bool      `json:"doc_values,omitempty"`
	Enabled        *bool      `json:"enabled,omitempty"`
	Format         string     `json:"format,omitempty"`
	Index          *bool      `json:"index,omitempty"`
	SearchAnalyzer string     `json:"search_analyzer,omitempty"`
	Dimension      *int       `json:"dimension,omitempty"`
	Method         *KNNMethod `json:"method,omitempty"`
}

func (s *ElasticSchema) Get(path ...string) *ElasticSchema {
	current := s
	for _, p := range path {
		if current == nil {
			return nil
		}
		current = current.Properties[p]
	}
	return current
}

type ElasticSchemaer interface {
	ElasticSchema() (*ElasticSchema, error)
}

func Generate(t any) (*ElasticSchema, error) {
	return generate(reflect.TypeOf(t), make([]string, 0, 10))
}

func generate(t reflect.Type, path []string) (*ElasticSchema, error) {
	if eskimaer, ok := reflect.NewAt(t, nil).Interface().(ElasticSchemaer); ok {
		s, err := eskimaer.ElasticSchema()
		if err != nil {
			return nil, fmt.Errorf("at %s: error from ElasticSchemaer: %w", strings.Join(path, "."), err)
		}
		return s, nil
	}

	switch t.Kind() {
	case reflect.Struct:
		return generateStruct(t, path)
	case reflect.Ptr:
		return generate(t.Elem(), path)
	case reflect.String:
		return &ElasticSchema{Type: "text"}, nil
	case reflect.Int:
		// 32 bit is a good default for int
		return &ElasticSchema{Type: "integer"}, nil
	case reflect.Int8:
		return &ElasticSchema{Type: "byte"}, nil
	case reflect.Int16:
		return &ElasticSchema{Type: "short"}, nil
	case reflect.Int32:
		return &ElasticSchema{Type: "integer"}, nil
	case reflect.Int64:
		return &ElasticSchema{Type: "long"}, nil
	case reflect.Uint64:
		return &ElasticSchema{Type: "unsigned_long"}, nil
	case reflect.Float64, reflect.Float32:
		// TODO should they be separated? Usage of float64 for everything is commonplace but not needed in indexing.
		return &ElasticSchema{Type: "float"}, nil
	case reflect.Bool:
		return &ElasticSchema{Type: "boolean"}, nil
	case reflect.Slice:
		return generate(t.Elem(), path)
	default:
		return nil, fmt.Errorf("at %s: unsupported type %s", strings.Join(path, "."), t.Kind())
	}
}

func generateStruct(t reflect.Type, path []string) (*ElasticSchema, error) {
	properties := map[string]*ElasticSchema{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if commaIndex := strings.Index(jsonTag, ","); commaIndex != -1 {
			jsonTag = jsonTag[:commaIndex]
		}
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		eskimaTag := field.Tag.Get("eskima")
		if strings.HasPrefix(eskimaTag, "-") {
			continue
		}
		if strings.HasPrefix(eskimaTag, "!") {
			f := false
			properties[jsonTag] = &ElasticSchema{Type: "object", Enabled: &f}
			continue
		}
		if eskimaTag != "" {
			properties[jsonTag] = &ElasticSchema{Type: eskimaTag}
			continue
		}

		inner, err := generate(field.Type, append(path, jsonTag))
		if err != nil {
			return nil, err
		}
		properties[jsonTag] = inner
	}
	return &ElasticSchema{Properties: properties}, nil
}
