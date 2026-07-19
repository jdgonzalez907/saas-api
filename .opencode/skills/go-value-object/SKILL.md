---
name: go-value-object
description: Use when the user wants to create, modify, or refactor Go value objects. Generates immutable value objects with New constructors, own error definitions, getters, Equals (by value), and ToDTO methods. Handles single-field (primitive types) and multi-field (structs) value objects. Creates or updates functional file and table-driven tests for 100% coverage.
---

# Go Value Object Generator

Generates immutable value objects in Go following these rules:

## Structure Rules

| Fields | Implementation |
|--------|----------------|
| 1 field | `type Xxx primitiveType` (no DTO, ToDTO returns primitive) |
| 2+ fields | `type Xxx struct { ... }` with `type XxxDTO struct { ... }` |

## Calculated Fields

- Fields computed from other fields
- NOT included in DTO
- Only accessible via getter method

```go
type Money struct {
    amount   int64
    currency string
}

// Calculated field - not in DTO
func (m Money) IsPositive() bool {
    return m.amount > 0
}

func (m Money) ToDTO() MoneyDTO {
    return MoneyDTO{Amount: m.amount, Currency: m.currency}
    // IsPositive NOT in DTO
}
```

## Required Elements

### File: `<name>.go`

```go
package <package>

import "errors"

// Errors
var (
    Err<Name>Invalid = errors.New("<name> is invalid")
    // Add specific errors as needed
)

// Value Object
type <Name> struct {
    field1 type1
    field2 type2
}

// DTO
type <Name>DTO struct {
    Field1 type1 `json:"field1"`
    Field2 type2 `json:"field2"`
}

// Constructor - always returns (VO, error)
func New<Name>(field1 type1, field2 type2) (<Name>, error) {
    // Validation here
    if invalid {
        return <Name>{}, Err<Name>Invalid
    }
    return <Name>{field1: field1, field2: field2}, nil
}

// Getters
func (v <Name>) Field1() type1 { return v.field1 }
func (v <Name>) Field2() type2 { return v.field2 }

// Equals - compares by value, field by field
func (v <Name>) Equals(other <Name>) bool {
    return v.field1 == other.field1 && v.field2 == other.field2
}

// ToDTO
func (v <Name>) ToDTO() <Name>DTO {
    return <Name>DTO{Field1: v.field1, Field2: v.field2}
}
```

### File: `<name>_test.go`

```go
package <package>

import "testing"

func TestNew<Name>(t *testing.T) {
    tests := []struct {
        name    string
        // params
        wantErr error
    }{
        {
            name:    "success",
            // valid params
            wantErr: nil,
        },
        {
            name:    "error - invalid",
            // invalid params
            wantErr: Err<Name>Invalid,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            v, err := New<Name>(/* params */)
            if err != tt.wantErr {
                t.Errorf("New<Name>() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if err == nil {
                // Verify fields
            }
        })
    }
}

func TestEquals(t *testing.T) {
    tests := []struct {
        name string
        v1   <Name>
        v2   <Name>
        want bool
    }{
        // Test cases: equal, different field1, different field2, both different
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.v1.Equals(tt.v2); got != tt.want {
                t.Errorf("<Name>.Equals() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test<Name>_ToDTO(t *testing.T) {
    // Create VO, call ToDTO(), verify fields match
}
```

## Single Field Example (primitive type)

### email.go

```go
package valueobject

import (
    "errors"
    "net/mail"
)

var (
    ErrEmailInvalid = errors.New("email is invalid")
    ErrEmailEmpty   = errors.New("email cannot be empty")
)

type Email string

func NewEmail(email string) (Email, error) {
    if email == "" {
        return "", ErrEmailEmpty
    }

    if _, err := mail.ParseAddress(email); err != nil {
        return "", ErrEmailInvalid
    }

    return Email(email), nil
}

func (v Email) Equals(other Email) bool {
    return v == other
}

func (v Email) ToDTO() string {
    return string(v)
}
```

### email_test.go

```go
package valueobject

import "testing"

func TestNewEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        want    Email
        wantErr error
    }{
        {
            name:    "success - valid email",
            email:   "user@example.com",
            want:    "user@example.com",
            wantErr: nil,
        },
        {
            name:    "success - email with subdomain",
            email:   "user@sub.example.com",
            want:    "user@sub.example.com",
            wantErr: nil,
        },
        {
            name:    "error - empty email",
            email:   "",
            want:    "",
            wantErr: ErrEmailEmpty,
        },
        {
            name:    "error - invalid format",
            email:   "not-an-email",
            want:    "",
            wantErr: ErrEmailInvalid,
        },
        {
            name:    "error - missing @",
            email:   "userexample.com",
            want:    "",
            wantErr: ErrEmailInvalid,
        },
        {
            name:    "error - missing domain",
            email:   "user@",
            want:    "",
            wantErr: ErrEmailInvalid,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NewEmail(tt.email)
            if err != tt.wantErr {
                t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("NewEmail() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestEquals(t *testing.T) {
    tests := []struct {
        name string
        v1   Email
        v2   Email
        want bool
    }{
        {
            name: "equal - same email",
            v1:   "user@example.com",
            v2:   "user@example.com",
            want: true,
        },
        {
            name: "equal - same value different case",
            v1:   "User@Example.com",
            v2:   "user@example.com",
            want: false,
        },
        {
            name: "not equal - different emails",
            v1:   "user1@example.com",
            v2:   "user2@example.com",
            want: false,
        },
        {
            name: "not equal - one empty",
            v1:   "user@example.com",
            v2:   "",
            want: false,
        },
        {
            name: "equal - both empty",
            v1:   "",
            v2:   "",
            want: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.v1.Equals(tt.v2); got != tt.want {
                t.Errorf("Email.Equals() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestEmail_ToDTO(t *testing.T) {
    tests := []struct {
        name  string
        email Email
        want  string
    }{
        {
            name:  "success",
            email: "user@example.com",
            want:  "user@example.com",
        },
        {
            name:  "empty email",
            email: "",
            want:  "",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.email.ToDTO(); got != tt.want {
                t.Errorf("Email.ToDTO() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Multi Field Example (struct)

### name.go

```go
package valueobject

