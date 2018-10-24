package templates

import (
	"bytes"
	"html/template"
)

func ComposeVerificationEmail(url, username string) (string, error) {
	var (
		err error
	)

	data := struct {
		Url      string
		Username string
	}{
		Url:      url,
		Username: username,
	}

	t, err := template.ParseFiles("templates/verification_template.html")
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
