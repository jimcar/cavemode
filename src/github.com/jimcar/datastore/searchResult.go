package datastore

import (
  "encoding/json"
)

// ----------------------------------------------------------------------------
//  Name: newSearchResults
//  Desc:

func newSearchResults(name, queryStr string) (int, error) {
  // TraceMsg("==> Datastore:newSearchResults")

  var r       Ref
  var value   ResultValue

   // Process the query param.
  searchTerms, logicOp := getSearchTerms(queryStr)

  // Destroy any previously-created list for this query.
  queryID := queryName(name, queryStr)
  destroyList("AllQueryResults", queryID)

  // Prepare to search the collection.
  collection := getCollectionHandle(name)
  ro.SetFillCache(false)
  iterator := collection.NewIterator(ro)
  defer iterator.Close()
  iterator.SeekToFirst()

  refTable := getCollectionHandle("RefTable")
  count := 0

  for iterator = iterator; iterator.Valid(); iterator.Next() {

    // Get ref value from collection/key/ordinal.
    ref := string(iterator.Value())

    // Unmarshal the json data/metadata extracted from the ref table
    r, value = Ref{}, ResultValue{}
    rdata, _ := refTable.Get(ro, []byte(ref))
    json.Unmarshal(rdata, &r)
    json.Unmarshal([]byte(r.Value), &value)
    // TraceMsg("value: %v", value)

    // If search terms match, add ref to query results.
    if searchTermsMatch(value, searchTerms, logicOp) {
      addItemToList("AllQueryResults", queryID, ref)
      count++
    }
  }
  ro.SetFillCache(true)

  var err error = nil
  if err = iterator.GetError(); err != nil {
    logError("datastore:newSearchResults", err)
  }

  return count, err
}

