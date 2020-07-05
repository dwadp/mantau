# Mantau
Mantau is a golang library for transforming data. Mantau is used for transforming struct, map and/or slice of struct by providing a schema of how the data will be transformed.

### Installation
To start using mantau, first you need to install the library by running the following command:
```bash
go get -u github.com/dwadp/mantau
```

### A Brief Explanation
Creating a mantau instance

```go
mantau.New("")
```

When creating a mantau instance, you just need to pass one argument which will be used to look for struct tag for the given data. You can pass an empty string to set the default tag which is a `json` tag.

When transforming a data, mantau will match with the provied schema.

```go
// A schema describing of how a data should be transformed
Schema map[string]SchemaField

SchemaField struct {
    // The key should be match with struct tag when initializing mantau instance
    // If the data is a type of map, the key should match with the map key
    Key string

    // Value will be use for nested schema when transforming nested data
    // You can pass a Schema when dealing with nested data or leave it to nil
    Value interface{}
}
```

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
    m := mantau.New("json")

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
    m := mantau.New("")

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

----------
**THIS LIBRARY IS STILL UNDER DEVELOPMENT AND CANNOT BE USED YET**