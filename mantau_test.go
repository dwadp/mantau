package mantau

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type (
	KindTest struct {
		Name   string
		Result interface{}
		Want   interface{}
	}

	TransformTest struct {
		Name   string
		Schema Schema
		Data   interface{}
		Want   interface{}
	}

	User struct {
		Name        string                   `json:"name"`
		Email       string                   `json:"email"`
		Phone       string                   `json:"phone"`
		IsActive    *bool                    `json:"is_active"`
		Address     UserAddress              `json:"user_address"`
		Permissions []Permission             `json:"permissions"`
		Products    []map[string]interface{} `json:"products"`
	}

	UserAddress struct {
		PostalCode string `json:"postal_code"`
		Address    string `json:"address"`
	}

	Permission struct {
		PermissionName string `json:"permission_name"`
		PermissionCode int    `json:"permission_code"`
	}

	CustomTag struct {
		ProductName  string  `schema:"product_name"`
		ProductPrice float32 `schema:"product_price"`
		ProductQty   int     `schema:"product_qty"`
	}

	Author struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	Book struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Price       float64  `json:"price"`
		Tags        []string `json:"tags"`
		Author      *Author  `json:"author"`
	}
)

// Test for mantau.getKind method
func TestDataKind(t *testing.T) {
	m := New()

	tests := []KindTest{
		{
			Name:   "StructShouldReturnStructKind",
			Result: m.getKind(User{}),
			Want:   Struct,
		},
		{
			Name:   "PointerOfStructShouldReturnPointerKind",
			Result: m.getKind(&User{}),
			Want:   Pointer,
		},
		{
			Name:   "MapShouldReturnMapKind",
			Result: m.getKind(map[string]interface{}{}),
			Want:   Map,
		},
		{
			Name:   "PointerOfMapShouldReturnPointerKind",
			Result: m.getKind(&map[string]interface{}{}),
			Want:   Pointer,
		},
		{
			Name:   "SliceShouldReturnSliceKind",
			Result: m.getKind([]User{}),
			Want:   Slice,
		},
		{
			Name:   "PointerOfSliceShouldReturnPointerKind",
			Result: m.getKind(&[]User{}),
			Want:   Pointer,
		},
		{
			Name:   "ArrayShouldReturnArrayKind",
			Result: m.getKind([3]User{}),
			Want:   Array,
		},
		{
			Name:   "PointerOfArrayShouldReturnPointerKind",
			Result: m.getKind(&[3]User{}),
			Want:   Pointer,
		},
		{
			Name:   "IntegerShouldReturnOtherKind",
			Result: m.getKind(1),
			Want:   Other,
		},
		{
			Name:   "StringShouldReturnOtherKind",
			Result: m.getKind("hello"),
			Want:   Other,
		},
		{
			Name:   "BooleanShouldReturnOtherKind",
			Result: m.getKind(true),
			Want:   Other,
		},
		{
			Name:   "FloatShouldReturnOtherKind",
			Result: m.getKind(0.95),
			Want:   Other,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.Want, test.Result)
		})
	}
}

