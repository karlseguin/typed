package typed

import (
	"bytes"
	"encoding/json"
	"sort"
	"testing"
	"time"
)

func Test_Must(t *testing.T) {
	typed := Must([]byte(`{"power": 9001}`))
	equal(t, typed.Int("power"), 9001)

	defer mustTest(t, "unexpected EOF")
	Must([]byte(`{`))
	t.FailNow()
}

func Test_JsonArray(t *testing.T) {
	typed, err := JsonArray([]byte(`[{"id":1},{"id":2}]`))
	equal(t, err, nil)
	equal(t, typed[0].Int("id"), 1)
	equal(t, typed[1].Int("id"), 2)

	_, err = JsonArray([]byte(`{}`))
	equal(t, err.Error(), "json: cannot unmarshal object into Go value of type []interface {}")

	typed, err = JsonArray([]byte(`[]`))
	equal(t, err, nil)
	equal(t, len(typed), 0)
}

func Test_JsonFileArray(t *testing.T) {
	typed, err := JsonFileArray("test_array.json")
	equal(t, err, nil)
	equal(t, typed[0].String("name"), "goku")
	equal(t, typed[1].Int("power"), 9001)

	typed, err = JsonFileArray("invalid2.json")
	equal(t, err.Error(), "open invalid2.json: no such file or directory")
}

func Test_JsonReader(t *testing.T) {
	json := []byte(`{"power": 8988876781182962205}`)
	stream := bytes.NewBuffer(json)

	typed, err := JsonReader(stream)
	equal(t, typed.Int("power"), 8988876781182962205)
	equal(t, err, nil)
}

func Test_Json(t *testing.T) {
	typed, err := Json([]byte(`{"power": 8988876781182962205}`))
	equal(t, typed.Int("power"), 8988876781182962205)
	equal(t, err, nil)
}

func Test_JsonFile(t *testing.T) {
	typed, err := JsonFile("test.json")
	equal(t, typed.String("name"), "leto")
	equal(t, err, nil)

	typed, err = JsonFile("invalid.json")
	equal(t, err.Error(), "open invalid.json: no such file or directory")
}

func Test_Keys(t *testing.T) {
	typed := New(build("name", "leto", "type", []int{1, 2, 3}, "number", 1))
	keys := typed.Keys()
	sort.Strings(keys)
	equalList(t, keys, []string{"name", "number", "type"})
}

func Test_Bool(t *testing.T) {
	typed := New(build("log", true, "ace", false, "nope", 99))
	equal(t, typed.Bool("log"), true)
	equal(t, typed.BoolOr("log", false), true)
	equal(t, typed.Bool("other"), false)
	equal(t, typed.BoolOr("other", true), true)

	equal(t, typed.BoolMust("log"), true)
	equal(t, typed.BoolMust("ace"), false)

	value, exists := typed.BoolIf("nope")
	equal(t, value, false)
	equal(t, exists, false)

	defer mustTest(t, "expected boolean value for fail")
	typed.BoolMust("fail")
	t.FailNow()
}

func Test_Int(t *testing.T) {
	typed := New(build("port", 84, "string", "30", "i16", int16(1), "i32", int32(2), "i64", int64(3), "f64", float64(4), "number", json.Number("5"), "nope", true))
	equal(t, typed.Int("port"), 84)
	equal(t, typed.IntOr("port", 11), 84)
	value, exists := typed.IntIf("port")
	equal(t, value, 84)
	equal(t, exists, true)

	equal(t, typed.Int("string"), 30)
	equal(t, typed.IntOr("string", 11), 30)
	value, exists = typed.IntIf("string")
	equal(t, value, 30)
	equal(t, exists, true)

	equal(t, typed.Int("other"), 0)
	equal(t, typed.IntOr("other", 33), 33)
	value, exists = typed.IntIf("other")
	equal(t, value, 0)
	equal(t, exists, false)

	value, exists = typed.IntIf("nope")
	equal(t, value, 0)
	equal(t, exists, false)

	equal(t, typed.Int("i16"), 1)
	equal(t, typed.Int("i32"), 2)
	equal(t, typed.Int("i64"), 3)
	equal(t, typed.Int("f64"), 4)
	equal(t, typed.Int("number"), 5)

	equal(t, typed.IntMust("port"), 84)

	defer mustTest(t, "expected int value for fail")
	typed.IntMust("fail")
	t.FailNow()
}

