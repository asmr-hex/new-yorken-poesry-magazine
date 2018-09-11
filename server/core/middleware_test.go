package core

import (
	"os"
	"testing"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/utils"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func MainTest(m *testing.M) {
	retcode := m.Run()

	os.Exit(retcode)
}

// move this test to utils
func TestIsValidUUIDV4_WithValidUUID(t *testing.T) {
	id := uuid.NewV4().String()

	assert.True(t, utils.IsValidUUIDV4(id))
}

func TestIsValidUUIDV4_WithInValidUUID(t *testing.T) {
	id := uuid.NewV4().String()

	// say no to sql injections
	assert.False(t, utils.IsValidUUIDV4(id+" OR 1=1 --see u in h3ll"))
}

func TestValidateParams_NoParams(t *testing.T) {
	params := map[string]string{}

	assert.NoError(t, ValidateParams(params))
}

func TestValidateParams_UnknownParam(t *testing.T) {
	params := map[string]string{"Haunt": "Me"}

	assert.Error(t, ValidateParams(params))
}

func TestValidateParams_IdParam_ValidUUID(t *testing.T) {
	params := map[string]string{API_ID_PATH_PARAM: uuid.NewV4().String()}

	assert.NoError(t, ValidateParams(params))
}

func TestValidateParams_IdParam_InValidUUID(t *testing.T) {
	params := map[string]string{API_ID_PATH_PARAM: "/bin/sh"}

	assert.Error(t, ValidateParams(params))
}
