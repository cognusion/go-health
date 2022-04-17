package health

import (
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func Test_StatusRegistry(t *testing.T) {

	Convey("When a StatusRegistry is created and the functions jogged, everything is as-expected", t, func() {
		sr := NewStatusRegistry()
		sr.Add("Bob", "Ok", "Smith", "Smith")

		keys := sr.Keys()
		So(keys, ShouldResemble, []string{"Bob"})

		bob, err := sr.Get("Bob")
		So(err, ShouldBeNil)
		So(bob, ShouldNotResemble, Status{})

		sr.Remove("Bob")
		keys2 := sr.Keys()
		So(keys2, ShouldResemble, []string{})

		_, err2 := sr.Get("Bob")
		So(err2, ShouldNotBeNil)

	})
}
