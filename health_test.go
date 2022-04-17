package health

import (
	. "github.com/smartystreets/goconvey/convey"
	jschema "github.com/xeipuuv/gojsonschema"

	"encoding/json"
	"io/ioutil"
	"testing"
)

var apiJSON = []byte(`{
	"overallStatus":"OK",
	"properties":{},
	"services":[],
	"systems":[
		{"status":"OK","name":"DBConnection","timeStamp":1538157610214,"timeout":null,"message":null,"error":null},
		{"status":"OK","name":"DBConnection","timeStamp":1538157610216,"timeout":null,"message":null,"error":null}
		],
		"metrics":[],
		"isPaused":false
}`)

var workerJSON = []byte(`{
  "overallStatus" : "OK",
  "properties" : {
    "queues" : [ "https://sqs.us-east-1.amazonaws.com/290913789/policy_complete", "https://sqs.us-east-1.amazonaws.com/290913789/response_qa" ],
    "topics" : [ "arn:aws:sns:us-east-1:290913789:topic_system-events", "arn:aws:sns:us-east-1:290913789:topic_healthcheck" ],
    "name" : "workers.system.dataset",
    "host" : "worker-ateam-0201"
  },
  "services" : [ {
    "name" : "com.workers.health.AwsSqsHealthIndicator",
    "status" : "UP",
    "value" : { },
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154003,
    "timeout" : null,
    "error" : null,
    "message" : null
  }, {
    "name" : "com.workers.health.AwsSnsHealthIndicator",
    "status" : "UP",
    "value" : { },
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007,
    "timeout" : null,
    "error" : null,
    "message" : null
  }, {
    "name" : "org.springframework.boot.actuate.health.DiskSpaceHealthIndicator",
    "status" : "UP",
    "value" : {
      "total" : 10726932480,
      "free" : 9218899968,
      "threshold" : 10485760
    },
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007,
    "timeout" : null,
    "error" : null,
    "message" : null
  }, {
    "name" : "org.springframework.boot.actuate.health.MongoHealthIndicator",
    "status" : "UP",
    "value" : {
      "version" : "3.2.17"
    },
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007,
    "timeout" : null,
    "error" : null,
    "message" : null
  } ],
  "systems" : [ ],
  "metrics" : [ {
    "status" : "OK",
    "name" : "mem",
    "value" : 845643,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "mem.free",
    "value" : 324294,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "processors",
    "value" : 2,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "instance.uptime",
    "value" : 10294270,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "uptime",
    "value" : 10357719,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "systemload.average",
    "value" : 0.01,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "heap.committed",
    "value" : 726528,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "heap.init",
    "value" : 129024,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "heap.used",
    "value" : 402233,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "heap",
    "value" : 1817088,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "nonheap.committed",
    "value" : 121240,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "nonheap.init",
    "value" : 2496,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "nonheap.used",
    "value" : 119115,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "nonheap",
    "value" : 0,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "threads.peak",
    "value" : 146,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "threads.daemon",
    "value" : 25,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "threads.totalStarted",
    "value" : 155,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "threads",
    "value" : 143,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "classes",
    "value" : 13082,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "classes.loaded",
    "value" : 13082,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "classes.unloaded",
    "value" : 0,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "gc.ps_scavenge.count",
    "value" : 45,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "gc.ps_scavenge.time",
    "value" : 1654,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "gc.ps_marksweep.count",
    "value" : 3,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "gc.ps_marksweep.time",
    "value" : 1698,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538154007
  }, {
    "status" : "OK",
    "name" : "counter.job.start.total",
    "value" : 1,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538144099
  }, {
    "status" : "OK",
    "name" : "counter.job.start.active",
    "value" : 0,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "timeStamp" : 1538144102123
  }, {
    "name" : "broken.test.value",
    "value" : 6,
    "expectedValue" : null,
    "minValue" : null,
    "maxValue" : null,
    "warnOver": 0,
    "badOver": 5,
    "timeStamp" : 1538144102
  } ],
  "isPaused" : false
}
`)

