package orchio

// ----------------------------------------------------------------------------
//  Type definitions for error message body

type Item struct {
  Collection string      `json:"collection"`
  Key        string      `json:"key"`
  Type       string      `json:"type,omitempty"`         // event only
  Ref        string      `json:"ref,omitempty"`          // ref only
  Timestamp  string      `json:"timestamp,omitempty"`    // event only
  Ordinal    int         `json:"ordinal,omitempty"`      // event only
}

type Detail struct {
  Items      []Item      `json:"items"`                  // 404 only
}

type ErrMsgBody struct {
  Message    string      `json:"message"`
  Details    Detail      `json:"details,omitempty"`      // 404 only
  Code       string      `json:"code"`
}

