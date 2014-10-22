package typed

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func Test_Json(t *testing.T) {
	spec := gspec.New(t)
	typed, err := Json([]byte(`{"power": 9002}`))
	spec.Expect(typed.Int("power")).ToEqual(9002)
	spec.Expect(err).ToBeNil()
}

func Test_JsonFile(t *testing.T) {
	spec := gspec.New(t)
	typed, err := JsonFile("test.json")
	spec.Expect(typed.String("name")).ToEqual("leto")
	spec.Expect(err).ToBeNil()
}

func Test_Bool(t *testing.T) {
	spec := gspec.New(t)
	typed := New(build("log", true))
	spec.Expect(typed.Bool("log")).ToEqual(true)
	spec.Expect(typed.BoolOr("log", false)).ToEqual(true)
	spec.Expect(typed.Bool("other")).ToEqual(false)
	spec.Expect(typed.BoolOr("other", true)).ToEqual(true)
}

func Test_Int(t *testing.T) {
	spec := gspec.New(t)
	typed := New(build("port", 84))
	spec.Expect(typed.Int("port")).ToEqual(84)
	spec.Expect(typed.IntOr("port", 11)).ToEqual(84)
	spec.Expect(typed.Int("other")).ToEqual(0)
	spec.Expect(typed.IntOr("other", 33)).ToEqual(33)
}

func Test_Float(t *testing.T) {
	spec := gspec.New(t)
	typed := New(build("pi", 3.14))
	spec.Expect(typed.Float("pi")).ToEqual(3.14)
	spec.Expect(typed.FloatOr("pi", 11.3)).ToEqual(3.14)
	spec.Expect(typed.Float("other")).ToEqual(0.0)
	spec.Expect(typed.FloatOr("other", 11.3)).ToEqual(11.3)
}

func Test_String(t *testing.T) {
	spec := gspec.New(t)
	typed := New(build("host", "localhost"))
	spec.Expect(typed.String("host")).ToEqual("localhost")
	spec.Expect(typed.StringOr("host", "openmymind.net")).ToEqual("localhost")
	spec.Expect(typed.String("other")).ToEqual("")
	spec.Expect(typed.StringOr("other", "openmymind.net")).ToEqual("openmymind.net")

}

func Test_Object(t *testing.T) {
	spec := gspec.New(t)
	typed := New(build("server", build("port", 32)))
	spec.Expect(typed.Object("server").Int("port")).ToEqual(32)
	spec.Expect(typed.ObjectOr("server", build("a", "b")).Int("port")).ToEqual(32)
	spec.Expect(len(typed.Object("other"))).ToEqual(0)
	spec.Expect(typed.ObjectOr("other", build("x", "y")).String("x")).ToEqual("y")
}


func Test_Bools(t *testing.T) {
	spec := gspec.New(t)
	typed := New(build("boring", []interface{}{true, false}))
	spec.Expect(typed.Bools("boring")).ToEqual([]bool{true, false})
	spec.Expect(typed.Bools("other")).ToEqual([]bool{})
}

func Test_Ints(t *testing.T) {
	spec := gspec.New(t)
	typed := New(build("scores", []interface{}{2, 1}))
	spec.Expect(typed.Ints("scores")).ToEqual([]int{2, 1})
	spec.Expect(typed.Ints("other")).ToEqual([]int{})
}

func Test_Ints_WithFloats(t *testing.T) {
	spec := gspec.New(t)
	typed := New(build("scores", []interface{}{2.1, 7.39}))
	spec.Expect(typed.Ints("scores")).ToEqual([]int{2, 7})
}

func Test_Floats(t *testing.T) {
	spec := gspec.New(t)
	typed := New(build("ranks", []interface{}{2.1, 1.2}))
	spec.Expect(typed.Floats("ranks")).ToEqual([]float64{2.1, 1.2})
	spec.Expect(typed.Floats("other")).ToEqual([]float64{})
}

func Test_Strings(t *testing.T) {
	spec := gspec.New(t)
	typed := New(build("names", []interface{}{"a", "b"}))
	spec.Expect(typed.Strings("names")).ToEqual([]string{"a", "b"})
	spec.Expect(typed.Strings("other")).ToEqual([]string{})
	spec.Expect(typed.StringsOr("names", []string{"c", "d"})).ToEqual([]string{"a", "b"})
	spec.Expect(typed.StringsOr("other", []string{"c", "d"})).ToEqual([]string{"c", "d"})
}

func Test_Objects(t *testing.T) {
	spec := gspec.New(t)
	typed := New(build("names", []interface{}{build("first", 1), build("second", 2)}))
	spec.Expect(typed.Objects("names")[0].Int("first")).ToEqual(1)
}

func Test_Maps(t *testing.T) {
	spec := gspec.New(t)
	typed := New(build("names", []interface{}{build("first", 1), build("second", 2)}))
	spec.Expect(typed.Maps("names")[1]["second"]).ToEqual(2)
}

func Test_StringBool(t *testing.T) {
	spec := gspec.New(t)
	typed, _ := JsonString(`{"blocked":{"a":true,"b":false}}`)
	m := typed.StringBool("blocked")
	spec.Expect(m["a"]).ToEqual(true)
	spec.Expect(m["b"]).ToEqual(false)
}

func Test_StringInt(t *testing.T) {
	spec := gspec.New(t)
	typed, _ := JsonString(`{"count":{"a":123,"b":43}}`)
	m := typed.StringInt("count")
	spec.Expect(m["a"]).ToEqual(123)
	spec.Expect(m["b"]).ToEqual(43)
	spec.Expect(m["xxz"]).ToEqual(0)
}


func Test_StringFloat(t *testing.T) {
	spec := gspec.New(t)
	typed, _ := JsonString(`{"rank":{"aa":3.4,"bz":4.2}}`)
	m := typed.StringFloat("rank")
	spec.Expect(m["aa"]).ToEqual(3.4)
	spec.Expect(m["bz"]).ToEqual(4.2)
	spec.Expect(m["xx"]).ToEqual(0.0)
}

func Test_StringString(t *testing.T) {
	spec := gspec.New(t)
	typed, _ := JsonString(`{"atreides":{"leto":"ghanima","paul":"alia"}}`)
	m := typed.StringString("atreides")
	spec.Expect(m["leto"]).ToEqual("ghanima")
	spec.Expect(m["paul"]).ToEqual("alia")
	spec.Expect(m["vladimir"]).ToEqual("")
}

func Test_StringObject(t *testing.T) {
	spec := gspec.New(t)
	typed, _ := JsonString(`{"atreides":{"leto":{"sister": "ghanima"}, "goku": {"power": 9001}}}`)
	m := typed.StringObject("atreides")
	spec.Expect(m["leto"].String("sister")).ToEqual("ghanima")
	spec.Expect(m["goku"].Int("power")).ToEqual(9001)
}

func build(values ...interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(values))
	for i := 0; i < len(values); i += 2 {
		m[values[i].(string)] = values[i+1]
	}
	return m
}
