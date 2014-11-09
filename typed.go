package typed

import (
	"encoding/json"
	"io/ioutil"
)

type Typed map[string]interface{}

func New(m map[string]interface{}) Typed {
	return Typed(m)
}

func Json(data []byte) (Typed, error) {
	var m map[string]interface{}
	err := json.Unmarshal(data, &m)
	return Typed(m), err
}

func JsonString(data string) (Typed, error) {
	return Json([]byte(data))
}

func JsonFile(path string) (Typed, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Json(data)
}

func (t Typed) Bool(key string) bool {
	return t.BoolOr(key, false)
}

func (t Typed) BoolOr(key string, d bool) bool {
	value, exists := t[key]
	if exists == false {
		return d
	}
	if n, ok := value.(bool); ok {
		return n
	}
	return d
}

func (t Typed) Int(key string) int {
	return t.IntOr(key, 0)
}

func (t Typed) IntOr(key string, d int) int {
	value, exists := t[key]
	if exists == false {
		return d
	}
	if n, ok := value.(int); ok {
		return n
	}
	if f, ok := value.(float64); ok {
		return int(f)
	}
	return d
}

func (t Typed) Float(key string) float64 {
	return t.FloatOr(key, 0)
}

func (t Typed) FloatOr(key string, d float64) float64 {
	value, exists := t[key]
	if exists == false {
		return d
	}
	if n, ok := value.(float64); ok {
		return n
	}
	return d
}

func (t Typed) String(key string) string {
	return t.StringOr(key, "")
}

func (t Typed) StringOr(key string, d string) string {
	value, exists := t[key]
	if exists == false {
		return d
	}
	if n, ok := value.(string); ok {
		return n
	}
	return d
}

func (t Typed) Object(key string) Typed {
	o := t.ObjectOr(key, nil)
	if o == nil {
		return Typed(nil)
	}
	return o
}

func (t Typed) ObjectOr(key string, d map[string]interface{}) Typed {
	value, exists := t[key]
	if exists == false {
		return Typed(d)
	}
	if n, ok := value.(map[string]interface{}); ok {
		return Typed(n)
	}
	return Typed(d)
}

func (t Typed) Map(key string) map[string]interface{} {
	return t.MapOr(key, nil)
}

func (t Typed) MapOr(key string, d map[string]interface{}) map[string]interface{} {
	value, exists := t[key]
	if exists == false {
		return d
	}
	if n, ok := value.(map[string]interface{}); ok {
		return n
	}
	return d
}

func (t Typed) Bools(key string) []bool {
	return t.BoolsOr(key, nil)
}

func (t Typed) BoolsOr(key string, d []bool) []bool {
	value, exists := t[key]
	if exists == false {
		return d
	}
	if n, ok := value.([]bool); ok {
		return n
	}
	if a, ok := value.([]interface{}); ok {
		l := len(a)
		n := make([]bool, l)
		for i := 0; i < l; i++ {
			n[i] = a[i].(bool)
		}
		return n
	}
	return d
}

func (t Typed) Ints(key string) []int {
	return t.IntsOr(key, nil)
}

func (t Typed) IntsOr(key string, d []int) []int {
	value, exists := t[key]
	if exists == false {
		return d
	}
	if n, ok := value.([]int); ok {
		return n
	}
	if a, ok := value.([]interface{}); ok {
		l := len(a)
		if l == 0 {
			return d
		}

		n := make([]int, l)
		switch a[0].(type) {
		case float64:
			for i := 0; i < l; i++ {
				n[i] = int(a[i].(float64))
			}
		case int:
			for i := 0; i < l; i++ {
				n[i] = a[i].(int)
			}
		}
		return n
	}
	return d
}

func (t Typed) Floats(key string) []float64 {
	return t.FloatsOr(key, nil)
}

func (t Typed) FloatsOr(key string, d []float64) []float64 {
	value, exists := t[key]
	if exists == false {
		return d
	}
	if n, ok := value.([]float64); ok {
		return n
	}
	if a, ok := value.([]interface{}); ok {
		l := len(a)
		n := make([]float64, l)
		for i := 0; i < l; i++ {
			n[i] = a[i].(float64)
		}
		return n
	}
	return nil
}

func (t Typed) Strings(key string) []string {
	return t.StringsOr(key, nil)
}

func (t Typed) StringsOr(key string, d []string) []string {
	value, exists := t[key]
	if exists == false {
		return d
	}
	if n, ok := value.([]string); ok {
		return n
	}
	if a, ok := value.([]interface{}); ok {
		l := len(a)
		n := make([]string, l)
		for i := 0; i < l; i++ {
			n[i] = a[i].(string)
		}
		return n
	}
	return d
}

func (t Typed) Objects(key string) []Typed {
	value, exists := t[key]
	if exists == true {
		if a, ok := value.([]interface{}); ok {
			l := len(a)
			n := make([]Typed, l)
			for i := 0; i < l; i++ {
				n[i] = Typed(a[i].(map[string]interface{}))
			}
			return n
		}
	}
	return nil
}

func (t Typed) Maps(key string) []map[string]interface{} {
	value, exists := t[key]
	if exists == true {
		if a, ok := value.([]interface{}); ok {
			l := len(a)
			n := make([]map[string]interface{}, l)
			for i := 0; i < l; i++ {
				n[i] = a[i].(map[string]interface{})
			}
			return n
		}
	}
	return nil
}

func (t Typed) StringBool(key string) map[string]bool {
	raw, ok := t.getmap(key)
	if ok == false {
		return nil
	}
	m := make(map[string]bool, len(raw))
	for k, value := range raw {
		m[k] = value.(bool)
	}
	return m
}

func (t Typed) StringInt(key string) map[string]int {
	raw, ok := t.getmap(key)
	if ok == false {
		return nil
	}
	m := make(map[string]int, len(raw))
	for k, value := range raw {
		switch t := value.(type) {
		case int:
			m[k] = t
		case float64:
			m[k] = int(t)
		}
	}
	return m
}

func (t Typed) StringFloat(key string) map[string]float64 {
	raw, ok := t.getmap(key)
	if ok == false {
		return nil
	}
	m := make(map[string]float64, len(raw))
	for k, value := range raw {
		m[k] = value.(float64)
	}
	return m
}

func (t Typed) StringString(key string) map[string]string {
	raw, ok := t.getmap(key)
	if ok == false {
		return nil
	}
	m := make(map[string]string, len(raw))
	for k, value := range raw {
		m[k] = value.(string)
	}
	return m
}

func (t Typed) StringObject(key string) map[string]Typed {
	raw, ok := t.getmap(key)
	if ok == false {
		return nil
	}
	m := make(map[string]Typed, len(raw))
	for k, value := range raw {
		m[k] = Typed(value.(map[string]interface{}))
	}
	return m
}

func (t Typed) getmap(key string) (raw map[string]interface{}, exists bool) {
	value, exists := t[key]
	if exists == false {
		return
	}
	raw, exists = value.(map[string]interface{})
	return
}
