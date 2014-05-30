package orchio

import (
  "fmt"
  "strings"
  "os"
  "code.google.com/p/gorest"
  "github.com/jimcar/datastore"
)

// ----------------------------------------------------------------------------
//  Name: getRelations
//  Desc: The underlying utility for services GetRelations 1 through n.

func getRelations(r *gorest.ResponseBuilder, name, key, kind string) string {
  datastore.TraceMsg(fmt.Sprintf("==> Orchio::getRelations: %s/%s/relations/%s", name, key, kind))

  data, _ := datastore.GetRelations(name, key, kind)

  response := updateResponse(r, "")
  completeTheResponse(response, 200, data)
  return ""
}

// ----------------------------------------------------------------------------
//  Name: GetRelations
//  Desc: Service for GET "/v0/$collection/$key/relations/$kind1/../$kind-n"
//        current max value for n is 3.

func (serv OrchioReadService) GetRelations1(name, key, kind string) string {
  return getRelations(serv.ResponseBuilder(), name, key, kind)
}

func (serv OrchioReadService) GetRelations2(name, key, kind1, kind2 string) string {
  kind := strings.Join([]string{kind1, kind2}, "/")
  return getRelations(serv.ResponseBuilder(), name, key, kind)
}

func (serv OrchioReadService) GetRelations3(name, key, kind1, kind2, kind3 string) string {
  kind := strings.Join([]string{kind1, kind2, kind3}, "/")
  return getRelations(serv.ResponseBuilder(), name, key, kind)
}


// ----------------------------------------------------------------------------
//
//  Extending the max graph depth, aka "number of hops".
//
//  The default maximum graph depth is three hops. The max can be doubled by
//  setting env var CAVEMODE_EXTEND_GRAPH_DEPTH=true.
//
//  TODO:
//  A more elegant solution for extending depth support might be created by
//  directly leveraging the net/http mux capabilities.  But for now, the
//  copy-paste-and-edit method can be used to extend depth support to as many
//  levels as needed - in about the same time that it took to write this comment.
//
//  How many hops is enough?

func registerExtendedGraphServices() {
  if os.Getenv("CAVEMODE_EXTEND_GRAPH_DEPTH") == "true" {
    gorest.RegisterServiceOnPath(orchioRoot, new(GraphServiceExt))
  }
}

// -------------------------------------------
//  Orchio Extended Graph Service Definitions

type GraphServiceExt struct {

  gorest.RestService `root:"/"`

  // Get
  getRelations4 gorest.EndPoint `method:"GET" path:"/{name:string}/{key:string}/relations/{kind1:string}/{kind2:string}/{kind3:string}/{kind4:string}" output:"string"`
  getRelations5 gorest.EndPoint `method:"GET" path:"/{name:string}/{key:string}/relations/{kind1:string}/{kind2:string}/{kind3:string}/{kind4:string}/{kind5:string}" output:"string"`
  getRelations6 gorest.EndPoint `method:"GET" path:"/{name:string}/{key:string}/relations/{kind1:string}/{kind2:string}/{kind3:string}/{kind4:string}/{kind5:string}/{kind6:string}" output:"string"`

  // Head
  getRelations4Head gorest.EndPoint `method:"HEAD" path:"/{name:string}/{key:string}/relations/{kind1:string}/{kind2:string}/{kind3:string}/{kind4:string}"`
  getRelations5Head gorest.EndPoint `method:"HEAD" path:"/{name:string}/{key:string}/relations/{kind1:string}/{kind2:string}/{kind3:string}/{kind4:string}/{kind5:string}"`
  getRelations6Head gorest.EndPoint `method:"HEAD" path:"/{name:string}/{key:string}/relations/{kind1:string}/{kind2:string}/{kind3:string}/{kind4:string}/{kind5:string}/{kind6:string}"`
}

// ----------------------------------------------------------------------------
//  Name: GetRelations
//  Desc: Extension services for GET "/v0/$collection/$key/relations/$kind-m/../$kind-n"
//        where [m:n] = [4:6]

func (serv GraphServiceExt) GetRelations4(name, key, kind1, kind2, kind3, kind4 string) string {
  kind := strings.Join([]string{kind1, kind2, kind3, kind4}, "/")
  return getRelations(serv.ResponseBuilder(), name, key, kind)
}

func (serv GraphServiceExt) GetRelations5(name, key, kind1, kind2, kind3, kind4, kind5 string) string {
  kind := strings.Join([]string{kind1, kind2, kind3, kind4, kind5}, "/")
  return getRelations(serv.ResponseBuilder(), name, key, kind)
}

func (serv GraphServiceExt) GetRelations6(name, key, kind1, kind2, kind3, kind4, kind5, kind6 string) string {
  kind := strings.Join([]string{kind1, kind2, kind3, kind4, kind5, kind6}, "/")
  return getRelations(serv.ResponseBuilder(), name, key, kind)
}

// ----------------------------------------------------------------------------
//  Name: GetRelationsHead 4 thru 6
//  Desc: Extension services for HEAD "/v0/$collection/$key/relations/$kind-m/../$kind-n"
//        where [m:n] = [4:6]

func (serv GraphServiceExt) GetRelations4Head(name, key, kind1, kind2, kind3, kind4 string) {
  response := updateResponse(serv.ResponseBuilder(), "")
  completeTheResponse(response, 200, "")
  return
}

func (serv GraphServiceExt) GetRelations5Head(name, key, kind1, kind2, kind3, kind4, kind5 string) {
  response := updateResponse(serv.ResponseBuilder(), "")
  completeTheResponse(response, 200, "")
  return
}

func (serv GraphServiceExt) GetRelations6Head(name, key, kind1, kind2, kind3, kind4, kind5, kind6 string) {
  response := updateResponse(serv.ResponseBuilder(), "")
  completeTheResponse(response, 200, "")
  return
}


