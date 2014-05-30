package datastore

import (
  "fmt"
  "code.google.com/p/go-leveldb"
  "os"
)

// ---------------- Global vars -----------------------------------------------

var ro = leveldb.NewReadOptions()
var wo = leveldb.NewWriteOptions()

// ----------------------------------------------------------------------------
//  Name: dbDirectory
//  Desc: Location for leveldb directories/files
//        Checks env var CAVEMODE_DB_DIR
//        Default is "$HOME/.cavemode/"
//        Executes only on first invocation; otherwise returns immediately.

var leveldbDir string = ""

func dbDirectory() string {

  // Return immediately if we've already done this!
  if leveldbDir != "" {
    return leveldbDir
  }

  // Check env var
  leveldbDir = os.Getenv("CAVEMODE_DB_DIR")

  // Use default if env var is not set or empty.
  if leveldbDir == "" {
    leveldbDir = os.Getenv("HOME") + "/.cavemode/leveldb-files"
  }

  // Create the directory if it doesn't exist.
  if _, err := os.Stat(leveldbDir); err != nil {
    os.MkdirAll(leveldbDir, 0755)
  }

  fmt.Printf("datastore:dbDirectory() %s\n", leveldbDir)

  return leveldbDir
}

// ----------------------------------------------------------------------------
//  Name: dbName
//  Desc:

func dbName(name string) string {
  return dbDirectory() + "/" + name
}

// ----------------------------------------------------------------------------
//  Name: dbOptions
//  Desc:

func dbOptions() *leveldb.Options {
  opts := leveldb.NewOptions()
  opts.SetCache(leveldb.NewLRUCache(3<<30))
  opts.SetCreateIfMissing(true)
  return opts
}

// ----------------------------------------------------------------------------
//  Name: getCollectionHandle
//  Desc:

var Collections = make(map[string]*leveldb.DB)

func getCollectionHandle(name string) (*leveldb.DB) {
  var handle *leveldb.DB = nil

  if db, ok := Collections[name]; ok {
    handle = db
  } else {
    opts := dbOptions()
    if db, err := leveldb.Open(dbName(name), opts); err == nil {
      Collections[name] = db
      handle = db
    } else {
      logError("getCollectionHandle", err)
    }
  }
  return handle
}

// ----------------------------------------------------------------------------
//  Name: destroyCollection
//  Desc:

func destroyCollection(name string) error {

  return leveldb.DestroyDatabase(dbName(name), dbOptions())

}



