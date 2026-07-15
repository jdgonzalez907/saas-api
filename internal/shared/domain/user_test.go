package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"jdgonzalez907/saas-api/internal/shared/domain"
)

func TestNewUser(t *testing.T) {
	user := domain.NewUser(
		1,
		"John Doe",
		"57",
		"1234567890",
		nil,
	)

	userDTO := user.ToDTO()

	assert.Equal(t, int64(1), user.ID())
	assert.Equal(t, "John Doe", user.FullName())
	assert.Equal(t, "57", user.PhoneCountryCode())
	assert.Equal(t, "1234567890", user.PhoneNumber())
	assert.Nil(t, user.Email())

	assert.Equal(t, int64(1), userDTO.ID)
	assert.Equal(t, "John Doe", userDTO.FullName)
	assert.Equal(t, "57", userDTO.PhoneCountryCode)
	assert.Equal(t, "1234567890", userDTO.PhoneNumber)
	assert.Nil(t, userDTO.Email)
}
