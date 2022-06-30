package xfer

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestMixin_UnmarshalStep(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/step-input.yaml")
	require.NoError(t, err)

	var action Action
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)
	assert.Equal(t, "install", action.Name)
	require.Len(t, action.Steps, 1)

	step := action.Steps[0]
	assert.Equal(t, "File Transfer", step.Description)
	assert.Equal(t, "This will likely work", step.Command)
}
