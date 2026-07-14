package controllers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"jdgonzalez907/saas-api/internal/posts/infrastructure/controllers"
)

func TestNewRouter(t *testing.T) {
	t.Run("empty params", func(t *testing.T) {
		params := controllers.RouterParams{}
		r := controllers.NewRouter(params)
		assert.NotNil(t, r)
	})

	t.Run("full params", func(t *testing.T) {
		params := controllers.RouterParams{
			CreatePost:         controllers.NewCreatePostController(nil),
			FindPostByID:       controllers.NewFindPostByIDController(nil),
			UpdatePost:         controllers.NewUpdatePostController(nil),
			DeletePost:         controllers.NewDeletePostController(nil),
			FindPostsPaginated: controllers.NewFindPostsPaginatedController(nil),
		}
		r := controllers.NewRouter(params)
		assert.NotNil(t, r)
	})
}