// Test for mantau.Transform method
func TestTransformMethod(t *testing.T) {
	m := New()
	isActive := true

	tests := []TransformTest{
		{
			Name: "TransformStruct",
			Data: User{
				Name:     "John doe",
				Email:    "johndoe@example.com",
				Phone:    "911",
				IsActive: &isActive,
				Address: UserAddress{
					Address:    "Street",
					PostalCode: "809120",
				},
				Permissions: []Permission{
					{"Admin", 0},
					{"Customer", 1},
					{"Seller", 2},
				},
				Products: []map[string]interface{}{
					{"product_name": "Apple", "product_price": 5, "product_qty": 1},
					{"product_name": "Orange", "product_price": 10, "product_qty": 2},
					{"product_name": "Lemon", "product_price": 10, "product_qty": 2},
				},
			},
			Schema: Schema{
				"useremail": Field{
					Key: "email",
				},
				"username": Field{
					Key: "name",
				},
				"active": Field{
					Key: "is_active",
				},
				"address": Field{
					Key: "user_address",
					Value: Schema{
						"code": Field{
							Key: "postal_code",
						},
						"address": Field{
							Key: "address",
						},
					},
				},
				"user_permissions": Field{
					Key: "permissions",
					Value: Schema{
						"code": Field{
							Key: "permission_code",
						},
						"name": Field{
							Key: "permission_name",
						},
					},
				},
				"products": Field{
					Key: "products",
					Value: Schema{
						"name": Field{
							Key: "product_name",
						},
						"price": Field{
							Key: "product_price",
						},
					},
				},
			},
			Want: Result{
				"useremail": "johndoe@example.com",
				"username":  "John doe",
				"active":    isActive,
				"address": Result{
					"address": "Street",
					"code":    "809120",
				},
				"user_permissions": []Result{
					{"name": "Admin", "code": 0},
					{"name": "Customer", "code": 1},
					{"name": "Seller", "code": 2},
				},
				"products": []Result{
					{"name": "Apple", "price": 5},
					{"name": "Orange", "price": 10},
					{"name": "Lemon", "price": 10},
				},
			},
		},

		{
			Name: "TransformSliceOfStruct",
			Data: []Permission{
				{"Admin", 0},
				{"Customer", 1},
				{"Seller", 2},
			},
			Schema: Schema{
				"name": Field{
					Key: "permission_name",
				},
			},
			Want: []Result{
				{"name": "Admin"},
				{"name": "Customer"},
				{"name": "Seller"},
			},
		},

		{
			Name: "TransformSliceOfMap",
			Data: []map[string]interface{}{
				{"product_name": "Apple", "product_price": 1.50, "product_qty": 50},
				{"product_name": "Banana", "product_price": 2.50, "product_qty": 20},
				{"product_name": "Peach", "product_price": 0.5, "product_qty": 100},
				{"product_name": "Coconut", "product_price": 3.25, "product_qty": 10},
			},
			Schema: Schema{
				"name": Field{
					Key: "product_name",
				},
				"price": Field{
					Key: "product_price",
				},
			},
			Want: []Result{
				{"name": "Apple", "price": 1.50},
				{"name": "Banana", "price": 2.50},
				{"name": "Peach", "price": 0.5},
				{"name": "Coconut", "price": 3.25},
			},
		},

		{
			Name: "TransformArrayOfStruct",
			Data: [3]Permission{
				{"Admin", 0},
				{"Customer", 1},
				{"Seller", 2},
			},
			Schema: Schema{
				"name": Field{
					Key: "permission_name",
				},
			},
			Want: []Result{
				{"name": "Admin"},
				{"name": "Customer"},
				{"name": "Seller"},
			},
		},

		{
			Name: "TransformArrayOfMap",
			Data: [4]map[string]interface{}{
				{"product_name": "Apple", "product_price": 1.50, "product_qty": 50},
				{"product_name": "Banana", "product_price": 2.50, "product_qty": 20},
				{"product_name": "Peach", "product_price": 0.5, "product_qty": 100},
				{"product_name": "Coconut", "product_price": 3.25, "product_qty": 10},
			},
			Schema: Schema{
				"name": Field{
					Key: "product_name",
				},
				"price": Field{
					Key: "product_price",
				},
			},
			Want: []Result{
				{"name": "Apple", "price": 1.50},
				{"name": "Banana", "price": 2.50},
				{"name": "Peach", "price": 0.5},
				{"name": "Coconut", "price": 3.25},
			},
		},

		{
			Name: "TransformMap",
			Data: map[string]interface{}{
				"name":        "Apple",
				"qty":         10,
				"price":       5,
				"description": "A fresh apple",
				"user": &User{
					Name:  "John Doe",
					Email: "johndoe@example.com",
					Phone: "911",
				},
			},
			Schema: Schema{
				"product_name": Field{
					Key: "name",
				},
				"product_qty": Field{
					Key: "qty",
				},
				"buyer": Field{
					Key: "user",
					Value: Schema{
						"username": Field{
							Key: "name",
						},
						"useremail": Field{
							Key: "email",
						},
					},
				},
			},
			Want: Result{
				"product_name": "Apple",
				"product_qty":  10,
				"buyer": Result{
					"username":  "John Doe",
					"useremail": "johndoe@example.com",
				},
			},
		},

		{
			Name:   "TransformNilValue",
			Data:   nil,
			Schema: Schema{},
			Want:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Helper()

			result, err := m.Transform(test.Data, test.Schema)

			assert.NoError(t, err, "Should not return any error")
			assert.Equal(t, test.Want, result, "The result do not match")
		})
	}
}

