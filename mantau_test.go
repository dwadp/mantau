package mantau

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type (
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

// Test data kind checking
func TestDataKindChecking(t *testing.T) {
	m := New()

	// Should pass a data of struct type either of: a struct, a pointer of struct
	t.Run("ShouldReturnStructAndPointerKind", func(t *testing.T) {
		t.Helper()

		src := User{}
		result := m.getDataKind(src)
		resultPtr := m.getDataKind(&src)

		want := Struct
		wantPtr := Pointer
		dontWant := Other

		assert.NotEqual(t, dontWant, result)
		assert.NotEqual(t, dontWant, resultPtr)

		assert.Equal(t, want, result)
		assert.Equal(t, wantPtr, resultPtr)
	})

	// Should pass a data of map type either of: a map, a pointer of map
	t.Run("ShouldReturnMapAndPointerKind", func(t *testing.T) {
		t.Helper()

		src := map[string]interface{}{}

		result := m.getDataKind(src)
		resultPtr := m.getDataKind(&src)

		want := Map
		wantPtr := Pointer
		dontWant := Other

		assert.NotEqual(t, dontWant, result)
		assert.NotEqual(t, dontWant, resultPtr)

		assert.Equal(t, want, result)
		assert.Equal(t, wantPtr, resultPtr)
	})

	// Should pass a data of slice type either of: a slice, a pointer of slice,
	// a pointer of struct slice
	t.Run("ShouldReturnSliceAndPointerKind", func(t *testing.T) {
		t.Helper()

		src := make([]User, 3)
		srcPtrOfStruct := make([]*User, 3)

		result := m.getDataKind(src)
		resultPtr := m.getDataKind(&src)

		resultStruct := m.getDataKind(srcPtrOfStruct)
		resultStructPtr := m.getDataKind(&srcPtrOfStruct)

		want := Slice
		wantPtr := Pointer
		dontWant := Other

		assert.NotEqual(t, dontWant, result)
		assert.NotEqual(t, dontWant, resultPtr)

		assert.NotEqual(t, dontWant, resultStruct)
		assert.NotEqual(t, dontWant, resultStructPtr)

		assert.Equal(t, want, result)
		assert.Equal(t, wantPtr, resultPtr)

		assert.Equal(t, want, resultStruct)
		assert.Equal(t, wantPtr, resultStructPtr)
	})
}

// Should return error if we pass a type other than, Struct, Map, Slice or Array
func TestTransform(t *testing.T) {
	m := New()

	_, err := m.Transform(0, nil)

	assert.Error(t, err, "Should return an error")

	_, err = m.Transform("john doe", nil)

	assert.Error(t, err, "Should return an error")

	_, err = m.Transform(true, nil)

	assert.Error(t, err, "Should return an error")
}

func TestTransformData(t *testing.T) {
	m := New()

	testStructTransforming(t, m)
	testSliceTransforming(t, m)
	testArrayTransforming(t, m)
	testMapTransforming(t, m)
	testWithCustomTag(t)
}

func testStructTransforming(t *testing.T, m *mantau) {
	t.Run("ShouldPassStructTransforming", func(t *testing.T) {
		t.Helper()

		isActive := true

		result, err := m.Transform(User{
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
		}, Schema{
			"useremail": SchemaField{
				Key: "email",
			},
			"username": SchemaField{
				Key: "name",
			},
			"active": SchemaField{
				Key: "is_active",
			},
			"address": SchemaField{
				Key: "user_address",
				Value: Schema{
					"code": SchemaField{
						Key: "postal_code",
					},
					"address": SchemaField{
						Key: "address",
					},
				},
			},
			"user_permissions": SchemaField{
				Key: "permissions",
				Value: Schema{
					"code": SchemaField{
						Key: "permission_code",
					},
					"name": SchemaField{
						Key: "permission_name",
					},
				},
			},
			"products": SchemaField{
				Key: "products",
				Value: Schema{
					"name": SchemaField{
						Key: "product_name",
					},
					"price": SchemaField{
						Key: "product_price",
					},
				},
			},
		})

		want := Result{
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
		}

		assert.NoError(t, err, "Should not return any error")
		assert.NotNil(t, result, "The result shoud not be a nil value")
		assert.Equal(t, want, result, "The result do not match")
	})
}

func testSliceTransforming(t *testing.T, m *mantau) {
	t.Run("ShouldPassSliceOfStructTransforming", func(t *testing.T) {
		t.Helper()

		sliceOfStruct, err := m.Transform([]Permission{
			{"Admin", 0},
			{"Customer", 1},
			{"Seller", 2},
		}, Schema{
			"name": SchemaField{
				Key: "permission_name",
			},
		})

		sliceOfStructResult := []Result{}
		sliceOfStructWant := []Result{
			{"name": "Admin"},
			{"name": "Customer"},
			{"name": "Seller"},
		}

		for i := 0; i < reflect.ValueOf(sliceOfStruct).Len(); i++ {
			value := reflect.ValueOf(sliceOfStruct).Index(i).Interface()
			result := value.(Result)

			sliceOfStructResult = append(sliceOfStructResult, result)
		}

		assert.NoError(t, err, "Should not return any error")
		assert.NotNil(t, sliceOfStruct, "The result should not be a nil value")
		assert.Equal(t, sliceOfStructWant, sliceOfStructResult, "The result do not match")
	})

	t.Run("ShouldPassSliceOfMapTransforming", func(t *testing.T) {
		t.Helper()

		sliceOfMap, err := m.Transform([]map[string]interface{}{
			{"product_name": "Apple", "product_price": 1.50, "product_qty": 50},
			{"product_name": "Banana", "product_price": 2.50, "product_qty": 20},
			{"product_name": "Peach", "product_price": 0.5, "product_qty": 100},
			{"product_name": "Coconut", "product_price": 3.25, "product_qty": 10},
		}, Schema{
			"name": SchemaField{
				Key: "product_name",
			},
			"price": SchemaField{
				Key: "product_price",
			},
		})

		sliceOfMapResult := []Result{}
		sliceOfMapWant := []Result{
			{"name": "Apple", "price": 1.50},
			{"name": "Banana", "price": 2.50},
			{"name": "Peach", "price": 0.5},
			{"name": "Coconut", "price": 3.25},
		}

		for i := 0; i < reflect.ValueOf(sliceOfMap).Len(); i++ {
			value := reflect.ValueOf(sliceOfMap).Index(i).Interface()
			result := value.(Result)

			sliceOfMapResult = append(sliceOfMapResult, result)
		}

		assert.NoError(t, err, "Should not return any error")
		assert.NotNil(t, sliceOfMap, "The resultl shoud not be a nil value")
		assert.Equal(t, sliceOfMapWant, sliceOfMapResult, "The result do not match")
	})
}

func testArrayTransforming(t *testing.T, m *mantau) {
	t.Run("ShouldPassArrayOfStructTransforming", func(t *testing.T) {
		t.Helper()

		arrayOfStruct, err := m.Transform([3]Permission{
			{"Admin", 0},
			{"Customer", 1},
			{"Seller", 2},
		}, Schema{
			"name": SchemaField{
				Key: "permission_name",
			},
		})

		arrayOfStructResult := [3]Result{}
		arrayOfStructWant := [3]Result{
			{"name": "Admin"},
			{"name": "Customer"},
			{"name": "Seller"},
		}

		for i := 0; i < reflect.ValueOf(arrayOfStruct).Len(); i++ {
			value := reflect.ValueOf(arrayOfStruct).Index(i).Interface()
			result := value.(Result)

			arrayOfStructResult[i] = result
		}

		assert.NoError(t, err, "Should not return any error")
		assert.NotNil(t, arrayOfStruct, "The result should not be a nil value")
		assert.Equal(t, arrayOfStructWant, arrayOfStructResult, "The result do not match")
	})

	t.Run("ShouldPassArrayOfMapTransforming", func(t *testing.T) {
		t.Helper()

		arrayOfMap, err := m.Transform([4]map[string]interface{}{
			{"product_name": "Apple", "product_price": 1.50, "product_qty": 50},
			{"product_name": "Banana", "product_price": 2.50, "product_qty": 20},
			{"product_name": "Peach", "product_price": 0.5, "product_qty": 100},
			{"product_name": "Coconut", "product_price": 3.25, "product_qty": 10},
		}, Schema{
			"name": SchemaField{
				Key: "product_name",
			},
			"price": SchemaField{
				Key: "product_price",
			},
		})

		arrayOfMapResult := [4]Result{}
		arrayOfMapWant := [4]Result{
			{"name": "Apple", "price": 1.50},
			{"name": "Banana", "price": 2.50},
			{"name": "Peach", "price": 0.5},
			{"name": "Coconut", "price": 3.25},
		}

		for i := 0; i < reflect.ValueOf(arrayOfMap).Len(); i++ {
			value := reflect.ValueOf(arrayOfMap).Index(i).Interface()
			result := value.(Result)

			arrayOfMapResult[i] = result
		}

		assert.NoError(t, err, "Should not return any error")
		assert.NotNil(t, arrayOfMap, "The result should not be a nil value")
		assert.Equal(t, arrayOfMapWant, arrayOfMapResult, "The result do not match")
	})
}

func testMapTransforming(t *testing.T, m *mantau) {
	t.Run("ShouldPassMapTransforming", func(t *testing.T) {
		t.Helper()

		result, err := m.Transform(map[string]interface{}{
			"name":        "Apple",
			"qty":         10,
			"price":       5,
			"description": "A fresh apple",
			"user": &User{
				Name:  "John Doe",
				Email: "johndoe@example.com",
				Phone: "911",
			},
		}, Schema{
			"product_name": SchemaField{
				Key: "name",
			},
			"product_qty": SchemaField{
				Key: "qty",
			},
			"buyer": SchemaField{
				Key: "user",
				Value: Schema{
					"username": SchemaField{
						Key: "name",
					},
					"useremail": SchemaField{
						Key: "email",
					},
				},
			},
		})

		want := Result{
			"product_name": "Apple",
			"product_qty":  10,
			"buyer": Result{
				"username":  "John Doe",
				"useremail": "johndoe@example.com",
			},
		}

		assert.NoError(t, err, "Should not return any error")
		assert.NotNil(t, result, "The result should not be a nil value")
		assert.Equal(t, want, result, "The result do not match")
	})
}

func testWithCustomTag(t *testing.T) {
	t.Run("ShouldPassTransformingACustomTag", func(t *testing.T) {
		t.Helper()

		var price float32 = 2769.99

		m := New()
		m.SetOpt(&Options{
			StructTag: "schema",
		})

		data := CustomTag{
			ProductName:  "Apple",
			ProductPrice: price,
			ProductQty:   100,
		}

		result, err := m.Transform(data, Schema{
			"name": SchemaField{
				Key: "product_name",
			},
			"price": SchemaField{
				Key: "product_price",
			},
			"qty": SchemaField{
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
	})
}

func TestTransformWithNil(t *testing.T) {
	m := New()

	t.Run("TransformShouldReturnNilIfTheValueIsNil", func(t *testing.T) {
		t.Helper()

		result, err := m.Transform(nil, Schema{
			"test": SchemaField{
				Key: "test",
			},
			"test2": SchemaField{
				Key: "test2",
			},
		})

		assert.NoError(t, err, "Should not return any error")
		assert.Nil(t, result, "The result should be a nil value")
	})

	t.Run("TransformStructShouldReturnNilIfTheValueIsNil", func(t *testing.T) {
		t.Helper()

		result, err := m.transformStruct(nil, Schema{
			"test": SchemaField{
				Key: "test",
			},
			"test2": SchemaField{
				Key: "test2",
			},
		})

		assert.NoError(t, err, "Should not return any error")
		assert.Nil(t, result, "The result should be a nil value")
	})

	t.Run("TransformCollectionsShouldReturnNilIfTheValueIsNil", func(t *testing.T) {
		t.Helper()

		result, err := m.transformCollections(nil, Schema{
			"test": SchemaField{
				Key: "test",
			},
			"test2": SchemaField{
				Key: "test2",
			},
		})

		assert.NoError(t, err, "Should not return any error")
		assert.Nil(t, result, "The result should be a nil value")
	})

	t.Run("TransformMapShouldReturnNilIfTheValueIsNil", func(t *testing.T) {
		t.Helper()

		result, err := m.transformMap(nil, Schema{
			"test": SchemaField{
				Key: "test",
			},
			"test2": SchemaField{
				Key: "test2",
			},
		})

		assert.NoError(t, err, "Should not return any error")
		assert.Nil(t, result, "The result should be a nil value")
	})
}

func TestTransformStruct(t *testing.T) {
	m := New()

	t.Run("SouldTransformAStruct", func(t *testing.T) {
		t.Helper()

		result, err := m.transformStruct(Book{
			Title:       "A new book",
			Description: "Description of a new book",
			Price:       35.87,
			Tags:        []string{"Romance", "Family"},
			Author: &Author{
				FirstName: "John",
				LastName:  "Doe",
			},
		}, Schema{
			"book_title": SchemaField{
				Key: "title",
			},
			"book_description": SchemaField{
				Key: "description",
			},
			"book_tags": SchemaField{
				Key: "tags",
			},
			"book_author": SchemaField{
				Key: "author",
				Value: Schema{
					"first": SchemaField{
						Key: "first_name",
					},
				},
			},
		})

		want := Result{
			"book_title":       "A new book",
			"book_description": "Description of a new book",
			"book_tags":        []string{"Romance", "Family"},
			"book_author": Result{
				"first": "John",
			},
		}

		assert.NoError(t, err, "Should not return any error")
		assert.NotNil(t, result, "The result should not be a nil value")
		assert.Equal(t, want, result, "The result do not match")
	})

	t.Run("SouldReturnNilWhenTransformingNilValue", func(t *testing.T) {
		t.Helper()

		result, err := m.transformStruct(nil, Schema{
			"test": SchemaField{
				Key: "test",
			},
		})

		assert.NoError(t, err, "Should not return any error")
		assert.Nil(t, result, "The result should be a nil value")
	})

	t.Run("SouldTransformANilValueInsideStruct", func(t *testing.T) {
		t.Helper()

		result, err := m.transformStruct(Book{
			Title:       "A new book",
			Description: "Description of a new book",
			Price:       35.87,
			Tags:        []string{"Romance", "Family"},
		}, Schema{
			"book_title": SchemaField{
				Key: "title",
			},
			"book_description": SchemaField{
				Key: "description",
			},
			"book_tags": SchemaField{
				Key: "tags",
			},
			"book_author": SchemaField{
				Key: "author",
				Value: Schema{
					"first": SchemaField{
						Key: "first_name",
					},
					"last": SchemaField{
						Key: "last_name",
					},
				},
			},
		})

		want := Result{
			"book_title":       "A new book",
			"book_description": "Description of a new book",
			"book_tags":        []string{"Romance", "Family"},
			"book_author":      nil,
		}

		assert.NoError(t, err, "Should not return any error")
		assert.NotNil(t, result, "The result shoud not be a nil value")
		assert.Equal(t, want, result, "The result do not match")
	})
}

func TestTransformCollections(t *testing.T) {
	m := New()

	t.Run("SouldTransformASliceOfStruct", func(t *testing.T) {
		t.Helper()

		result, err := m.transformCollections([]Book{
			{
				Title:       "A new book",
				Description: "Description of a new book",
				Price:       35.87,
				Tags:        []string{"Romance", "Family"},
				Author: &Author{
					FirstName: "John",
					LastName:  "Doe",
				},
			},
			{
				Title:       "An old book",
				Description: "Description of an old book",
				Price:       100.99,
				Tags:        []string{"History"},
				Author: &Author{
					FirstName: "Jane",
					LastName:  "Doe",
				},
			},
		}, Schema{
			"book_title": SchemaField{
				Key: "title",
			},
			"book_description": SchemaField{
				Key: "description",
			},
			"book_tags": SchemaField{
				Key: "tags",
			},
			"book_author": SchemaField{
				Key: "author",
				Value: Schema{
					"first": SchemaField{
						Key: "first_name",
					},
					"last": SchemaField{
						Key: "last_name",
					},
				},
			},
		})

		want := []Result{
			{
				"book_title":       "A new book",
				"book_description": "Description of a new book",
				"book_tags":        []string{"Romance", "Family"},
				"book_author": Result{
					"first": "John",
					"last":  "Doe",
				},
			},
			{
				"book_title":       "An old book",
				"book_description": "Description of an old book",
				"book_tags":        []string{"History"},
				"book_author": Result{
					"first": "Jane",
					"last":  "Doe",
				},
			},
		}

		assert.NoError(t, err, "Should not return any error")
		assert.NotNil(t, result, "The result should not be a nil value")
		assert.Equal(t, want, result, "The result do not match")
	})

	t.Run("SouldTransformAnArrayOfStruct", func(t *testing.T) {
		t.Helper()

		result, err := m.transformCollections([2]Book{
			{
				Title:       "A new book",
				Description: "Description of a new book",
				Price:       35.87,
				Tags:        []string{"Romance", "Family"},
				Author: &Author{
					FirstName: "John",
					LastName:  "Doe",
				},
			},
			{
				Title:       "An old book",
				Description: "Description of an old book",
				Price:       100.99,
				Tags:        []string{"History"},
				Author: &Author{
					FirstName: "Jane",
					LastName:  "Doe",
				},
			},
		}, Schema{
			"book_title": SchemaField{
				Key: "title",
			},
			"book_description": SchemaField{
				Key: "description",
			},
			"book_tags": SchemaField{
				Key: "tags",
			},
			"book_author": SchemaField{
				Key: "author",
				Value: Schema{
					"first": SchemaField{
						Key: "first_name",
					},
					"last": SchemaField{
						Key: "last_name",
					},
				},
			},
		})

		want := [2]Result{
			{
				"book_title":       "A new book",
				"book_description": "Description of a new book",
				"book_tags":        []string{"Romance", "Family"},
				"book_author": Result{
					"first": "John",
					"last":  "Doe",
				},
			},
			{
				"book_title":       "An old book",
				"book_description": "Description of an old book",
				"book_tags":        []string{"History"},
				"book_author": Result{
					"first": "Jane",
					"last":  "Doe",
				},
			},
		}

		results := [2]Result{}

		for k, v := range result {
			results[k] = v
		}

		assert.NoError(t, err, "Should not return any error")
		assert.NotNil(t, result, "The result should not be a nil value")
		assert.Equal(t, want, results, "The result do not match")
	})

	t.Run("SouldTransformASliceOfMap", func(t *testing.T) {
		t.Helper()

		result, err := m.transformCollections([]map[string]interface{}{
			{
				"title":       "A new book",
				"description": "Description of a new book",
				"tags":        []string{"Romance", "Family"},
				"author": map[string]string{
					"first_name": "John",
					"last_name":  "Doe",
				},
			},
			{
				"title":       "An old book",
				"description": "Description of an old book",
				"tags":        []string{"History"},
				"author": map[string]string{
					"first_name": "Jane",
					"last_name":  "Doe",
				},
			},
		}, Schema{
			"book_title": SchemaField{
				Key: "title",
			},
			"book_description": SchemaField{
				Key: "description",
			},
			"book_tags": SchemaField{
				Key: "tags",
			},
			"book_author": SchemaField{
				Key: "author",
				Value: Schema{
					"first": SchemaField{
						Key: "first_name",
					},
					"last": SchemaField{
						Key: "last_name",
					},
				},
			},
		})

		want := []Result{
			{
				"book_title":       "A new book",
				"book_description": "Description of a new book",
				"book_tags":        []string{"Romance", "Family"},
				"book_author": Result{
					"first": "John",
					"last":  "Doe",
				},
			},
			{
				"book_title":       "An old book",
				"book_description": "Description of an old book",
				"book_tags":        []string{"History"},
				"book_author": Result{
					"first": "Jane",
					"last":  "Doe",
				},
			},
		}

		assert.NoError(t, err, "Should not return any error")
		assert.NotNil(t, result, "The result should not be a nil value")
		assert.Equal(t, want, result, "The result do not match")
	})

	t.Run("SouldTransformAnArrayOfMap", func(t *testing.T) {
		t.Helper()

		result, err := m.transformCollections([2]map[string]interface{}{
			{
				"title":       "A new book",
				"description": "Description of a new book",
				"tags":        []string{"Romance", "Family"},
				"author": map[string]string{
					"first_name": "John",
					"last_name":  "Doe",
				},
			},
			{
				"title":       "An old book",
				"description": "Description of an old book",
				"tags":        []string{"History"},
				"author": map[string]string{
					"first_name": "Jane",
					"last_name":  "Doe",
				},
			},
		}, Schema{
			"book_title": SchemaField{
				Key: "title",
			},
			"book_description": SchemaField{
				Key: "description",
			},
			"book_tags": SchemaField{
				Key: "tags",
			},
			"book_author": SchemaField{
				Key: "author",
				Value: Schema{
					"first": SchemaField{
						Key: "first_name",
					},
					"last": SchemaField{
						Key: "last_name",
					},
				},
			},
		})

		want := [2]Result{
			{
				"book_title":       "A new book",
				"book_description": "Description of a new book",
				"book_tags":        []string{"Romance", "Family"},
				"book_author": Result{
					"first": "John",
					"last":  "Doe",
				},
			},
			{
				"book_title":       "An old book",
				"book_description": "Description of an old book",
				"book_tags":        []string{"History"},
				"book_author": Result{
					"first": "Jane",
					"last":  "Doe",
				},
			},
		}

		results := [2]Result{}

		for k, v := range result {
			results[k] = v
		}

		assert.NoError(t, err, "Should not return any error")
		assert.NotNil(t, result, "The result should not be a nil value")
		assert.Equal(t, want, results, "The result do not match")
	})

	t.Run("SouldReturnNilWhenTransformingNilValue", func(t *testing.T) {
		t.Helper()

		result, err := m.transformCollections(nil, Schema{
			"test": SchemaField{
				Key: "test",
			},
		})

		assert.NoError(t, err, "Should not return any error")
		assert.Nil(t, result, "The result should be a nil value")
	})
}

func TestTransformMap(t *testing.T) {
	m := New()

	t.Run("ShouldTransformAMap", func(t *testing.T) {
		t.Helper()

		movieReleaseDate := time.Date(2019, 12, 13, 20, 0, 0, 0, time.UTC)

		result, err := m.transformMap(map[string]interface{}{
			"name":         "6 Underground",
			"release_date": &movieReleaseDate,
			"platform":     "netflix",
			"running_time": 128,
			"country":      "United States",
			"budget":       150000000,
		}, Schema{
			"movieName": SchemaField{
				Key: "name",
			},
			"movieReleaseDate": SchemaField{
				Key: "release_date",
			},
			"movieRunningTime": SchemaField{
				Key: "running_time",
			},
			"movieBudget": SchemaField{
				Key: "budget",
			},
		})

		want := Result{
			"movieName":        "6 Underground",
			"movieReleaseDate": movieReleaseDate,
			"movieRunningTime": 128,
			"movieBudget":      150000000,
		}

		assert.NoError(t, err, "Should not return any error")
		assert.NotNil(t, result, "The result should not be a nil value")
		assert.Equal(t, want, result, "The result do not match")
	})

	t.Run("ShouldReturnNilIfValueIsNil", func(t *testing.T) {
		t.Helper()

		result, err := m.transformMap(nil, Schema{
			"testKey": SchemaField{
				Key: "testKeySource",
			},
		})

		assert.NoError(t, err, "Should not return any error")
		assert.Nil(t, result, "The result should be a nil value")
	})
}
