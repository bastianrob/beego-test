package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_APIError(t *testing.T) {
	err := &APIError{Code: 400, Message: "Bad Request"}
	assert.Error(t, err, "API error is derived from error object")
	assert.Equal(t, "Bad Request", err.Error())
}
