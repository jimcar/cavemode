package datastore

import (
  "fmt"
  "strings"
  "strconv"
)

// ----------------------------------------------------------------------------
//  Name: GetRelations
//  Desc:

func GetRelations(name, key, kind string) (string, error) {

  kinds := strings.Split(kind, "/")
  depth := len(kinds)

  relations := getRelations(name, key, kinds[0])
  nextRelations := []string{}

  for currentDepth := 1; currentDepth < depth; currentDepth++ {

    for _, relation := range relations {
      nextRelations = append(nextRelations, relation)
    }

    relations = []string{}
    for _, relation := range nextRelations {
      from := strings.Split(relation, "/")
      for _, r := range getRelations(from[0], from[1], kinds[currentDepth]) {
        relations = appendRelation(relations, r)
      }
    }

    nextRelations = []string{}
  }

  // Create the response body with whatever relations dropped out.
  results, err := relationResults(relations)
  body := listResponse(ResponseBody{len(results), results, "", "", 0})

  return body, err
}

// ----------------------------------------------------------------------------
//  Name: PutRelation
//  Desc:

func PutRelation(name, key, kind, to_c, to_k string) error {

  // Assign the next ordinal value
  ordinal := strconv.Itoa(getNextOrdinalValue(name, key))

  // Update collection/key/kind with ordinal value.
  graphTable := graphTableName(name, key)
  addItemToList(graphTable, kind, ordinal)

  // Update collection/key/ordinal with relation data
  relation := fmt.Sprintf("%s/%s", to_c, to_k)
  relationTable := getCollectionHandle(relationTableName(name, key))
  relationTable.Put(wo, []byte(ordinal), []byte(relation))

  return error(nil)
}

// ----------------------------------------------------------------------------
//  Name: DeleteRelation
//  Desc:

func DeleteRelation(name, key, kind, to_c, to_k string) error {

  relTableName := relationTableName(name, key)
  relationTable := getCollectionHandle(relTableName)
  deleteMe := fmt.Sprintf("%s/%s", to_c, to_k)

  graphTable := graphTableName(name, key)
  ordinals := getItems(graphTable, kind)

  for _, ordinal := range ordinals {
    relation, _ := relationTable.Get(ro, []byte(ordinal))
    if deleteMe == string(relation) {
      relationTable.Delete(wo, []byte(ordinal))
      deleteItemFromList(graphTable, kind, ordinal)
    }
  }

  return error(nil)
}

