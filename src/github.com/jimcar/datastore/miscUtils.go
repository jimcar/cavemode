package datastore

import (
  "fmt"
  "strings"
  "strconv"
  "encoding/json"
  "os"
)

// ----------------------------------------------------------------------------
//  Name: getNextOrdinalValue
//  Desc: Get and increment the ordinal value.
//        Update maxOrdinals/collection/key with new maxOrdinal

func getNextOrdinalValue(name, key string) int {
  maxOrdinals := getCollectionHandle("MaxOrdinals")
  myKey := createKey(name, key)
  currentMax, _ := maxOrdinals.Get(ro, []byte(myKey))
  ordinal, _ := strconv.Atoi(string(currentMax))
  ordinal++
  maxOrdinals.Put(wo, []byte(myKey), []byte(strconv.Itoa(ordinal)))
  return ordinal
}

// ----------------------------------------------------------------------------
//  Name: isDuplicateData
//  Desc:

func isDuplicateData(data, name, key string) (bool, string) {

  allRefsKey := createKey(name, key)
  db := getCollectionHandle("AllRefsTable")

  if slice, err := db.Get(ro, []byte(allRefsKey)); err == nil {
    if slice != nil {
      refTable := getCollectionHandle("RefTable")
      refs := strings.Split(string(slice), ":")
      for i := range refs {
        rdata, _ := refTable.Get(ro, []byte(refs[i]))
        var r Ref
        json.Unmarshal(rdata, &r)
        if data == string(r.Value) {
          return true, refs[i]
        }
      }
    }
  }
  return false, ""
}

// ----------------------------------------------------------------------------
//  Name: createKey
//  Desc:

func createKey(values ...string) string {
  key := values[0]
  for i, _ := range values {
    if i == 0 {
      continue
    }
    key = strings.Join([]string{key, values[i]}, "-")
  }
  return key
}

func createName(values ...string) string {
  key := values[0]
  for i, _ := range values {
    if i == 0 {
      continue
    }
    key = strings.Join([]string{key, values[i]}, "-")
  }
  return key
}

func queryName(name, queryStr string) string {
  return fmt.Sprintf("%s-query-%s", name, queryStr)
}

func eventTableName(name, key string) string {
  return strings.Join([]string{name, key, "events"}, "-")
}

func eventTsTableName(name, key, etype string) string {
  return strings.Join([]string{name, key, "events", etype}, "-")
}

func graphTableName(name, key string) string {
  return strings.Join([]string{name, key, "graph"}, "-")
}

func relationTableName(name, key string) string {
  return strings.Join([]string{name, key, "relations"}, "-")
}

// ----------------------------------------------------------------------------
//  Name: outputError
//  Desc:

func outputError(name string, err error) {
  if err != nil {
    fmt.Printf("%s: err = %v\n", name, err)
  }
}

// ----------------------------------------------------------------------------
//  Name: logError
//  Desc:

func logError(name string, err error) {
  outputError(name, err)
}

// ----------------------------------------------------------------------------
//  Name: traceMsg
//  Desc:

var TraceMode string = ""

func traceMode() bool {

  // Return immediately if we've already done this!
  if TraceMode != "" {
    return TraceMode == "true"
  }

  TraceMode = os.Getenv("CAVEMODE_TRACE")
  if TraceMode != "true" {
    TraceMode = "false"
  }

  return TraceMode == "true"
}


func TraceMsg(msg string) {
  if traceMode() {
    fmt.Printf("%s\n", msg)
  }
}