// Test with a custom hook tag
func TestCustomTagHook(t *testing.T) {
	var price float32 = 2769.99

	m := New()
	m.SetOpt(&Options{
		Hook: "schema",
	})

	data := CustomTag{
		ProductName:  "Apple",
		ProductPrice: price,
		ProductQty:   100,
	}

	result, err := m.Transform(data, Schema{
		"name": Field{
			Key: "product_name",
		},
		"price": Field{
			Key: "product_price",
		},
		"qty": Field{
			Key: "product_qty",
		},
	})

	want := Result{
		"name":  "Apple",
		"price": price,
		"qty":   100,
	}

	assert.NoError(t, err, "Should not return any error")
	assert.NotNil(t, result, "The result should not be a nil value")
	assert.Equal(t, want, result, "The result do not match")
}

func TestShouldSkipTransform(t *testing.T) {
	data := []interface{}{
		new(time.Time), new(string), new(bool), new(int), new(int8), new(int16), new(int32), new(int64),
		new(uint), new(uint8), new(uint16), new(uint32), new(uint64), new(float32), new(float64), new(complex64),
		new(complex128), new(byte), []time.Time{}, []string{}, []bool{}, []int{}, []int8{}, []int16{},
		[]int32{}, []int64{}, []uint{}, []uint8{}, []uint16{}, []uint32{}, []uint64{}, []float32{},
		[]float64{}, []complex64{}, []complex128{}, []byte{},
	}

	m := New()

	for _, val := range data {
		result := m.shouldSkipTransform(val)

		assert.Equal(t, true, result, "The result do not match")
	}
}

func TestValue(t *testing.T) {
	emptyKey := Value{Key: "", Value: 1}

	assert.Equal(t, emptyKey.IsEmpty(), true, "IsEmpty should return true")

	emptyValue := Value{Key: "something", Value: nil}

	assert.Equal(t, emptyValue.IsEmpty(), true, "IsEmpty should return true")
}

func TestGetType(t *testing.T) {
	val := 1

	m := New()

	valType := m.getType(val)
	valTypePtr := m.getType(&val)

	assert.NotNil(t, valType, "Non pointer should not return nil")
	assert.NotNil(t, valTypePtr, "Pointer should not return nil")
}

func TestGetPtrValue(t *testing.T) {
	zeroValues := []interface{}{"", 0, false, nil}

	m := New()

	for _, v := range zeroValues {
		result := m.getPtrValue(&v)

		assert.Nil(t, result, "Zero value should return nil")
	}

	values := []interface{}{"hello", 1, true}

	for _, v := range values {
		ptrResult := m.getPtrValue(&v)
		nonPtrResult := m.getPtrValue(v)

		assert.NotNil(t, ptrResult, "Pointer value should not return nil")
		assert.Nil(t, nonPtrResult, "Non pointer should return nil")
		assert.Equal(t, v, ptrResult, "Pointer should return it's original value")
	}

	value := 1

	assert.Nil(t, m.getPtrValue(nil), "Nil value should return nil")
	assert.Equal(t, value, m.getPtrValue(&value), "Should equal to original value")
}

