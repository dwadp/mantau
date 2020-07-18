# Mantau
Mantau is a golang library for transforming data. Mantau can be used for transforming struct, map and/or slice of struct by providing a schema of how the data will be transformed.

----------
**THIS LIBRARY IS STILL UNDER DEVELOPMENT AND CANNOT BE USED YET**

### Installation
To start using mantau, first you need to install the library by running the following command:
```bash
go get -u github.com/dwadp/mantau
```

### A Brief Explanation
#### Creating instance
When creating a mantau instance, a default options will be passed for initialization. The default struct tag will be set as `json`.
```go
mantau.New()
```

#### Overriding options
You can override mantau default options by calling `SetOpt` function and pass a `mantau.Options`.
```go
m := mantau.New()

// Override the default mantau options
m.SetOpt(&mantau.Options{
    StructTag: "acustomstructtag"
})
```

#### Schema
When transforming a data, mantau will match the data with the provied schema. For examples:

```go
package main

import "github.com/dwadp/mantau"

type AStruct struct {
    FieldOne string `schema:"field_one"`
    FieldTwo string `schema:"field_two"`
}

func main() {
    m := mantau.New()

    // Setting the struct tag to a custom tag 'schema'
    m.SetOpt(&mantau.Options{
        StructTag: "schema",
    })

    // Transforming a struct
    m.Transform(AStruct{
        FieldOne: "Struct value one",
        FieldTwo: "Struct value two",
    }, mantau.Schema{
        "one": SchemaField{
            // Key should match with the struct tag 'schema'
            Key: "field_one"
        },
        "two": SchemaField{
            // Key should match with the struct tag 'schema'
            Key: "field_two"
        },
    })

    // Transforming a map
    m.Transform(map[string]interface{}{
        "key_one": "Map value one",
        "key_two": "Map value two",
    }, mantau.Schema{
        "one": mantau.SchemaField{
            // Key should match with the 'map key'
            Key: "key_one",
        },
        "two": mantau.SchemaField{
            // Key should match with the 'map key'
            Key: "key_two",
        },
    })

    // Transforming nested data structure
    m.Transform(map[string]interface{}{
        "key_one": "Value one",
        "key_two": "Value two",
        "nested": AStruct{
            FieldOne: "Struct value one",
            FieldTwo: "Struct value two",
        },
    }, mantau.Schema{
        "one": mantau.SchemaField{
            // Key should match with the 'map key'
            Key: "key_one",
        },
        "two": mantau.SchemaField{
            // Key should match with the 'map key'
            Key: "key_two",
        },
        "three": mantau.SchemaField{
            // Key should match with the 'map key'
            Key: "nested",
            // Use 'Value' for a schema of a nested data structure
            Value: mantau.Schema{
                "one": mantau.SchemaField{
                    // Key should match with the 'map key'
                    Key: "key_one",
                },
                "two": mantau.SchemaField{
                    // Key should match with the 'map key'
                    Key: "key_two",
                },
            },
        },
    })
}
```

In `mantau.SchemaField` you can leave the `Value` field to nil or omit it if you are not dealing with a nested data structure.

### Examples
Below are some examples on how to use this library.

#### Transforming a struct
Below are example transforming a struct

```go
package main

import (
    "fmt"
    "github.com/dwadp/mantau"
)

type User struct {
    Name        string                   `json:"name"`
    Email       string                   `json:"email"`
    Phone       string                   `json:"phone"`
    IsActive    *bool                    `json:"is_active"`
}

func main() {
    m := mantau.New()

    isActive := true

    result, _ := m.Transform(User{
        Name:     "John doe",
        Email:    "johndoe@example.com",
        Phone:    "911",
        IsActive: &isActive,
    }, mantau.Schema{
        "username": SchemaField{
            Key: "name",
        },
        "useremail": SchemaField{
            Key: "email",
        },
        "active": SchemaField{
            Key: "is_active",
        },
    })

    fmt.Println(result)
}
```

will result:
```
{
    "username": "John doe",
    "useremail": "johndoe@example.com",
    "active": false
}
```

#### Transforming a slice
Below are example transforming a slice

```go
package main

import (
    "fmt"
    "github.com/dwadp/mantau"
)

type User struct {
    Name        string                   `json:"name"`
    Email       string                   `json:"email"`
    Phone       string                   `json:"phone"`
    IsActive    *bool                    `json:"is_active"`
}

func main() {
    m := mantau.New()

    active := true
    inactive := false

    result, _ := m.Transform([]User{
        User{
            Name:     "John doe",
            Email:    "johndoe@example.com",
            Phone:    "911",
            IsActive: &inactive,
        },
        User{
            Name:     "Jane doe",
            Email:    "janedoe@example.com",
            Phone:    "912",
            IsActive: &active,
        }
    }, mantau.Schema{
        "name": mantau.SchemaField{
            Key: "name",
        },
        "email": mantau.SchemaField{
            Key: "email",
        },
        "active": mantau.SchemaField{
            Key: "is_active",
        },
        "phone": mantau.SchemaField{
            Key: "phone",
        },
    })

    fmt.Println(result)
}
```

will result:
```
[
    {
        "name": "John doe",
        "email": "johndoe@example.com",
        "active": false,
        "phone: "911"
    },
    {
        "name": "Jane doe",
        "email": "janedoe@example.com",
        "active": true,
        "phone: "912"
    }
]
```

# TODO
- Write more tests
- Write documentation
- Fix more bugs
