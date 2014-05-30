package orchio

import (
  "code.google.com/p/gorest"
  "net/http"
)

const orchioRoot = "v0"

func OrchioMain() {
  gorest.RegisterServiceOnPath(orchioRoot, new(OrchioReadService))
  gorest.RegisterServiceOnPath(orchioRoot, new(OrchioWriteService))
  registerExtendedGraphServices()
  http.Handle("/", gorest.Handle())
  http.ListenAndServe(":8787", nil)
}

// ---------------------------------
//  Orchio Read Service Definitions
// ---------------------------------

type OrchioReadService struct {

  gorest.RestService `root:"/"`

  // ----------- GET -----------
  getCollection     gorest.EndPoint `method:"GET" path:"/{name:string}" output:"string"`
  getKey            gorest.EndPoint `method:"GET" path:"/{name:string}/{key:string}" output:"string"`
  getByRef          gorest.EndPoint `method:"GET" path:"/{name:string}/{key:string}/refs/{ref:string}" output:"string"`
  listRefs          gorest.EndPoint `method:"GET" path:"/{name:string}/{key:string}/refs" output:"string"`
  getEvent          gorest.EndPoint `method:"GET" path:"/{name:string}/{key:string}/events/{etype:string}/{timestamp:string}/{ordinal:int}" output:"string"`
  listEvents        gorest.EndPoint `method:"GET" path:"/{name:string}/{key:string}/events/{etype:string}" output:"string"`
  getRelations1     gorest.EndPoint `method:"GET" path:"/{name:string}/{key:string}/relations/{kind:string}" output:"string"`
  getRelations2     gorest.EndPoint `method:"GET" path:"/{name:string}/{key:string}/relations/{kind1:string}/{kind2:string}" output:"string"`
  getRelations3     gorest.EndPoint `method:"GET" path:"/{name:string}/{key:string}/relations/{kind1:string}/{kind2:string}/{kind3:string}" output:"string"`

  // ----------- HEAD ----------
  getCollectionHead gorest.EndPoint `method:"HEAD" path:"/{name:string}"`
  getKeyHead        gorest.EndPoint `method:"HEAD" path:"/{name:string}/{key:string}"`
  getByRefHead      gorest.EndPoint `method:"HEAD" path:"/{name:string}/{key:string}/refs/{ref:string}"`
  getEventHead      gorest.EndPoint `method:"HEAD" path:"/{name:string}/{key:string}/events/{etype:string}/{timestamp:string}/{ordinal:int}"`
  getEventsHead     gorest.EndPoint `method:"HEAD" path:"/{name:string}/{key:string}/events/{etype:string}"`
  getRelations1Head gorest.EndPoint `method:"HEAD" path:"/{name:string}/{key:string}/relations/{kind:string}"`
  getRelations2Head gorest.EndPoint `method:"HEAD" path:"/{name:string}/{key:string}/relations/{kind1:string}/{kind2:string}"`
  getRelations3Head gorest.EndPoint `method:"HEAD" path:"/{name:string}/{key:string}/relations/{kind1:string}/{kind2:string}/{kind3:string}"`
}

// ----------------------------------
//  Orchio Write Service Definitions
// ----------------------------------

type OrchioWriteService struct {

  gorest.RestService `root:"/"`

  // ----------- POST ----------
  postEvent         gorest.EndPoint `method:"POST" path:"/{name:string}/{key:string}/events/{etype:string}" postdata:"string"`
  postEventTs       gorest.EndPoint `method:"POST" path:"/{name:string}/{key:string}/events/{etype:string}/{timestamp:string}" postdata:"string"`

  // ----------- PUT -----------
  putKey            gorest.EndPoint `method:"PUT" path:"/{name:string}/{key:string}/" postdata:"string" output:"string"`
  putEvent          gorest.EndPoint `method:"PUT" path:"/{name:string}/{key:string}/events/{etype:string}/{timestamp:string}/{ordinal:int}" postdata:"string"`
  putRelation       gorest.EndPoint `method:"PUT" path:"/{name:string}/{key:string}/relation/{kind:string}/{to_c:string}/{to_k:string}" postdata:"string"`

  // ----------- DELETE -----------
  deleteCollection  gorest.EndPoint `method:"DELETE" path:"/{name:string}"`
  deleteKey         gorest.EndPoint `method:"DELETE" path:"/{name:string}/{key:string}"`
  deleteEvent       gorest.EndPoint `method:"DELETE" path:"/{name:string}/{key:string}/events/{etype:string}/{timestamp:string}/{ordinal:int}"`
  deleteRelation    gorest.EndPoint `method:"DELETE" path:"/{name:string}/{key:string}/relation/{kind:string}/{to_c:string}/{to_k:string}"`
}