func TestMapWithSchema(t *testing.T) {
	m := New()

	sample := struct {
		SomeField string `anyhing:"not_found"`
	}{}

	result, err := m.mapWithSchema("not_found", sample, Schema{
		"something": Field{Key: "something"},
	})

	assert.Error(t, err, "Not found struct field should return error")
	assert.True(t, result.IsEmpty(), "Not found struct field should return empty")
}

func TestTransformStruct(t *testing.T) {
	m := New()

	sample := struct {
		SomeField string `anyhing:"not_found"`
	}{}

	nilResult, err := m.transformStruct(nil, Schema{
		"something": Field{Key: "something"},
	})

	assert.Nil(t, nilResult, "Nil should return nil")
	assert.NoError(t, err, "Nil should not return any error")

	result, err := m.transformStruct(sample, Schema{
		"not_found": Field{Key: "not_found"},
	})

	assert.Error(t, err, "If struct field cannot be found, it should return error")
	assert.Nil(t, result, "If struct field cannot be found, the result should be nil")
}

// func TestTransformWithNil(t *testing.T) {
// 	m := New()

// 	t.Run("TransformShouldReturnNilIfTheValueIsNil", func(t *testing.T) {
// 		t.Helper()

// 		result, err := m.Transform(nil, Schema{
// 			"test": Field{
// 				Key: "test",
// 			},
// 			"test2": Field{
// 				Key: "test2",
// 			},
// 		})

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.Nil(t, result, "The result should be a nil value")
// 	})

// 	t.Run("TransformStructShouldReturnNilIfTheValueIsNil", func(t *testing.T) {
// 		t.Helper()

// 		result, err := m.transformStruct(nil, Schema{
// 			"test": Field{
// 				Key: "test",
// 			},
// 			"test2": Field{
// 				Key: "test2",
// 			},
// 		})

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.Nil(t, result, "The result should be a nil value")
// 	})

// 	t.Run("TransformCollectionsShouldReturnNilIfTheValueIsNil", func(t *testing.T) {
// 		t.Helper()

// 		result, err := m.transformCollections(nil, Schema{
// 			"test": Field{
// 				Key: "test",
// 			},
// 			"test2": Field{
// 				Key: "test2",
// 			},
// 		})

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.Nil(t, result, "The result should be a nil value")
// 	})

// 	t.Run("TransformMapShouldReturnNilIfTheValueIsNil", func(t *testing.T) {
// 		t.Helper()

// 		result, err := m.transformMap(nil, Schema{
// 			"test": Field{
// 				Key: "test",
// 			},
// 			"test2": Field{
// 				Key: "test2",
// 			},
// 		})

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.Nil(t, result, "The result should be a nil value")
// 	})
// }

// func TestTransformStruct(t *testing.T) {
// 	m := New()

// 	t.Run("SouldTransformAStruct", func(t *testing.T) {
// 		t.Helper()

// 		result, err := m.transformStruct(Book{
// 			Title:       "A new book",
// 			Description: "Description of a new book",
// 			Price:       35.87,
// 			Tags:        []string{"Romance", "Family"},
// 			Author: &Author{
// 				FirstName: "John",
// 				LastName:  "Doe",
// 			},
// 		}, Schema{
// 			"book_title": Field{
// 				Key: "title",
// 			},
// 			"book_description": Field{
// 				Key: "description",
// 			},
// 			"book_tags": Field{
// 				Key: "tags",
// 			},
// 			"book_author": Field{
// 				Key: "author",
// 				Value: Schema{
// 					"first": Field{
// 						Key: "first_name",
// 					},
// 				},
// 			},
// 		})

