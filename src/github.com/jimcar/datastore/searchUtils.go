package datastore

import (
  "strings"
  "regexp"
)

// ----------------------------------------------------------------------------
//  Name: getSearchTerms
//  Desc:

func getSearchTerms(queryStr string) ([]string, string) {
  // ("==> datastore:getSearchTerms")

  r := strings.NewReplacer("%20", " ")
  queryStr = r.Replace(queryStr)

  sep, logicOp := "", ""
  if strings.Contains(queryStr, " AND ") {
    sep, logicOp = " AND ", "and"
  } else if strings.Contains(queryStr, " OR ") {
    sep, logicOp = " OR ", "or"
  }

  var searchTerms []string
  if sep == "" {
    searchTerms = []string{queryStr}
  } else {
    searchTerms = strings.Split(queryStr, sep)
  }
  // TraceMsg("searchTerms: %v, logicOp: %s", searchTerms, logicOp)

  return searchTerms, logicOp
}

// ----------------------------------------------------------------------------
//  Name: searchTermsMatch
//  Desc: Check whether the data (value) matches the search terms.

func searchTermsMatch(value ResultValue, searchTerms []string, logicOp string) bool {
  // TraceMsg("==> datastore:searchTermsMatch" )

  searchTermCount := len(searchTerms)
  matchCount := 0
  matchFound := false

  for _, term := range searchTerms {

    if strings.Contains(term, ":") { // specified field

      kvpair := strings.Split(term, ":")
      if len(kvpair) != 2 {
        continue
      }
      k, v := kvpair[0], kvpair[1]

      if logicOp == "and" {
        if !isMatch(value[k], v) {
          break
        }
        matchCount++
        if matchCount == searchTermCount {
          matchFound = true
          break
        }
      } else {
        if isMatch(value[k], v) {
          matchFound = true
          break
        }
      }

    } else {  // default, or unspecified, field

      for _, v := range value {
        if isMatch(v, term) {
          if logicOp == "and" {
            matchCount++
            if matchCount == searchTermCount {
              matchFound = true
            }
          } else {
            matchFound = true
          }
          break
        }
      }

      if matchFound == true {
        break
      }
    }
  }

  return matchFound
}

// ----------------------------------------------------------------------------
//  Name: isMatch
//  Desc: case-insensitive, wildcards, regexp matching, field matching

func isMatch(s1, s2 string) bool {

  // Return immediate failure for empty value.
  if s1 == "" {
    return false
  }

  // Return immediate success for wildcard.
  if s2 == "*" {
    return true
  }

  // Make case-insensitive
  s1 = strings.ToLower(s1)
  s2 = strings.ToLower(s2)

  // Return immediate success for full match (case-insensitive).
  if s1 == s2 {
    return true
  }

  // Replace ":", ";", "," and "/" chars with whitespace.
  // Split s1 into whitespace separated fields.
  r := strings.NewReplacer(":", " ", ";", " ", ",", " ", "/", " ")
  fields := strings.Fields(r.Replace(s1))

  // Adjust any wildcards.
  if strings.ContainsAny(s2, "?*") {
    r := strings.NewReplacer("?", ".?", "*", ".*")
    s2 = r.Replace(s2)
  }

  // Default return value.
  retval := false

  // Create the regexp match checker (type *regexp.Regexp).
  checkMatch := regexp.MustCompile(s2)

  // Match the search term (s2) against each of s1's fields.
  for _, field := range fields {
    if checkMatch.MatchString(field) {
      retval = true
      break
    }
  }

  return retval
}



