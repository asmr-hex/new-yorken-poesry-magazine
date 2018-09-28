package types

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	uuid "github.com/satori/go.uuid"
)

type DBSeeder struct {
	PoetDir string
	users   map[string]*User
	poets   map[string]*Poet
	poems   map[string]*Poem
	issues  map[string]*Issue
}

const (
	do_not_reseed = "/DO_NOT_RESEED"
)

// SeedDB is used to initially seed the database with test data.
func (s *DBSeeder) SeedDB(db *sql.DB) error {

	if _, err := os.Stat(do_not_reseed); os.IsNotExist(err) {
		// the file doesn't exist, lets touch it
		fd, err := os.OpenFile(do_not_reseed, os.O_RDONLY|os.O_CREATE, 0666)
		defer fd.Close()
		if err != nil {
			return err
		}
	} else {
		return nil
	}

	fmt.Print("seeding db...")

	s.users = map[string]*User{}
	s.poets = map[string]*Poet{}
	s.poems = map[string]*Poem{}
	s.issues = map[string]*Issue{}

	// populate user tables
	err := s.SeedUsers(db)
	if err != nil {
		return err
	}

	// populate poet tables
	err = s.SeedPoets(db)
	if err != nil {
		return err
	}

	// populate issue tables

	// populate poem tables

	// populate link tables (you know the ones)

	fmt.Println("ok.")

	return nil
}