func Test_Float(t *testing.T) {
	typed := New(build("pi", 3.14, "string", "30.14", "number", json.Number("32e-005"), "nope", true))
	equal(t, typed.Float("pi"), 3.14)
	equal(t, typed.FloatOr("pi", 11.3), 3.14)
	value, exists := typed.FloatIf("pi")
	equal(t, value, 3.14)
	equal(t, exists, true)

	equal(t, typed.Float("string"), 30.14)
	equal(t, typed.FloatOr("string", 11.3), 30.14)
	value, exists = typed.FloatIf("string")
	equal(t, value, 30.14)
	equal(t, exists, true)

	equal(t, typed.Float("other"), 0.0)
	equal(t, typed.FloatOr("other", 11.3), 11.3)
	value, exists = typed.FloatIf("other")
	equal(t, value, 0.0)
	equal(t, exists, false)

	value, exists = typed.FloatIf("nope")
	equal(t, value, 0.0)
	equal(t, exists, false)

	equal(t, typed.Float("number"), 0.00032)

	equal(t, typed.FloatMust("pi"), 3.14)

	defer mustTest(t, "expected float value for fail")
	typed.FloatMust("fail")
	t.FailNow()
}

func Test_String(t *testing.T) {
	typed := New(build("host", "localhost", "nope", 1))
	equal(t, typed.String("host"), "localhost")
	equal(t, typed.StringOr("host", "openmymind.net"), "localhost")
	value, exists := typed.StringIf("host")
	equal(t, value, "localhost")
	equal(t, exists, true)

	equal(t, typed.String("other"), "")
	equal(t, typed.StringOr("other", "openmymind.net"), "openmymind.net")
	value, exists = typed.StringIf("other")
	equal(t, value, "")
	equal(t, exists, false)

	value, exists = typed.StringIf("nope")
	equal(t, value, "")
	equal(t, exists, false)

	equal(t, typed.StringMust("host"), "localhost")

	defer mustTest(t, "expected string value for fail")
	typed.StringMust("fail")
	t.FailNow()
}

func Test_Object(t *testing.T) {
	typed := New(build("server", build("port", 32), "nope", "a"))
	equal(t, typed.Object("server").Int("port"), 32)
	equal(t, typed.ObjectOr("server", build("a", "b")).Int("port"), 32)

	equal(t, len(typed.Object("other")), 0)
	equal(t, typed.ObjectOr("other", build("x", "y")).String("x"), "y")

	value, exists := typed.ObjectIf("other")
	equal(t, len(value), 0)
	equal(t, exists, false)

	value, exists = typed.ObjectIf("nope")
	equal(t, len(value), 0)
	equal(t, exists, false)

	equal(t, typed.ObjectMust("server").Int("port"), 32)

	defer mustTest(t, "expected map for fail")
	typed.ObjectMust("fail")
	t.FailNow()
}

func Test_ObjectType(t *testing.T) {
	typed := New(build("server", Typed(build("port", 32))))
	equal(t, typed.Object("server").Int("port"), 32)
}

func Test_Interface(t *testing.T) {
	typed := New(build("host", "localhost"))
	equal(t, typed.Interface("host").(string), "localhost")
	equal(t, typed.InterfaceOr("host", "openmymind.net").(string), "localhost")
	value, exists := typed.InterfaceIf("host")
	equal(t, value.(string), "localhost")
	equal(t, exists, true)

	equal(t, typed.Interface("other"), nil)
	equal(t, typed.InterfaceOr("other", "openmymind.net").(string), "openmymind.net")
	value, exists = typed.InterfaceIf("other")
	equal(t, value, nil)
	equal(t, exists, false)

	equal(t, typed.InterfaceMust("host").(string), "localhost")

	defer mustTest(t, "expected map for fail")
	typed.InterfaceMust("fail")
	t.FailNow()
}

func Test_Bools(t *testing.T) {
	typed := New(build("boring", []interface{}{true, false}, "fail", []interface{}{true, "goku"}, "bools", []bool{false, false, true}, "nope", 1))
	equalList(t, typed.Bools("boring"), []bool{true, false})
	equal(t, len(typed.Bools("other")), 0)
	equalList(t, typed.BoolsOr("boring", []bool{false, true}), []bool{true, false})
	equalList(t, typed.BoolsOr("other", []bool{false, true}), []bool{false, true})
	equalList(t, typed.Bools("bools"), []bool{false, false, true})

	values, exists := typed.BoolsIf("fail")
	equalList(t, values, []bool{true, false})
	equal(t, exists, false)

	values, exists = typed.BoolsIf("nope")
	equal(t, len(values), 0)
	equal(t, exists, false)
}

