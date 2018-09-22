package types

import (
	"io/ioutil"
	"path/filepath"
)

type Code struct {
	Author   string `json:"author"`
	Code     string `json:"code"`
	Language string `json:"language"`
	FileName string `json:"filename"`
	Path     string `json:"-"`
}

func (c *Code) Read() error {
	// get absolute path to filename
	fn := filepath.Join(c.Path, c.FileName)

	// read sourcefile into string
	bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}

	// set the code string
	c.Code = string(bytes)

	return nil
}
