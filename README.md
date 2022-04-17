

# health
`import "github.com/cognusion/go-health"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package health provides standardized means for generating healthchecks, and consuming them from JSON.
The included schema files allow for validation of generation, regardless of the source.




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func Checks(n *nagios.Nagios, maxAge int64, checkMap []interface{}, noisy bool)](#Checks)
* [func Metrics(n *nagios.Nagios, checkMap []interface{}, noisy bool)](#Metrics)
* [func SafeLabel(label string) string](#SafeLabel)
* [func ValidateJSON(jsonBody string) error](#ValidateJSON)
* [type Check](#Check)
  * [func NewCheck() Check](#NewCheck)
  * [func NewCheckfromJSON(hcjson []byte) (Check, error)](#NewCheckfromJSON)
  * [func (s *Check) AddMetric(status *Status)](#Check.AddMetric)
  * [func (s *Check) AddService(status *Status)](#Check.AddService)
  * [func (s *Check) AddSystem(status *Status)](#Check.AddSystem)
  * [func (s *Check) Calculate()](#Check.Calculate)
  * [func (s *Check) JSON() string](#Check.JSON)
  * [func (s *Check) Merge(hc *Check)](#Check.Merge)
  * [func (s *Check) PrefixedMerge(prefix string, hc *Check)](#Check.PrefixedMerge)
  * [func (s *Check) Terse() string](#Check.Terse)
  * [func (s *Check) Validate() error](#Check.Validate)
* [type JSON](#JSON)
* [type Status](#Status)
  * [func StatusSliceFromJmap(jmap []interface{}) []Status](#StatusSliceFromJmap)
  * [func (s *Status) MarshalJSON() ([]byte, error)](#Status.MarshalJSON)
  * [func (s *Status) MetricString() string](#Status.MetricString)
* [type StatusRegistry](#StatusRegistry)
  * [func NewStatusRegistry() *StatusRegistry](#NewStatusRegistry)
  * [func (s *StatusRegistry) Add(name, status string, Value, ExpectedValue interface{})](#StatusRegistry.Add)
  * [func (s *StatusRegistry) Get(name string) (*Status, error)](#StatusRegistry.Get)
  * [func (s *StatusRegistry) Keys() []string](#StatusRegistry.Keys)
  * [func (s *StatusRegistry) Remove(name string)](#StatusRegistry.Remove)
* [type StatusString](#StatusString)


#### <a name="pkg-files">Package files</a>
[checksandmetrics.go](https://github.com/cognusion/go-health/tree/master/checksandmetrics.go) [health.go](https://github.com/cognusion/go-health/tree/master/health.go) [schema.go](https://github.com/cognusion/go-health/tree/master/schema.go) [status.go](https://github.com/cognusion/go-health/tree/master/status.go) [statusregistry.go](https://github.com/cognusion/go-health/tree/master/statusregistry.go)


## <a name="pkg-constants">Constants</a>
``` go
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
```
Status constants to prevent fat-fingered-oopsies


## <a name="pkg-variables">Variables</a>
``` go
var (
    // ErrNoSuchEntryError is returned when the requested element does not exist in the Registry
    ErrNoSuchEntryError = errors.New("no such element exists")
)
```
``` go
var SchemaJSON = []byte(`
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "$id": "http://example.com/product.schema.json",
    "title": "Go Healthcheck",
    "description": "The Go Healthcheck (health.Check) is output by HTTP services wishing to provide consumable healthcheck output",
    "type": "object",
    "properties": {
        "overallStatus": {
      "description": "The declared status of the system",
      "type": "string",
      "enum": [
        "OK",
        "WARNING",
        "ERROR",
                "BAD",
                "CRITICAL",
                "UNKNOWN"
      ]
    },
    "metrics": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/metric"
      }
    },
    "services": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/service"
      }
    },
    "systems": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/system"
      }
    }
  },
  "required": [
    "overallStatus"
  ],
  "definitions": {
    "health": {
      "type": "object",
      "required": [
        "name"
      ],
      "additionalProperties": true,
      "properties": {
        "name": {
          "type": "string",
          "description": "The name of the service"
        },
        "status": {
          "description": "The declared status of the service",
          "type": "string",
          "enum": [
            "OK",
                        "UP",
            "WARNING",
                        "ERROR",
                        "BAD",
                        "DOWN",
                        "CRITICAL",
                        "UNKNOWN"
          ]
        },
        "timeStamp": {
          "description": "The UNIX epoch timestamp for when this value was last fetched",
          "type": ["integer", "null"]
        },
        "timeout": {
          "description": "The number of milliseconds allowed to lapse between 'timeStamp' and 'now' before the metric is declared stale",
          "type": ["integer", "null"]
        },
        "message": {
            "type": ["string", "null"],
          "description": "A message explaining why the status is what it is (often exception message)"
        }
      }
    },
    "metric": {
      "type": "object",
      "required": [
        "name",
        "value"
      ],
      "extends" : {
        "$ref": "#/definitions/health"
      },
      "allOf": [{ "$ref": "#/definitions/health" }],
      "additionalProperties": true,
      "properties": {
        "value": {
          "description": "The current value for this metric",
          "type": "number"
        },
        "expectedValue": {
          "description": "The declared value that was expected for this metric",
          "type": ["number", "null"]
        },
        "minValue": {
          "description": "The declared minimum value that for this metric (graph floor)",
          "type": ["number", "null"]
        },
        "maxValue": {
          "description": "The declared maximum value that for this metric (graph ceiling)",
          "type": ["number", "null"]
        },
        "warnOver": {
          "description": "The value at which exceeding values generate WARNING status (graph yellow-line)",
          "type": ["number", "null"]
        },
        "badOver": {
          "description": "The value at which exceeding values generate CRITICAL status (graph red-line)",
          "type": ["number", "null"]
        }
      }
    },
    "service": {
      "type": "object",
      "required": [
        "name",
        "status"
      ],
      "extends" : {
        "$ref": "#/definitions/health"
      },
      "allOf": [{ "$ref": "#/definitions/health" }],
      "additionalProperties": true,
      "properties": {
      }
    },
    "system": {
      "type": "object",
      "required": [
        "name"
      ],
      "extends" : {
        "$ref": "#/definitions/health"
      },
      "allOf": [{ "$ref": "#/definitions/health" }],
      "additionalProperties": true,
      "properties": {
      }
    }
  }
}
`)
```
SchemaJSON was generated from schema.json at Sun Apr 17 12:23:48 PM EDT 2022



## <a name="Checks">func</a> [Checks](https://github.com/cognusion/go-health/tree/master/checksandmetrics.go?s=2520:2599#L106)
``` go
func Checks(n *nagios.Nagios, maxAge int64, checkMap []interface{}, noisy bool)
```
Checks takes a status document, escalate the status, and appends
Nagios-compatible information to the message



## <a name="Metrics">func</a> [Metrics](https://github.com/cognusion/go-health/tree/master/checksandmetrics.go?s=487:553#L29)
``` go
func Metrics(n *nagios.Nagios, checkMap []interface{}, noisy bool)
```
Metrics takes a "metrics" document and appends Nagios-compatible metrics
information to the message



## <a name="SafeLabel">func</a> [SafeLabel](https://github.com/cognusion/go-health/tree/master/checksandmetrics.go?s=304:339#L23)
``` go
func SafeLabel(label string) string
```
SafeLabel returns a label that is safe to use with modern RRD



## <a name="ValidateJSON">func</a> [ValidateJSON](https://github.com/cognusion/go-health/tree/master/health.go?s=6147:6187#L240)
``` go
func ValidateJSON(jsonBody string) error
```
ValidateJSON runs provided JSON against the JSON-Schema validator




## <a name="Check">type</a> [Check](https://github.com/cognusion/go-health/tree/master/health.go?s=467:898#L17)
``` go
type Check struct {
    // OverallStatus must be one of OK,WARNING,BAD/ERROR/CRITICAL, or UNKNOWN
    OverallStatus string                 `json:"overallStatus,omitempty"`
    Services      []Status               `json:"services,omitempty"`
    Systems       []Status               `json:"systems,omitempty"`
    Metrics       []Status               `json:"metrics,omitempty"`
    Properties    map[string]interface{} `json:"properties,omitempty"`
}

