package orchio

import (
  "strconv"
)

// ----------------------------------------------------------------------------
//  Name: GetCollectionHead
//  Desc: service for HEAD "/v0/$collection"

func (serv OrchioReadService) GetCollectionHead(name string) {
  response := updateResponse(serv.ResponseBuilder(), "")
  completeTheResponse(response, 200, "")
  return
}

// ----------------------------------------------------------------------------
//  Name: GetKeyHead
//  Desc: service for HEAD "/v0/$collection/$key"

func (serv OrchioReadService) GetKeyHead(name, key string) {

  response := updateResponse(serv.ResponseBuilder(), "")
  ref := getRefValue(name, key)
  if ref == "" {
    completeTheResponse(response, 404, "")
    return
  }

  response.AddHeader("Content-Location", locationString(name, key, ref))
  response.AddHeader("Etag", strconv.Quote(ref))
  completeTheResponse(response, 200, "")
  return
}

// ----------------------------------------------------------------------------
//  Name: GetByRefHead
//  Desc: service for HEAD "/v0/$collection/$key/refs/$ref"

func (serv OrchioReadService) GetByRefHead(name, key, ref string) {

  response := updateResponse(serv.ResponseBuilder(), "")
  responseCode := checkRefValue(name, key, ref)
  if responseCode != 0 {
    completeTheResponse(response, responseCode, "")
    return
  }

  response.AddHeader("Content-Location", locationString(name, key, ref))
  response.AddHeader("Etag", strconv.Quote(ref))
  completeTheResponse(response, 200, "")
  return
}

// ----------------------------------------------------------------------------
//  Name: GetEventHead
//  Desc: service for HEAD "/v0/$collection/$key/events/$etype/$timestamp/$ordinal"

func (serv OrchioReadService) GetEventHead(name, key, etype, timestamp string, ordinal int) {
  response := updateResponse(serv.ResponseBuilder(), "")
  response.AddHeader("Location", eventLocationString(name, key, etype, timestamp, ordinal))
  completeTheResponse(response, 200, "")
  return
}

// ----------------------------------------------------------------------------
//  Name: GetEventsHead
//  Desc: service for HEAD "/v0/$collection/$key/events/$etype"

func (serv OrchioReadService) GetEventsHead(name, key, etype string) {
  response := updateResponse(serv.ResponseBuilder(), "")
  completeTheResponse(response, 200, "")
  return
}

// ----------------------------------------------------------------------------
//  Name: GetRelationsHead1
//  Desc: service for HEAD "/v0/$collection/$key/relations/$kind"

func (serv OrchioReadService) GetRelations1Head(name, key, kind string) {
  response := updateResponse(serv.ResponseBuilder(), "")
  completeTheResponse(response, 200, "")
  return
}

// ----------------------------------------------------------------------------
//  Name: GetRelationsHead2
//  Desc: service for HEAD "/v0/$collection/$key/relations/$kind1/$kind2"

func (serv OrchioReadService) GetRelations2Head(name, key, kind1, kind2 string) {
  response := updateResponse(serv.ResponseBuilder(), "")
  completeTheResponse(response, 200, "")
  return
}

// ----------------------------------------------------------------------------
//  Name: GetRelationsHead3
//  Desc: service for HEAD "/v0/$collection/$key/relations/$kind1/$kind2/$kind3"

func (serv OrchioReadService) GetRelations3Head(name, key, kind1, kind2, kind3 string) {
  response := updateResponse(serv.ResponseBuilder(), "")
  completeTheResponse(response, 200, "")
  return
}

