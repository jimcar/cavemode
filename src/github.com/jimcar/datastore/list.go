package datastore

import (
  "strconv"
)

const(
  listLimitDefault =  10
  listLimitMax     = 100
)

func ListLimit(paramValues map[string]string) int {
  limit := listLimitDefault
  if tmplimit, ok := paramValues["limit"]; ok {
    limit, _ = strconv.Atoi(tmplimit)
    if limit > listLimitMax {
      limit = listLimitMax
    } else if limit <= 0 {
      limit = listLimitDefault
    }
  }
  return limit
}

