package fsrouter

import (
	"html/template"
	"io"
	"os"
)

type templateCache struct {
	// reload forces templates to be reloaded each time
	reload bool

	// templates is a map of loaded templates
	templates map[string]*template.Template
}

func (tc templateCache) loadTemplate(filePath string) (*template.Template, error) {
	if tmpl, ok := tc.templates[filePath]; !tc.reload && ok {
		return tmpl, nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("").Parse(string(content))
	if err != nil {
		return nil, err
	}

	tc.templates[filePath] = tmpl

	return tmpl, nil
}

// Render the template for "view" into "w" with "data"
func (tc templateCache) Render(view string, w io.Writer, data any) error {
	tmpl, err := tc.loadTemplate(view)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}

type TemplateEngine interface {
	Render(view string, w io.Writer, data any) error
}

func NewTemplateCache(reload bool) TemplateEngine {
	return templateCache{
		reload:    reload,
		templates: map[string]*template.Template{},
	}
}