// 		want := Result{
// 			"book_title":       "A new book",
// 			"book_description": "Description of a new book",
// 			"book_tags":        []string{"Romance", "Family"},
// 			"book_author": Result{
// 				"first": "John",
// 			},
// 		}

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.NotNil(t, result, "The result should not be a nil value")
// 		assert.Equal(t, want, result, "The result do not match")
// 	})

// 	t.Run("SouldReturnNilWhenTransformingNilValue", func(t *testing.T) {
// 		t.Helper()

// 		result, err := m.transformStruct(nil, Schema{
// 			"test": Field{
// 				Key: "test",
// 			},
// 		})

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.Nil(t, result, "The result should be a nil value")
// 	})

// 	t.Run("SouldTransformANilValueInsideStruct", func(t *testing.T) {
// 		t.Helper()

// 		result, err := m.transformStruct(Book{
// 			Title:       "A new book",
// 			Description: "Description of a new book",
// 			Price:       35.87,
// 			Tags:        []string{"Romance", "Family"},
// 		}, Schema{
// 			"book_title": Field{
// 				Key: "title",
// 			},
// 			"book_description": Field{
// 				Key: "description",
// 			},
// 			"book_tags": Field{
// 				Key: "tags",
// 			},
// 			"book_author": Field{
// 				Key: "author",
// 				Value: Schema{
// 					"first": Field{
// 						Key: "first_name",
// 					},
// 					"last": Field{
// 						Key: "last_name",
// 					},
// 				},
// 			},
// 		})

// 		want := Result{
// 			"book_title":       "A new book",
// 			"book_description": "Description of a new book",
// 			"book_tags":        []string{"Romance", "Family"},
// 			"book_author":      nil,
// 		}

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.NotNil(t, result, "The result shoud not be a nil value")
// 		assert.Equal(t, want, result, "The result do not match")
// 	})
// }

// func TestTransformCollections(t *testing.T) {
// 	m := New()

// 	t.Run("SouldTransformASliceOfStruct", func(t *testing.T) {
// 		t.Helper()

// 		result, err := m.transformCollections([]Book{
// 			{
// 				Title:       "A new book",
// 				Description: "Description of a new book",
// 				Price:       35.87,
// 				Tags:        []string{"Romance", "Family"},
// 				Author: &Author{
// 					FirstName: "John",
// 					LastName:  "Doe",
// 				},
// 			},
// 			{
// 				Title:       "An old book",
// 				Description: "Description of an old book",
// 				Price:       100.99,
// 				Tags:        []string{"History"},
// 				Author: &Author{
// 					FirstName: "Jane",
// 					LastName:  "Doe",
// 				},
// 			},
// 		}, Schema{
// 			"book_title": Field{
// 				Key: "title",
// 			},
// 			"book_description": Field{
// 				Key: "description",
// 			},
// 			"book_tags": Field{
// 				Key: "tags",
// 			},
// 			"book_author": Field{
// 				Key: "author",
// 				Value: Schema{
// 					"first": Field{
// 						Key: "first_name",
// 					},
// 					"last": Field{
// 						Key: "last_name",
// 					},
// 				},
// 			},
// 		})

// 		want := []Result{
// 			{
// 				"book_title":       "A new book",
// 				"book_description": "Description of a new book",
// 				"book_tags":        []string{"Romance", "Family"},
// 				"book_author": Result{
// 					"first": "John",
// 					"last":  "Doe",
// 				},
// 			},
// 			{
// 				"book_title":       "An old book",
// 				"book_description": "Description of an old book",
// 				"book_tags":        []string{"History"},
// 				"book_author": Result{
// 					"first": "Jane",
// 					"last":  "Doe",
// 				},
// 			},
// 		}

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.NotNil(t, result, "The result should not be a nil value")
// 		assert.Equal(t, want, result, "The result do not match")
// 	})

// 	t.Run("SouldTransformAnArrayOfStruct", func(t *testing.T) {
// 		t.Helper()