func Benchmark_SchemaSmallFull(b *testing.B) {

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		schemaLoader := jschema.NewReferenceLoader("file://./schema.json")
		schema, _ := jschema.NewSchema(schemaLoader)
		loader := jschema.NewStringLoader(string(apiJSON))
		result, err := schema.Validate(loader)
		if err != nil || !result.Valid() {
			b.Errorf("Bomb: %v %+v\n", err, result.Errors())
		}
	}
}

func Benchmark_SchemaSmallFast(b *testing.B) {
	schemaLoader := jschema.NewReferenceLoader("file://./schema.json")
	schema, _ := jschema.NewSchema(schemaLoader)
	w := string(apiJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loader := jschema.NewStringLoader(w)
		result, err := schema.Validate(loader)
		if err != nil || !result.Valid() {
			b.Errorf("Bomb: %v %+v\n", err, result.Errors())
		}
	}
}

func Test_Schemas(t *testing.T) {

	Convey("When the reference schema is loaded", t, func() {
		schemaLoader := jschema.NewReferenceLoader("file://./schema.json")
		schema, err := jschema.NewSchema(schemaLoader)
		So(err, ShouldBeNil) // Trap bad schema earlier

		Convey("and the 'refJSON' blob is checked, all is good", FailureContinues, func() {
			loader := jschema.NewReferenceLoader("file://./ref.json")

			result, err := schema.Validate(loader)
			So(err, ShouldBeNil)
			So(result.Valid(), ShouldBeTrue)
			if !result.Valid() {
				Println()
				for _, desc := range result.Errors() {
					Printf("- %s\n", desc)
				}
			}
		})

		Convey("and the 'workerJSON' blob is checked, all is good", FailureContinues, func() {
			loader := jschema.NewStringLoader(string(workerJSON))

			result, err := schema.Validate(loader)
			So(err, ShouldBeNil)
			So(result.Valid(), ShouldBeTrue)
			if !result.Valid() {
				Println()
				for _, desc := range result.Errors() {
					Printf("- %s\n", desc)
				}
			}
		})
	})
}

func Test_NewCheckfromJSON(t *testing.T) {

	Convey("When NewCheckfromJSON is called on known good JSON", t, func() {

		Convey("workerJSON, the result is a valid Check", func() {
			hc, err := NewCheckfromJSON(workerJSON)

			So(err, ShouldBeNil)
			So(hc, ShouldNotBeNil)
			So(hc, ShouldNotEqual, Check{})
		})

		Convey("refJSON, the result is a valid Check", func() {
			buf, err := ioutil.ReadFile("./ref.json")
			So(err, ShouldBeNil)

			hc, err := NewCheckfromJSON(buf)

			So(err, ShouldBeNil)
			So(hc, ShouldNotBeNil)
			So(hc, ShouldNotEqual, Check{})
		})

	})
}

func Test_StatusMetric(t *testing.T) {

	Convey("When an Check has a metric that is critical, the resulting status is correct", t, func() {

		hc, err := NewCheckfromJSON(workerJSON)

		So(err, ShouldBeNil)
		So(hc, ShouldNotBeNil)
		So(hc, ShouldNotEqual, Check{})

		Convey("The resulting status and MetricString are correct", func() {
			So(hc.OverallStatus, ShouldEqual, "CRITICAL")
			So(hc.Metrics[len(hc.Metrics)-1].MetricString(), ShouldEqual, "'broken.test.value'=6;0;5;;")
		})
	})
}

