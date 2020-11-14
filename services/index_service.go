package services

import (
	"strings"
)

type IndexService struct {
	Cache       map[string][]IndexServiceSearchResult
	Index       map[string][]int
	LookupTable map[int]interface{}
}

type IndexServiceSearchResult struct {
	Id         int
	Item       interface{}
	Query      string
	SearchTerm string
}

///////////////////////////////////////////////////////////////////////////////
// FILE MARKER - METHODS - EXPORTED ///////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// This adds an item to the index -- allowing them to be searchable by the
// given search terms. When a search is conducted via .Search(), the query
// passed into the .Search() method is used to match against all search terms.
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

// This searches the index for items that match the given query.
func (i IndexService) Search(query string) []IndexServiceSearchResult {
	if i.Cache[query] != nil {
		return i.Cache[query]
	}

	results := []IndexServiceSearchResult{}

	for i1, ids := range i.Index {
		if strings.Contains(i1, query) {
			for i2 := range ids {
				result := IndexServiceSearchResult{
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
