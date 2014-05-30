package datastore

import (
  "fmt"
)

// ----------------------------------------------------------------------------
//  Name: offsetResultPage
//  Desc: Returns "next" or "prev" result page strings for search and listRefs.

func offsetResultPage(name string, limit int, offset string) string {
  return fmt.Sprintf("/%s/%s?limit=%v&offset=%v", "v0", name, limit, offset)
}

// ----------------------------------------------------------------------------
//  Name: keyResultPage
//  Desc: Returns "next" result page string for listKeys.

func keyResultPage(name string, limit int, afterKey string) string {
  return fmt.Sprintf("/%s/%s?limit=%v&afterKey=%s", "v0", name, limit, afterKey)
}

// ----------------------------------------------------------------------------
//  Name: eventResultPage
//  Desc: Returns "next" result page string for listEvents.

func eventResultPage(name, key, etype string, limit int, beforeEvent string) string {
  return fmt.Sprintf("/%s/%s/%s/events/%s?limit=%v&beforeEvent=%s",
                     "v0", name, key, etype, limit, beforeEvent)
}

