---
name: go-entity
description: Use when the user wants to create, modify, or refactor Go entities. Generates mutable entities with ID, explicit mutation methods (named by intention), getters, Equals (by ID), and ToDTO for infra only. Creates or updates functional file and table-driven tests for 100% coverage.
---

# Go Entity Generator

Generates mutable entities in Go following these rules:

## Structure Rules

| Element | Implementation |
|---------|----------------|
| ID | Required, type specified by developer |
| Fields | Private, accessed via getters |
| Mutations | Explicit methods named by business intention |
| DTO | For infra only, never for calculations |

## Calculated Fields

- Fields computed from other fields or business logic
- NOT included in DTO
- Only accessible via getter method

```go
type Order struct {
    id         int64
    items      []OrderItem
    discount   Money
}

// Calculated field - not in DTO
func (o *Order) Total() Money {
    sum := Money{}
    for _, item := range o.items {
        sum = sum.Add(item.Subtotal())
    }
    return sum.Subtract(o.discount)
}

func (o *Order) ToDTO() OrderDTO {
    return OrderDTO{ID: o.id, Items: o.items, Discount: o.discount}
    // Total NOT in DTO
}
```

## Required Elements

### File: `<name>.go`

```go
package <package>

import (
    "errors"
    // ID type import if needed (e.g., "github.com/google/uuid")
)

// Errors
var (
    Err<Entity>Invalid       = errors.New("<entity> is invalid")
    Err<Entity>IDRequired    = errors.New("<entity> ID is required")
    // Add specific errors as needed
)

// Entity
type <Entity> struct {
    id   <IDType>
    field1 type1
    field2 type2
    // Add fields as needed
}

// DTO - FOR INFRA ONLY, NEVER USE FOR CALCULATIONS
type <Entity>DTO struct {
    ID      <IDType> `json:"id"`
    Field1  type1    `json:"field1"`
    Field2  type2    `json:"field2"`
}

// Constructor - returns pointer and error
func New<Entity>(id <IDType>, field1 type1, field2 type2) (*<Entity>, error) {
    if id == zero {
        return nil, Err<Entity>IDRequired
    }
    // Additional validation
    return &<Entity>{id: id, field1: field1, field2: field2}, nil
}

// Getters
func (e *<Entity>) ID() <IDType> { return e.id }
func (e *<Entity>) Field1() type1 { return e.field1 }
func (e *<Entity>) Field2() type2 { return e.field2 }

// Equals - compares by ID only
func (e *<Entity>) Equals(other *<Entity>) bool {
    if other == nil {
        return false
    }
    return e.id == other.id
}

// Mutation methods - named by business intention
func (e *<Entity>) UpdateField1(newValue type1) error {
    // Validation if needed
    e.field1 = newValue
    return nil
}

// ToDTO - FOR INFRA ONLY
func (e *<Entity>) ToDTO() <Entity>DTO {
    return <Entity>DTO{
        ID:     e.id,
        Field1: e.field1,
        Field2: e.field2,
    }
}
```

### File: `<name>_test.go`

```go
package <package>

import "testing"

func TestNew<Entity>(t *testing.T) {
    tests := []struct {
        name    string
        id      <IDType>
        field1  type1
        field2  type2
        wantErr error
    }{
        {
            name:    "success",
            id:      // valid ID
            field1:  // valid value
            field2:  // valid value
            wantErr: nil,
        },
        {
            name:    "error - missing ID",
            id:      zero,
            field1:  // any value
            field2:  // any value
            wantErr: Err<Entity>IDRequired,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e, err := New<Entity>(tt.id, tt.field1, tt.field2)
            if err != tt.wantErr {
                t.Errorf("New<Entity>() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if err == nil {
                if e.ID() != tt.id {
                    t.Errorf("New<Entity>().ID() = %v, want %v", e.ID(), tt.id)
                }
                // Verify other fields
            }
        })
    }
}

func Test<Entity>_Equals(t *testing.T) {
    tests := []struct {
        name string
        e1   *<Entity>
        e2   *<Entity>
        want bool
    }{
        {
            name: "equal - same ID",
            e1:   // entity with ID "1"
            e2:   // entity with ID "1"
            want: true,
        },
        {
            name: "not equal - different ID",
            e1:   // entity with ID "1"
            e2:   // entity with ID "2"
            want: false,
        },
        {
            name: "not equal - nil other",
            e1:   // entity
            e2:   nil,
            want: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.e1.Equals(tt.e2); got != tt.want {
                t.Errorf("<Entity>.Equals() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test<Entity>_<MutationMethod>(t *testing.T) {
    tests := []struct {
        name      string
        entity    *<Entity>
        // mutation params
        wantErr   error
        wantField // expected value
    }{
        {
            name:    "success",
            entity:  // create entity
            // valid params
            wantErr: nil,
            wantField: // expected value
        },
        // Add error cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.entity.<MutationMethod>(/* params */)
            if err != tt.wantErr {
                t.Errorf("<Entity>.<MutationMethod>() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if err == nil {
                if got := tt.entity.Field1(); got != tt.wantField {
                    t.Errorf("<Entity>.Field1() = %v, want %v", got, tt.wantField)
                }
            }
        })
    }
}

func Test<Entity>_ToDTO(t *testing.T) {
    // Create entity, call ToDTO(), verify fields match
    // Note: ToDTO is for infra only, do not use for business logic
}
```

