package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetComments(t *testing.T) {
	resp, err := GetComments()
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

func TestGetPosts(t *testing.T) {
	resp, err := GetPosts()
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

func TestGetPost(t *testing.T) {
	resp, err := GetPost(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}
