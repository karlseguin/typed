package typed

import (
	"bytes"
	"testing"

	. "github.com/karlseguin/expect"
)

type TypedTests struct{}

func Test_Typed(t *testing.T) {
	Expectify(new(TypedTests), t)
}

func (_ TypedTests) JsonReader() {
	json := []byte(`{"power": 8988876781182962205}`)
	stream := bytes.NewBuffer(json)

	typed, err := JsonReader(stream)
	Expect(typed.Int("power")).To.Equal(8988876781182962205)
	Expect(err).To.Equal(nil)
}

func (_ TypedTests) Json() {
	typed, err := Json([]byte(`{"power": 8988876781182962205}`))
	Expect(typed.Int("power")).To.Equal(8988876781182962205)
	Expect(err).To.Equal(nil)
}

func (_ TypedTests) JsonFile() {
	typed, err := JsonFile("test.json")
	Expect(typed.String("name")).To.Equal("leto")
	Expect(err).To.Equal(nil)
}

func (_ TypedTests) Keys() {
	typed := New(build("name", "leto", "type", []int{1, 2, 3}, "number", 1))
	Expect(typed.Keys()).To.Equal([]string{"name", "type", "number"})
}

func (_ TypedTests) Bool() {
	typed := New(build("log", true))
	Expect(typed.Bool("log")).To.Equal(true)
	Expect(typed.BoolOr("log", false)).To.Equal(true)
	Expect(typed.Bool("other")).To.Equal(false)
	Expect(typed.BoolOr("other", true)).To.Equal(true)
}

func (_ TypedTests) Int() {
	typed := New(build("port", 84, "string", "30"))
	Expect(typed.Int("port")).To.Equal(84)
	Expect(typed.IntOr("port", 11)).To.Equal(84)
	Expect(typed.IntIf("port")).To.Equal(84, true)

	Expect(typed.Int("string")).To.Equal(30)
	Expect(typed.IntOr("string", 11)).To.Equal(30)
	Expect(typed.IntIf("string")).To.Equal(30, true)

	Expect(typed.Int("other")).To.Equal(0)
	Expect(typed.IntOr("other", 33)).To.Equal(33)
	Expect(typed.IntIf("other")).To.Equal(0, false)
}

func (_ TypedTests) Float() {
	typed := New(build("pi", 3.14, "string", "30.14"))
	Expect(typed.Float("pi")).To.Equal(3.14)
	Expect(typed.FloatOr("pi", 11.3)).To.Equal(3.14)
	Expect(typed.FloatIf("pi")).To.Equal(3.14, true)

	Expect(typed.Float("string")).To.Equal(30.14)
	Expect(typed.FloatOr("string", 11.3)).To.Equal(30.14)
	Expect(typed.FloatIf("string")).To.Equal(30.14, true)

	Expect(typed.Float("other")).To.Equal(0.0)
	Expect(typed.FloatOr("other", 11.3)).To.Equal(11.3)
	Expect(typed.FloatIf("other")).To.Equal(0.0, false)

}

func (_ TypedTests) String() {
	typed := New(build("host", "localhost"))
	Expect(typed.String("host")).To.Equal("localhost")
	Expect(typed.StringOr("host", "openmymind.net")).To.Equal("localhost")
	Expect(typed.StringIf("host")).To.Equal("localhost", true)

	Expect(typed.String("other")).To.Equal("")
	Expect(typed.StringOr("other", "openmymind.net")).To.Equal("openmymind.net")
	Expect(typed.StringIf("other")).To.Equal("", false)
}

func (_ TypedTests) Object() {
	typed := New(build("server", build("port", 32)))
	Expect(typed.Object("server").Int("port")).To.Equal(32)
	Expect(typed.ObjectOr("server", build("a", "b")).Int("port")).To.Equal(32)

	Expect(len(typed.Object("other"))).To.Equal(0)
	Expect(typed.ObjectOr("other", build("x", "y")).String("x")).To.Equal("y")
	Expect(typed.ObjectIf("other")).To.Equal(nil, false)
}

func (_ TypedTests) ObjectType() {
	typed := New(build("server", Typed(build("port", 32))))
	Expect(typed.Object("server").Int("port")).To.Equal(32)
}

