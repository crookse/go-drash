package services

import (
	"strings"
)

type SearchResult struct {
	Id         int
	Item       interface{}
	Query      string
	SearchTerm string
}

type IndexService struct {
	Cache       map[string][]SearchResult // e.g., ["query" []SearchResult]
	Index       map[string][]int       // e.g., ["query", [1,2,3,4,5]]
	LookupTable map[int]interface{}      // e.g., [1, SomeType]
}

func (i IndexService) AddItem(searchTerms []string, item interface{}) {
	id := len(i.LookupTable)

	i.LookupTable[id] = item

	for iSt := range searchTerms {
		query := searchTerms[iSt]
		var ids = i.Index[query]
		if ids == nil {
			ids = []int{}
		}
		ids = append(ids, id)
		i.Index[query] = ids
	}
}

func (i IndexService) Search(query string) []SearchResult {
	if i.Cache[query] != nil {
		return i.Cache[query]
	}

	results := []SearchResult{}

	for i1, ids := range i.Index {
		if strings.Contains(i1, query) {
			for i2 := range ids {
				result := SearchResult{
					Id: ids[i2],
					Item: i.LookupTable[ids[i2]],
					Query: query,
				}
				results = append(results, result)
			}
		}
	}
	
	i.Cache[query] = results

	return results
}