func Test_Ints(t *testing.T) {
	typed := New(build("scores", []interface{}{2, 1, "3"}, "fail1", []interface{}{2, "nope"}, "fail2", []interface{}{2, true}, "ints", []int{9, 8}, "empty", []interface{}{}))
	equalList(t, typed.Ints("scores"), []int{2, 1, 3})
	equalList(t, typed.Ints("ints"), []int{9, 8})
	equal(t, len(typed.Ints("empty")), 0)
	equal(t, len(typed.Ints("other")), 0)
	equalList(t, typed.IntsOr("scores", []int{3, 4, 5}), []int{2, 1, 3})
	equalList(t, typed.IntsOr("other", []int{3, 4, 5}), []int{3, 4, 5})

	values, exists := typed.IntsIf("fail1")
	equalList(t, values, []int{2, 0})
	equal(t, exists, false)

	values, exists = typed.IntsIf("fail2")
	equalList(t, values, []int{2, 0})
	equal(t, exists, false)
}

func Test_Ints64(t *testing.T) {
	typed := New(build("scores", []interface{}{2, 1, "3", 2.0, int64(939292992929292), json.Number("882")}, "fail1", []interface{}{2, "nope"}, "fail2", []interface{}{2, true}, "fail3", "nope"))
	equalList(t, typed.Ints64("scores"), []int64{2, 1, 3, 2.0, 939292992929292, 882})
	equal(t, len(typed.Ints64("other")), 0)
	equalList(t, typed.Ints64Or("scores", []int64{3, 4, 5}), []int64{2, 1, 3, 2.0, 939292992929292, 882})
	equalList(t, typed.Ints64Or("other", []int64{3, 4, 5}), []int64{3, 4, 5})

	values, exists := typed.Ints64If("fail1")
	equalList(t, values, []int{2, 0})
	equal(t, exists, false)

	values, exists = typed.Ints64If("fail2")
	equalList(t, values, []int{2, 0})
	equal(t, exists, false)

	values, exists = typed.Ints64If("fail3")
	equal(t, len(values), 0)
	equal(t, exists, false)
}

func Test_Ints_WithFloats(t *testing.T) {
	typed := New(build("scores", []interface{}{2.1, 7.39}))
	equalList(t, typed.Ints("scores"), []int{2, 7})
}

func Test_Floats(t *testing.T) {
	typed := New(build("ranks", []interface{}{2.1, 1.2, "3.0", json.Number("44.45")}, "fail1", []interface{}{"a"}, "floats", []float64{9.0}))
	equalList(t, typed.Floats("floats"), []float64{9.0})
	equalList(t, typed.Floats("ranks"), []float64{2.1, 1.2, 3.0, 44.45})
	equal(t, len(typed.Floats("other")), 0)
	equalList(t, typed.FloatsOr("ranks", []float64{3.1, 4.2, 5.3}), []float64{2.1, 1.2, 3.0, 44.45})
	equalList(t, typed.FloatsOr("other", []float64{3.1, 4.2, 5.3}), []float64{3.1, 4.2, 5.3})

	values, exists := typed.FloatsIf("fail1")
	equalList(t, values, []float64{0})
	equal(t, exists, false)
}

func Test_Strings(t *testing.T) {
	typed := New(build("names", []interface{}{"a", "b"}, "strings", []string{"s1", "s2"}, "fail1", []interface{}{1, true}, "fail2", "2"))
	equalList(t, typed.Strings("names"), []string{"a", "b"})
	equalList(t, typed.Strings("strings"), []string{"s1", "s2"})
	equal(t, len(typed.Strings("other")), 0)
	equalList(t, typed.StringsOr("names", []string{"c", "d"}), []string{"a", "b"})
	equalList(t, typed.StringsOr("other", []string{"c", "d"}), []string{"c", "d"})

	values, exists := typed.StringsIf("fail1")
	equalList(t, values, []string{"", ""})
	equal(t, exists, false)

	values, exists = typed.StringsIf("fail2")
	equal(t, len(values), 0)
	equal(t, exists, false)
}