## Mutation Methods

Name by business intention, NOT by field name:

| Bad | Good |
|-----|------|
| `SetName(string)` | `Rename(string)` |
| `SetStatus(string)` | `Activate()` / `Deactivate()` |
| `SetEmail(string)` | `ChangeEmail(string)` |
| `SetAmount(float64)` | `Deposit(float64)` / `Withdraw(float64)` |

Each mutation method should:
- Validate input if needed
- Return error if validation fails
- Mutate the entity
- Return nil on success

## Workflow

1. **Read existing files**: If `<name>.go` exists, read it and its test
2. **Create or modify**: Generate/modify the entity following rules
3. **Generate tests**: Create table-driven tests covering 100%
4. **Verify**: Ensure tests pass with `go test -v -run Test<Entity>`

## Scope

- **Create**: Generate new entity from scratch
- **Edit**: Modify existing entity (add fields, methods, validations)
- **Refactor**: Improve existing entity without changing behavior
- **Context**: ONLY `<name>.go` and `<name>_test.go`
- **External errors**: Developer must explicitly ask to fix

## Audit Fields

Ask the developer which audit fields the entity needs:

| Field | Description | Mutator |
|-------|-------------|---------|
| `createdAt` | Auto-set on creation | None (readonly) |
| `updatedAt` | Auto-set on modification | Auto-updated by mutators |
| `createdBy` | User who created | Set on creation |
| `updatedBy` | User who last modified | Auto-updated by mutators |

**Example with audit fields:**

```go
type User struct {
    id        int64
    name      Name
    email     Email
    createdAt time.Time
    updatedAt time.Time
    createdBy int64
    updatedBy int64
}

// Mutators update audit fields automatically
func (u *User) ChangeEmail(email Email, modifiedBy int64) error {
    if email == "" {
        return ErrUserEmailInvalid
    }
    u.email = email
    u.updatedAt = time.Now()
    u.updatedBy = modifiedBy
    return nil
}
```

**Questions to ask:**
- Does this entity need `createdAt`?
- Does this entity need `updatedAt`?
- Does this entity need `createdBy`?
- Does this entity need `updatedBy`?

## Context Limits

- ONLY modify `<name>.go` and `<name>_test.go`
- Do NOT fix compilation errors in other files
- Do NOT modify callers or external code
- If other files have errors due to entity changes, developer must explicitly ask to fix them

## DTO Rules

- DTO is for INFRA ONLY (serialization, API responses, etc.)
- NEVER use DTO for business calculations
- NEVER pass DTO to domain methods
- DTO is a data representation, not a domain object
- **VO fields in DTO**: Call `ToDTO()` on VOs, never raw values

```go
// Entity with VOs
type User struct {
    id    int64
    name  Name   // VO
    email Email  // VO
}

type UserDTO struct {
    ID    int64   `json:"id"`
    Name  NameDTO `json:"name"`   // NameDTO from Name.ToDTO()
    Email string  `json:"email"`  // string from Email.ToDTO()
}

// ToDTO must call ToDTO() on VOs
func (u *User) ToDTO() UserDTO {
    return UserDTO{
        ID:    u.id,
        Name:  u.name.ToDTO(),    // Required: call VO's ToDTO()
        Email: u.email.ToDTO(),   // Required: call VO's ToDTO()
    }
}
```

## Naming Convention

- **Package**: Always in `domain` package (or developer-specified package)
- File: `<name>.go` and `<name>_test.go`
- Type: `<Entity>` (PascalCase, no suffix)
- DTO: `<Entity>DTO`
- Constructor: `New<Entity>`
- Errors: `Err<Entity><Reason>`
- Getters: `FieldName()` (PascalCase)
- Mutators: Business intention (`Rename`, `Activate`, etc.)