// 		result, err := m.transformCollections([2]Book{
// 			{
// 				Title:       "A new book",
// 				Description: "Description of a new book",
// 				Price:       35.87,
// 				Tags:        []string{"Romance", "Family"},
// 				Author: &Author{
// 					FirstName: "John",
// 					LastName:  "Doe",
// 				},
// 			},
// 			{
// 				Title:       "An old book",
// 				Description: "Description of an old book",
// 				Price:       100.99,
// 				Tags:        []string{"History"},
// 				Author: &Author{
// 					FirstName: "Jane",
// 					LastName:  "Doe",
// 				},
// 			},
// 		}, Schema{
// 			"book_title": Field{
// 				Key: "title",
// 			},
// 			"book_description": Field{
// 				Key: "description",
// 			},
// 			"book_tags": Field{
// 				Key: "tags",
// 			},
// 			"book_author": Field{
// 				Key: "author",
// 				Value: Schema{
// 					"first": Field{
// 						Key: "first_name",
// 					},
// 					"last": Field{
// 						Key: "last_name",
// 					},
// 				},
// 			},
// 		})

// 		want := [2]Result{
// 			{
// 				"book_title":       "A new book",
// 				"book_description": "Description of a new book",
// 				"book_tags":        []string{"Romance", "Family"},
// 				"book_author": Result{
// 					"first": "John",
// 					"last":  "Doe",
// 				},
// 			},
// 			{
// 				"book_title":       "An old book",
// 				"book_description": "Description of an old book",
// 				"book_tags":        []string{"History"},
// 				"book_author": Result{
// 					"first": "Jane",
// 					"last":  "Doe",
// 				},
// 			},
// 		}

// 		results := [2]Result{}

// 		for k, v := range result {
// 			results[k] = v
// 		}

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.NotNil(t, result, "The result should not be a nil value")
// 		assert.Equal(t, want, results, "The result do not match")
// 	})

// 	t.Run("SouldTransformASliceOfMap", func(t *testing.T) {
// 		t.Helper()

// 		result, err := m.transformCollections([]map[string]interface{}{
// 			{
// 				"title":       "A new book",
// 				"description": "Description of a new book",
// 				"tags":        []string{"Romance", "Family"},
// 				"author": map[string]string{
// 					"first_name": "John",
// 					"last_name":  "Doe",
// 				},
// 			},
// 			{
// 				"title":       "An old book",
// 				"description": "Description of an old book",
// 				"tags":        []string{"History"},
// 				"author": map[string]string{
// 					"first_name": "Jane",
// 					"last_name":  "Doe",
// 				},
// 			},
// 		}, Schema{
// 			"book_title": Field{
// 				Key: "title",
// 			},
// 			"book_description": Field{
// 				Key: "description",
// 			},
// 			"book_tags": Field{
// 				Key: "tags",
// 			},
// 			"book_author": Field{
// 				Key: "author",
// 				Value: Schema{
// 					"first": Field{
// 						Key: "first_name",
// 					},
// 					"last": Field{
// 						Key: "last_name",
// 					},
// 				},
// 			},
// 		})

// 		want := []Result{
// 			{
// 				"book_title":       "A new book",
// 				"book_description": "Description of a new book",
// 				"book_tags":        []string{"Romance", "Family"},
// 				"book_author": Result{
// 					"first": "John",
// 					"last":  "Doe",
// 				},
// 			},
// 			{
// 				"book_title":       "An old book",
// 				"book_description": "Description of an old book",
// 				"book_tags":        []string{"History"},
// 				"book_author": Result{
// 					"first": "Jane",
// 					"last":  "Doe",
// 				},
// 			},
// 		}

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.NotNil(t, result, "The result should not be a nil value")
// 		assert.Equal(t, want, result, "The result do not match")
// 	})

// 	t.Run("SouldTransformAnArrayOfMap", func(t *testing.T) {
// 		t.Helper()

