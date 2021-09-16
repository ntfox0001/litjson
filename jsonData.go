package litjson

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	TypeNone = iota
	TypeBool
	TypeString
	TypeDouble
	TypeList
	TypeMap
)

type JSONData struct {
	data      interface{}
	valueType int
}

func NewJSONData() *JSONData {
	return &JSONData{}
}
func NewJSONDataByType(jdType int) *JSONData {
	jd := &JSONData{
		valueType: jdType,
	}

	switch jdType {
	case TypeBool:
		jd.data = false
	case TypeString:
		jd.data = ""
	case TypeDouble:
		jd.data = 0.0
	case TypeList:
		jd.data = make([]*JSONData, 0)
	case TypeMap:
		jd.data = make(map[string]*JSONData)
	}

	return jd
}

func NewJSONDataFromJSON(js string) *JSONData {
	jd := &JSONData{}
	if err := jd.InitByJSON(js); err != nil {
		return nil
	}

	return jd
}
func NewJSONDataFromObject(obj interface{}) *JSONData {
	jd := &JSONData{}
	if err := jd.InitByObject(obj); err != nil {
		return nil
	}

	return jd
}

func (jd *JSONData) InitByObject(obj interface{}) error {
	if obj == nil {
		jd.valueType = TypeNone
		return nil
	}

	switch obj := obj.(type) {
	case bool:
		jd.data = obj
		jd.valueType = TypeBool
	case int32:
		jd.data = float64(obj)
		jd.valueType = TypeDouble
	case int64:
		jd.data = float64(obj)
		jd.valueType = TypeDouble
	case uint32:
		jd.data = float64(obj)
		jd.valueType = TypeDouble
	case uint64:
		jd.data = float64(obj)
		jd.valueType = TypeDouble
	case float32:
		jd.data = float64(obj)
		jd.valueType = TypeDouble
	case float64:
		jd.data = obj
		jd.valueType = TypeDouble
	case int:
		jd.data = float64(obj)
		jd.valueType = TypeDouble
	case uint:
		jd.data = float64(obj)
		jd.valueType = TypeDouble
	case string:
		jd.data = obj
		jd.valueType = TypeString
	case []interface{}:
		nlist := make([]*JSONData, 0)
		for _, v := range obj {
			nlist = append(nlist, NewJSONDataFromObject(v))
		}

		jd.data = nlist
		jd.valueType = TypeList
	case map[string]interface{}:
		nmap := make(map[string]*JSONData)
		for k, v := range obj {
			nmap[k] = NewJSONDataFromObject(v)
		}

		jd.data = nmap
		jd.valueType = TypeMap
	case *JSONData:
		jd.data = obj.data
		jd.valueType = obj.valueType
	default:
		bytejs, err := json.Marshal(obj)
		if err != nil {
			return err
		}

		return jd.InitByJSON(string(bytejs))
	}

	return nil
}
func (jd *JSONData) InitByJSON(js string) error {
	var obj interface{}
	if err := json.Unmarshal([]byte(js), &obj); err != nil {
		return err
	}

	return jd.InitByObject(obj)
}

// 确保是指定类型
func (jd *JSONData) ensure(valueType int) bool {
	// 如果还没有初始化，那么就用第一次调用的类型
	if jd.valueType == TypeNone {
		jd.valueType = valueType

		switch valueType {
		case TypeList:
			jd.data = make([]*JSONData, 0)
		case TypeMap:
			jd.data = make(map[string]*JSONData)
		}

		return true
	}

	return jd.valueType == valueType
}

// 确认key是否存在
func (jd *JSONData) Confirm(keys ...string) (string, bool) {
	if !jd.ensure(TypeMap) {
		return "not-map", false
	}

	rt := ""

	for _, k := range keys {
		if !jd.hasKey(k) {
			rt = fmt.Sprintf("%s, %s", rt, k)
		}
	}

	if rt != "" {
		return rt[2:], false
	}

	return "", true
}

