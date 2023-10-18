package fsrouter

import (
	"html/template"
	"io"
	"os"
	"sync"
)

type TemplateEngine interface {
	Render(w io.Writer, view string, data any) error
}

type templateCache struct {
	// mutex
	mu sync.RWMutex

	// reload forces templates to be reloaded each time
	reload bool

	// templates is a map of loaded templates
	templates map[string]*template.Template
}

func NewTemplateCache(reload bool) TemplateEngine {
	return &templateCache{
		reload:    reload,
		templates: map[string]*template.Template{},
	}
}

func (tc *templateCache) retriveFromCache(filePath string) (*template.Template, bool) {
	tc.mu.RLock()
	defer tc.mu.Unlock()

	if tc.reload {
		return nil, false
	}

	if tmpl, ok := tc.templates[filePath]; ok {
		return tmpl, true
	}

	return nil, false
}

func (tc *templateCache) loadTemplate(filePath string) (*template.Template, error) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

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
func (tc *templateCache) Render(w io.Writer, view string, data any) error {
	tmpl, ok := tc.retriveFromCache(view)
	if ok {
		return tmpl.Execute(w, data)
	}

	tmpl, err := tc.loadTemplate(view)
	if err == nil {
		return tmpl.Execute(w, data)
	}

	return err
}
