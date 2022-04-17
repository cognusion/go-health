package health

import (
	nagios "github.com/cognusion/go-nagios-checks"
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func Test_SafeLabel(t *testing.T) {

	Convey("When a string is Nagios-unsafe, it is rendered safe", t, func() {
		badString := "This is a very unsafe-string, 100% of the time."
		goodString := "This_is_a_very_unsafe_string,_100perc_of_the_time:"
		So(SafeLabel(badString), ShouldEqual, goodString)
	})
}

func Test_Metrics(t *testing.T) {
	metric := map[string]interface{}{
		"status":        "OK",
		"name":          "mem",
		"value":         845643,
		"expectedValue": nil,
		"minValue":      nil,
		"maxValue":      nil,
		"timeStamp":     1538154007,
	}
	metrics := make([]interface{}, 1)
	metrics[0] = metric
	var empty nagios.Nagios
	Convey("When a metrics document is added to an empty Nagios struct, it is properly updated", t, func() {
		var newNag nagios.Nagios
		Metrics(&newNag, metrics, false)
		So(newNag, ShouldNotResemble, empty)
		So(newNag.Status(), ShouldEqual, nagios.OK)
	})
}

func Test_MetricsCompute(t *testing.T) {
	metric := map[string]interface{}{
		"name":          "mem",
		"value":         845643,
		"expectedValue": 840000,
		"warnover":      840000,
		"maxValue":      850000,
		"timeStamp":     1538154007,
	}
	metrics := make([]interface{}, 1)
	metrics[0] = metric
	var empty nagios.Nagios
	Convey("When a metrics document is added to an empty Nagios struct, it is properly updated", t, func() {
		var newNag nagios.Nagios
		Metrics(&newNag, metrics, false)
		So(newNag, ShouldNotResemble, empty)
		So(newNag.Status(), ShouldEqual, nagios.WARNING)
	})
}

func Test_Checks(t *testing.T) {
	check := map[string]interface{}{
		"name":          "com.workers.health.AwsSqsHealthIndicator",
		"status":        "UP",
		"value":         nil,
		"expectedValue": nil,
		"minValue":      nil,
		"maxValue":      nil,
		"timeStamp":     1538154003,
		"timeout":       nil,
		"error":         nil,
		"message":       nil,
	}
	checks := make([]interface{}, 1)
	checks[0] = check
	var empty nagios.Nagios
	Convey("When a checks document is added to an empty Nagios struct, it is properly updated", t, func() {
		var newNag nagios.Nagios
		Checks(&newNag, 0, checks, false)
		So(newNag, ShouldNotResemble, empty)
		So(newNag.Status(), ShouldEqual, nagios.WARNING)
	})
}

func Test_ChecksEscalate(t *testing.T) {
	check := map[string]interface{}{
		"name":          "com.workers.health.AwsSqsHealthIndicator",
		"status":        "DOWN",
		"value":         nil,
		"expectedValue": nil,
		"minValue":      nil,
		"maxValue":      nil,
		"timeStamp":     1538154003,
		"timeout":       nil,
		"error":         nil,
		"message":       nil,
	}
	checks := make([]interface{}, 1)
	checks[0] = check
	var empty nagios.Nagios
	Convey("When a checks document is added to an empty Nagios struct, it is properly updated", t, func() {
		var newNag nagios.Nagios
		Checks(&newNag, 0, checks, false)
		So(newNag, ShouldNotResemble, empty)
		So(newNag.Status(), ShouldEqual, nagios.CRITICAL)
	})
}