func (_ TypedTests) Interface() {
	typed := New(build("host", "localhost"))
	Expect(typed.Interface("host").(string)).To.Equal("localhost")
	Expect(typed.InterfaceOr("host", "openmymind.net").(string)).To.Equal("localhost")
	Expect(typed.InterfaceIf("host")).To.Equal(interface{}("localhost"), true)

	Expect(typed.Interface("other")).To.Equal(nil)
	Expect(typed.InterfaceOr("other", "openmymind.net").(string)).To.Equal("openmymind.net")
	Expect(typed.InterfaceIf("other")).To.Equal(nil, false)
}

func (_ TypedTests) Bools() {
	typed := New(build("boring", []interface{}{true, false}))
	Expect(typed.Bools("boring")).To.Equal([]bool{true, false})
	Expect(len(typed.Bools("other"))).To.Equal(0)
	Expect(typed.BoolsOr("boring", []bool{false, true})).To.Equal([]bool{true, false})
	Expect(typed.BoolsOr("other", []bool{false, true})).To.Equal([]bool{false, true})
}

func (_ TypedTests) Ints() {
	typed := New(build("scores", []interface{}{2, 1, "3"}))
	Expect(typed.Ints("scores")).To.Equal([]int{2, 1, 3})
	Expect(len(typed.Ints("other"))).To.Equal(0)
	Expect(typed.IntsOr("scores", []int{3, 4, 5})).To.Equal([]int{2, 1, 3})
	Expect(typed.IntsOr("other", []int{3, 4, 5})).To.Equal([]int{3, 4, 5})
}

func (_ TypedTests) Ints64() {
	typed := New(build("scores", []interface{}{2, 1, "3"}))
	Expect(typed.Ints64("scores")).To.Equal([]int64{2, 1, 3})
	Expect(len(typed.Ints64("other"))).To.Equal(0)
	Expect(typed.Ints64Or("scores", []int64{3, 4, 5})).To.Equal([]int64{2, 1, 3})
	Expect(typed.Ints64Or("other", []int64{3, 4, 5})).To.Equal([]int64{3, 4, 5})
}

func (_ TypedTests) Ints_WithFloats() {
	typed := New(build("scores", []interface{}{2.1, 7.39}))
	Expect(typed.Ints("scores")).To.Equal([]int{2, 7})
}

func (_ TypedTests) Floats() {
	typed := New(build("ranks", []interface{}{2.1, 1.2, "3.0"}))
	Expect(typed.Floats("ranks")).To.Equal([]float64{2.1, 1.2, 3.0})
	Expect(len(typed.Floats("other"))).To.Equal(0)
	Expect(typed.FloatsOr("ranks", []float64{3.1, 4.2, 5.3})).To.Equal([]float64{2.1, 1.2, 3.0})
	Expect(typed.FloatsOr("other", []float64{3.1, 4.2, 5.3})).To.Equal([]float64{3.1, 4.2, 5.3})
}

func (_ TypedTests) Strings() {
	typed := New(build("names", []interface{}{"a", "b"}))
	Expect(typed.Strings("names")).To.Equal([]string{"a", "b"})
	Expect(len(typed.Strings("other"))).To.Equal(0)
	Expect(typed.StringsOr("names", []string{"c", "d"})).To.Equal([]string{"a", "b"})
	Expect(typed.StringsOr("other", []string{"c", "d"})).To.Equal([]string{"c", "d"})
}

func (_ TypedTests) Objects() {
	typed := New(build("names", []interface{}{build("first", 1), build("second", 2)}))
	Expect(typed.Objects("names")[0].Int("first")).To.Equal(1)
}

func (_ TypedTests) ObjectsIf() {
	typed := New(build("names", []interface{}{build("first", 1), build("second", 2)}))
	objects, exists := typed.ObjectsIf("names")
	Expect(objects[0].Int("first")).To.Equal(1)
	Expect(exists).To.Equal(true)

	objects, exists = typed.ObjectsIf("non_existing")
	Expect(objects).To.Equal(nil)
	Expect(exists).To.Equal(false)
}

func (_ TypedTests) ObjectsMust() {
	typed := New(build("names", []interface{}{build("first", 1), build("second", 2)}))
	objects := typed.ObjectsMust("names")
	Expect(objects[0].Int("first")).To.Equal(1)

	paniced := false
	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				paniced = true
			}
		}()

		typed.ObjectsMust("non_existing")
	}()

	Expect(paniced).To.Equal(true)
}

func (_ TypedTests) ObjectsAsMap() {
	typed := New(build("names", []map[string]interface{}{build("first", 1), build("second", 2)}))
	Expect(typed.Objects("names")[0].Int("first")).To.Equal(1)
}

func (_ TypedTests) Maps() {
	typed := New(build("names", []interface{}{build("first", 1), build("second", 2)}))
	Expect(typed.Maps("names")[1]["second"]).To.Equal(2)
}