func (jd *JSONData) Get(key string) *JSONData {
	if jd.ensure(TypeMap) {
		if v, ok := jd.data.(map[string]*JSONData)[key]; ok {
			return v
		}
	}

	return nil
}

func (jd *JSONData) SafeGet(key string) (*JSONData, error) {
	if jd.ensure(TypeMap) {
		if v, ok := jd.data.(map[string]*JSONData)[key]; ok {
			return v, nil
		}

		return nil, errors.New("no found item")
	}

	return nil, errors.New("type error")
}

func (jd *JSONData) SetKey(key string, value interface{}) {
	if jd.ensure(TypeMap) {
		jd.data.(map[string]*JSONData)[key] = jd.isJSONData(value)
	}
}

func (jd *JSONData) RemoveKey(key string) {
	if jd.ensure(TypeMap) {
		delete(jd.data.(map[string]*JSONData), key)
	}
}

func (jd *JSONData) HasKey(key string) bool {
	if jd.ensure(TypeMap) {
		return jd.hasKey(key)
	}

	return false
}

func (jd *JSONData) hasKey(key string) bool {
	_, ok := jd.data.(map[string]*JSONData)[key]
	return ok
}

func (jd *JSONData) GetType() int {
	return jd.valueType
}

func (jd *JSONData) Index(id int) *JSONData {
	if jd.ensure(TypeList) {
		if jd.Len() < id {
			return nil
		}

		return jd.data.([]*JSONData)[id]
	}

	return nil
}

func (jd *JSONData) SafeIndex(id int) (*JSONData, error) {
	if jd.ensure(TypeList) {
		if jd.Len() >= id {
			return nil, errors.New("id overflow")
		}

		return jd.data.([]*JSONData)[id], nil
	}

	return nil, errors.New("type error")
}

func (jd *JSONData) SetIndex(id int, value interface{}) {
	if jd.ensure(TypeList) {
		if jd.Len() <= id {
			return
		}

		jd.data.([]*JSONData)[id] = jd.isJSONData(value)
	}
}

func (jd *JSONData) Append(value interface{}) {
	if jd.ensure(TypeList) {
		jd.data = append(jd.data.([]*JSONData), jd.isJSONData(value))
	}
}
func (jd *JSONData) RemoveID(id int) {
	if jd.ensure(TypeList) {
		jd.data = sliceDel(jd.data.([]*JSONData), id)
	}
}

func sliceDel(s []*JSONData, id int) []*JSONData {
	if s == nil {
		return nil
	}

	if len(s) <= id {
		return s
	}

	t := append(s[:id], s[id+1:]...)

	return t
}

func (jd *JSONData) isJSONData(value interface{}) *JSONData {
	switch value := value.(type) {
	case *JSONData:
		return value
	default:
		return NewJSONDataFromObject(value)
	}
}

func (jd *JSONData) Len() int {
	switch jd.valueType {
	case TypeMap:
		return len(jd.data.(map[string]*JSONData))
	case TypeList:
		return len(jd.data.([]*JSONData))
	}

	return 0
}

func (jd *JSONData) GetString() string {
	if jd.ensure(TypeString) {
		return jd.data.(string)
	}

	return ""
}

func (jd *JSONData) GetFloat32() float32 {
	if jd.ensure(TypeDouble) {
		return float32(jd.data.(float64))
	}

	return 0
}

func (jd *JSONData) GetFloat64() float64 {
	if jd.ensure(TypeDouble) {
		return jd.data.(float64)
	}

	return 0
}

func (jd *JSONData) GetInt32() int32 {
	if jd.ensure(TypeDouble) {
		return int32(jd.data.(float64))
	}

	return 0
}

func (jd *JSONData) GetInt() int {
	if jd.ensure(TypeDouble) {
		return int(jd.data.(float64))
	}

	return 0
}

func (jd *JSONData) GetInt64() int64 {
	if jd.ensure(TypeDouble) {
		return int64(jd.data.(float64))
	}

	return 0
}

