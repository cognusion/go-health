package health

import (
	"github.com/spf13/cast"

	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Status constants to prevent fat-fingered-oopsies
const (
	OK       = StatusString("OK")
	UP       = StatusString("UP")
	WARNING  = StatusString("WARNING")
	BAD      = StatusString("BAD")
	DOWN     = StatusString("DOWN")
	CRITICAL = StatusString("CRITICAL")
	UNKNOWN  = StatusString("UNKNOWN")
	ERROR    = StatusString("ERROR")
)

// StatusString is a string type for static string consistency
type StatusString = string

// Status is a type used to convey status-related information about a Service, System, or Metric
type Status struct {
	// Name is univerally required
	Name string `json:"name,omitempty"`
	// Status is universally required, one of: OK,WARNING,
	// BAD/ERROR/CRITICAL, or UNKNOWN
	Status string `json:"status,omitempty"`
	// Value is optional for all but Metrics, and is used to convey a
	// numeric-type representation
	Value interface{} `json:"value,omitempty"`
	// ExpectedValue is optional, generally used with Services to represent what
	// Value should be, to understand deviations
	ExpectedValue interface{} `json:"expectedValue,omitempty"`
	// WarnOver is only for Metrics, and is used to represent the Value at which a
	// WARNING state will be triggered
	WarnOver interface{} `json:"warnOver,omitempty"`
	// BadOver is ony for Metrics, and is used to represent the Value at which a
	// CRITICAL state will be triggered
	BadOver interface{} `json:"badOver,omitempty"`
	// TimeStamp is optional, and is used to convey the time the Status or Value
	// was retrieved
	TimeStamp *time.Time `json:"timestamp,omitempty"`
	// TimeOut is optional, and is used to set the amount of time TimeStamp can
	// drift before it is considered stale, tiggering an alert
	TimeOut *time.Duration `json:"timeout,omitempty"`
	// Suffix is optional, and is used appended to Value for metrics
	Suffix string `json:"suffix,omitempty"`
}

// rawStatus is the Status struct without the higher-level time.Time and time.Duration used in
// Status, for marshalling back to JSON
type rawStatus struct {
	// Name is univerally required
	Name string `json:"name,omitempty"`
	// Status is universally required, one of: OK,WARNING,
	// BAD/ERROR/CRITICAL, or UNKNOWN
	Status string `json:"status,omitempty"`
	// Value is optional for all but Metrics, and is used to convey a
	// numeric-type representation
	Value interface{} `json:"value,omitempty"`
	// ExpectedValue is optional, generally used with Services to represent what
	// Value should be, to understand deviations
	ExpectedValue interface{} `json:"expectedValue,omitempty"`
	// WarnOver is only for Metrics, and is used to represent the Value at which a
	// WARNING state will be triggered
	WarnOver interface{} `json:"warnOver,omitempty"`
	// BadOver is ony for Metrics, and is used to represent the Value at which a
	// CRITICAL state will be triggered
	BadOver interface{} `json:"badOver,omitempty"`
	// TimeStamp is optional, and is used to convey the time the Status or Value
	// was retrieved
	TimeStamp *int64 `json:"timestamp,omitempty"`
	// TimeOut is optional, and is used to set the amount of time TimeStamp can
	// drift before it is considered stale, tiggering an alert
	TimeOut *int64 `json:"timeout,omitempty"`
	// Suffix is optional, and is used appended to Value for metrics
	Suffix string `json:"suffix,omitempty"`
}

// MarshalJSON is a custom marshaller for JSON encoding,
// to output TimeStamp and TimeOut as numbers instead of pretty strings.
func (s *Status) MarshalJSON() ([]byte, error) {
	newS := rawStatus{
		Name:          s.Name,
		Status:        s.Status,
		Value:         s.Value,
		ExpectedValue: s.ExpectedValue,
		WarnOver:      s.WarnOver,
		BadOver:       s.BadOver,
		Suffix:        s.Suffix,
	}

	if s.TimeStamp != nil {
		t := s.TimeStamp.UnixMilli()
		newS.TimeStamp = &t
	}

	if s.TimeOut != nil {
		t := s.TimeOut.Milliseconds()
		newS.TimeOut = &t
	}

	return json.Marshal(&newS)
}

// MetricString returns a Nagios Performance Data -compatible representation of Status
func (s *Status) MetricString() string {
	value := cast.ToString(s.Value)
	if s.Suffix != "" {
		value = fmt.Sprintf("%s%s", value, s.Suffix)
	}

	// note we are missing min and max
	return fmt.Sprintf("'%s'=%s;%s;%s;%s;%s", s.Name,
		value, cast.ToString(s.WarnOver),
		cast.ToString(s.BadOver), "", "")
}

// StatusSliceFromJmap is a hacky function that might take a slice of interfaces, and return a same-sized slice of Status
func StatusSliceFromJmap(jmap []interface{}) []Status {
	var statuses = make([]Status, len(jmap))

	c := 0
	for _, r := range jmap {
		jr := lcKeys(cast.ToStringMap(r))

		var (
			ts *time.Time
			to *time.Duration
		)

		if m, ok := jr["timestamp"]; ok {

			mi := cast.ToInt64(m)
			tsx := cast.ToTime(mi)

			if tsx.After(time.Now().Add(time.Hour * 24)) {
				// Safety valve for some Java stamps...
				mi /= 1000
				tsx = cast.ToTime(mi)
			}

			ts = &tsx
		}

		if m, ok := jr["timeout"]; ok {
			tox := cast.ToDuration(m)
			to = &tox
		}

		// Fix a pair of spec messups
		if m, ok := jr["warnvalue"]; ok {
			jr["warnover"] = m
		}
		if m, ok := jr["badvalue"]; ok {
			jr["badover"] = m
		}

		s := Status{
			Name:          cast.ToString(jr["name"]),
			Status:        cast.ToString(jr["status"]),
			Value:         jr["value"],
			ExpectedValue: jr["expectedvalue"],
			WarnOver:      jr["warnover"],
			BadOver:       jr["badover"],
			TimeStamp:     ts,
			TimeOut:       to,
		}
		statuses[c] = s
		c++
	}
	return statuses
}

// Lowercases all of the keys in a map[string]interface{}
func lcKeys(mcmap map[string]interface{}) map[string]interface{} {
	newmap := make(map[string]interface{})
	for k, v := range mcmap {
		newmap[strings.ToLower(k)] = v
	}
	return newmap
}
