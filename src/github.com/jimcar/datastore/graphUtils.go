package datastore

import (
  "strings"
  "encoding/json"
)

// ----------------------------------------------------------------------------
//  Name: relationResults
//  Desc:

func relationResults(relations []string) ([]Result, error) {

  results := []Result{}

  refs := getCollectionHandle("RefTable")

  for _, relation := range relations {
    // TraceMsg(fmt.Sprintf("\tFINAL to_relation: %s", relation))

    jmc := strings.Split(relation, "/")
    name, key := jmc[0], jmc[1]

    collection := getCollectionHandle(name)
    tmpref, _ := collection.Get(ro, []byte(key))
    ref := string(tmpref)

    r, value := Ref{}, ResultValue{}
    rdata, _ := refs.Get(ro, []byte(ref))
    json.Unmarshal(rdata, &r)
    json.Unmarshal([]byte(r.Value), &value)

    path := PathValue{r.Mdata.Collection, r.Mdata.Key, "", r.Mdata.Ref, 0, 0, false}
    results = append(results, Result{path, value, 0, 0, 0})
  }

  return results, error(nil)
}


// ----------------------------------------------------------------------------
//  Name: getRelations
//  Desc:

func getRelations(name, key, kind string) []string {

  graphTable := graphTableName(name, key)
  ordinals := getItems(graphTable, kind)

  relations := []string{}
  relationTable := getCollectionHandle(relationTableName(name, key))

  for _, ordinal := range ordinals {
    relation, _ := relationTable.Get(ro, []byte(ordinal))
    relations = append(relations, string(relation))
  }

  return relations
}

// ----------------------------------------------------------------------------
//  Name: appendRelation
//  Desc:

func appendRelation(relations []string, newRelation string) []string {

  for _, relation := range relations {
    if newRelation == relation {
      return relations
    }
  }

  relations = append(relations, newRelation)
  return relations
}



