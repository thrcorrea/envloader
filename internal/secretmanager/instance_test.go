package secretmanager_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thrcorrea/envloader/internal/secretmanager"
)

func Test_NewInstance(t *testing.T) {
	t.Run("should return a valid instance", func(t *testing.T) {
		sm, err := secretmanager.NewInstance("us-east-1")
		if err != nil {
			t.Fatal(err)
		}
		assert.NotEmpty(t, sm)
		assert.Empty(t, err)
	})
}
