package datastore

import (
  // "fmt"
  "encoding/json"
)

// ----------------------------------------------------------------------------
//  Name: GetKey
//  Desc:

func GetKey(name, key, ref string) (string, string, error) {

  var body, mdata string = "", ""

  if ref == "" {
    collection := getCollectionHandle(name)
    tmpref, _ := collection.Get(ro, []byte(key))
    ref = string(tmpref)
  }

  // Verify that this ref is included in the list for collection/key
  if ok := isValidKeyRef(name, key, ref); ok {
    refTable := getCollectionHandle("RefTable")
    rdata, _ := refTable.Get(ro, []byte(ref))
    body, mdata = itemResponse(rdata)
  }

  return body, mdata, error(nil)
}

// ----------------------------------------------------------------------------
//  Name: PutKey
//  Desc:

func PutKey(data, name, key string) string {

  // For duplicate data, just use the existing ref value.
  isDup, ref := isDuplicateData(data, name, key)

  if isDup == false {
    // Generate a new ref value
    ref = generateRefValue()

    // Add ref value to AllRefsTable {key: "collection-key"}
    allRefsKey := createKey(name, key)
    addItemToList("AllRefsTable", allRefsKey, ref)
  }

  // Update collection table with ref data
  collection := getCollectionHandle(name)
  collection.Put(wo, []byte(key), []byte(ref))

  // Generate datestamp
  datestamp := generateTimestamp()     // aka "reftime" (ms since unix epoch)

  // Initialize ref instance, update RefTable
  refdata, _ := json.Marshal(Ref{Metadata{name, key, ref, datestamp, "", "", 0, false}, data})
  db := getCollectionHandle("RefTable")
  db.Put(wo, []byte(ref), refdata)

  return ref
}

// ----------------------------------------------------------------------------
//  Name: DeleteKey
//  Desc:

func DeleteKey(name, key string, purge bool) error {

  var ref []byte
  var err error

  // Get the ref value from (the head of) collection/key
  collection := getCollectionHandle(name)
  ref, _ = collection.Get(ro, []byte(key))

  // Delete the ref value from (the head of) collection/key
  err = collection.Delete(wo, []byte(key))

  refTable := getCollectionHandle("RefTable")

  if purge {

    // Delete the data from refTable.
    err = refTable.Delete(wo, ref)

    // Delete the ref value from allRefsTable
    allRefsKey := createKey(name, key)
    allRefsTable := getCollectionHandle("AllRefsTable")
    err = allRefsTable.Delete(wo, []byte(allRefsKey))

  } else {  // purge == false

    // Set the tombstone flag to true
    refData, _ := refTable.Get(ro, ref)
    var r Ref
    json.Unmarshal(refData, &r)
    r.Mdata.Tombstone = true
    r.Mdata.Datestamp = generateTimestamp()
    refData, _ = json.Marshal(r)
    err = refTable.Put(wo, ref, refData)
  }

  return err
}


