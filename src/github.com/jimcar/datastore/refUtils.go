package datastore

import (
  "fmt"
  "strconv"
  "math/rand"
  "time"
)

// ----------------------------------------------------------------------------
//  Name: generateRandomRefValue
//  Desc:

func generateRandomRefValue() string {

  var r *rand.Rand
  sval := ""

  for len(sval) < 16 {
    r = rand.New(rand.NewSource(int64(time.Now().UnixNano())))
    sval = fmt.Sprintf("%x", r.Int63())
  }
  return sval
}

// ----------------------------------------------------------------------------
//  Name: generateRefValue
//  Desc:

func generateRefValue() string {
  ref := ""
  refTable := getCollectionHandle("RefTable")
  eventRefTable := getCollectionHandle("EventRefTable")

  isUnique := false
  for isUnique == false {
    ref = generateRandomRefValue()
    rdata, _ := refTable.Get(ro, []byte(ref))
    edata, _ := eventRefTable.Get(ro, []byte(ref))
    if rdata == nil && edata == nil {
      isUnique = true
    }
  }
  return ref
}

// ----------------------------------------------------------------------------
//  Name: isValidKeyRef
//  Desc: Verify that ref is included in AllRefsTable/collection/key

func isValidKeyRef(name, key, ref string) bool {
  allRefsKey := createKey(name, key)
  return isValidItem("AllRefsTable", allRefsKey, ref)
}

// ----------------------------------------------------------------------------
//  Name: IsValidKeyRef
//  Desc: Publicly available; unquotes ref val before calling isValidKeyRef

func IsValidKeyRef(name, key, qref string) bool {
  if ref, err := strconv.Unquote(qref); err == nil {
    return isValidKeyRef(name, key, ref)
  }
  return false
}
