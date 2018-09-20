package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestMain(m *testing.M) {
	retCode := m.Run()

	os.Exit(retCode)
}

func TestStrCheckSuite(t *testing.T) {
	suite.Run(t, &StrCheckTestSuite{})
}

func TestRandomSuite(t *testing.T) {
	suite.Run(t, &RandomTestSuite{})
}
