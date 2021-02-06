package mantau

import (
	"errors"
	"reflect"
	"time"
)

type (
	// Mantau type
	mantau struct {
		opt *Options
	}

	// Mantau options
	Options struct {
		// Hook with determine how mantau take individual field and transform it
		// Based on the given schema
		Hook string
	}
)

type (
	// A schema describing how the data should be transformed
	Schema map[string]Field

	// A field describe the matching key or tag from source data
	Field struct {
		// The result mapped key
		Key string

		// Value could be nil or a schema
		Value interface{}
	}

	// A value will store the schema field name and corresponding value after it's being transformed
	Value struct {
		// Key will store the schema key
		Key string

		// Value will store the transformed value
		Value interface{}
	}

	// Result will store the final result of the data after it's being transformed
	Result map[string]interface{}

	// Kind is just a string type aliase for this package
	Kind string
)

// Data kinds
var (
	Struct  Kind = "struct"
	Map     Kind = "map"
	Slice   Kind = "slice"
	Array   Kind = "array"
	Pointer Kind = "pointer"
	Other   Kind = "other"
	Nil     Kind = "nil"
)

// IsEmpty will check if the Key or Value field is empty
// This will prevent an empty value result being added to the mapped result
func (v *Value) IsEmpty() bool {
	if v.Key == "" {
		return true
	}

	if v.Value == nil {
		return true
	}

	return false
}

// Create a new mantau instance and set the default options
func New() *mantau {
	return &mantau{
		opt: &Options{
			Hook: "json",
		},
	}
}

// SetOpt will override the default options with the given options
func (m *mantau) SetOpt(opt *Options) {
	m.opt = opt
}

// Transform data with the given schema
func (m *mantau) Transform(src interface{}, schema Schema) (interface{}, error) {
	return m.serialize(src, schema)
}

// Get the input data kind based on given value
func (m *mantau) getKind(src interface{}) Kind {
	if src == nil {
		return Nil
	}

	switch reflect.TypeOf(src).Kind() {
	case reflect.Struct:
		return Struct
	case reflect.Map:
		return Map
	case reflect.Slice:
		return Slice
	case reflect.Array:
		return Array
	case reflect.Ptr:
		return Pointer
	default:
		return Other
	}
}

// Check if the type of the given value other than a struct, map, array or slice
// If so, we should not transform it
func (m *mantau) shouldSkipTransform(src interface{}) bool {
	value := m.getValue(src).Interface()

	switch value.(type) {
	case time.Time:
		return true
	case string:
		return true
	case bool:
		return true
	case int, int8, int16, int32, int64:
		return true
	case uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64:
		return true
	case complex64, complex128:
		return true
	case []time.Time:
		return true
	case []string:
		return true
	case []bool:
		return true
	case []int, []int8, []int16, []int32, []int64:
		return true
	case []uint, []uint8, []uint16, []uint32, []uint64:
		return true
	case []float32, []float64:
		return true
	case []complex64, []complex128:
		return true
	}

	return false
}

// getValue will retrieve the value from the given source regardless the source is pointer or other
// and return reflect.Value
func (m *mantau) getValue(src interface{}) reflect.Value {
	val := reflect.ValueOf(src)

	if reflect.TypeOf(src).Kind() == reflect.Ptr {
		return val.Elem()
	}

	return val
}

// getType will retrieve the type of the given source regardless the source is pointer or other
// and return reflect.Type
func (m *mantau) getType(src interface{}) reflect.Type {
	val := reflect.TypeOf(src)

	if reflect.TypeOf(src).Kind() == reflect.Ptr {
		return val.Elem()
	}

	return val
}

// getPtrValue will retrieve the actual value from pointer and return an interface{}
func (m *mantau) getPtrValue(src interface{}) interface{} {
	if src == nil {
		return nil
	}

	if m.getKind(src) != Pointer {
		return nil
	}

	value := reflect.ValueOf(src).Elem()

	if value.Interface() == nil {
		return nil
	}

	if reflect.ValueOf(value.Interface()).IsZero() {
		return nil
	}

	// If the type of value is struct then return it directly because the next step is to check if the value is nil
	if m.getKind(value.Interface()) == Struct {
		return value.Interface()
	}

	return value.Interface()
}

