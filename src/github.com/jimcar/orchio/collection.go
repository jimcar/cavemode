package orchio

import (
  "fmt"
  "github.com/jimcar/datastore"
)

// Orchio Collection Services

// ----------------------------------------------------------------------------
//  Name: GetCollection
//  Desc: service for GET "/v0/$collection"

func (serv OrchioReadService) GetCollection(name string) string {
  datastore.TraceMsg(fmt.Sprintf("==> Orchio::GetCollection: %s\n", name))

  req := serv.Context.Request()
  paramValues := getParamValues(req.URL.RawQuery)

  responseCode := 200
  response := updateResponse(serv.ResponseBuilder(), "")

  data, next, prev, err := datastore.GetCollection(name, paramValues)
  outputError("Orchio::GetCollection:", err)

  if next != "" {
    response.AddHeader("Link", fmt.Sprintf("<%s>; rel=\"next\"", next))
  }

  if _, ok := paramValues["query"]; ok && prev != "" {
    response.AddHeader("Link", fmt.Sprintf("<%s>; rel=\"prev\"", prev))
  }

  completeTheResponse(response, responseCode, data)
  return ""
}

// ----------------------------------------------------------------------------
//  Name: DeleteCollection
//  Desc: service for DELETE "/v0/$collection"

func (serv OrchioWriteService) DeleteCollection(name string) {
  datastore.TraceMsg("==> Orchio::DeleteCollection:")

  req := serv.Context.Request()
  response := updateResponse(serv.ResponseBuilder(), "delete")

  // check for "?force=true"
  force := getParamBooleanValue("force", req.URL.RawQuery)
  datastore.DeleteCollection(name, force)
  completeTheResponse(response, 204, "")
  return
}

