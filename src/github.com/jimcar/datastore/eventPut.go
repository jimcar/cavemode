package datastore

import (
  "strconv"
  "encoding/json"
)

// ----------------------------------------------------------------------------
//  Name: PutEvent
//  Desc:

func PutEvent(data, name, key, etype, timestamp string, ordinal int) (string, string, int) {

  var ref string

  if isDup, mdata := isDuplicateEventData(data, name, key, etype, timestamp); isDup == true {

    // For duplicate data, just use the existing ref, timestamp and ordinal values.
    ref, timestamp, ordinal = mdata.Ref, mdata.Timestamp, mdata.Ordinal

  } else {

    ref = generateRefValue()

    if timestamp == "" {
      timestamp = generateTimestamp()
    } else {
      timestamp = getTimestamp(timestamp)
    }

    // Assign the next ordinal value
    if ordinal == 0 {
      ordinal = getNextOrdinalValue(name, key)
    }

    // Update collection/key/etype/timestamp with ordinal value.
    etypeTsListName := eventTsTableName(name, key, etype)
    addItemToList(etypeTsListName, timestamp, strconv.Itoa(ordinal))

    // Update collection/key/etype/refs with current ref.
    refsKey := createKey(name, key, etype)
    addItemToList("AllEventRefsTable", refsKey, ref)
  }

  // Update collection/key/ordinal with ref value.
  myKey := strconv.Itoa(ordinal)
  collectionEvents := getCollectionHandle(eventTableName(name, key))
  collectionEvents.Put(wo, []byte(myKey), []byte(ref))

  // Generate datestamp.
  datestamp := generateDatestamp()

  // Initialize event ref instance.
  r := Ref{Metadata{name, key, ref, datestamp, etype, timestamp, ordinal, false}, data}

  // Update EventRefTable[ref] with eventRef data.
  refdata, _ := json.Marshal(r)
  eventRefTable := getCollectionHandle("EventRefTable")
  eventRefTable.Put(wo, []byte(ref), refdata)

  return ref, timestamp, ordinal
}


