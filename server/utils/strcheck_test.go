package utils

import "github.com/stretchr/testify/suite"

type StrCheckTestSuite struct {
	suite.Suite
}

func (s *StrCheckTestSuite) TestValidateUsername_English() {
	username := "wintermute"

	err := ValidateUsername(username)
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_English_OneChar() {
	username := "c"

	err := ValidateUsername(username)
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_English_NoChar() {
	username := ""

	err := ValidateUsername(username)
	s.Error(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_English_Underscore() {
	username := "n_n"

	err := ValidateUsername(username)
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_English_UnderscoreRepeated() {
	username := "n__n"

	err := ValidateUsername(username)
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_English_SymbolEnd() {
	err := ValidateUsername("n_")
	s.Error(err)

	err = ValidateUsername("n.")
	s.Error(err)

	err = ValidateUsername("n-")
	s.Error(err)

	err = ValidateUsername("n ")
	s.Error(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_English_SymbolBegining() {
	err := ValidateUsername("_n")
	s.Error(err)

	err = ValidateUsername(".n")
	s.Error(err)

	err = ValidateUsername("-n")
	s.Error(err)

	err = ValidateUsername(" n")
	s.Error(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_English_Hyphen() {
	username := "e-e"

	err := ValidateUsername(username)
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_English_HyphenRepeated() {
	username := "e--e"

	err := ValidateUsername(username)
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_English_Period() {
	username := "x.x"

	err := ValidateUsername(username)
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_English_PeriodRepeated() {
	username := "x..x"

	err := ValidateUsername(username)
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_English_space() {
	username := "h h"

	err := ValidateUsername(username)
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_English_spaceRepeated() {
	username := "h  h"

	err := ValidateUsername(username)
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_Japanese() {
	err := ValidateUsername("暗い眠り")
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_Chinese() {
	err := ValidateUsername("黑暗的睡眠")
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_Cyrillic() {
	err := ValidateUsername("темный сон")
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateUsername_Arabic() {
	err := ValidateUsername("الظلام")
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateEmail() {
	err := ValidateEmail("geomancer666@benevolent.magik")
	s.NoError(err)
}

func (s *StrCheckTestSuite) TestValidateEmail_InvalidEmail() {
	err := ValidateEmail("ur@somewherere bad")
	s.Error(err)
}
