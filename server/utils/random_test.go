package utils

import "github.com/stretchr/testify/suite"

type RandomTestSuite struct {
	suite.Suite
}

func (s *RandomTestSuite) TestNRandomUniqueInts() {
	for i := 0; i < 100; i++ {
		ns, err := NRandomUniqueInts(5, 10, 20)
		s.NoError(err)
		for _, n := range ns {
			s.True(n >= 10 && n <= 20)
		}
	}
}

func (s *RandomTestSuite) TestNRandomUniqueInts_OutOfBounds() {
	for i := 0; i < 100; i++ {
		_, err := NRandomUniqueInts(5, 10, 13)
		s.Error(err)
	}
}
