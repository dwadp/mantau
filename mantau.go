package mantau

import (
	"errors"
	"reflect"
	"time"
)

type (
	// Mantau instance
	mantau struct {
		opt *Options
	}

	// Mantau options
	Options struct {
		// Struct tag specify a custom tag to match with the schema when transforming a struct
		StructTag string `json:"struct_tag"`
	}
)

type (
	// A schema describing of how a data should be transformed
	Schema map[string]SchemaField

	// A SchemaField describe the matching key or tag from source data
	SchemaField struct {
		// The result mapped key
		Key string

		// Value could be nil if not using a nested schema
		Value interface{}
	}

	// The result after transforming data
	Result map[string]interface{}

	// Available data types
	DataType string
)

// Data type declarations
var (
	Struct  DataType = "struct"
	Map     DataType = "map"
	Slice   DataType = "slice"
	Array   DataType = "array"
	Pointer DataType = "pointer"
	Other   DataType = "other"
	Nil     DataType = "nil"
)

// Create a new mantau instance
// And set the default options
func New() *mantau {
	return &mantau{
		opt: &Options{
			StructTag: "json",
		},
	}
}

// Override the default mantau options
func (m *mantau) SetOpt(opt *Options) {
	m.opt = opt
}

// Transform a data with the given schema
func (m *mantau) Transform(src interface{}, schema Schema) (interface{}, error) {
	dataType := m.getDataType(src)

	if dataType == Other {
		return nil, errors.New("Uknown source type")
	}

	return m.execute(src, dataType, schema)
}

// Get the input data type for further processing
func (m *mantau) getDataType(src interface{}) DataType {
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

// Get the original value of pointer type
func (m *mantau) getValue(src interface{}) reflect.Value {
	val := reflect.ValueOf(src)

	if reflect.TypeOf(src).Kind() == reflect.Ptr {
		return val.Elem()
	}

	return val
}

// Get the original type of the source
func (m *mantau) getType(src interface{}) reflect.Type {
	val := reflect.TypeOf(src)

	if reflect.TypeOf(src).Kind() == reflect.Ptr {
		return val.Elem()
	}

	return val
}

// Get the original value of pointer
func (m *mantau) getPtrValue(src interface{}) interface{} {
	if reflect.ValueOf(src).IsZero() {
		return nil
	}

	return reflect.ValueOf(src).Elem().Interface()
}

// Check if source is the type of built-in struct type such as time.Time
func (m *mantau) isCustomStruct(src interface{}) bool {
	switch src.(type) {
	case time.Time:
		return true
	}

	return false
}

// Begin the transforming process
func (m *mantau) execute(src interface{}, dataType DataType, schema Schema) (interface{}, error) {
	if src == nil {
		return nil, nil
	}

	switch dataType {
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

// Transform individual value with a given schema
func (m *mantau) transformValue(s interface{}, t DataType, f Schema) (interface{}, error) {
	if t == Other {
		return s, nil
	}

	switch t {
	case Struct:
		if m.isCustomStruct(s) {
			return s, nil
		}

		return m.transformStruct(s, f)
	case Slice:
		return m.transformCollections(s, f)
	case Array:
		return m.transformCollections(s, f)
	case Map:
		return m.transformMap(s, f)
	case Pointer:
		value := m.getPtrValue(s)

		return m.transformValue(value, m.getDataType(value), f)
	case Other:
		return s, nil
	}

	return nil, nil
}

// Transforming a struct with the given schema
func (m *mantau) transformStruct(src interface{}, schema Schema) (Result, error) {
	if src == nil {
		return nil, nil
	}

	result := Result{}
	srcValue := m.getValue(src)
	srcType := m.getType(src)

	for i := 0; i < srcValue.NumField(); i++ {
		for k, v := range schema {
			tag, err := m.tagLookup(srcValue.Type(), srcType.Field(i).Name)

			if err != nil {
				return nil, err
			}

			if v.Key == tag {
				fieldValue := srcValue.Field(i).Interface()
				fieldValueType := m.getDataType(fieldValue)

				if fieldValueType == Nil {
					continue
				}

				schemaValue := schema

				if s, ok := v.Value.(Schema); ok {
					schemaValue = s
				}

				value, err := m.transformValue(fieldValue, fieldValueType, schemaValue)

				if err != nil {
					return nil, err
				}

				result[k] = value
			}
		}
	}

	return result, nil
}

// Transforming a collections of data could be slice or array
func (m *mantau) transformCollections(src interface{}, schema Schema) ([]interface{}, error) {
	if src == nil {
		return nil, nil
	}

	result := make([]interface{}, 0)
	srcValue := m.getValue(src)

	for i := 0; i < srcValue.Len(); i++ {
		fieldValue := srcValue.Index(i).Interface()
		fieldValueType := m.getDataType(fieldValue)

		if fieldValueType == Nil {
			continue
		}

		value, err := m.transformValue(fieldValue, fieldValueType, schema)

		if err != nil {
			return nil, err
		}

		if fieldValueType == Other {
			result = append(result, value)
			continue
		}

		res, ok := value.(Result)

		if !ok {
			continue
		}

		result = append(result, res)
	}

	return result, nil
}

// Transforming a map with the given schema
func (m *mantau) transformMap(src interface{}, schema Schema) (Result, error) {
	if src == nil {
		return nil, nil
	}

	result := Result{}
	srcValue := m.getValue(src)

	for _, mapValue := range srcValue.MapKeys() {
		fieldValue := srcValue.MapIndex(mapValue).Interface()
		fieldValueType := m.getDataType(fieldValue)

		if fieldValueType == Nil {
			continue
		}

		for k, v := range schema {
			if v.Key == mapValue.String() {
				schemaValue := schema

				if s, ok := v.Value.(Schema); ok {
					schemaValue = s
				}

				value, err := m.transformValue(fieldValue, fieldValueType, schemaValue)

				if err != nil {
					return nil, err
				}

				result[k] = value
			}
		}
	}

	return result, nil
}

// Find struct tag based on the field name
// This is used when matching a struct field with a schema
func (m *mantau) tagLookup(t reflect.Type, fieldName string) (string, error) {
	field, ok := t.FieldByName(fieldName)

	if !ok {
		return "", errors.New("Cannot find the field")
	}

	tag, ok := field.Tag.Lookup(m.opt.StructTag)

	if tag == "" {
		return "", nil
	}

	if !ok {
		return "", errors.New("Cannot find tag")
	}

	return tag, nil
}
