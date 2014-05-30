package datastore

import (
  // "fmt"
  "strconv"
)

// ----------------------------------------------------------------------------
//  Name: GetEvent
//  Desc:

func GetEvent(name, key, etype, timestamp string, ordinal int) (string, string, error) {

  // Return immediately if the ordinal is not associated with the timestamp
  etypeTsTableName := eventTsTableName(name, key, etype)
  if ok := isValidItem(etypeTsTableName, timestamp, strconv.Itoa(ordinal)); !ok {
    return "", "", error(nil)
  }

  // Get ref value from collection/key/ordinal.
  collectionEvents := getCollectionHandle(eventTableName(name, key))
  ref, _ := collectionEvents.Get(ro, []byte(strconv.Itoa(ordinal)))

  // Get data/metadata from event ref table
  eventRefTable := getCollectionHandle("EventRefTable")
  eventRefData, _ := eventRefTable.Get(ro, []byte(ref))
  body, mdata := itemResponse(eventRefData)

  return body, mdata, error(nil)
}


