package orchio

import (
  "fmt"
  "strconv"
  "github.com/jimcar/datastore"
)

// ----------------------------------------------------------------------------
//  Name: PostEvent (create)
//  Desc: service for POST "/v0/$collection/$key/events/$type"

func (serv OrchioWriteService) PostEvent(data, name, key, etype string) {
  datastore.TraceMsg("==> Orchio::PostEvent:")
  putEvent(serv, data, name, key, etype, "", 0)
  return
}

// ----------------------------------------------------------------------------
//  Name: PostEventTs (create w/timestamp)
//  Desc: service for POST "/v0/$collection/$key/events/$type/$timestamp"

func (serv OrchioWriteService) PostEventTs(data, name, key, etype, timestamp string) {
  datastore.TraceMsg("==> Orchio::PostEventTS:")
  putEvent(serv, data, name, key, etype, timestamp, 0)
  return
}

// ----------------------------------------------------------------------------
//  Name: PutEvent (update)
//  Desc: service for PUT "/v0/$collection/$key/events/$type/$timestamp/$ordinal"

func (serv OrchioWriteService) PutEvent(data, name, key, etype, timestamp string, ordinal int) {
  datastore.TraceMsg("==> Orchio::PutEvent:")
  putEvent(serv, data, name, key, etype, timestamp, ordinal)
  return
}

// ----------------------------------------------------------------------------
//  Name: putEvent
//  Desc: underlying utility for PostEvent, PostEventTS and PutEvent services

func putEvent(serv OrchioWriteService, data, name, key, etype, timestamp string, ordinal int) {

  req := serv.Context.Request()
  response := updateResponse(serv.ResponseBuilder(), "")

  var responseCode int
  var msgCode      string
  var ok           bool
  if responseCode, msgCode, ok = checkContentType(req.Header, "post"); ok {
    responseCode = 201
    // Check conditional headers for PUT (update)
    if timestamp != "" && ordinal != 0 {
      if responseCode, msgCode, ok = checkConditionalHeaders(req.Header, name, key); ok {
        responseCode = 204
      }
    }
  }

  if !ok {
    errorMsgBody := eventErrorMsgBody(responseCode, msgCode, name, key, etype, timestamp, ordinal)
    completeTheResponse(response, responseCode, errorMsgBody)
    return
  }

  ref, timestamp, ordinal := datastore.PutEvent(data, name, key, etype, timestamp, ordinal)
  response.AddHeader("Location", eventLocationString(name, key, etype, timestamp, ordinal))
  response.AddHeader("Etag", strconv.Quote(ref))
  completeTheResponse(response, responseCode, "OrchioService::putEvent")
  return
}


// ----------------------------------------------------------------------------
//  Name: GetEvent
//  Desc: service for GET "/v0/$collection/$key/events/$type/$timestamp/$ordinal"

func (serv OrchioReadService) GetEvent(name, key, etype, timestamp string, ordinal int) string {
  datastore.TraceMsg("==> Orchio::GetEvent")

  data, mdata, err := datastore.GetEvent(name, key, etype, timestamp, ordinal)
  logError("Orchio::GetEvent", err)

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
    errorMsgBody := eventErrorMsgBody(responseCode, msgCode, name, key, etype, timestamp, ordinal)
    completeTheResponse(response, responseCode, errorMsgBody)
    return ""
  }

  completeTheResponse(response, responseCode, data)
  return ""
}

// ----------------------------------------------------------------------------
//  Name: ListEvents
//  Desc: service for GET "/v0/$collection/$key/events/$type"

func (serv OrchioReadService) ListEvents(name, key, etype string) string {
  datastore.TraceMsg("==> Orchio::GetEvents")

  // Check for range params (timestamp/ordinal)
  req := serv.Context.Request()
  paramValues := getParamValues(req.URL.RawQuery)

  response := updateResponse(serv.ResponseBuilder(), "")

  data, next, _, _ := datastore.ListEvents(name, key, etype, paramValues)

  if next != "" {
    response.AddHeader("Link", fmt.Sprintf("<%s>", next))
  }
  completeTheResponse(response, 200, data)
  return ""
}

// ----------------------------------------------------------------------------
//  Name: DeleteEvent
//  Desc: service for DELETE "/v0/$collection/$key/events/$etype/$timestamp/$ordinal"

func (serv OrchioWriteService) DeleteEvent(name, key, etype, timestamp string, ordinal int) {
  datastore.TraceMsg("==> Orchio::DeleteEvent:")

  req := serv.Context.Request()
  response := updateResponse(serv.ResponseBuilder(), "delete")

  if responseCode, msgCode, ok := checkConditionalHeaders(req.Header, name, key); !ok {
    errorMsgBody := eventErrorMsgBody(responseCode, msgCode, name, key, etype, timestamp, ordinal)
    completeTheResponse(response, responseCode, errorMsgBody)
    return
  }

  // check for "?purge=true"
  purge := getParamBooleanValue("purge", req.URL.RawQuery)
  datastore.DeleteEvent(name, key, etype, timestamp, ordinal, purge)

  completeTheResponse(response, 204, "")
  return
}




