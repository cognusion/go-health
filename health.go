// Package health provides standardized means for generating healthchecks, and consuming them from JSON.
// The included schema files allow for validation of generation, regardless of the source.
package health

import (
	"github.com/spf13/cast"
	jschema "github.com/xeipuuv/gojsonschema"

	"encoding/json"
	"fmt"
)

// JSON is an encapsulating type for "jmap"-based structures
type JSON map[string]interface{}

// Check is a type used to define healthcheck statuses
type Check struct {
	// OverallStatus must be one of OK,WARNING,BAD/ERROR/CRITICAL, or UNKNOWN
	OverallStatus string                 `json:"overallStatus,omitempty"`
	Services      []Status               `json:"services,omitempty"`
	Systems       []Status               `json:"systems,omitempty"`
	Metrics       []Status               `json:"metrics,omitempty"`
	Properties    map[string]interface{} `json:"properties,omitempty"`
}

// NewCheck returns an empty Check
func NewCheck() Check {
	return Check{
		OverallStatus: "UNKNOWN",
		Properties:    make(map[string]interface{}),
	}
}

// NewCheckfromJSON returns an Check populated from an Check JSON
func NewCheckfromJSON(hcjson []byte) (Check, error) {
	var hc = NewCheck()

	// Convert the body into a JSON map
	jmap, err := jsonToMap(hcjson)
	if err != nil {
		return hc, err
	}

	hc.Services = StatusSliceFromJmap(cast.ToSlice(jmap["services"]))
	hc.Systems = StatusSliceFromJmap(cast.ToSlice(jmap["systems"]))
	hc.Metrics = StatusSliceFromJmap(cast.ToSlice(jmap["metrics"]))
	if props, ok := jmap["properties"]; ok {
		hc.Properties = cast.ToStringMap(props)
	}

	hc.Calculate()

	return hc, nil
}

// Merge integrates one Check into this Check. Does not merge metrics. Does not dedupe
func (s *Check) Merge(hc *Check) {
	if len(hc.Services) > 0 {
		s.Services = append(s.Services, hc.Services...)
	}
	if len(hc.Systems) > 0 {
		s.Systems = append(s.Systems, hc.Systems...)
	}

	// Properties are not merged, because they are free-form

	/* Metrics are not merged, as they are not named uniquely across healthchecks
	if len(hc.Metrics) > 0 {
		s.Metrics = append(s.Metrics, hc.Metrics...)
	}
	*/

	// Update the overallstatus maybe
	s.Calculate()
}

// PrefixedMerge integrates one Check into this Check, prefixing all named items in the source hc before merging.
func (s *Check) PrefixedMerge(prefix string, hc *Check) {

	if len(hc.Services) > 0 {
		for i := range hc.Services {
			hc.Services[i].Name = SafeLabel(fmt.Sprintf("%s_%s", prefix, hc.Services[i].Name))
		}
		s.Services = append(s.Services, hc.Services...)
	}
	if len(hc.Systems) > 0 {
		for i := range hc.Systems {
			hc.Systems[i].Name = SafeLabel(fmt.Sprintf("%s_%s", prefix, hc.Systems[i].Name))
		}
		s.Systems = append(s.Systems, hc.Systems...)
	}
	if len(hc.Metrics) > 0 {
		for i := range hc.Metrics {
			hc.Metrics[i].Name = SafeLabel(fmt.Sprintf("%s_%s", prefix, hc.Metrics[i].Name))
		}
		s.Metrics = append(s.Metrics, hc.Metrics...)
	}
	// Properties are not merged, because they are free-form

	// Update the overallstatus maybe
	s.Calculate()
}

// AddService adds the provided Status to the Services array
func (s *Check) AddService(status *Status) {
	s.Services = append(s.Services, *status)
}

// AddSystem adds the provided Status to the Systems array
func (s *Check) AddSystem(status *Status) {
	s.Systems = append(s.Systems, *status)
}

// AddMetric adds the provided Status to the Metrics array
func (s *Check) AddMetric(status *Status) {
	s.Metrics = append(s.Metrics, *status)
}

