package datastore

import (
  "sort"
)

type By func(r1, r2 *Result) bool

func (by By) Sort(results []Result) {
  rs := &resultSorter{
    results: results,
    by:      by, // The Sort method's receiver is the function (closure)
  }              // that defines the sort order.
  sort.Sort(rs)
}

type resultSorter struct {
  results []Result
  by      func(r1, r2 *Result) bool // Closure used in the Less method.
}

func (s *resultSorter) Len() int {
  return len(s.results)
}

func (s *resultSorter) Swap(i, j int) {
  s.results[i], s.results[j] = s.results[j], s.results[i]
}

func (s *resultSorter) Less(i, j int) bool {
  return s.by(&s.results[i], &s.results[j])
}

