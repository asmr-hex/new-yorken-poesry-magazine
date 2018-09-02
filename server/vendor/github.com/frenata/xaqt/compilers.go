package xaqt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// ExecutionDetails specifies how to execute certain code.
type ExecutionDetails struct {
	Compiler           string `json:"compiler"`
	SourceFile         string `json:"sourceFile"`
	OptionalExecutable string `json:"optionalExecutable"`
	CompilerFlags      string `json:"compilerFlags"`
	Disabled           bool   `json:"disabled"`
}

// CompositionDetails specifies how to write code in a given language
type CompositionDetails struct {
	Boilerplate   string `json:"boilerplate"`
	CommentPrefix string `json:"commentPrefix"`
}

// CompilerDetails contains everything XAQT knows about handling a certain language
type CompilerDetails struct {
	ExecutionDetails
	CompositionDetails
}

// Compilers maps language names to the details of how to execute code in that language.
type Compilers map[string]CompilerDetails

// availableLanguages returns a list of currently supported languages.
func (c Compilers) availableLanguages() map[string]CompositionDetails {
	fmt.Printf("Received languages request...")
	langs := make(map[string]CompositionDetails)

	// make a list of currently supported languages
	for k, v := range c {
		if !v.Disabled {
			langs[k] = v.CompositionDetails
		}
	}

	log.Printf("currently supporting %d of %d known languages\n", len(langs), len(c))
	return langs
}

// GetCompilers reads a compilers map from a file.
func GetCompilers(filenames ...string) Compilers {
	var (
		filename    string
		compilerMap = make(Compilers, 0)
		err         error
	)

	if len(filenames) > 0 {
		// just use the first filename
		filename = filenames[0]
	} else {
		// return default compilers if no filename is provided
		return DEFAULT_COMPILERS
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Fatal: failed to read language file: %s", err)
	}

	err = json.Unmarshal(bytes, &compilerMap)
	if err != nil {
		log.Fatalf("Fatal: failed to parse JSON: %s", err)
	}

	return compilerMap
}

// write a compilers map to a json file.
//
func (c Compilers) ExportToJSON(filename string) error {
	var (
		err error
	)

	bytes, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	// write json to file
	err = ioutil.WriteFile(filename, bytes, 0777)
	if err != nil {
		return err
	}

	return nil
}
