package orchio

import (
  "code.google.com/p/gorest"
  // "fmt"
  "encoding/json"
  "os"
)

// ----------------------------------------------------------------------------
//  Name: updateResponse
//  Desc:

func updateResponse(r *gorest.ResponseBuilder, method string) *gorest.ResponseBuilder {
  r.ConnectionKeepAlive()
  r.AddHeader("X-Orchestrate-Req-Id", reqIdString())
  if method != "delete" {
    r.SetContentType("application/json")
  }
  return r
}

// ----------------------------------------------------------------------------
//  Name: completeTheResponse
//  Desc:

func completeTheResponse(r *gorest.ResponseBuilder, code int, msgBody string) *gorest.ResponseBuilder {
  r.SetResponseCode(code)
  if msgBody != "" {
    r.WriteAndOveride([]byte(msgBody))
  }
  return r
}

// ----------------------------------------------------------------------------
//  Name: keyErrorMsgBody
//  Desc:

func keyErrorMsgBody(responseCode int, msgCode, name, key string) string {
  return createErrorMsgBody(responseCode, msgCode, name, key, "", "", "", 0)
}

// ----------------------------------------------------------------------------
//  Name: refErrorMsgBody
//  Desc:

func refErrorMsgBody(responseCode int, msgCode, name, key, ref string) string {
  return createErrorMsgBody(responseCode, msgCode, name, key, "", ref, "", 0)
}

// ----------------------------------------------------------------------------
//  Name: eventErrorMsgBody
//  Desc:

func eventErrorMsgBody(responseCode int, msgCode, name, key, etype, timestamp string, ordinal int) string {
  return createErrorMsgBody(responseCode, msgCode, name, key, etype, "", timestamp, ordinal)
}

// ----------------------------------------------------------------------------
//  Name: createErrorMsgBody
//  Desc:

func createErrorMsgBody(responseCode int, msgCode, name, key, etype, ref, timestamp string, ordinal int) string {

  msgMap := map[string]string{
    /* 400 */ "api_bad_request":       "Invalid value for header ''If-Match''.",
    /* 400 */ "item_ref_malformed":    "Invalid value for header ''If-None-Match''.",
    /* 400 */ "invalid_content_type":  "Invalid value for header ''Content-Type''.",
    /* 404 */ "items_not_found":       "The requested items could not be found.",
    /* 412 */ "item_version_mismatch": "The version of the item does not match.",
    /* 412 */ "item_already_present":  "The item is already present.",
    /* 500 */ "internal_error":        "Internal error.",
  }

  var msgBody ErrMsgBody
  if msg, ok := msgMap[msgCode]; ok {
    if responseCode == 404 {
      var items []Item
      item := Item{name, key, etype, ref, timestamp, ordinal}
      msgBody = ErrMsgBody{msg, Detail{append(items, item)}, msgCode}
    } else {
      msgBody = ErrMsgBody{msg, Detail{}, msgCode}
    }
  } else {
    msgBody = ErrMsgBody{"Unknown message code.", Detail{}, msgCode}
  }

  var jsonMsgBody []byte
  if os.Getenv("CAVEMODE_JSON_INDENT") == "true" {
    jsonMsgBody, _ = json.MarshalIndent(msgBody, "", " ")
  } else {
    jsonMsgBody, _ = json.Marshal(msgBody)
  }

  return string(jsonMsgBody)
}