// transformMap will take a map as an input and transform it's value based on the given schema
// and return mantau.Result as the final result
func (m *mantau) transformMap(src interface{}, schema Schema) (Result, error) {
	if src == nil {
		return nil, nil
	}

	result := Result{}
	value := m.getValue(src)

	for _, val := range value.MapKeys() {
		v, err := m.mapWithSchema(
			val.String(),
			value.MapIndex(val).Interface(),
			schema,
		)

		if err != nil {
			return nil, err
		}

		if v.IsEmpty() {
			continue
		}

		result[v.Key] = v.Value
	}

	return result, nil
}

// mapWithSchema will iterates the given schema and find the corresponding data based on the given value
// and return mantau.Value as the final result
func (m *mantau) mapWithSchema(field string, value interface{}, schema Schema) (Value, error) {
	for key, val := range schema {
		if val.Key == field {
			schemaValue := schema

			if s, ok := val.Value.(Schema); ok {
				schemaValue = s
			}

			v, err := m.transformValue(value, schemaValue)

			if err != nil {
				return Value{}, err
			}

			return Value{Key: key, Value: v}, nil
		}
	}

	return Value{}, nil
}

// tagLookup is used specifically for struct
// tagLookup will find the struct tag on a struct field
// the tag is used to map the struct value with the schema
func (m *mantau) tagLookup(t reflect.Type, fieldName string) (string, error) {
	field, ok := t.FieldByName(fieldName)

	if !ok {
		return "", errors.New("Cannot find the field")
	}

	tag, ok := field.Tag.Lookup(m.opt.Hook)

	if tag == "" || !ok {
		return "", errors.New("Cannot find tag")
	}

	return tag, nil
}

// serialize will check for the given value and determine which process need to take
// based on the given value and the given schema
func (m *mantau) serialize(src interface{}, schema Schema) (interface{}, error) {
	kind := m.getKind(src)

	if kind == Other {
		return nil, errors.New("Source type is not allowed")
	}

	if kind == Nil {
		return nil, nil
	}

	switch kind {
	case Struct:
		return m.transformStruct(src, schema)
	case Slice:
		return m.transformCollections(src, schema)
	case Array:
		return m.transformCollections(src, schema)
	case Map:
		return m.transformMap(src, schema)
	}

	return nil, nil
}

// transformValue will check for individual value after it's being transformed,
// if the given value contains nested data structure it will determine which process to take
// to get the final result
func (m *mantau) transformValue(src interface{}, schema Schema) (interface{}, error) {

	// Check if the value cannot be transformed. If so, then just return it
	if m.shouldSkipTransform(src) {
		return m.getValue(src).Interface(), nil
	}

	kind := m.getKind(src)

	switch kind {
	case Struct:
		return m.transformStruct(src, schema)
	case Slice:
		return m.transformCollections(src, schema)
	case Array:
		return m.transformCollections(src, schema)
	case Map:
		return m.transformMap(src, schema)
	case Pointer:
		value := m.getPtrValue(src)

		return m.transformValue(
			value,
			schema,
		)
	}

	return nil, nil
}

// transformCollections will take an array or slice as an input and transform
// it's value based on the given schema and return mantau.Result as the final result
func (m *mantau) transformCollections(src interface{}, schema Schema) ([]Result, error) {
	if src == nil {
		return nil, nil
	}

	result := make([]Result, 0)
	value := m.getValue(src)

	for i := 0; i < value.Len(); i++ {
		v, err := m.transformValue(value.Index(i).Interface(), schema)

		if err != nil {
			return nil, err
		}

		res, ok := v.(Result)

		if !ok {
			continue
		}

		result = append(result, res)
	}

	return result, nil
}

// transformStruct will take a struct as an input and transform it's value
// based on the given schema and return mantau.Result as the final result
func (m *mantau) transformStruct(src interface{}, schema Schema) (Result, error) {
	if src == nil {
		return nil, nil
	}

	result := Result{}
	value := m.getValue(src)
	dataType := m.getType(src)

	for i := 0; i < value.NumField(); i++ {
		tag, err := m.tagLookup(value.Type(), dataType.Field(i).Name)

		if err != nil {
			return nil, err
		}

		v, err := m.mapWithSchema(tag, value.Field(i).Interface(), schema)

		if err != nil {
			return nil, err
		}

		if v.IsEmpty() {
			continue
		}

		result[v.Key] = v.Value
	}

	return result, nil
}