func (_ TypedTests) StringBool() {
	typed, _ := JsonString(`{"blocked":{"a":true,"b":false}}`)
	m := typed.StringBool("blocked")
	Expect(m["a"]).To.Equal(true)
	Expect(m["b"]).To.Equal(false)
}

func (_ TypedTests) StringInt() {
	typed, _ := JsonString(`{"count":{"a":123,"b":43,"c":"55"}}`)
	m := typed.StringInt("count")
	Expect(m["a"]).To.Equal(123)
	Expect(m["b"]).To.Equal(43)
	Expect(m["c"]).To.Equal(55)
	Expect(m["xxz"]).To.Equal(0)
}

func (_ TypedTests) StringFloat() {
	typed, _ := JsonString(`{"rank":{"aa":3.4,"bz":4.2,"cc":"5.5"}}`)
	m := typed.StringFloat("rank")
	Expect(m["aa"]).To.Equal(3.4)
	Expect(m["bz"]).To.Equal(4.2)
	Expect(m["cc"]).To.Equal(5.5)
	Expect(m["xx"]).To.Equal(0.0)
}

func (_ TypedTests) StringString() {
	typed, _ := JsonString(`{"atreides":{"leto":"ghanima","paul":"alia"}}`)
	m := typed.StringString("atreides")
	Expect(m["leto"]).To.Equal("ghanima")
	Expect(m["paul"]).To.Equal("alia")
	Expect(m["vladimir"]).To.Equal("")
}

func (_ TypedTests) StringObject() {
	typed, _ := JsonString(`{"atreides":{"leto":{"sister": "ghanima"}, "goku": {"power": 9001}}}`)
	m := typed.StringObject("atreides")
	Expect(m["leto"].String("sister")).To.Equal("ghanima")
	Expect(m["goku"].Int("power")).To.Equal(9001)
}

func (_ TypedTests) ToBytes() {
	typed, _ := JsonString(`{"atreides":{"leto":{"sister": "ghanima"}}, "goku": {"power": 9001}}`)
	m, err := typed.ToBytes("goku")
	Expect(err).To.Equal(nil)
	Expect(string(m)).To.Equal(`{"power":9001}`)
}

func (_ TypedTests) ToBytesNullHandling() {
	typed, _ := JsonString(`{"atreides":null}`)
	m, err := typed.ToBytes("atreides")
	Expect(err).To.Equal(nil)
	Expect(m).To.Equal(nil)
}

func (_ TypedTests) MustBytes() {
	typed, _ := JsonString(`{"atreides":{"leto":{"sister": "ghanima"}}, "goku": {"power": 9001}}`)
	m := typed.MustBytes("goku")
	Expect(string(m)).To.Equal(`{"power":9001}`)
}

func (_ TypedTests) ToBytesSelf() {
	typed, _ := JsonString(`{"atreides":{"leto":{"sister": "ghanima"}}, "goku": {"power": 9001}}`)
	m, err := typed.Object("atreides").ToBytes("")
	Expect(err).To.Equal(nil)
	Expect(string(m)).To.Equal(`{"leto":{"sister":"ghanima"}}`)
}

func (_ TypedTests) ToBytesNotFound() {
	typed, _ := JsonString(`{"atreides":{"leto":{"sister": "ghanima"}}, "goku": {"power": 9001}}`)
	m, err := typed.ToBytes("other")
	Expect(err).To.Equal(KeyNotFound)
	Expect(string(m)).To.Equal("")
}

func (_ TypedTests) RootArray() {
	typed, _ := JsonStringArray(`[{"id":1},{"id":2}]`)
	Expect(typed[0].Int("id")).To.Equal(1)
	Expect(typed[1].Int("id")).To.Equal(2)
}

func (_ TypedTests) RootArrayPrimitives() {
	typed, _ := JsonStringArray(`[-41, {"id":1}, 2]`)
	Expect(typed[0].Int("0")).To.Equal(-41)
	Expect(typed[1].Int("id")).To.Equal(1)
	Expect(typed[2].Int("0")).To.Equal(2)
}

func (_ TypedTests) Exists() {
	typed := New(build("power", 9001))
	Expect(typed.Exists("power")).To.Equal(true)
	Expect(typed.Exists("spice")).To.Equal(false)
}

func build(values ...interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(values))
	for i := 0; i < len(values); i += 2 {
		m[values[i].(string)] = values[i+1]
	}
	return m
}