// Calculate walks the tree and updates OverallStatus if applicable
func (s *Check) Calculate() {
	ostatus := OK

FLOOP:
	for _, service := range s.Services {
		switch service.Status {
		case OK:
		case UP:
		case WARNING:
			if ostatus == OK || ostatus == UNKNOWN {
				ostatus = WARNING
			}
		case BAD:
			fallthrough
		case ERROR:
			fallthrough
		case DOWN:
			fallthrough
		case CRITICAL:
			ostatus = CRITICAL
			break FLOOP
		case UNKNOWN:
			if ostatus != CRITICAL && ostatus != WARNING {
				ostatus = UNKNOWN
			}
		}
	}

	if ostatus != CRITICAL {
	NCFLOOP:
		for _, system := range s.Systems {
			switch system.Status {
			case OK:
			case UP:
			case WARNING:
				if ostatus == OK || ostatus == UNKNOWN {
					ostatus = WARNING
				}
			case BAD:
				fallthrough
			case ERROR:
				fallthrough
			case DOWN:
				fallthrough
			case CRITICAL:
				ostatus = CRITICAL
				break NCFLOOP
			case UNKNOWN:
				if ostatus != CRITICAL && ostatus != WARNING {
					ostatus = UNKNOWN
				}
			}
		}
	}

	if ostatus != CRITICAL {
	MFLOOP:
		for _, metric := range s.Metrics {
			//fmt.Printf("Metric %s = %v\n", metric.Name, metric.Value)
			if metric.Status != "" {
				//fmt.Printf("\thas status '%s'\n", metric.Status)
				switch metric.Status {
				case OK:
				case UP:
				case WARNING:
					if ostatus == OK || ostatus == UNKNOWN {
						ostatus = WARNING
					}
				case BAD:
					fallthrough
				case ERROR:
					fallthrough
				case DOWN:
					fallthrough
				case CRITICAL:
					ostatus = CRITICAL
					break MFLOOP
				}
			} else if v, ok := isNumericGimme(cast.ToString(metric.Value)); ok {
				// No status declared, but value is a number, so lets see what we got with the other numbers
				//fmt.Printf("\tis %+v\n", metric)
				if cv, ok := isNumericGimme(cast.ToString(metric.BadOver)); ok && v > cv {
					ostatus = CRITICAL
				} else if wv, ok := isNumericGimme(cast.ToString(metric.WarnOver)); ok && v > wv {
					ostatus = WARNING
				}
			}
		}
	}

	s.OverallStatus = ostatus
}

// JSON returns the JSON-encoded version of the Check
func (s *Check) JSON() string {
	j, err := json.Marshal(s)
	if err != nil {
		return "{}"
	}
	return string(j)
}

// Terse returns the JSON-encoded version of just the overall status of the Check
func (s *Check) Terse() string {
	ts := make(map[string]string)
	ts["overallStatus"] = s.OverallStatus

	j, err := json.Marshal(&ts)
	if err != nil {
		return "{}"
	}
	return string(j)
}

// Validate runs the Check.JSON() output against a JSON-Schema validator
func (s *Check) Validate() error {
	return ValidateJSON(s.JSON())
}

// ValidateJSON runs provided JSON against the JSON-Schema validator
func ValidateJSON(jsonBody string) error {
	schemaLoader := jschema.NewStringLoader(string(SchemaJSON))
	schema, err := jschema.NewSchema(schemaLoader)
	if err != nil {
		return fmt.Errorf("error creating reference schema: %w", err)
	}

	loader := jschema.NewStringLoader(jsonBody)
	result, err := schema.Validate(loader)
	if err != nil {
		return fmt.Errorf("cannot validate response: %v", err)
	}
	if !result.Valid() {
		var resp string
		for _, desc := range result.Errors() {
			resp += fmt.Sprintf(" %s ", desc)
		}
		return fmt.Errorf("JSON is not valid: %s", resp)
	}

	return nil
}

// Converts a JSON text block into a JSON type
func jsonToMap(j []byte) (JSON, error) {
	newMap := make(JSON)
	err := json.Unmarshal(j, &newMap)
	if err != nil {
		return nil, err
	}
	return newMap, nil
}
