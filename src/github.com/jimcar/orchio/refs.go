package orchio

import (
  "fmt"
  "github.com/jimcar/datastore"
)

// Orchio List Refs Services

// ----------------------------------------------------------------------------
//  Name: ListRefs
//  Desc: service for GET "/v0/$collection/$key/refs/"

func (serv OrchioReadService) ListRefs(name, key string) string {
  datastore.TraceMsg(fmt.Sprintf("==> Orchio::ListRefs: %s", name))

  req := serv.Context.Request()
  paramValues := getParamValues(req.URL.RawQuery)

  responseCode := 200
  response := updateResponse(serv.ResponseBuilder(), "")

  data, next, prev, err := datastore.ListRefs(name, key, paramValues)
  outputError("Orchio::ListRefs:", err)

  if next != "" {
    response.AddHeader("Link", fmt.Sprintf("<%s>; rel=\"next\"", next))
  }

  if prev != "" {
    response.AddHeader("Link", fmt.Sprintf("<%s>; rel=\"prev\"", prev))
  }

  completeTheResponse(response, responseCode, data)
  return ""
}


