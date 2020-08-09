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

	// Data kind mapping
	DataKind string
)

// Data kind declarations
var (
	Struct  DataKind = "struct"
	Map     DataKind = "map"
	Slice   DataKind = "slice"
	Array   DataKind = "array"
	Pointer DataKind = "pointer"
	Other   DataKind = "other"
	Nil     DataKind = "nil"
)

// Create a new mantau instance and set the default options
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
	dataKind := m.getDataKind(src)

	if dataKind == Other {
		return nil, errors.New("Uknown source type")
	}

	return m.begin(src, dataKind, schema)
}

// Get the input data kind for further processing
func (m *mantau) getDataKind(src interface{}) DataKind {
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

// Determine if the value doesn't need for further transforming
func (m *mantau) shouldSkipTransform(src interface{}) bool {
	switch src.(type) {
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

// Begin the transforming process
func (m *mantau) begin(src interface{}, dataKind DataKind, schema Schema) (interface{}, error) {
	if src == nil {
		return nil, nil
	}

	switch dataKind {
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
func (m *mantau) transformValue(src interface{}, dataKind DataKind, schema Schema) (interface{}, error) {
	if dataKind == Other {
		return src, nil
	}

	if m.shouldSkipTransform(src) {
		return src, nil
	}

	switch dataKind {
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
		kind := m.getDataKind(value)

		return m.transformValue(value, kind, schema)
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
		fieldValue := srcValue.Field(i).Interface()
		fieldValueKind := m.getDataKind(fieldValue)

		if fieldValueKind == Nil {
			continue
		}

		tag, err := m.tagLookup(srcValue.Type(), srcType.Field(i).Name)

		if err != nil {
			return nil, err
		}

		for k, v := range schema {
			if v.Key == tag {
				schemaValue := schema

				if s, ok := v.Value.(Schema); ok {
					schemaValue = s
				}

				value, err := m.transformValue(fieldValue, fieldValueKind, schemaValue)

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
func (m *mantau) transformCollections(src interface{}, schema Schema) ([]Result, error) {
	if src == nil {
		return nil, nil
	}

	result := make([]Result, 0)
	srcValue := m.getValue(src)

	for i := 0; i < srcValue.Len(); i++ {
		fieldValue := srcValue.Index(i).Interface()
		fieldValueKind := m.getDataKind(fieldValue)

		if fieldValueKind == Nil {
			continue
		}

		value, err := m.transformValue(fieldValue, fieldValueKind, schema)

		if err != nil {
			return nil, err
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
		fieldValueKind := m.getDataKind(fieldValue)

		if fieldValueKind == Nil {
			continue
		}

		for k, v := range schema {
			if v.Key == mapValue.String() {
				schemaValue := schema

				if s, ok := v.Value.(Schema); ok {
					schemaValue = s
				}

				value, err := m.transformValue(fieldValue, fieldValueKind, schemaValue)

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
