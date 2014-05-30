package datastore

import (
  "strconv"
  "encoding/json"
)

// ----------------------------------------------------------------------------
//  Name: searchCollection
//  Desc:

func searchCollection(name string, params map[string]string) (string, string, string, error) {

  // Get the search results.
  results, limit, next, prev, totalCount, err := searchResults(name, params)

  if next != "" {
    next = offsetResultPage(name, limit, next)
  }

  if prev != "" {
    prev = offsetResultPage(name, limit, prev)
  }

  body := listResponse(ResponseBody{len(results), results, next, prev, totalCount})

  // Return json response body, next, error
  return body, next, prev, err
}

// ----------------------------------------------------------------------------
//  Name: searchResults
//  Desc:

func searchResults(name string, params map[string]string) ([]Result, int, string, string, int, error) {
  // TraceMsg("==> Datastore:searchResults")

  var r       Ref
  var value   ResultValue
  var results []Result

  // query param
  queryStr, _ := params["query"]
  queryID := queryName(name, queryStr)

  // limit param
  limit := ListLimit(params)

  // offset param - to specify start of key range
  offset := 0
  if offsetStr, ok := params["offset"]; ok {
    offset, _ = strconv.Atoi(offsetStr)
  }

  // Determine whether there are stored results for this query.
  qresults := getItems("AllQueryResults", queryID)
  totalCount := len(qresults)

  // Generate new results when either: 1) offset is 0, or
  //                                   2) there are no stored results.
  if offset == 0 || totalCount == 0 {
    totalCount, _ = newSearchResults(name, queryStr)
    qresults = getItems("AllQueryResults", queryID)
  }

  // The table that contains the data, indexed by ref value.
  refTable := getCollectionHandle("RefTable")

  next := ""
  count, nextOffset := 0, 0

  for i, ref := range qresults {

    nextOffset++

    if i < offset {
      continue
    }

    // Unmarshal the json data/metadata extracted from the ref table
    r, value = Ref{}, ResultValue{}
    rdata, _ := refTable.Get(ro, []byte(ref))
    json.Unmarshal(rdata, &r)
    json.Unmarshal([]byte(r.Value), &value)

    // Create new result and append to results.
    path := PathValue{r.Mdata.Collection, r.Mdata.Key, "", r.Mdata.Ref, 0, 0, false}
    results = append(results, Result{path, value, 0, 0, 0})
    count++

    // Check limit, set "next" offset if limit reached before end of range.
    if count == limit {
      if nextOffset < totalCount {
        next = strconv.Itoa(nextOffset)
      }
      break
    }
  }

  // Set "prev" offset.
  prev := ""
  if offset >= limit {
    prev = strconv.Itoa(offset - limit)
  }

  return results, limit, next, prev, totalCount, error(nil)
}



