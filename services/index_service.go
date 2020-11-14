package services

import (
	"regexp"
	"strings"

	"github.com/drashland/go-drash/http"
)

type SearchResult struct {
	ID         int
	Item       interface
	Query      string
	SearchTerm string
}

type IndexService struct {
	Cache       map[int][]SearchResult // e.g., [1, []SearchResult]
	Index       map[string][]int       // e.g., ["query", [1,2,3,4,5]]
	LookupTable map[int]interface      // e.g., [1, SomeType]
}

func (i IndexService) AddItem(searchTerms []string, item T) {
	id := len(i.LookupTable)

	i.LookupTable[id] = item

	for i := range searchTerms {
		ids := i.Index[searchTerms[i]]
		if ids == nil {
			ids := []int{}
		}
		ids = append(ids, id)
		s.Index[searchTerm] = ids
	}
}

func (i IndexService) Search(query string) map[int][]SearchResult {
	if i.Cache[query] != nil {
		return i.Cache[query]
	}

	results map[int][]SearchResult
	for key, ids := range i.Index {
		if key.Contains(query) {
			for i, id := range ids {
				results[id] = SearchResult{
					Id: id,
					Item: i.LookupTable[id]
					SearchTerm: key,
					Query: query,
				}
			}
		}
	}

	i.Cache[query] = results

	return results
}
