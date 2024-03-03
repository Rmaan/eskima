package eskima

type Keyword string

func (k *Keyword) ElasticSchema() (*ElasticSchema, error) {
	return &ElasticSchema{Type: "keyword"}, nil
}
