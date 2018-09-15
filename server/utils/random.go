package utils

import (
	"fmt"
	"math/rand"
)

func NRandomUniqueInts(n, min, max int) ([]int, error) {
	var (
		l []int
	)

	// ensure that this is even possible!
	if (max - min) < n {
		return nil, fmt.Errorf(
			"can't pick %d unique ints from [%d, %d]",
			n, min, max,
		)
	}

	// get a randomly permuted array in normalized range
	p := rand.Perm(max - min)

	// choose first n
	l = p[0:n]

	// if min is non-zero, then we need to re-add it
	if min != 0 {
		for idx, v := range l {
			l[idx] = v + min
		}
	}

	return l, nil
}
