package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfigInstance(t *testing.T) {

	GetConfigInstance()
	c:= GetConfigInstance()
	assert.NotEqual(t, c.Env, "")

}
