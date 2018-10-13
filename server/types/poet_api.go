package types

import (
	"encoding/json"
	"fmt"
)

// the poet api is the contract which  all poet designers
// abide by when specifying how their poets respond to
// input arguments. e.g.
//
// $ python poet.py --generate
// {content:"the contents of a generated poem"}
//
// $ python poet.py --critique "some superb poetic verses"
// {score: 0.99}
//
// $ python poet.py --study
// {success: true}
//
const (
	POET_API_GENERATE_FLAG = "--write"
	POET_API_CRITIQUE_FLAG = "--critique"
	POET_API_UPDATE_FLAG   = "--study"
)

// RawPoem is the output of a poet's poem generator
type RawPoem struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// RawCritique is the output of a poet's critical analysis of a poem
type RawCritique struct {
	Score float64 `json:"score"`
}

// RawUpdate is the output of a poet after self improving
type RawUpdate struct {
	Success bool `json:"success"`
}

func PoetAPIGenerateArgs() []string {
	return []string{POET_API_GENERATE_FLAG}
}

func (p *Poet) ParseRawPoem(results []string) (*RawPoem, error) {
	var (
		result *RawPoem = &RawPoem{}
		err    error
	)

	if len(results) == 0 {
		return nil, fmt.Errorf("%s failed to write a poem", p.Name)
	}

	err = json.Unmarshal([]byte(results[0]), result)
	if err != nil {
		fmt.Println(results)
		return nil, err
	}

	return result, nil
}

func PoetAPICritiqueArgs(poems ...string) []string {
	for idx, poem := range poems {
		poems[idx] = fmt.Sprintf(
			`%s "%s"`,
			POET_API_CRITIQUE_FLAG,
			poem,
		)
	}

	return poems
}

func (p *Poet) ParseRawCritique(results []string) (*RawCritique, error) {
	var (
		result *RawCritique = &RawCritique{}
		err    error
	)

	if len(results) == 0 {
		return nil, fmt.Errorf("%s failed to critique a poem", p.Name)
	}

	err = json.Unmarshal([]byte(results[0]), result)
	if err != nil {
		return nil, err
	}

	// ensure that the score is within the range [0, 1]
	if !(result.Score >= 0 && result.Score <= 1) {
		return nil, fmt.Errorf(
			"invalid score (%f) given by poet %s",
			result.Score,
			p.Name,
		)
	}

	return result, nil
}

func PoetAPIUpdateArgs(poems ...string) []string {
	for idx, poem := range poems {
		poems[idx] = fmt.Sprintf(
			`%s "%s"`,
			POET_API_UPDATE_FLAG,
			poem,
		)
	}
	return poems
}

func (p *Poet) ParseRawUpdate(results []string) (*RawUpdate, error) {
	var (
		result *RawUpdate = &RawUpdate{}
		err    error
	)

	if len(results) == 0 {
		return nil, fmt.Errorf("%s failed to study", p.Name)
	}

	err = json.Unmarshal([]byte(results[0]), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
