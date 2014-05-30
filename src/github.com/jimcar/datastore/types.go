package datastore

// ----------------------------------------------------------------------------
//  Type definition for KEY or EVENT data

type ResultValue map[string]string

// ----------------------------------------------------------------------------
//  Type definitions for LIST response (key or event data)

type PathValue struct {
  Collection string      `json:"collection"`
  Key        string      `json:"key"`
  Type       string      `json:"type,omitempty"`         // event only
  Ref        string      `json:"ref"`
  Timestamp  int         `json:"timestamp,omitempty"`    // event only
  Ordinal    int         `json:"ordinal,omitempty"`      // event only
  Tombstone  bool        `json:"tombstone,omitempty"`    // refs only
}

type Result struct {
  Path       PathValue   `json:"path"`
  Value      ResultValue `json:"value,omitempty"`
  Timestamp  int         `json:"timestamp,omitempty"`    // event only
  Ordinal    int         `json:"ordinal,omitempty"`      // event only
  Reftime    int         `json:"reftime,omitempty"`      // refs only
}

type ResponseBody struct {
  Count      int         `json:"count"`
  Results    []Result    `json:"results"`
  Next       string      `json:"next,omitempty"`
  Prev       string      `json:"prev,omitempty"`         // refs, search only
  TotalCount int         `json:"total_count,omitempty"`  // search only
}

// ----------------------------------------------------------------------------
//  Type definitions for ref data/metadata

type Metadata struct {
  Collection string      `json:"collection"`
  Key        string      `json:"key"`
  Ref        string      `json:"ref"`
  Datestamp  string      `json:"datestamp"`
  Type       string      `json:"type,omitempty"`         // event only
  Timestamp  string      `json:"timestamp,omitempty"`    // event only
  Ordinal    int         `json:"ordinal,omitempty"`      // event only
  Tombstone  bool        `json:"tombstone,omitempty"`    // refs only
}

type Ref struct {
  Mdata      Metadata    `json:"mdata"`
  Value      string      `json:"value"`
}


