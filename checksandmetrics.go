package health

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	nagios "github.com/cognusion/go-nagios-checks"
	"github.com/spf13/cast"
)

var (
	safeReplacer = strings.NewReplacer(
		".", ":",
		" ", "_",
		"-", "_",
		"%", "perc",
	)
)

// SafeLabel returns a label that is safe to use with modern RRD
func SafeLabel(label string) string {
	return safeReplacer.Replace(label)
}

// Metrics takes a "metrics" document and appends Nagios-compatible metrics
// information to the message
func Metrics(n *nagios.Nagios, checkMap []interface{}, noisy bool) {
	for _, r := range checkMap {
		jr := lcKeys(cast.ToStringMap(r))
		// We check to see if there is a value, and that it is
		// numeric before adding it
		if _, ok := jr["value"]; ok {
			val := jr["value"]

			switch val := val.(type) {
			default:
				value := cast.ToString(val)

				// Fix null values
				if value == "" {
					value = "0"
				}

				// Casts
				var (
					warn string
					crit string
					min  string
					max  string
					name string
				)

				if _, ok := jr["warnover"]; ok {
					warn = cast.ToString(jr["warnover"])
				} else {
					warn = cast.ToString(jr["warnvalue"])
				}

				if _, ok := jr["badover"]; ok {
					crit = cast.ToString(jr["badover"])
				} else {
					crit = cast.ToString(jr["badvalue"])
				}

				min = cast.ToString(jr["minvalue"])
				max = cast.ToString(jr["maxvalue"])
				name = cast.ToString(jr["name"])

				n.AddMetricNumbers(name, value, warn, crit, min, max)

				if _, ok := jr["status"]; ok {
					status := cast.ToString(jr["status"])

					n.AddMessageIfBool(fmt.Sprintf(" %s %s=%s", status, name, value), noisy || status != "OK")
					switch status {
					case WARNING:
						n.EscalateIf(nagios.WARNING)
					case BAD, ERROR, CRITICAL, DOWN:
						n.EscalateIf(nagios.CRITICAL)
					}

				} else if v, ok := isNumericGimme(value); ok {
					// No status declared, but value is a number, so lets see what we got with the other numbers

					if cv, ok := isNumericGimme(crit); ok && v > cv {
						n.AddMessage(fmt.Sprintf(" %s %s=%s", CRITICAL, name, value))
						n.EscalateIf(nagios.CRITICAL)
					} else if wv, ok := isNumericGimme(warn); ok && v > wv {
						n.AddMessage(fmt.Sprintf(" %s %s=%s", WARNING, name, value))
						n.EscalateIf(nagios.WARNING)
					}
				}
			case map[string]interface{}:
				// TODO: recurse maps?
				// TODO: Handle geopoints?
			}

		}
	}
}

// Checks takes a status document, escalate the status, and appends
// Nagios-compatible information to the message
func Checks(n *nagios.Nagios, maxAge int64, checkMap []interface{}, noisy bool) {

	now := time.Now().UnixMilli()
	for _, r := range checkMap {
		jr := lcKeys(cast.ToStringMap(r))

		// Cache the name for easier use later
		checkName := cast.ToString(jr["name"])

		// For an item, if a timeout has been defined, use it. Else use maxage
		lMaxAge := maxAge * 1000
		if timeout, ok := jr["timeout"]; ok && cast.ToInt64(timeout) != 0 {
			lMaxAge = cast.ToInt64(timeout)
		}
		//fmt.Printf("%s %d %d\n", checkName, maxAge, lMaxAge)

		// Check the freshness
		if timestamp, ok := jr["timestamp"]; ok && cast.ToInt64(timestamp) != 0 {
			old := now - cast.ToInt64(timestamp)
			if old > lMaxAge {
				n.EscalateIf(nagios.WARNING)
				n.AddMessage(fmt.Sprintf(" %s: STALE (%d seconds old) ", checkName, old))
			}
		}

		// Check the status
		value := cast.ToString(jr["value"])
		status := cast.ToString(jr["status"])
		errorMessage := cast.ToString(jr["error"])
		message := cast.ToString(jr["message"])

		if message == "" && value != "" {
			message = value
		}

		// We're muting OK messages normally
		n.AddMessageIfBool(fmt.Sprintf(" %s: %s", checkName, status), noisy || !isOK(status))

		switch status {
		case WARNING:
			n.EscalateIf(nagios.WARNING)
			n.AddMessageIf(fmt.Sprintf(" %s", errorMessage), errorMessage)
			n.AddMessageIf(fmt.Sprintf(" (%s)", message), message)
		case BAD:
			fallthrough
		case ERROR:
			fallthrough
		case DOWN:
			fallthrough
		case CRITICAL:
			n.EscalateIf(nagios.CRITICAL)
			n.AddMessageIf(fmt.Sprintf(" %s", errorMessage), errorMessage)
			n.AddMessageIf(fmt.Sprintf(" (%s)", message), message)
		case UP:
			fallthrough
		case OK:
			n.EscalateIf(nagios.OK)
			n.AddMessageIfBool(fmt.Sprintf(" (%s)", message), noisy && message != "")
		case UNKNOWN:
			n.EscalateIf(nagios.UNKNOWN)
			n.AddMessage(" Unknown state! ")
			n.AddMessageIf(fmt.Sprintf(" %s", errorMessage), errorMessage)
			n.AddMessageIf(fmt.Sprintf(" (%s)", message), message)
		default:
			n.AddMessageIf(message, message)
		}
	}
}

// Test to see if a string is a valid number type, and return a float64 of it if so
func isNumericGimme(s string) (float64, bool) {
	v, err := strconv.ParseFloat(s, 64)
	return v, err == nil
}

func isOK(status string) bool {
	if status == OK || status == UP {
		return true
	}
	return false
}
