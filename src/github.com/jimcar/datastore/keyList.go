package datastore

import (
  "encoding/json"
)

// ----------------------------------------------------------------------------
//  Name: listKeys
//  Desc:

func listKeys(name string, params map[string]string) (string, string, string, error) {

  // Get the list results.
  results, limit, next, err := listResults(name, params)

  if next != "" {
    next = keyResultPage(name, limit, next)
  }

  body := listResponse(ResponseBody{len(results), results, next, "", 0})

  // Return the json response body
  return body, next, "", err
}

// ----------------------------------------------------------------------------
//  Name: listResults
//  Desc:

func listResults(name string, params map[string]string) ([]Result, int, string, error) {
  // TraceMsg("==> Datastore:listResults")

  var r       Ref
  var value   ResultValue
  var results []Result

  collection := getCollectionHandle(name)
  refs := getCollectionHandle("RefTable")

  ro.SetFillCache(false)
  iterator := collection.NewIterator(ro)
  defer iterator.Close()

  // limit param
  limit := ListLimit(params)

  // startKey (inclusive) and afterKey (exclusive) to specify start of key range
  startKey, startOk := params["startKey"]
  afterKey, afterOk := params["afterKey"]

  if startOk && afterOk {
    // illegal combo, get actual error from oio
    return []Result{}, 0, "", error(nil)
  }

  if startOk {
    iterator.Seek([]byte(startKey))
  } else if afterOk {
    iterator.Seek([]byte(afterKey))
    iterator.Next()
  } else {
    iterator.SeekToFirst()
  }

  // beforeKey (exclusive) and endKey (inclusive) to specify end of key range
  beforeKey, beforeOk := params["beforeKey"]
  endKey,    endOk    := params["endKey"]

  if beforeOk && endOk {
    // illegal combo, get actual error from oio
    return []Result{}, 0, "", error(nil)
  }

  next := ""
  count := 0
  for iterator = iterator; iterator.Valid(); iterator.Next() {

    currentKey := string(iterator.Key())

    if beforeOk && beforeKey == currentKey {
      break
    }

    // Get ref value from collection/key/ordinal.
    ref := string(iterator.Value())

    // Unmarshal the json data/metadata extracted from the ref table
    r, value = Ref{}, ResultValue{}
    rdata, _ := refs.Get(ro, []byte(ref))
    json.Unmarshal(rdata, &r)
    json.Unmarshal([]byte(r.Value), &value)

    // Create new result and append to results.
    path := PathValue{r.Mdata.Collection, r.Mdata.Key, "", r.Mdata.Ref, 0, 0, false}
    results = append(results, Result{path, value, 0, 0, 0})

    if endOk {
      if endKey == currentKey {
        break
      }
    }

    count++
    if count == limit {
      lastKey := ""
      if endOk {
        iterator.Seek([]byte(endKey))
        if iterator.Valid() {
          lastKey = string(iterator.Key())
        }
      } else if beforeOk {
        iterator.Seek([]byte(beforeKey))
        if iterator.Valid() {
          iterator.Prev()
          if iterator.Valid() {
            lastKey = string(iterator.Key())
          }
        }
      } else {
        iterator.SeekToLast()
        lastKey = string(iterator.Key())
      }

      if lastKey != currentKey {
        next = currentKey
      }
      break
    }
  }
  ro.SetFillCache(true)

  var estat error = nil
  if err := iterator.GetError(); err != nil {
    estat = err
  }

  return results, limit, next, estat
}


