package room

import "github.com/olivere/elastic/v7"

const (
	defaultSize  = 100
	idField      = "id"
	nameField    = "name"
	defaultIndex = "rooms"
)

type Request struct {
	ids      []int32
	names    []string
	size     int32
	password string
	index    string
}

func New() *Request {
	return &Request{
		size:  defaultSize,
		index: defaultIndex,
	}
}

func (r *Request) Ids(ids []int32) *Request {
	r.ids = ids
	return r
}

func (r *Request) Size(size int32) *Request {
	r.size = size
	return r
}

func (r *Request) Names(names []string) *Request {
	r.names = names
	return r
}

func (r *Request) Index(index string) *Request {
	r.index = index
	return r
}

func (r *Request) Password(pass string) *Request {
	r.password = pass
	return r
}

func (r *Request) buildSearchQuery() *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	if len(r.ids) > 0 {
		query.Must(elastic.NewTermsQuery(idField, r.ids))
	}
	if len(r.names) > 0 {
		query.Must(elastic.NewTermsQuery(nameField, r.names))
	}
	return query
}