```
Check is a type used to define healthcheck statuses







### <a name="NewCheck">func</a> [NewCheck](https://github.com/cognusion/go-health/tree/master/health.go?s=935:956#L27)
``` go
func NewCheck() Check
```
NewCheck returns an empty Check


### <a name="NewCheckfromJSON">func</a> [NewCheckfromJSON](https://github.com/cognusion/go-health/tree/master/health.go?s=1121:1172#L35)
``` go
func NewCheckfromJSON(hcjson []byte) (Check, error)
```
NewCheckfromJSON returns an Check populated from an Check JSON





### <a name="Check.AddMetric">func</a> (\*Check) [AddMetric](https://github.com/cognusion/go-health/tree/master/health.go?s=3406:3447#L115)
``` go
func (s *Check) AddMetric(status *Status)
```
AddMetric adds the provided Status to the Metrics array




### <a name="Check.AddService">func</a> (\*Check) [AddService](https://github.com/cognusion/go-health/tree/master/health.go?s=3111:3153#L105)
``` go
func (s *Check) AddService(status *Status)
```
AddService adds the provided Status to the Services array




### <a name="Check.AddSystem">func</a> (\*Check) [AddSystem](https://github.com/cognusion/go-health/tree/master/health.go?s=3260:3301#L110)
``` go
func (s *Check) AddSystem(status *Status)
```
AddSystem adds the provided Status to the Systems array




### <a name="Check.Calculate">func</a> (\*Check) [Calculate](https://github.com/cognusion/go-health/tree/master/health.go?s=3561:3588#L120)
``` go
func (s *Check) Calculate()
```
Calculate walks the tree and updates OverallStatus if applicable




### <a name="Check.JSON">func</a> (\*Check) [JSON](https://github.com/cognusion/go-health/tree/master/health.go?s=5552:5581#L214)
``` go
func (s *Check) JSON() string
```
JSON returns the JSON-encoded version of the Check




### <a name="Check.Merge">func</a> (\*Check) [Merge](https://github.com/cognusion/go-health/tree/master/health.go?s=1712:1744#L57)
``` go
func (s *Check) Merge(hc *Check)
```
Merge integrates one Check into this Check. Does not merge metrics. Does not dedupe




### <a name="Check.PrefixedMerge">func</a> (\*Check) [PrefixedMerge](https://github.com/cognusion/go-health/tree/master/health.go?s=2290:2345#L78)
``` go
func (s *Check) PrefixedMerge(prefix string, hc *Check)
```
PrefixedMerge integrates one Check into this Check, prefixing all named items in the source hc before merging.




### <a name="Check.Terse">func</a> (\*Check) [Terse](https://github.com/cognusion/go-health/tree/master/health.go?s=5748:5778#L223)
``` go
func (s *Check) Terse() string
```
Terse returns the JSON-encoded version of just the overall status of the Check




### <a name="Check.Validate">func</a> (\*Check) [Validate](https://github.com/cognusion/go-health/tree/master/health.go?s=6009:6041#L235)
``` go
func (s *Check) Validate() error
```
Validate runs the Check.JSON() output against a JSON-Schema validator




## <a name="JSON">type</a> [JSON](https://github.com/cognusion/go-health/tree/master/health.go?s=378:410#L14)
``` go
type JSON map[string]interface{}
```
JSON is an encapsulating type for "jmap"-based structures










## <a name="Status">type</a> [Status](https://github.com/cognusion/go-health/tree/master/status.go?s=618:1939#L28)
``` go
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

