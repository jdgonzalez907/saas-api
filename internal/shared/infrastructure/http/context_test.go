package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrUnauthenticated(t *testing.T) {
	assert.Equal(t, "unauthenticated user", ErrUnauthenticated.Error())
}
