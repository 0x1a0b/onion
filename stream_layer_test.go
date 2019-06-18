package onion

import (
	"bytes"
	"io"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const validJSON = `
{
	"string" : "str",
	"number" : 100,
	"nested" : {
		"bool" : "true"
	}
}
`

func TestNewStreamLayer(t *testing.T) {
	Convey("Stream layer test", t, func() {
		buf := bytes.NewBufferString(validJSON)
		l, err := NewStreamLayer(buf, "json")
		So(err, ShouldBeNil)
		o := New(l)
		So(o.GetString("string"), ShouldEqual, "str")
		So(o.GetInt("number"), ShouldEqual, 100)
		So(o.GetBool("nested.bool"), ShouldBeTrue)
	})
}

type dummyDecoder struct {
	data map[string]interface{}
	err  error
}

func (d *dummyDecoder) Decode(io.Reader) (map[string]interface{}, error) {
	return d.data, d.err
}

func TestRegisterDecoder(t *testing.T) {
	Convey("Test dummy decoder", t, func() {
		RegisterDecoder(&dummyDecoder{
			data: map[string]interface{}{
				"hi": 10,
			},
		}, "dummy")
		l, err := NewStreamLayer(nil, "dummy")
		So(err, ShouldBeNil)
		o := New(l)
		So(o.GetInt("hi"), ShouldEqual, 10)
	})

	Convey("Fail decode", t, func() {
		_, err := NewStreamLayer(nil, "hi_i_am_not_a_format")
		So(err, ShouldBeError)

		buf := bytes.NewBufferString(`{INVALID}`)
		_, err = NewStreamLayer(buf, "json")
		So(err, ShouldBeError)
	})
}
