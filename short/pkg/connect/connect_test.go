package connect

import (
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {

	c.Convey("normal example", t, func() {
		url := "http://www.baidu.com"
		got := Get(url)
		c.ShouldBeTrue(got)
	})

	c.Convey("fatal example:", t, func() {
		url := "posts/rabbit/12345/"
		got := Get(url)
		c.ShouldBeFalse(got)
	})
}
