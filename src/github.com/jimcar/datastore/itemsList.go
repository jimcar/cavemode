package datastore

import (
  "strings"
)

//  Name: itemsList.go
//
//  Desc: Lists of simple items, implemented as ":" delimited strings.
//        Examples of simple items are values for ordinals and refs.
//        Each list is stored in an underlying leveldb "items table."


// ----------------------------------------------------------------------------
//  Name: isValidItem
//  Desc: could also be called isItemInList

func isValidItem(itemsTable, itemKey, item string) bool {

  retval := false

  db := getCollectionHandle(itemsTable)
  if list, err := db.Get(ro, []byte(itemKey)); err == nil {
    if list != nil {
      items := strings.Split(string(list), ":")
      for i := range items {
        if items[i] == item {
          retval = true
          break
        }
      }
    }
  }

  return retval
}

// ----------------------------------------------------------------------------
//  Name: getItems
//  Desc: returns all items in the specified list

func getItems(itemsTable, itemKey string) []string {

  retval := []string{}

  db := getCollectionHandle(itemsTable)
  if list, err := db.Get(ro, []byte(itemKey)); err == nil {
    if list != nil {
      retval = strings.Split(string(list), ":")
    }
  }

  return retval
}

// ----------------------------------------------------------------------------
//  Name: destroyList
//  Desc: Deletes all items from list

func destroyList(itemsTable, itemsKey string) error {

  db := getCollectionHandle(itemsTable)
  err := db.Delete(wo, []byte(itemsKey))

  return err
}

// ----------------------------------------------------------------------------
//  Name: deleteItemFromList
//  Desc: Remove item from a ":" delimited list

func deleteItemFromList(itemsTable, key, item string) error {

  err := error(nil)

  db := getCollectionHandle(itemsTable)
  if list, err := db.Get(ro, []byte(key)); err == nil {
    oldItems := strings.Split(string(list), ":")

    var newItems string = ""
    for _, oldItem := range oldItems {
      if oldItem != item {
        if newItems == "" {
          newItems = oldItem
        } else {
          newItems = strings.Join([]string{newItems, oldItem}, ":")
        }
      }
    }
    // TraceMsg(fmt.Sprintf("\tnewItems = %s", newItems))

    err = db.Put(wo, []byte(key), []byte(newItems))
  }

  return err
}

// ----------------------------------------------------------------------------
//  Name: addItemToList
//  Desc: Add new item to list of items (implemented as ":" delimited string)

func addItemToList(itemsTable, key, item string) error {

  err := error(nil)

  db := getCollectionHandle(itemsTable)
  if list, err := db.Get(ro, []byte(key)); err == nil {
    items := item
    if list != nil {
      items = string(list) + ":" + item
    }
    // TraceMsg(fmt.Sprintf("\titems = %s", items))

    err = db.Put(wo, []byte(key), []byte(items))
  }

  return err
}

