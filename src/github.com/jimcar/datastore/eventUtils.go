package datastore

import (
  "fmt"
  "strings"
  "strconv"
  "encoding/json"
)

// ----------------------------------------------------------------------------
//  Name: getRangeEvent
//  Desc:

func getRangeEvent(ts string, ord int) string {
  if ord == 0 {
    return ts
  }
  return fmt.Sprintf("%s/%v", ts, ord)
}

// ----------------------------------------------------------------------------
//  Name: splitRangeParam
//  Desc:

func splitRangeParam(event string) (string, int) {
  var ord int = 0
  tmpslice := strings.Split(event, "/")
  if len(tmpslice) == 2 {
    ord, _ = strconv.Atoi(tmpslice[1])
  }
  return getTimestamp(tmpslice[0]), ord
}

// ----------------------------------------------------------------------------
//  Name: eventTimestamp
//  Desc:

func eventTimestamp(ref string) string {

  if ref == "" {
    return generateTimestamp()
  }

  db := getCollectionHandle("EventRefTable")
  tmpEventRef, _ := db.Get(ro, []byte(ref))
  var r Ref
  json.Unmarshal(tmpEventRef, &r)
  return r.Mdata.Timestamp
}

// ----------------------------------------------------------------------------
//  Name: isDuplicateEventData
//  Desc:

func isDuplicateEventData(data, name, key, etype, timestamp string) (bool, Metadata) {

  refsKey := createKey(name, key, etype)
  db := getCollectionHandle("AllEventRefsTable")

  if slice, err := db.Get(ro, []byte(refsKey)); err == nil {
    if slice != nil {
      var r Ref
      eventRefTable := getCollectionHandle("EventRefTable")
      refs := strings.Split(string(slice), ":")
      for i := range refs {
        rdata, _ := eventRefTable.Get(ro, []byte(refs[i]))
        json.Unmarshal(rdata, &r)
        if data == string(r.Value) && timestamp == r.Mdata.Timestamp {
          return true, r.Mdata
        }
      }
    }
  }
  return false, Metadata{}
}

// // ----------------------------------------------------------------------------
// //  Name: isValidEventRef
// //  Desc: Verify that ref is included in AllRefsTable/collection/key/events/etype

// func isValidEventRef(name, key, etype, ref string) bool {
//   refsKey := createKey(name, key, etype)
//   return isValidItem("AllEventRefsTable", refsKey, ref)
// }

// // ----------------------------------------------------------------------------
// //  Name: IsValidEventRef
// //  Desc: Publicly available; unquotes ref val before calling isValidEventRef

// func IsValidEventRef(name, key, etype, qref string) bool {
//   if ref, err := strconv.Unquote(qref); err == nil {
//     return isValidEventRef(name, key, etype, ref)
//   }
//   return false
// }


