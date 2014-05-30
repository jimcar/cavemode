package datastore

import (
  "fmt"
  "time"
  "strings"
  "strconv"
)

// ----------------------------------------------------------------------------
//  Name: generateDatestamp
//  Desc: returns current datestamp string

func generateDatestamp() string {
  r := strings.NewReplacer("UTC", "GMT")
  return r.Replace(time.Now().UTC().Format(time.RFC1123))
}

// ----------------------------------------------------------------------------
//  Name: getDatestamp
//  Desc: accepts timestamp string (milliseconds since Unix epoch)
//        returns datestamp string

func getDatestamp(timestamp string) string {
  msec, _ := strconv.Atoi(timestamp)
  mytime := time.Unix(int64(msec/1000), 0)
  r := strings.NewReplacer("UTC", "GMT")
  return r.Replace(mytime.UTC().Format(time.RFC1123))
}

// ----------------------------------------------------------------------------
//  Name: generateTimestamp
//  Desc: returns current timestamp string (milliseconds since Unix epoch)

func generateTimestamp() string {
  return strconv.FormatInt(time.Now().UTC().UnixNano()/1000000, 10)
}

// ----------------------------------------------------------------------------
//  Name: getTimestamp
//  Desc: accepts timestamp string (various supported formats)
//        returns timestamp string (milliseconds since Unix epoch)

func getTimestamp(timestamp string) string {
  return strconv.FormatInt(convertTime(timestamp), 10)
}

// ----------------------------------------------------------------------------
//  Name: convertTime
//  Desc: accepts timestamp string (various supported formats)
//        returns timestamp int64  (milliseconds since Unix epoch)

func convertTime(timestr string) int64 {

  // Return immediately if timestr represents ms since Unix epoch.
  if mytime, err := strconv.Atoi(timestr); err == nil {
    return int64(mytime)
  }

  // Try parsing timestr with all of the supported formats.
  // Return immediately if timestr is successfully converted.
  var mytime time.Time
  var err error
  for j, _ := range timeFormats {
    if mytime, err = time.Parse(timeFormats[j], timestr); err == nil {
      return mytime.UnixNano()/1000000
    }
  }

  // Return failure.
  fmt.Printf("Error: unable to parse '%s'\n", timestr)
  return 0
}

// ----------------------------------------------------------------------------
//  Supported time formats.

var timeFormats = []string{
  ISO8601basic,
  ISO8601basicZ,
  ISO8601basicMs,
  ISO8601basicMsZ,
  ISO8601extended,
  ISO8601extendedZ,
  ISO8601extendedMs,
  ISO8601extendedMsZ,
  time.RFC1123,       // "Mon, 02 Jan 2006 15:04:05 MST"
  time.RFC1123Z,      // "Mon, 02 Jan 2006 15:04:05 -0700"
  RFC1036,
  RFC1036Z,
  time.RFC850,        // "Monday, 02-Jan-06 15:04:05 MST"
  RFC850Z,
  time.ANSIC,         // "Mon Jan _2 15:04:05 2006"
}


// ----------------------------------------------------------------------------
//  Define constants for the layout formats that aren't included in pkg "time".

const(
  // ISO8601 Basic
  ISO8601basic       = "2006102T150405Z"
  ISO8601basicZ      = "20060102T150405-0700"
  ISO8601basicMs     = "20060102T150405.000Z"
  ISO8601basicMsZ    = "20060102T150405.000-0700"

  // ISO8601 Extended
  ISO8601extended    = "2006-01-02T15:04:05Z"
  ISO8601extendedZ   = "2006-01-02T15:04:05-07:00"
  ISO8601extendedMs  = "2006-01-02T15:04:05.000Z"
  ISO8601extendedMsZ = "2006-01-02T15:04:05.000-07:00"

  // RFC1036
  RFC1036            = "Monday, 02-Jan-2006 15:04:05 MST"
  RFC1036Z           = "Monday, 02-Jan-2006 15:04:05 -0700"

  // RFC850Z
  RFC850Z            = "Monday, 02-Jan-06 15:04:05 -0700"

  // ASCTIME (UTC Enforced)
  ASCTIME            = "Mon Jan _2 15:04:05 2006"
)


