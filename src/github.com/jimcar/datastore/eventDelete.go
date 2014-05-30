package datastore

import (
  "strconv"
)

// ----------------------------------------------------------------------------
//  Name: DeleteEvent
//  Desc:

func DeleteEvent(name, key, etype, timestamp string, ordinal int, purge bool) {

  ord := strconv.Itoa(ordinal)

  // Verify that the ordinal is associated with the timestamp by checking
  // collection/key/etype/timestamp for the ordinal value.
  etypeTsTableName := eventTsTableName(name, key, etype)
  if ok := isValidItem(etypeTsTableName, timestamp, ord); !ok {
    return
  }
  deleteItemFromList(etypeTsTableName, timestamp, ord)

  // Get ref value from collection/key/ordinal.
  jmctable := getCollectionHandle(eventTableName(name, key))
  ref, _ := jmctable.Get(ro, []byte(ord))
  if err := jmctable.Delete(wo, []byte(ord)); err != nil {
    logError("datastore:DeleteEvent", err)
    return
  }

  // Verify that the ref is associated with etype for collection/key.
  refsKey := createKey(name, key, etype)
  if ok := isValidItem("AllEventRefsTable", refsKey, string(ref)); !ok {
    return
  }
  deleteItemFromList("AllEventRefsTable", refsKey, string(ref))

  // Delete ref value from collection/key/ordinal.
  myKey := strconv.Itoa(ordinal)
  collectionEvents := getCollectionHandle(createName(name, key, "eventRefs"))
  if err := collectionEvents.Delete(wo, []byte(myKey)); err != nil {
    logError("datastore:DeleteEvent", err)
    return
  }

  // Delete eventRef data from EventRefTable[ref].
  eventRefTable := getCollectionHandle("EventRefTable")
  if err := eventRefTable.Delete(wo, ref); err != nil {
    logError("datastore:DeleteEvent", err)
    return
  }

  return
}


