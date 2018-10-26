package app

import (
	"bytes"
	"context"
	"gitlab.petrichor.io/sandstone/v1/errors/stackerr"
	"io"
	"sort"
	"text/template"
)

func init() {
	Register(Command{
		Keyword:     "help",
		Description: "Available commands and their options/flags",
		Fn:          help,
	})
}

func help(ctx context.Context) (err error) {
	var rawTemplate = "\033[1m\033[38;5;26m{{.Name}} v{{.Version}}\033[0m" + `

` + "\033[1mDESCRIPTION:\033[0m" + `
` + "\033[38;5;249m {{.Description}}\033[0m" + `

` + "\033[1mCOMMANDS:\033[0m" + `
{{range .Commands -}}` +
		"\033[38;5;249m" + ` {{printf "%s %s" .Keyword .Usage | printf "%-40s"}} {{.Description}}` + "\033[0m\n" +
		"\033[38;5;249m" + ` {{printf "" | printf "%-45s"}} {{.UsageFlags}}` + "\033[0m\n" + `
{{end}}`

	t, _ := template.New("help display").Parse(rawTemplate)

	// For some reason, if we write directly to an output interface, it messes up.
	buf := &bytes.Buffer{}

	sort.Sort(a.commands)

	var vars = struct {
		Commands    Commands
		Description string
		Name        string
		Version     string
	}{
		Commands:    a.commands,
		Description: a.About.Description,
		Name:        a.About.Name,
		Version:     a.About.Version,
	}

	if err := t.Execute(buf, vars); err != nil {
		return stackerr.Wrap(err, "Could not execute template")
	}

	io.Copy(a.Log.NormalOut, buf)

	return nil
}
