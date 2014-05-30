package orchio

import (
  "github.com/jimcar/datastore"
  "encoding/json"
  "strconv"
)

// Orchio Collection/Key Services

// ----------------------------------------------------------------------------
//  Name: GetKey
//  Desc: service for GET "/v0/$collection/$key"

func (serv OrchioReadService) GetKey(name, key string) string {
  datastore.TraceMsg("==> Orchio::GetKey:")
  return getKey(serv, name, key, "")
}

// ----------------------------------------------------------------------------
//  Name: GetByRef
//  Desc: service for GET "/v0/$collection/$key/refs/$ref"

func (serv OrchioReadService) GetByRef(name, key, ref string) string {
  datastore.TraceMsg("==> Orchio::GetByRef:")
  return getKey(serv, name, key, ref)
}

// ----------------------------------------------------------------------------
//  Name: getKey
//  Desc: underlying utility for the GetKey and GetByRef services

func getKey(serv OrchioReadService, name, key, ref string) string {

  data, mdata, err := datastore.GetKey(name, key, ref)
  logError("Orchio::getKey:", err)

  responseCode := 200
  response := updateResponse(serv.ResponseBuilder(), "")

  // Guard clause for 404 and 500 errors.
  if mdata == "" {
    var msgCode string
    if data == "" {
      responseCode = 404
      msgCode = "items_not_found"
    } else {
      responseCode = 500
      msgCode = "internal_error"
    }
    errorMsgBody := refErrorMsgBody(responseCode, msgCode, name, key, ref)
    completeTheResponse(response, responseCode, errorMsgBody)
    return ""
  }

  if ref == "" {
    var tmpMdata datastore.Metadata
    if err := json.Unmarshal([]byte(mdata), &tmpMdata); err == nil {
      ref = tmpMdata.Ref
    }
  }
  response.AddHeader("Content-Location", locationString(name, key, ref))
  response.AddHeader("Etag", strconv.Quote(ref))
  completeTheResponse(response, responseCode, data)
  return ""
}

// ----------------------------------------------------------------------------
//  Name: PutKey
//  Desc: service for PUT "/v0/$collection/$key"

func (serv OrchioWriteService) PutKey(data, name, key string) {
  datastore.TraceMsg("==> Orchio::PutKey:")

  req := serv.Context.Request()
  response := updateResponse(serv.ResponseBuilder(), "")

  var responseCode int
  var msgCode      string
  var ok           bool
  if responseCode, msgCode, ok = checkContentType(req.Header, "put"); ok {
    responseCode, msgCode, ok = checkConditionalHeaders(req.Header, name, key)
  }

  if !ok {
    errorMsgBody := keyErrorMsgBody(responseCode, msgCode, name, key)
    completeTheResponse(response, responseCode, errorMsgBody)
    return
  }

  ref := datastore.PutKey(data, name, key)
  response.AddHeader("Location", locationString(name, key, ref))
  response.AddHeader("Etag", strconv.Quote(ref))
  completeTheResponse(response, 201, "")
  return
}

// ----------------------------------------------------------------------------
//  Name: DeleteKey
//  Desc: service for DELETE "/v0/$collection/$key"

func (serv OrchioWriteService) DeleteKey(name, key string) {
  datastore.TraceMsg("==> Orchio::DeleteKey:")

  req := serv.Context.Request()
  response := updateResponse(serv.ResponseBuilder(), "delete")

  if responseCode, msgCode, ok := checkConditionalHeaders(req.Header, name, key); !ok {
    errorMsgBody := keyErrorMsgBody(responseCode, msgCode, name, key)
    completeTheResponse(response, responseCode, errorMsgBody)
    return
  }

  // check for "?purge=true"
  purge := getParamBooleanValue("purge", req.URL.RawQuery)
  datastore.DeleteKey(name, key, purge)
  completeTheResponse(response, 204, "")
  return
}