func Test_Merge(t *testing.T) {

	Convey("When NewCheckfromJSON is called on known good JSON, the result is a valid Check", t, func() {

		hc, err := NewCheckfromJSON(apiJSON)

		sysLen := len(hc.Systems)
		srvLen := len(hc.Services)
		metLen := len(hc.Metrics)

		So(err, ShouldBeNil)
		So(hc, ShouldNotBeNil)
		So(hc, ShouldNotEqual, Check{})

		Convey("When another NewCheckfromJSON is called on known good JSON, the result is a valid Check", func() {
			whc, err := NewCheckfromJSON(workerJSON)

			wsysLen := len(whc.Systems)
			wsrvLen := len(whc.Services)

			So(err, ShouldBeNil)
			So(whc, ShouldNotBeNil)
			So(whc, ShouldNotEqual, Check{})

			Convey("When the second Check is merged into the first Check, the result is a valid Check", func() {
				hc.Merge(&whc)

				So(err, ShouldBeNil)
				So(hc, ShouldNotBeNil)
				So(hc, ShouldNotEqual, Check{})

				So(len(hc.Systems), ShouldEqual, sysLen+wsysLen)
				So(len(hc.Services), ShouldEqual, srvLen+wsrvLen)
				So(len(hc.Metrics), ShouldEqual, metLen) // We don't merge metrics
			})
		})

		//Printf("HC: %+v\n", hc)

	})
}

func Test_PrefixedMerge(t *testing.T) {

	Convey("When NewCheckfromJSON is called on known good JSON, the result is a valid Check", t, func() {

		hc, err := NewCheckfromJSON(apiJSON)

		sysLen := len(hc.Systems)
		srvLen := len(hc.Services)
		metLen := len(hc.Metrics)

		So(err, ShouldBeNil)
		So(hc, ShouldNotBeNil)
		So(hc, ShouldNotEqual, Check{})

		Convey("When another NewCheckfromJSON is called on known good JSON, the result is a valid Check", func() {
			whc, err := NewCheckfromJSON(workerJSON)

			wsysLen := len(whc.Systems)
			wsrvLen := len(whc.Services)
			wmetLen := len(whc.Metrics)

			So(err, ShouldBeNil)
			So(whc, ShouldNotBeNil)
			So(whc, ShouldNotEqual, Check{})

			Convey("When the second Check is prefix-merged into the first Check, the result is a valid Check", func() {
				hc.PrefixedMerge("workers.system.dataset", &whc)

				So(err, ShouldBeNil)
				So(hc, ShouldNotBeNil)
				So(hc, ShouldNotEqual, Check{})

				So(len(hc.Systems), ShouldEqual, sysLen+wsysLen)
				So(len(hc.Services), ShouldEqual, srvLen+wsrvLen)
				So(len(hc.Metrics), ShouldEqual, metLen+wmetLen)
			})
		})

		//Printf("HC: %+v\n", hc)

	})
}

func Test_JSONOutput(t *testing.T) {

	Convey("When NewCheckfromJSON is called on known good JSON, the result is a valid Check", t, func() {

		hc, err := NewCheckfromJSON(apiJSON)

		sysLen := len(hc.Systems)
		srvLen := len(hc.Services)
		metLen := len(hc.Metrics)

		So(err, ShouldBeNil)
		So(hc, ShouldNotBeNil)
		So(hc, ShouldNotEqual, Check{})

		Convey("When another NewCheckfromJSON is called on known good JSON, the result is a valid Check", func() {
			whc, err := NewCheckfromJSON(workerJSON)

			wsysLen := len(whc.Systems)
			wsrvLen := len(whc.Services)
			wmetLen := len(whc.Metrics)

			So(err, ShouldBeNil)
			So(whc, ShouldNotBeNil)
			So(whc, ShouldNotEqual, Check{})

			Convey("When the second Check is prefix-merged into the first Check, the result is a valid Check", func() {
				hc.PrefixedMerge("workers.system.dataset", &whc)

				So(err, ShouldBeNil)
				So(hc, ShouldNotBeNil)
				So(hc, ShouldNotEqual, Check{})

				So(len(hc.Systems), ShouldEqual, sysLen+wsysLen)
				So(len(hc.Services), ShouldEqual, srvLen+wsrvLen)
				So(len(hc.Metrics), ShouldEqual, metLen+wmetLen)

				Convey("When the resulting Check is marhalled to JSON, the result is valid JSON", func() {
					var v interface{}
					err := json.Unmarshal([]byte(hc.JSON()), &v)
					So(err, ShouldBeNil)
				})
			})
		})

		//Printf("%+s\n", hc.JSON())

	})
}