// 		result, err := m.transformCollections([2]map[string]interface{}{
// 			{
// 				"title":       "A new book",
// 				"description": "Description of a new book",
// 				"tags":        []string{"Romance", "Family"},
// 				"author": map[string]string{
// 					"first_name": "John",
// 					"last_name":  "Doe",
// 				},
// 			},
// 			{
// 				"title":       "An old book",
// 				"description": "Description of an old book",
// 				"tags":        []string{"History"},
// 				"author": map[string]string{
// 					"first_name": "Jane",
// 					"last_name":  "Doe",
// 				},
// 			},
// 		}, Schema{
// 			"book_title": Field{
// 				Key: "title",
// 			},
// 			"book_description": Field{
// 				Key: "description",
// 			},
// 			"book_tags": Field{
// 				Key: "tags",
// 			},
// 			"book_author": Field{
// 				Key: "author",
// 				Value: Schema{
// 					"first": Field{
// 						Key: "first_name",
// 					},
// 					"last": Field{
// 						Key: "last_name",
// 					},
// 				},
// 			},
// 		})

// 		want := [2]Result{
// 			{
// 				"book_title":       "A new book",
// 				"book_description": "Description of a new book",
// 				"book_tags":        []string{"Romance", "Family"},
// 				"book_author": Result{
// 					"first": "John",
// 					"last":  "Doe",
// 				},
// 			},
// 			{
// 				"book_title":       "An old book",
// 				"book_description": "Description of an old book",
// 				"book_tags":        []string{"History"},
// 				"book_author": Result{
// 					"first": "Jane",
// 					"last":  "Doe",
// 				},
// 			},
// 		}

// 		results := [2]Result{}

// 		for k, v := range result {
// 			results[k] = v
// 		}

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.NotNil(t, result, "The result should not be a nil value")
// 		assert.Equal(t, want, results, "The result do not match")
// 	})

// 	t.Run("SouldReturnNilWhenTransformingNilValue", func(t *testing.T) {
// 		t.Helper()

// 		result, err := m.transformCollections(nil, Schema{
// 			"test": Field{
// 				Key: "test",
// 			},
// 		})

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.Nil(t, result, "The result should be a nil value")
// 	})
// }

// func TestTransformMap(t *testing.T) {
// 	m := New()

// 	t.Run("ShouldTransformAMap", func(t *testing.T) {
// 		t.Helper()

// 		movieReleaseDate := time.Date(2019, 12, 13, 20, 0, 0, 0, time.UTC)

// 		result, err := m.transformMap(map[string]interface{}{
// 			"name":         "6 Underground",
// 			"release_date": &movieReleaseDate,
// 			"platform":     "netflix",
// 			"running_time": 128,
// 			"country":      "United States",
// 			"budget":       150000000,
// 		}, Schema{
// 			"movieName": Field{
// 				Key: "name",
// 			},
// 			"movieReleaseDate": Field{
// 				Key: "release_date",
// 			},
// 			"movieRunningTime": Field{
// 				Key: "running_time",
// 			},
// 			"movieBudget": Field{
// 				Key: "budget",
// 			},
// 		})

// 		want := Result{
// 			"movieName":        "6 Underground",
// 			"movieReleaseDate": movieReleaseDate,
// 			"movieRunningTime": 128,
// 			"movieBudget":      150000000,
// 		}

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.NotNil(t, result, "The result should not be a nil value")
// 		assert.Equal(t, want, result, "The result do not match")
// 	})

// 	t.Run("ShouldReturnNilIfValueIsNil", func(t *testing.T) {
// 		t.Helper()

// 		result, err := m.transformMap(nil, Schema{
// 			"testKey": Field{
// 				Key: "testKeySource",
// 			},
// 		})

// 		assert.NoError(t, err, "Should not return any error")
// 		assert.Nil(t, result, "The result should be a nil value")
// 	})
// }
