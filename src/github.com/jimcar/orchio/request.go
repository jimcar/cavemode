package orchio

import (
  "net/http"
  "fmt"
  "github.com/jimcar/datastore"
  "strings"
  "strconv"
)


// ---------------- examineRequest ------------------------------------------------

func printRequest(r *http.Request) *http.Request {
  fmt.Printf("\theader     = %v\n", r.Header)
  fmt.Printf("\turl        = %s\n", r.URL)
  fmt.Printf("\traw_query  = %v\n", r.URL.RawQuery)
  fmt.Printf("\tremoteAddr = %v\n", strings.Split(r.RemoteAddr, ":")[0])
  return r
}

// ---------------- getParams -------------------------------------------------

func getParamValues(query string) map[string]string {
  var kv = map[string]string{}
  if query != "" {
    params := strings.Split(query, "&")
    for i := range params {
      pair := strings.Split(params[i], "=")
      kv[pair[0]] = pair[1]
    }
  }
  return kv
}

func getParamBooleanValue(param, query string) bool {
  params := strings.Split(query, "=")
  if param == params[0] {
    return params[1] == "true"
  }
  return false
}

// ---------------- checkContentType Header  ----------------------------------

func checkContentType(header http.Header, methodType string) (int, string, bool) {
  if contentType, ok := header["Content-Type"]; ok {
    retval := contentType[0] == "application/json"
    if methodType == "get" {
      retval = retval || (contentType[0] == "*/*")
    }
    if retval {
      return 0, "", true
    }
  }
  // responseCode, errorMsgCode, ok
  return 400, "invalid_content_type", false
}

// ---------------- Conditional Header Support ---------------------------------

func itemAlreadyPresent(name, key string) bool {
  if data, _, err := datastore.GetKey(name, key, ""); err == nil {
    return(data != "")
  }
  return false
}

func itemVersionMismatch(name, key, qref string) bool {
  if matchRef, err := strconv.Unquote(qref); err == nil {
    return matchRef != getRefValue(name, key)
  }
  return true
}

func checkConditionalHeaders(header http.Header, name, key string) (int, string, bool) {

  if match, ok := header["If-Match"]; ok {
    if datastore.IsValidKeyRef(name, key, match[0]) == false {
      return 400, "api_bad_request", false
    }
    if itemVersionMismatch(name, key, match[0]) {
      return 412, "item_version_mismatch", false
    }
  } else if match, ok := header["If-None-Match"]; ok {
    if match[0] != strconv.Quote("*") {
      return 400, "item_ref_malformed", false
    } else if itemAlreadyPresent(name, key) {
      return 412, "item_already_present", false
    }
  }

  // responseCode, errorMsgCode, ok
  return 0, "", true
}

// ----------------------------------------------------------------------------


