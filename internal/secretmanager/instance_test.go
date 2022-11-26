package secretmanager_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bavatech/envloader/internal/secretmanager"
	"github.com/stretchr/testify/assert"
)

func Test_NewInstance(t *testing.T) {
	t.Run("should return a valid instance", func(t *testing.T) {
		sm, err := secretmanager.NewInstance("us-east-1", aws.Config{})
		if err != nil {
			t.Fatal(err)
		}
		assert.NotEmpty(t, sm)
		assert.Empty(t, err)
	})
}