func (jd *JSONData) GetUInt32() uint32 {
	if jd.ensure(TypeDouble) {
		return uint32(jd.data.(float64))
	}

	return 0
}

func (jd *JSONData) GetUInt() uint {
	if jd.ensure(TypeDouble) {
		return uint(jd.data.(float64))
	}

	return 0
}

func (jd *JSONData) GetUInt64() uint64 {
	if jd.ensure(TypeDouble) {
		return uint64(jd.data.(float64))
	}

	return 0
}

func (jd *JSONData) GetBool() bool {
	if jd.ensure(TypeBool) {
		return jd.data.(bool)
	}

	return false
}

func (jd *JSONData) SetString(value string) {
	if jd.ensure(TypeString) {
		jd.data = value
	}
}

func (jd *JSONData) SetInt32(value int32) {
	if jd.ensure(TypeDouble) {
		jd.data = float64(value)
	}
}

func (jd *JSONData) SetInt(value int) {
	if jd.ensure(TypeDouble) {
		jd.data = float64(value)
	}
}

func (jd *JSONData) SetInt64(value int64) {
	if jd.ensure(TypeDouble) {
		jd.data = float64(value)
	}
}

func (jd *JSONData) SetUInt32(value uint32) {
	if jd.ensure(TypeDouble) {
		jd.data = float64(value)
	}
}

func (jd *JSONData) SetUInt(value uint) {
	if jd.ensure(TypeDouble) {
		jd.data = float64(value)
	}
}

func (jd *JSONData) SetUInt64(value uint64) {
	if jd.ensure(TypeDouble) {
		jd.data = float64(value)
	}
}

func (jd *JSONData) SetFloat32(value float32) {
	if jd.ensure(TypeDouble) {
		jd.data = float64(value)
	}
}

func (jd *JSONData) SetFloat64(value float64) {
	if jd.ensure(TypeDouble) {
		jd.data = value
	}
}

func (jd *JSONData) SetBool(value bool) {
	if jd.ensure(TypeBool) {
		jd.data = value
	}
}

func (jd *JSONData) Map() map[string]*JSONData {
	if jd.ensure(TypeMap) {
		return jd.data.(map[string]*JSONData)
	}

	return nil
}

func (jd *JSONData) List() []*JSONData {
	if jd.ensure(TypeList) {
		return jd.data.([]*JSONData)
	}

	return nil
}

func (jd *JSONData) GetPath(path []string) *JSONData {
	if len(path) == 0 {
		return nil
	}

	cursor := jd
	for i := range path {
		cursor = cursor.Get(path[i])

		if cursor == nil {
			return nil
		}
	}

	return cursor
}

func (jd *JSONData) ToObject() interface{} {
	switch jd.valueType {
	case TypeMap:
		nmap := make(map[string]interface{})
		for k, v := range jd.data.(map[string]*JSONData) {
			nmap[k] = v.ToObject()
		}

		return nmap
	case TypeList:
		nlist := make([]interface{}, 0)
		for _, v := range jd.data.([]*JSONData) {
			nlist = append(nlist, v.ToObject())
		}

		return nlist
	default:
		return jd.data
	}
}

// 实现JSON导出接口
func (jd *JSONData) MarshalJSON() ([]byte, error) {
	s := jd.ToJSON()
	return []byte(s), nil
}

func (jd *JSONData) UnmarshalJSON(js []byte) error {
	return jd.InitByJSON(string(js))
}

func (jd *JSONData) ToJSON() string {
	bytejs, err := json.Marshal(jd.ToObject())
	if err != nil {
		return ""
	}

	return string(bytejs)
}

func (jd *JSONData) Conv2Obj(objPtr interface{}) error {
	return json.Unmarshal([]byte(jd.ToJSON()), objPtr)
}

func UnmarshalString(js string, objPtr interface{}) error {
	return json.Unmarshal([]byte(js), objPtr)
}

func Unmarshal(bytes []byte, objPtr interface{}) error {
	return json.Unmarshal(bytes, objPtr)
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}
