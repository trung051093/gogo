package usermodel

import (
	"context"
	"fmt"
	"user_management/common"
	esprovider "user_management/components/elasticsearch"
)

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
"sort" : [ { "%s" : "%s" } ]
`

type UserEsQuery struct {
	Query     string
	LastIndex string // for pagination
	Paging    *common.Pagination
	Filter    *UserFilter
}

func GetUserESQuery(ctx context.Context, userEsQuery *UserEsQuery) string {
	q := fmt.Sprintf(
		ElasticSearchQuery,
		userEsQuery.Query,
		userEsQuery.Paging.Limit,
		userEsQuery.Filter.SortField,
		userEsQuery.Filter.SortName,
	)
	return esprovider.BuildQuery(
		ctx,
		q,
		userEsQuery.LastIndex,
	)
}
