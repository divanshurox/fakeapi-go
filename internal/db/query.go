package db

type Query struct {
	Database   string
	Collection string
	Filter     interface{}
	Object     interface{}
}

func NewQuery() *Query {
	return &Query{}
}

func (query *Query) WithDatabase(database string) *Query {
	query.Database = database
	return query
}

func (query *Query) WithCollection(collection string) *Query {
	query.Collection = collection
	return query
}

func (query *Query) WithFilter(filter interface{}) *Query {
	query.Filter = filter
	return query
}

func (query *Query) WithObject(object interface{}) *Query {
	query.Object = object
	return query
}