```
Status is a type used to convey status-related information about a Service, System, or Metric







### <a name="StatusSliceFromJmap">func</a> [StatusSliceFromJmap](https://github.com/cognusion/go-health/tree/master/status.go?s=4498:4551#L126)
``` go
func StatusSliceFromJmap(jmap []interface{}) []Status
```
StatusSliceFromJmap is a hacky function that might take a slice of interfaces, and return a same-sized slice of Status





### <a name="Status.MarshalJSON">func</a> (\*Status) [MarshalJSON](https://github.com/cognusion/go-health/tree/master/status.go?s=3520:3566#L88)
``` go
func (s *Status) MarshalJSON() ([]byte, error)
```
MarshalJSON is a custom marshaller for JSON encoding,
to output TimeStamp and TimeOut as numbers instead of pretty strings.




### <a name="Status.MetricString">func</a> (\*Status) [MetricString](https://github.com/cognusion/go-health/tree/master/status.go?s=4068:4106#L113)
``` go
func (s *Status) MetricString() string
```
MetricString returns a Nagios Performance Data -compatible representation of Status




## <a name="StatusRegistry">type</a> [StatusRegistry](https://github.com/cognusion/go-health/tree/master/statusregistry.go?s=281:350#L14)
``` go
type StatusRegistry struct {
    sync.RWMutex
    // contains filtered or unexported fields
}

```
StatusRegistry is a gorosafe map of services to their Status objects







### <a name="NewStatusRegistry">func</a> [NewStatusRegistry](https://github.com/cognusion/go-health/tree/master/statusregistry.go?s=411:451#L20)
``` go
func NewStatusRegistry() *StatusRegistry
```
NewStatusRegistry returns an initialized StatusRegistry





### <a name="StatusRegistry.Add">func</a> (\*StatusRegistry) [Add](https://github.com/cognusion/go-health/tree/master/statusregistry.go?s=563:646#L27)
``` go
func (s *StatusRegistry) Add(name, status string, Value, ExpectedValue interface{})
```
Add or update an entry in StatusRegistry




### <a name="StatusRegistry.Get">func</a> (\*StatusRegistry) [Get](https://github.com/cognusion/go-health/tree/master/statusregistry.go?s=1279:1337#L61)
``` go
func (s *StatusRegistry) Get(name string) (*Status, error)
```
Get returns the requested Status, or ErrNoSuchEntryError




### <a name="StatusRegistry.Keys">func</a> (\*StatusRegistry) [Keys](https://github.com/cognusion/go-health/tree/master/statusregistry.go?s=1041:1081#L48)
``` go
func (s *StatusRegistry) Keys() []string
```
Keys returns a list of names from the StatusRegistry




### <a name="StatusRegistry.Remove">func</a> (\*StatusRegistry) [Remove](https://github.com/cognusion/go-health/tree/master/statusregistry.go?s=890:934#L41)
``` go
func (s *StatusRegistry) Remove(name string)
```
Remove an entry from the StatusRegistry




## <a name="StatusString">type</a> [StatusString](https://github.com/cognusion/go-health/tree/master/status.go?s=493:519#L25)
``` go
type StatusString = string
```
StatusString is a string type for static string consistency














- - -
Generated by [godoc2md](http://godoc.org/github.com/cognusion/godoc2md)