func Test_Objects(t *testing.T) {
	typed := New(build("names", []interface{}{build("first", 1), build("second", 2)}))
	equal(t, typed.Objects("names")[0].Int("first"), 1)
}

func Test_ObjectsIf(t *testing.T) {
	typed := New(build("names", []interface{}{build("first", 1), build("second", 2)}))
	objects, exists := typed.ObjectsIf("names")
	equal(t, objects[0].Int("first"), 1)
	equal(t, exists, true)

	objects, exists = typed.ObjectsIf("non_existing")
	equal(t, len(objects), 0)
	equal(t, exists, false)
}

func Test_ObjectsMust(t *testing.T) {
	typed := New(build("names", []interface{}{build("first", 1), build("second", 2)}))
	objects := typed.ObjectsMust("names")
	equal(t, objects[0].Int("first"), 1)

	paniced := false
	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				paniced = true
			}
		}()

		typed.ObjectsMust("non_existing")
	}()

	equal(t, paniced, true)
}

func Test_ObjectsAsMap(t *testing.T) {
	typed := New(build("names", []map[string]interface{}{build("first", 1), build("second", 2)}))
	equal(t, typed.Objects("names")[0].Int("first"), 1)
}

func Test_Maps(t *testing.T) {
	typed := New(build("names", []interface{}{build("first", 1), build("second", 2)}))
	equal(t, typed.Maps("names")[1]["second"], 2)
}

func Test_StringBool(t *testing.T) {
	typed, _ := JsonString(`{"blocked":{"a":true,"b":false}}`)
	m := typed.StringBool("blocked")
	equal(t, m["a"], true)
	equal(t, m["b"], false)
}

func Test_StringInt(t *testing.T) {
	typed, _ := JsonString(`{"count":{"a":123,"c":"55"}}`)
	m := typed.StringInt("count")
	equal(t, m["a"], 123)
	equal(t, m["c"], 55)
	equal(t, m["xxz"], 0)
	equal(t, len(typed.StringInt("nope")), 0)

	typed = New(build("count", map[string]interface{}{"a": 99, "b": 8.0, "c": "9"}, "fail1", map[string]interface{}{"a": "nope"}))
	m = typed.StringInt("count")
	equal(t, m["a"], 99)
	equal(t, m["b"], 8)
	equal(t, m["c"], 9)

	equal(t, len(typed.StringInt("fail")), 0)
}

func Test_StringFloat(t *testing.T) {
	typed, _ := JsonString(`{"rank":{"aa":3.4,"bz":4.2,"cc":"5.5"}}`)
	m := typed.StringFloat("rank")
	equal(t, m["aa"], 3.4)
	equal(t, m["bz"], 4.2)
	equal(t, m["cc"], 5.5)
	equal(t, m["xx"], 0.0)
	equal(t, len(typed.StringFloat("nope")), 0)

	typed = New(build("count", map[string]interface{}{"a": 1.1, "b": "2.2"}, "fail1", map[string]interface{}{"a": "nope"}))
	m = typed.StringFloat("count")
	equal(t, m["a"], 1.1)
	equal(t, m["b"], 2.2)

	equal(t, len(typed.StringFloat("fail")), 0)
}

func Test_StringString(t *testing.T) {
	typed, _ := JsonString(`{"atreides":{"leto":"ghanima","paul":"alia"}}`)
	m := typed.StringString("atreides")
	equal(t, m["leto"], "ghanima")
	equal(t, m["paul"], "alia")
	equal(t, m["vladimir"], "")
}

func Test_StringObject(t *testing.T) {
	typed, _ := JsonString(`{"atreides":{"leto":{"sister": "ghanima"}, "goku": {"power": 9001}}}`)
	m := typed.StringObject("atreides")
	equal(t, m["leto"].String("sister"), "ghanima")
	equal(t, m["goku"].Int("power"), 9001)
}

func Test_ToBytes(t *testing.T) {
	typed, _ := JsonString(`{"atreides":{"leto":{"sister": "ghanima"}}, "goku": {"power": 9001}}`)
	m, err := typed.ToBytes("goku")
	equal(t, err, nil)
	equal(t, string(m), `{"power":9001}`)
}

func Test_ToBytesNullHandling(t *testing.T) {
	typed, _ := JsonString(`{"atreides":null}`)
	m, err := typed.ToBytes("atreides")
	equal(t, err, nil)
	equal(t, len(m), 0)
}

