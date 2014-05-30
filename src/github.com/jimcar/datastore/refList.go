package datastore

import (
  "strconv"
  "encoding/json"
)

// ----------------------------------------------------------------------------
//  Name: ListRefs
//  Desc:

func ListRefs(name, key string, params map[string]string) (string, string, string, error) {

  // Get the search results.
  results, limit, offsets, err := listRefs(name, key, params)

  next, _ := offsets["next"]
  if next != "" {
    next = offsetResultPage(name, limit, next)
  }

  prev, _ := offsets["prev"]
  if prev != "" {
    prev = offsetResultPage(name, limit, prev)
  }

  body := listResponse(ResponseBody{len(results), results, next, prev, 0})

  // Return the json response body
  return body, next, prev, err
}

// ----------------------------------------------------------------------------
//  Name: listRefs
//  Desc:

func listRefs(name, key string, params map[string]string) ([]Result, int, map[string]string, error) {
  // TraceMsg("==> Datastore:listRefs")

  var r       Ref
  var value   ResultValue
  var results []Result

  // limit param
  limit := ListLimit(params)

  // offset param - to specify start of key range
  offset := 0
  if offsetStr, ok := params["offset"]; ok {
    offset, _ = strconv.Atoi(offsetStr)
  }

  // values param
  includeValueInResults := false
  if valuesStr, ok := params["values"]; ok {
    includeValueInResults = (valuesStr == "true")
  }

  // The table that contains the data, indexed by ref value.
  refTable := getCollectionHandle("RefTable")

  next := ""
  nextOffset := 0
  count := 0

  // Get all of the refs associated with this key.
  keyRefs := getItems("AllRefsTable", createKey(name, key))

  for i, ref := range keyRefs {

    nextOffset++

    if i < offset {
      continue
    }

    // Unmarshal the json data/metadata extracted from the ref table
    rdata, _ := refTable.Get(ro, []byte(ref))
    r = Ref{}
    json.Unmarshal(rdata, &r)

    //
    resultRef := r.Mdata.Ref
    if r.Mdata.Tombstone {
      resultRef = ""
    }

    value = ResultValue{}
    if includeValueInResults && r.Mdata.Tombstone == false {
      json.Unmarshal([]byte(r.Value), &value)
    }

    // Create new result and append to results.
    path := PathValue{r.Mdata.Collection, r.Mdata.Key, "", resultRef, 0, 0, r.Mdata.Tombstone}
    reftime, _ := strconv.Atoi(r.Mdata.Datestamp)
    results = append(results, Result{path, value, 0, 0, reftime})
    count++

    // Check limit, set "next" if limit reached before end of range.
    if count == limit {
      if nextOffset < len(keyRefs) {
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

  offsets := map[string]string{"next":next, "prev":prev}

  return results, limit, offsets, error(nil)
}