## Complete Example: User Entity

### user.go

```go
package domain

import (
    "errors"
    "time"
)

var (
    ErrUserInvalid       = errors.New("user is invalid")
    ErrUserIDRequired    = errors.New("user ID is required")
    ErrUserEmailInvalid  = errors.New("user email is invalid")
    ErrUserNameRequired  = errors.New("user name is required")
)

type User struct {
    id        int64
    name      Name
    email     Email
    birthDate time.Time
    createdAt time.Time
}

type UserDTO struct {
    ID        int64     `json:"id"`
    Name      NameDTO   `json:"name"`
    Email     string    `json:"email"`
    BirthDate time.Time `json:"birthDate"`
    CreatedAt time.Time `json:"createdAt"`
}

func NewUser(id int64, name Name, email Email, birthDate time.Time) (*User, error) {
    if id == 0 {
        return nil, ErrUserIDRequired
    }

    if name == (Name{}) {
        return nil, ErrUserNameRequired
    }

    if email == "" {
        return nil, ErrUserEmailInvalid
    }

    return &User{
        id:        id,
        name:      name,
        email:     email,
        birthDate: birthDate,
        createdAt: time.Now(),
    }, nil
}

// Getters
func (u *User) ID() int64         { return u.id }
func (u *User) Name() Name        { return u.name }
func (u *User) Email() Email      { return u.email }
func (u *User) BirthDate() time.Time { return u.birthDate }
func (u *User) CreatedAt() time.Time { return u.createdAt }

// Equals - compares by ID only
func (u *User) Equals(other *User) bool {
    if other == nil {
        return false
    }
    return u.id == other.id
}

// Mutators - named by business intention

// ChangeEmail updates the user's email
func (u *User) ChangeEmail(email Email) error {
    if email == "" {
        return ErrUserEmailInvalid
    }
    u.email = email
    return nil
}

// UpdatePersonalInfo updates name and birth date
func (u *User) UpdatePersonalInfo(name Name, birthDate time.Time) error {
    if name == (Name{}) {
        return ErrUserNameRequired
    }
    u.name = name
    u.birthDate = birthDate
    return nil
}

// ToDTO - FOR INFRA ONLY
func (u *User) ToDTO() UserDTO {
    return UserDTO{
        ID:        u.id,
        Name:      u.name.ToDTO(),
        Email:     u.email.ToDTO(),
        BirthDate: u.birthDate,
        CreatedAt: u.createdAt,
    }
}
```

### user_test.go