func Test_MustBytes(t *testing.T) {
	typed, _ := JsonString(`{"atreides":{"leto":{"sister": "ghanima"}}, "goku": {"power": 9001}}`)
	m := typed.MustBytes("goku")
	equal(t, string(m), `{"power":9001}`)

	defer mustTest(t, "Key not found")
	typed.MustBytes("hi")
	t.FailNow()
}

func Test_ToBytesSelf(t *testing.T) {
	typed, _ := JsonString(`{"atreides":{"leto":{"sister": "ghanima"}}, "goku": {"power": 9001}}`)
	m, err := typed.Object("atreides").ToBytes("")
	equal(t, err, nil)
	equal(t, string(m), `{"leto":{"sister":"ghanima"}}`)
}

func Test_ToBytesNotFound(t *testing.T) {
	typed, _ := JsonString(`{"atreides":{"leto":{"sister": "ghanima"}}, "goku": {"power": 9001}}`)
	m, err := typed.ToBytes("other")
	equal(t, err, KeyNotFound)
	equal(t, string(m), "")
}

func Test_RootArray(t *testing.T) {
	typed, _ := JsonStringArray(`[{"id":1},{"id":2}]`)
	equal(t, typed[0].Int("id"), 1)
	equal(t, typed[1].Int("id"), 2)
}

func Test_RootArrayPrimitives(t *testing.T) {
	typed, _ := JsonStringArray(`[-41, {"id":1}, 2]`)
	equal(t, typed[0].Int("0"), -41)
	equal(t, typed[1].Int("id"), 1)
	equal(t, typed[2].Int("0"), 2)
}

func Test_Exists(t *testing.T) {
	typed := New(build("power", 9001))
	equal(t, typed.Exists("power"), true)
	equal(t, typed.Exists("spice"), false)
}

func Test_Map(t *testing.T) {
	typed := New(build("data", map[string]interface{}{"a": 1}, "nope", 11))
	m := typed.Map("data")
	equal(t, m["a"].(int), 1)

	m = typed.MapOr("data", nil)
	equal(t, m["a"].(int), 1)

	m = typed.MapOr("other", map[string]interface{}{"a": 2})
	equal(t, m["a"].(int), 2)

	m, exists := typed.MapIf("data")
	equal(t, m["a"].(int), 1)
	equal(t, exists, true)

	m, exists = typed.MapIf("other")
	equal(t, len(m), 0)
	equal(t, exists, false)

	m, exists = typed.MapIf("nope")
	equal(t, len(m), 0)
	equal(t, exists, false)
}

func Test_Time(t *testing.T) {
	now := time.Now().UTC()
	zero := time.Time{}
	typed := New(build("ts", now, "nope", true))
	equal(t, typed.Time("ts"), now)
	equal(t, typed.TimeOr("ts", zero), now)
	equal(t, typed.TimeOr("other", zero), zero)

	ts, exists := typed.TimeIf("ts")
	equal(t, ts, now)
	equal(t, exists, true)

	ts, exists = typed.TimeIf("other")
	equal(t, ts, zero)
	equal(t, exists, false)

	ts, exists = typed.TimeIf("nope")
	equal(t, ts, zero)
	equal(t, exists, false)

	equal(t, typed.TimeMust("ts"), now)

	defer mustTest(t, "expected time.Time value for other")
	typed.TimeMust("other")
	t.FailNow()
}

func build(values ...interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(values))
	for i := 0; i < len(values); i += 2 {
		m[values[i].(string)] = values[i+1]
	}
	return m
}

func mustTest(t *testing.T, expected string) {
	if err := recover(); err != nil {
		switch e := err.(type) {
		case string:
			equal(t, e, expected)
		case error:
			equal(t, e.Error(), expected)
		default:
			panic("unknown recover type")
		}
	}
}

func equal(t *testing.T, actual interface{}, expected interface{}) {
	t.Helper()
	if actual != expected {
		t.Errorf("expected '%v' to equal '%v", actual, expected)
		t.FailNow()
	}
}

// awful
func equalList(t *testing.T, actual interface{}, expected interface{}) {
	t.Helper()
	actualJSON, err := json.Marshal(actual)
	if err != nil {
		panic(err)
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		panic(err)
	}
	equal(t, string(actualJSON), string(expectedJSON))
}
