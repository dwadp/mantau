package mantau

import (
	"reflect"
	"testing"
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
)

// Test data type checking
func TestDataTypeChecking(t *testing.T) {
	mantauInstance := New("json")

	// Should pass a data of struct type either
	// a struct
	// a pointer of struct
	t.Run("ShouldPassStructType", func(t *testing.T) {
		src := User{}

		if _, err := mantauInstance.Transform(src, nil); err != nil {
			t.Error(err)
		}

		if _, err := mantauInstance.Transform(&src, nil); err != nil {
			t.Error(err)
		}
	})

	// Should pass a data of map type either
	// a map
	// a pointer of map
	t.Run("ShouldPassMapType", func(t *testing.T) {
		src := map[string]interface{}{}
		srcString := map[string]string{}

		if _, err := mantauInstance.Transform(src, nil); err != nil {
			t.Error(err)
		}

		if _, err := mantauInstance.Transform(&src, nil); err != nil {
			t.Error(err)
		}

		if _, err := mantauInstance.Transform(srcString, nil); err != nil {
			t.Error(err)
		}

		if _, err := mantauInstance.Transform(&srcString, nil); err != nil {
			t.Error(err)
		}
	})

	// Should pass a data of slice type either
	// a slice
	// a pointer of slice
	// a pointer of struct slice
	t.Run("ShouldPassSliceType", func(t *testing.T) {
		src := make([]User, 3)
		srcPtrOfStruct := make([]*User, 3)

		if _, err := mantauInstance.Transform(src, nil); err != nil {
			t.Error(err)
		}

		if _, err := mantauInstance.Transform(srcPtrOfStruct, nil); err != nil {
			t.Error(err)
		}

		if _, err := mantauInstance.Transform(&src, nil); err != nil {
			t.Error(err)
		}
	})

	// Should return error if we pass a type other than
	// Struct
	// Map
	// Slice or Array
	t.Run("ShouldReturnErrorOfUknownType", func(t *testing.T) {
		if _, err := mantauInstance.Transform(0, nil); err == nil {
			t.Error("Transforming a data of type int should return an error")
		}

		if _, err := mantauInstance.Transform("john doe", nil); err == nil {
			t.Error("Transforming a data of type string should return an error")
		}

		if _, err := mantauInstance.Transform(true, nil); err == nil {
			t.Error("Transforming a data of type boolean should return an error")
		}
	})
}

func TestTransformData(t *testing.T) {
	m := New("json")

	testStructTransforming(t, m)

	testSliceTransforming(t, m)

	testMapTransforming(t, m)

	testWithCustomTag(t)
}

func testStructTransforming(t *testing.T, m *mantau) {
	t.Run("ShouldPassStructTransforming", func(t *testing.T) {
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

		if err != nil {
			t.Error(err)
		}

		if result == nil {
			t.Errorf("Transform should return a data")
		}

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

		res := result.(Result)

		if !reflect.DeepEqual(res, want) {
			t.Errorf("The result do not match\nGOT:\n%+v\nWANT:\n%+v\n", res, want)
		}
	})
}

func testSliceTransforming(t *testing.T, m *mantau) {
	t.Run("ShouldPassSliceOfStructTransforming", func(t *testing.T) {
		sliceOfStruct, err := m.Transform([]Permission{
			{"Admin", 0},
			{"Customer", 1},
			{"Seller", 2},
		}, Schema{
			"name": SchemaField{
				Key: "permission_name",
			},
		})

		if err != nil {
			t.Error(err)
		}

		if sliceOfStruct == nil {
			t.Errorf("Transform should return a data")
		}

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

		if !reflect.DeepEqual(sliceOfStructResult, sliceOfStructResult) {
			t.Errorf(
				"Transformed result do not match\nGOT:%+v\nWANT:\n%+v\n",
				sliceOfStructResult,
				sliceOfStructWant)
		}
	})

	t.Run("ShouldPassSliceOfMapTransforming", func(t *testing.T) {
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

		if err != nil {
			t.Error(err)
		}

		if sliceOfMap == nil {
			t.Errorf("Transform should return a data")
		}

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

		if !reflect.DeepEqual(sliceOfMapResult, sliceOfMapWant) {
			t.Errorf(
				"The result do not match\nGOT:\n%+v\nWANT:\n%+v\n",
				sliceOfMapResult,
				sliceOfMapWant)
		}
	})
}

func testMapTransforming(t *testing.T, m *mantau) {
	t.Run("ShouldPassMapTransforming", func(t *testing.T) {
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

		if err != nil {
			t.Error(err)
		}

		// fmt.Printf("%+v\n", result)

		if result == nil {
			t.Errorf("Transform should return a data")
		}

		want := Result{
			"product_name": "Apple",
			"product_qty":  10,
			"buyer": Result{
				"username":  "John Doe",
				"useremail": "johndoe@example.com",
			},
		}

		if !reflect.DeepEqual(result, want) {
			t.Errorf("The result do not match\nGOT:\n%+v\nWANT:\n%+v\n", result, want)
		}
	})
}

func testWithCustomTag(t *testing.T) {
	t.Run("ShouldPassTransformingACustomTag", func(t *testing.T) {
		var price float32 = 2769

		m := New("schema")

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

		if err != nil {
			t.Error(err)
		}

		if result == nil {
			t.Error("Transform should return a data")
		}

		want := Result{
			"name":  "Apple",
			"price": price,
			"qty":   100,
		}

		res := result.(Result)
		mappedResult := Result{
			"name":  res["name"].(string),
			"price": res["price"].(float32),
			"qty":   res["qty"].(int),
		}

		if !reflect.DeepEqual(mappedResult, want) {
			t.Errorf("The result do not match\nGOT:\n%+v\nWANT:\n%+v\n", mappedResult, want)
		}
	})
}