import "errors"

var (
    ErrNameInvalid       = errors.New("name is invalid")
    ErrNameEmptyFirst    = errors.New("first name cannot be empty")
    ErrNameEmptyLast     = errors.New("last name cannot be empty")
)

type Name struct {
    first string
    last  string
}

type NameDTO struct {
    First string `json:"first"`
    Last  string `json:"last"`
}

func NewName(first, last string) (Name, error) {
    if first == "" {
        return Name{}, ErrNameEmptyFirst
    }

    if last == "" {
        return Name{}, ErrNameEmptyLast
    }

    return Name{first: first, last: last}, nil
}

func (v Name) First() string { return v.first }
func (v Name) Last() string  { return v.last }

func (v other Name) bool {
    return v.first == other.first && v.last == other.last
}

func (v Name) ToDTO() NameDTO {
    return NameDTO{First: v.first, Last: v.last}
}
```

### name_test.go

```go
package valueobject

import "testing"

func TestNewName(t *testing.T) {
    tests := []struct {
        name    string
        first   string
        last    string
        want    Name
        wantErr error
    }{
        {
            name:    "success",
            first:   "John",
            last:    "Doe",
            want:    Name{first: "John", last: "Doe"},
            wantErr: nil,
        },
        {
            name:    "success - unicode chars",
            first:   "José",
            last:    "García",
            want:    Name{first: "José", last: "García"},
            wantErr: nil,
        },
        {
            name:    "error - empty first name",
            first:   "",
            last:    "Doe",
            want:    Name{},
            wantErr: ErrNameEmptyFirst,
        },
        {
            name:    "error - empty last name",
            first:   "John",
            last:    "",
            want:    Name{},
            wantErr: ErrNameEmptyLast,
        },
        {
            name:    "error - both empty",
            first:   "",
            last:    "",
            want:    Name{},
            wantErr: ErrNameEmptyFirst,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NewName(tt.first, tt.last)
            if err != tt.wantErr {
                t.Errorf("NewName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("NewName() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestName_First(t *testing.T) {
    v := Name{first: "John", last: "Doe"}
    if got := v.First(); got != "John" {
        t.Errorf("Name.First() = %v, want %v", got, "John")
    }
}

func TestName_Last(t *testing.T) {
    v := Name{first: "John", last: "Doe"}
    if got := v.Last(); got != "Doe" {
        t.Errorf("Name.Last() = %v, want %v", got, "Doe")
    }
}

func TestEquals(t *testing.T) {
    tests := []struct {
        name string
        v1   Name
        v2   Name
        want bool
    }{
        {
            name: "equal - same values",
            v1:   Name{first: "John", last: "Doe"},
            v2:   Name{first: "John", last: "Doe"},
            want: true,
        },
        {
            name: "not equal - different first",
            v1:   Name{first: "John", last: "Doe"},
            v2:   Name{first: "Jane", last: "Doe"},
            want: false,
        },
        {
            name: "not equal - different last",
            v1:   Name{first: "John", last: "Doe"},
            v2:   Name{first: "John", last: "Smith"},
            want: false,
        },
        {
            name: "not equal - different both",
            v1:   Name{first: "John", last: "Doe"},
            v2:   Name{first: "Jane", last: "Smith"},
            want: false,
        },
        {
            name: "equal - both empty",
            v1:   Name{},
            v2:   Name{},
            want: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.v1.Equals(tt.v2); got != tt.want {
                t.Errorf("Name.Equals() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestName_ToDTO(t *testing.T) {
    tests := []struct {
        name string
        vo   Name
        want NameDTO
    }{
        {
            name: "success",
            vo:   Name{first: "John", last: "Doe"},
            want: NameDTO{First: "John", Last: "Doe"},
        },
        {
            name: "empty values",
            vo:   Name{},
            want: NameDTO{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.vo.ToDTO(); got != tt.want {
                t.Errorf("Name.ToDTO() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Workflow

1. **Read existing files**: If `<name>.go` exists, read it and its test
2. **Create or modify**: Generate/modify the value object following rules
3. **Generate tests**: Create table-driven tests covering 100%
4. **Verify**: Ensure tests pass with `go test -v -run Test<Name>`

## Scope

- **Create**: Generate new value object from scratch
- **Edit**: Modify existing value object (add fields, methods, validations)
- **Refactor**: Improve existing value object without changing behavior
- **Context**: ONLY `<name>.go` and `<name>_test.go`
- **External errors**: Developer must explicitly ask to fix

## Context Limits

- ONLY modify `<name>.go` and `<name>_test.go`
- Do NOT fix compilation errors in other files
- Do NOT modify callers or external code
- If other files have errors due to VO changes, developer must explicitly ask to fix them

## Naming Convention

- **Package**: Always in `domain` package (or developer-specified package)
- File: `<name>.go` and `<name>_test.go`
- Type: `<Name>` (PascalCase, no suffix)
- DTO: `<Name>DTO`
- Constructor: `New<Name>`
- Errors: `Err<Name><Reason>`
