package datastore

import (
  "encoding/json"
  "os"
)

// ----------------------------------------------------------------------------
//  Name: itemResponse
//  Desc: Returns json response body (value) along with the item's metadata.
//        Item can be key or event.

func itemResponse(refData []byte) (string, string) {

  var r Ref
  json.Unmarshal(refData, &r)
  mdata, _ := json.Marshal(r.Mdata)
  body := r.Value

  if os.Getenv("CAVEMODE_JSON_INDENT") == "true" {
    var rv ResultValue
    json.Unmarshal([]byte(r.Value), &rv)
    jstr, _ := json.MarshalIndent(rv, "", " ")
    body = string(jstr) + "\n"
  }

  return body, string(mdata)
}

// ----------------------------------------------------------------------------
//  Name: listResponse
//  Desc: Returns json response body generated from list results.

func listResponse(responseBody ResponseBody) string {

  var body string

  if os.Getenv("CAVEMODE_JSON_INDENT") == "true" {

    jstr, _ := json.MarshalIndent(responseBody, "", " ")
    body = string(jstr) + "\n"

  } else {

    jstr, _ := json.Marshal(responseBody)
    body = string(jstr)
    // TraceMsg(fmt.Sprintf("json s = %s\n", body))
  }

  return body
}


