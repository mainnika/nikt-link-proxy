package loader

//go:generate npm run build-loader

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"strings"
	"sync"
)

var tmplsLock = sync.RWMutex{}
var tmpls = make(map[string]*template.Template)

func register(f fs.FS, templatePath string) {

	buf := &strings.Builder{}

	tmplContent, err := f.Open(templatePath)
	if err != nil {
		panic(fmt.Errorf("cannot open template content, name:%s, %w", templatePath, err))
	}
	size, err := io.Copy(buf, tmplContent)
	if err != nil {
		panic(fmt.Errorf("cannot read template content, name:%s, %w", templatePath, err))
	}

	tmplsLock.Lock()
	defer tmplsLock.Unlock()

	tmpls[templatePath], err = template.New(templatePath).Parse(buf.String())
	if err != nil {
		panic(fmt.Errorf("cannot parse template, name:%s, %w", templatePath, err))
	}

	fmt.Printf("Found template: %s, size:%d\n", templatePath, size)
}

func LookupBytes(name []byte) (tmpl *template.Template, err error) {

	tmplsLock.RLock()
	defer tmplsLock.RUnlock()

	tmpl, hasTmpl := tmpls[string(name)]
	if !hasTmpl {
		err = fmt.Errorf("cannot find template, name:%s", name)
	}

	return
}

func Lookup(name string) (tmpl *template.Template, err error) {

	tmplsLock.RLock()
	defer tmplsLock.RUnlock()

	tmpl, hasTmpl := tmpls[name]
	if !hasTmpl {
		err = fmt.Errorf("cannot find template, name:%s", name)
	}

	return
}
