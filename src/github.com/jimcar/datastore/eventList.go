package datastore

import (
  "fmt"
  "strings"
  "strconv"
  "encoding/json"
)

// ----------------------------------------------------------------------------
//  Name: ListEvents
//  Desc:

func ListEvents(name, key, etype string, params map[string]string) (string, string, string, error) {

  // Get the event results.
  results, _ := eventResults(name, key, etype, params)

  // Sort results by descending ordinal values (newest first)
  descendingOrdinal := func(r1, r2 *Result) bool {
    return r1.Ordinal > r2.Ordinal
  }
  By(descendingOrdinal).Sort(results)

  // ----------------------------------------------------------------------
  //  The value used to specify the "next" field in the result as well a
  //  as the Link header with a rel="next".  This value is used to build
  //  the uri that can be used to fetch the next page of results.
  nextBeforeEvent := ""

  // --------------------------------------------------------------------------
  // limit param
  limit := ListLimit(params)

  if len(results) > limit {
    nextBeforeEvent = strconv.Itoa(results[limit-1].Timestamp)
    if results[limit].Ordinal > 0 {
      nextBeforeEvent = fmt.Sprintf("%v/%v", nextBeforeEvent, results[limit-1].Ordinal)
    }
    nextBeforeEvent = eventResultPage(name, key, etype, limit, nextBeforeEvent)
    results = results[0:limit]
  }

  // Return the json response body
  body := listResponse(ResponseBody{len(results), results, nextBeforeEvent, "", 0})
  return body, nextBeforeEvent, "", error(nil)
}

// ----------------------------------------------------------------------------
//  Name: eventResults
//  Desc:

func eventResults(name, key, etype string, params map[string]string) ([]Result, error) {

  var value   ResultValue
  var r       Ref
  var results []Result

  etypeTsTableName := eventTsTableName(name, key, etype)
  etypeTimestamps  := getCollectionHandle(etypeTsTableName)
  collectionEvents := getCollectionHandle(eventTableName(name, key))
  eventRefTable    := getCollectionHandle("EventRefTable")

  ro.SetFillCache(false)
  tsIterator := etypeTimestamps.NewIterator(ro)
  defer tsIterator.Close()

  // ----------------------------------------------------------------------
  // params to specify start of event range
  // startEvent (inclusive) and afterEvent (exclusive)
  var ts string
  var startOrdinal, afterOrdinal int = 0, 0

  startEvent, startOk := params["startEvent"]
  afterEvent, afterOk := params["afterEvent"]

  if afterOk && startOk {
    // illegal combo, get actual error from oio
    return []Result{}, error(nil)
  }

  if startOk {
    ts, startOrdinal = splitRangeParam(startEvent)
    tsIterator.Seek([]byte(ts))
  } else if afterOk {
    ts, afterOrdinal = splitRangeParam(afterEvent)
    tsIterator.Seek([]byte(ts))
    if afterOrdinal == 0 {
      tsIterator.Next()
    }
  } else {
    tsIterator.SeekToFirst()
  }

  // ----------------------------------------------------------------------
  // params to specify end of event range
  // endEvent (inclusive) and beforeEvent (exclusive)

  var endTs, beforeTs string = "", ""
  var endOrdinal, beforeOrdinal int = 0, 0
  var endEvent, beforeEvent string = "", ""

  endEventParam,       endOk := params["endEvent"]
  beforeEventParam, beforeOk := params["beforeEvent"]
  if endOk && beforeOk {
    // illegal combo, get actual error from oio
    return []Result{}, error(nil)
  }

  if endOk {
    endTs, endOrdinal = splitRangeParam(endEventParam)
    endEvent = getRangeEvent(endTs, endOrdinal)
  } else if beforeOk {
    beforeTs, beforeOrdinal = splitRangeParam(beforeEventParam)
    beforeEvent = getRangeEvent(beforeTs, beforeOrdinal)
  }

  // --------------------------------------------------------------------------
  //  Iterate through the timestamps associated with collection/key/etype

  count := 0

  for tsIterator = tsIterator; tsIterator.Valid(); tsIterator.Next() {

    var currentEvent string
    currentTimestamp := string(tsIterator.Key())

    // Get outta here, if this is the beforeEvent and no ordinal is specified
    if beforeOk && beforeTs == currentTimestamp && beforeOrdinal == 0 {
      break
    }

    // Get the ordinals associated with this timestamp
    ordinals := strings.Split(string(tsIterator.Value()), ":")

    // --------------------------------------------------------------
    //  Iterate through the ordinals associated with this timestamp

    for _, ord := range ordinals {

      currentEvent = fmt.Sprintf("%s/%s", currentTimestamp, ord)

      // Get outta here, if this is the beforeEvent
      if beforeOk && beforeEvent == currentEvent {
        break
      }

      // Skip ordinals until we: 1) reach the startEvent, or 2) get past the afterEvent
      if count == 0 && (startOk || afterOk) {
        ordinal, _ := strconv.Atoi(ord)
        if ordinal < startOrdinal || ordinal <= afterOrdinal {
          continue
        }
        count++
      }

      // Get ref value from collection/key/ordinal.
      ref, _ := collectionEvents.Get(ro, []byte(ord))

      // Unmarshal the json data/metadata extracted from the event ref table
      r, value = Ref{}, ResultValue{}
      rdata, _ := eventRefTable.Get(ro, []byte(ref))
      json.Unmarshal(rdata, &r)
      json.Unmarshal([]byte(r.Value), &value)

      // Create new event result and append to results.
      timestamp, _ := strconv.Atoi(r.Mdata.Timestamp)
      path := PathValue{
        r.Mdata.Collection, r.Mdata.Key, r.Mdata.Type, r.Mdata.Ref,
        timestamp, r.Mdata.Ordinal, false,
      }
      results = append(results, Result{path, value, timestamp, r.Mdata.Ordinal, 0})

      // Get outta here, if this is the endEvent.
      if endOk && endEvent == currentEvent {
        break
      }

    } // End ordinal iteration.

    // Get outta here, if this is the endEvent.
    if endOk && endTs == currentTimestamp {
      break
    }

    // Get outta here, if this is the beforeEvent.
    if beforeOk && beforeTs == currentTimestamp {
      break
    }

  } // End timestamp iteration

  ro.SetFillCache(true)
  err := tsIterator.GetError()

  return results, err
}



