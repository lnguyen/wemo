package wemo

import (
	"fmt"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSwitches(t *testing.T) {
	Convey("Given some integer with a starting value", t, func() {
		switches, err := Switches()
		So(err, ShouldBeNil)
		fmt.Println(switches)
		if len(switches) > 1 {
			switches[0].Off()
			switches[0].On()
		}
	})
}
