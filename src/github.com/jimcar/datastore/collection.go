package datastore


// ----------------------------------------------------------------------------
//  Name: GetCollection
//  Desc: search or list the collection, returns response body, next, error.

func GetCollection(name string, params map[string]string) (string, string, string, error) {

  if _, ok := params["query"]; ok {

    return searchCollection(name, params)

  } else {

    return listKeys(name, params)
  }
}

// ----------------------------------------------------------------------------
//  Name: DeleteCollection
//  Desc:

func DeleteCollection(name string, force bool) error {

  if force {
    return destroyCollection(name)
  }

  return error(nil)
}
