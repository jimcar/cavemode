package orchio

import (
  "github.com/jimcar/datastore"
)

// ----------------------------------------------------------------------------
//  Name: PutRelation
//  Desc:

func (serv OrchioWriteService) PutRelation(data, name, key, kind, to_c, to_k string) {
  datastore.TraceMsg("==> Orchio::PutRelation:")

  datastore.PutRelation(name, key, kind, to_c, to_k)

  response := updateResponse(serv.ResponseBuilder(), "")
  completeTheResponse(response, 204, "OrchioService::PutRelation")
  return
}

// ----------------------------------------------------------------------------
//  Name: DeleteRelation
//  Desc:

func (serv OrchioWriteService) DeleteRelation(name, key, kind, to_c, to_k string) {
  datastore.TraceMsg("==> Orchio::DeleteRelation:")

  // check for "?purge=true"
  req := serv.Context.Request()
  if getParamBooleanValue("purge", req.URL.RawQuery) {
    datastore.DeleteRelation(name, key, kind, to_c, to_k)
  }
  response := updateResponse(serv.ResponseBuilder(), "delete")
  completeTheResponse(response, 204, "OrchioService::DeleteRelation")
  return
}