```go
package domain

import (
    "testing"
    "time"
)

func TestNewUser(t *testing.T) {
    validName, _ := NewName("John", "Doe")
    validEmail, _ := NewEmail("john@example.com")
    validBirthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

    tests := []struct {
        name      string
        id        int64
        nameArg   Name
        emailArg  Email
        birthDate time.Time
        wantErr   error
    }{
        {
            name:      "success",
            id:        1,
            nameArg:   validName,
            emailArg:  validEmail,
            birthDate: validBirthDate,
            wantErr:   nil,
        },
        {
            name:    "error - zero ID",
            id:      0,
            nameArg: validName,
            emailArg: validEmail,
            birthDate: validBirthDate,
            wantErr: ErrUserIDRequired,
        },
        {
            name:    "error - empty name",
            id:      1,
            nameArg: Name{},
            emailArg: validEmail,
            birthDate: validBirthDate,
            wantErr: ErrUserNameRequired,
        },
        {
            name:    "error - empty email",
            id:      1,
            nameArg: validName,
            emailArg: "",
            birthDate: validBirthDate,
            wantErr: ErrUserEmailInvalid,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            u, err := NewUser(tt.id, tt.nameArg, tt.emailArg, tt.birthDate)
            if err != tt.wantErr {
                t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if err == nil {
                if u.ID() != tt.id {
                    t.Errorf("NewUser().ID() = %v, want %v", u.ID(), tt.id)
                }
                if u.Name() != tt.nameArg {
                    t.Errorf("NewUser().Name() = %v, want %v", u.Name(), tt.nameArg)
                }
                if u.Email() != tt.emailArg {
                    t.Errorf("NewUser().Email() = %v, want %v", u.Email(), tt.emailArg)
                }
                if u.BirthDate() != tt.birthDate {
                    t.Errorf("NewUser().BirthDate() = %v, want %v", u.BirthDate(), tt.birthDate)
                }
                if u.CreatedAt().IsZero() {
                    t.Errorf("NewUser().CreatedAt() should not be zero")
                }
            }
        })
    }
}

func TestUser_Equals(t *testing.T) {
    name, _ := NewName("John", "Doe")
    email, _ := NewEmail("john@example.com")
    birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

    tests := []struct {
        name string
        e1   *User
        e2   *User
        want bool
    }{
        {
            name: "equal - same ID",
            e1:   mustNewUser(t, 1, name, email, birthDate),
            e2:   mustNewUser(t, 1, name, email, birthDate),
            want: true,
        },
        {
            name: "not equal - different ID",
            e1:   mustNewUser(t, 1, name, email, birthDate),
            e2:   mustNewUser(t, 2, name, email, birthDate),
            want: false,
        },
        {
            name: "not equal - nil other",
            e1:   mustNewUser(t, 1, name, email, birthDate),
            e2:   nil,
            want: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.e1.Equals(tt.e2); got != tt.want {
                t.Errorf("User.Equals() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestUser_ChangeEmail(t *testing.T) {
    name, _ := NewName("John", "Doe")
    oldEmail, _ := NewEmail("john@example.com")
    newEmail, _ := NewEmail("john.doe@example.com")
    birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

    tests := []struct {
        name    string
        entity  *User
        email   Email
        wantErr error
    }{
        {
            name:    "success",
            entity:  mustNewUser(t, 1, name, oldEmail, birthDate),
            email:   newEmail,
            wantErr: nil,
        },
        {
            name:    "error - empty email",
            entity:  mustNewUser(t, 1, name, oldEmail, birthDate),
            email:   "",
            wantErr: ErrUserEmailInvalid,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.entity.ChangeEmail(tt.email)
            if err != tt.wantErr {
                t.Errorf("User.ChangeEmail() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if err == nil {
                if got := tt.entity.Email(); got != tt.email {
                    t.Errorf("User.Email() = %v, want %v", got, tt.email)
                }
            }
        })
    }
}

func TestUser_UpdatePersonalInfo(t *testing.T) {
    oldName, _ := NewName("John", "Doe")
    newName, _ := NewName("Jane", "Smith")
    email, _ := NewEmail("john@example.com")
    oldBirthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
    newBirthDate := time.Date(1985, 6, 15, 0, 0, 0, 0, time.UTC)

    tests := []struct {
        name      string
        entity    *User
        nameArg   Name
        birthDate time.Time
        wantErr   error
    }{
        {
            name:      "success",
            entity:    mustNewUser(t, 1, oldName, email, oldBirthDate),
            nameArg:   newName,
            birthDate: newBirthDate,
            wantErr:   nil,
        },
        {
            name:      "error - empty name",
            entity:    mustNewUser(t, 1, oldName, email, oldBirthDate),
            nameArg:   Name{},
            birthDate: newBirthDate,
            wantErr:   ErrUserNameRequired,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.entity.UpdatePersonalInfo(tt.nameArg, tt.birthDate)
            if err != tt.wantErr {
                t.Errorf("User.UpdatePersonalInfo() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if err == nil {
                if got := tt.entity.Name(); got != tt.nameArg {
                    t.Errorf("User.Name() = %v, want %v", got, tt.nameArg)
                }
                if got := tt.entity.BirthDate(); got != tt.birthDate {
                    t.Errorf("User.BirthDate() = %v, want %v", got, tt.birthDate)
                }
            }
        })
    }
}

func TestUser_ToDTO(t *testing.T) {
    name, _ := NewName("John", "Doe")
    email, _ := NewEmail("john@example.com")
    birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
    u := mustNewUser(t, 1, name, email, birthDate)

    dto := u.ToDTO()

    if dto.ID != 1 {
        t.Errorf("User.ToDTO().ID = %v, want %v", dto.ID, 1)
    }
    if dto.Name != name.ToDTO() {
        t.Errorf("User.ToDTO().Name = %v, want %v", dto.Name, name.ToDTO())
    }
    if dto.Email != string(email) {
        t.Errorf("User.ToDTO().Email = %v, want %v", dto.Email, string(email))
    }
    if dto.BirthDate != birthDate {
        t.Errorf("User.ToDTO().BirthDate = %v, want %v", dto.BirthDate, birthDate)
    }
}

// Helper function for tests
func mustNewUser(t *testing.T, id int64, name Name, email Email, birthDate time.Time) *User {
    t.Helper()
    u, err := NewUser(id, name, email, birthDate)
    if err != nil {
        t.Fatalf("mustNewUser() error = %v", err)
    }
    return u
}
```
