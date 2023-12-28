package usermodel

import (
	"fmt"
	"gogo/common"
	esprovider "gogo/components/elasticsearch"
	elasticsearchmodel "gogo/components/elasticsearch/model"
)

type UserEsSearchResult struct {
	elasticsearchmodel.SearchResults
}

type UserEsSearchDto struct {
	*common.PagePagination
	Query     string `json:"-" query:"query"`
	LastIndex string // for pagination
}

const ElasticSearchQuery = `
"query" : {
	"multi_match" : {
		"query" : %q,
		"fields" : ["firstName^10", "lastName^10", "email", "company", "address^100"],
		"operator" : "and"
	}
},
"highlight" : {
	"fields" : {
		"firstName" : { "number_of_fragments" : 0 },
		"lastName" : { "number_of_fragments" : 0 },
		"email" : { "number_of_fragments" : 5, "fragment_size" : 25 }
	}
},
"size" : %d,
`

func (userEsQuery *UserEsSearchDto) ToQuery() string {
	q := fmt.Sprintf(
		ElasticSearchQuery,
		userEsQuery.Query,
		userEsQuery.Limit,
	)
	return esprovider.BuildQuery(
		q,
		userEsQuery.LastIndex,
	)
}
