package orchio

import (
  "fmt"
  "github.com/jimcar/datastore"
  "encoding/json"
)

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
//  Name: locationString
//  Desc:

func locationString(name, key, ref string) string {
  return fmt.Sprintf("\"/%s/%s/%s/refs/%s\"", orchioRoot, name, key, ref)
}

// ----------------------------------------------------------------------------
//  Name: eventLocationString
//  Desc:

func eventLocationString(name, key, etype, timestamp string, ordinal int) string {
  return fmt.Sprintf("\"/%s/%s/%s/events/%s/%s/%d\"", orchioRoot, name, key, etype, timestamp, ordinal)
}

// ----------------------------------------------------------------------------
//  Name: reqIdString
//  Desc:

func reqIdString() string {
  return "abcdefg-hijk-lmnop-qrst-uvw-xyy"
}

// ----------------------------------------------------------------------------
//  Name: getRefValueFromMdata
//  Desc:

func getRefValueFromMdata(mdata string) string {
  var jmctmp datastore.Metadata
  if err := json.Unmarshal([]byte(mdata), &jmctmp); err != nil {
    return ""
  }
  return jmctmp.Ref
}

// ----------------------------------------------------------------------------
//  Name: getRefValue
//  Desc:

func getRefValue(name, key string) string {
  if _, mdata, err := datastore.GetKey(name, key, ""); err == nil {
    if mdata != "" {
      var jmctmp datastore.Metadata
      if err := json.Unmarshal([]byte(mdata), &jmctmp); err == nil {
        return jmctmp.Ref
      }
    }
  }
  return ""
}

// ----------------------------------------------------------------------------
//  Name: checkRefValue
//  Desc:

func checkRefValue(name, key, ref string) int {
  data, _, err := datastore.GetKey(name, key, ref)
  if err != nil {
    return 500
  } else if data == "" {
    return 404
  }
  return 0
}


