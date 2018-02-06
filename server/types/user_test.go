package types

import (
	"os"
	"testing"
)

func MainTest(m *testing.M) {
	returnCode := m.Run()

	os.Exit(returnCode)
}

func TestNothing(t *testing.T) {

}