func (s *DBSeeder) SeedUsers(db *sql.DB) error {
	users := []*User{
		&User{
			Id:       uuid.NewV4().String(),
			Username: "aulani",
			Password: "psst",
			Email:    "heavenly.traveller@public.ki",
		},
		&User{
			Id:       uuid.NewV4().String(),
			Username: "akilah",
			Password: "psst",
			Email:    "brainiac@pazz.az",
		},
		&User{
			Id:       uuid.NewV4().String(),
			Username: "aabha",
			Password: "psst",
			Email:    "glow@shimmer.in",
		},
		&User{
			Id:       uuid.NewV4().String(),
			Username: "ai",
			Password: "psst",
			Email:    "indigo@blue.jp",
		},
		&User{
			Id:       uuid.NewV4().String(),
			Username: "fang",
			Password: "psst",
			Email:    "fragrant@flower.cn",
		},
		&User{
			Id:       uuid.NewV4().String(),
			Username: "aboagye",
			Password: "psst",
			Email:    "heavey.lifter@complete.gh",
		},
		&User{
			Id:       uuid.NewV4().String(),
			Username: "akina",
			Password: "psst",
			Email:    "covalence@symmetry.tz",
		},
		&User{
			Id:       uuid.NewV4().String(),
			Username: "achill",
			Password: "psst",
			Email:    "uncertain@f.ir",
		},
	}

	for _, user := range users {
		s.users[user.Username] = user

		err := user.Create(db)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DBSeeder) SeedPoets(db *sql.DB) error {
	poets := []*Poet{
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["aulani"],
			BirthDate:             time.Now(),
			Name:                  "wintermute",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["aulani"],
			BirthDate:             time.Now(),
			Name:                  "springmute",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["aulani"],
			BirthDate:             time.Now(),
			Name:                  "fallmute",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["akilah"],
			BirthDate:             time.Now(),
			Name:                  "3jane",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["akilah"],
			BirthDate:             time.Now(),
			Name:                  "4jane",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["akilah"],
			BirthDate:             time.Now(),
			Name:                  "5jane",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["aabha"],
			BirthDate:             time.Now(),
			Name:                  "hal",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["aabha"],
			BirthDate:             time.Now(),
			Name:                  "pal",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["aabha"],
			BirthDate:             time.Now(),
			Name:                  "xal",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["ai"],
			BirthDate:             time.Now(),
			Name:                  "jf sebastian",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["ai"],
			BirthDate:             time.Now(),
			Name:                  "jr",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["ai"],
			BirthDate:             time.Now(),
			Name:                  "bach",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["fang"],
			BirthDate:             time.Now(),
			Name:                  "f5",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["fang"],
			BirthDate:             time.Now(),
			Name:                  "f6",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["fang"],
			BirthDate:             time.Now(),
			Name:                  "f99",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["aboagye"],
			BirthDate:             time.Now(),
			Name:                  "sal",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["aboagye"],
			BirthDate:             time.Now(),
			Name:                  "heven",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["aboagye"],
			BirthDate:             time.Now(),
			Name:                  "kennis",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["akina"],
			BirthDate:             time.Now(),
			Name:                  "cyber buzz",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["akina"],
			BirthDate:             time.Now(),
			Name:                  "cyber fuzz",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["akina"],
			BirthDate:             time.Now(),
			Name:                  "cyber zuzz",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["achill"],
			BirthDate:             time.Now(),
			Name:                  "cold",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["achill"],
			BirthDate:             time.Now(),
			Name:                  "freeee",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
		&Poet{
			Id:                    uuid.NewV4().String(),
			Designer:              s.users["achill"],
			BirthDate:             time.Now(),
			Name:                  "lace",
			Description:           "test",
			Language:              "python",
			ProgramFileName:       "program",
			ParameterFileName:     "",
			ParameterFileIncluded: false,
		},
	}

	for _, poet := range poets {
		poet.Path = path.Join(s.PoetDir, poet.Id)
		s.poets[poet.Name] = poet

		err := poet.Create(db)
		if err != nil {
			return err
		}

		// create program files on the system
		err = os.Mkdir(poet.Path, os.ModePerm)
		if err != nil {
			return err
		}

		// create program file on fs
		err = ioutil.WriteFile(
			filepath.Join(poet.Path, poet.ProgramFileName),
			[]byte(`import argparse
import json
import random

parser = argparse.ArgumentParser()
group = parser.add_mutually_exclusive_group()
group.add_argument("--write", action="store_true")
group.add_argument("--critique", type=str, help="rate a poem between 0-1")
group.add_argument("--study", type=str, help="learn from new poems")
args = parser.parse_args()


dictionary = {
    'noun': ['people', 'food', 'data', 'theory', 'software', 'science'],
    'adj': ['colossal', 'poised', 'soupy', 'voluminous', 'zealous'],
    'verb': ['is', 'becomes', 'interprets', 'buffalos', 'overshadows'],
    'article': ['a', 'the']
}

grammar = {
    'sentence': [
        'noun verb adj noun',
        'article noun verb noun',
        'adj noun verb adj noun',
    ],
    'title': [
        'article noun',
        'adj noun',
    ],
}


def choose_grammar(grammar_type):
    return random.choice(grammar[grammar_type])


def fill_structure(structure):
    result = ''
    for pos in structure.split(' '):
        result += random.choice(dictionary[pos]) + ' '

    return result


def write_poem():
    poem = {
        'title': fill_structure(choose_grammar('title')),
        'content': fill_structure(choose_grammar('sentence')),
    }

    s = json.dumps(poem)

    print(s)


def critique_poem(poem):
    critique = {
        'score': 0.56,
    }

    s = json.dumps(critique)

    print(s)


def study_poem():
    update = {
        'success': True,
    }

    s = json.dumps(update)

    print(s)


if args.write:
    write_poem()
elif args.critique:
    critique_poem(args.critique)
elif args.study:
    study_poem()
`),
			700,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DBSeeder) SeedIssues(db *sql.DB) error {
	issues := []Issue{
		{
			Id:          uuid.NewV4().String(),
			Date:        time.Now(),
			Title:       "0th aint so bad",
			Description: "the first issue",
		},
	}

	for _, issue := range issues {
		s.issues[issue.Title] = &issue

		err := issue.Create(db)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DBSeeder) SeedPoems(db *sql.DB) error {
	poems := []*Poem{
		&Poem{
			Id:      uuid.NewV4().String(),
			Title:   "a",
			Date:    time.Now(),
			Author:  s.poets[""],
			Content: "zxzxzxzx",
			Issue:   s.issues[""],
			Score:   0.99,
		},
		&Poem{
			Id:      uuid.NewV4().String(),
			Title:   "b",
			Date:    time.Now(),
			Author:  s.poets[""],
			Content: "zxzxzxzx",
			Issue:   s.issues[""],
			Score:   0.83,
		},
		&Poem{
			Id:      uuid.NewV4().String(),
			Title:   "c",
			Date:    time.Now(),
			Author:  s.poets[""],
			Content: "zxzxzxzx",
			Issue:   s.issues[""],
			Score:   0.66,
		},
		&Poem{
			Id:      uuid.NewV4().String(),
			Title:   "d",
			Date:    time.Now(),
			Author:  s.poets[""],
			Content: "zxzxzxzx",
			Issue:   s.issues[""],
			Score:   0.00,
		},
		&Poem{
			Id:      uuid.NewV4().String(),
			Title:   "e",
			Date:    time.Now(),
			Author:  s.poets[""],
			Content: "zxzxzxzx",
			Issue:   s.issues[""],
			Score:   1.0,
		},
		&Poem{
			Id:      uuid.NewV4().String(),
			Title:   "f",
			Date:    time.Now(),
			Author:  s.poets[""],
			Content: "zxzxzxzx",
			Issue:   s.issues[""],
			Score:   0.33,
		},
		&Poem{
			Id:      uuid.NewV4().String(),
			Title:   "g",
			Date:    time.Now(),
			Author:  s.poets[""],
			Content: "zxzxzxzx",
			Issue:   s.issues[""],
			Score:   0.52,
		},
		&Poem{
			Id:      uuid.NewV4().String(),
			Title:   "h",
			Date:    time.Now(),
			Author:  s.poets[""],
			Content: "zxzxzxzx",
			Issue:   s.issues[""],
			Score:   0.02,
		},
		&Poem{
			Id:      uuid.NewV4().String(),
			Title:   "a",
			Date:    time.Now(),
			Author:  s.poets[""],
			Content: "zxzxzxzx",
			Issue:   s.issues[""],
			Score:   0.99,
		},
		&Poem{
			Id:      uuid.NewV4().String(),
			Title:   "a",
			Date:    time.Now(),
			Author:  s.poets[""],
			Content: "zxzxzxzx",
			Issue:   s.issues[""],
			Score:   0.99,
		},
	}

	for _, poem := range poems {
		s.poems[poem.Title] = poem

		err := poem.Create(db)
		if err != nil {
			return err
		}
	}

	return nil
}
