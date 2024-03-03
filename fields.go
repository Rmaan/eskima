package eskima

type Keyword string

func (k *Keyword) ElasticSchema() (*ElasticSchema, error) {
	return &ElasticSchema{Type: "keyword"}, nil
}

// TODO these two types are not that useful to be in library code.

type Timestamp int64

func (m *Timestamp) ElasticSchema() (*ElasticSchema, error) {
	return &ElasticSchema{Type: "date"}, nil
}

type SecondTimestamp int64

func (s *SecondTimestamp) ElasticSchema() (*ElasticSchema, error) {
	return &ElasticSchema{Type: "date"}, nil
}
