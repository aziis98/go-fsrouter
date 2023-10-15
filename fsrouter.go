package fsrouter

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/mattn/go-zglob"
)

var (
	paramRegex = regexp.MustCompile(`\[((?:\.\.\.)?([a-zA-Z0-9]+))\]`)

	paramRegexSingle = regexp.MustCompile(`\[([a-zA-Z0-9]+)\]`)
	paramRegexNested = regexp.MustCompile(`\[\.\.\.([a-zA-Z0-9]+)\]`)
)

type Preset struct {
	NamedParamReplacement string
	WildcardReplacement   string
}

var FiberPreset = Preset{
	NamedParamReplacement: ":$1",
	WildcardReplacement:   "*",
}

var ChiPreset = Preset{
	NamedParamReplacement: "{$1}",
	WildcardReplacement:   "*",
}

type FSRouter struct {
	Root           string
	IncludePattern string

	Preset Preset
}

func New(rootDir string, preset Preset) *FSRouter {
	return &FSRouter{
		Root:           rootDir,
		IncludePattern: "**/*.html",

		Preset: preset,
	}
}

type RouteParam struct {
	Name   string
	Nested bool
}

type Route struct {
	Name       string
	ParamNames []RouteParam

	Path string
}

func (r Route) Realize(params map[string]string) string {
	path := r.Name

	for _, pn := range r.ParamNames {
		if pn.Nested {
			path = strings.ReplaceAll(path, fmt.Sprintf("[...%s]", pn.Name), params[pn.Name])
		} else {
			path = strings.ReplaceAll(path, fmt.Sprintf("[%s]", pn.Name), params[pn.Name])
		}
	}

	return path
}

func (r Route) ExtractMap(valueFn func(param string) string) map[string]string {
	m := map[string]string{}
	for _, pn := range r.ParamNames {
		m[pn.Name] = valueFn(pn.Name)
	}

	return m
}

func (fsr FSRouter) LoadRoutes() ([]Route, error) {
	matches, err := zglob.Glob(filepath.Join(fsr.Root, fsr.IncludePattern))
	if err != nil {
		return nil, err
	}

	sort.Strings(matches)

	routes := make([]Route, len(matches))
	for i, path := range matches {
		relPath, err := filepath.Rel(fsr.Root, path)
		if err != nil {
			return nil, err
		}

		routes[i] = fsr.parseRoute(relPath)
	}

	return routes, nil
}

func (fsr FSRouter) parseRoute(path string) Route {
	paramNames := []RouteParam{}

	allParamMatches := paramRegex.FindAllStringSubmatch(path, -1)

	for _, paramMatch := range allParamMatches {
		paramNames = append(paramNames, RouteParam{
			Name:   paramMatch[2],
			Nested: strings.HasPrefix(paramMatch[1], "..."),
		})
	}

	// apply route replacement using the current preset
	route := "/" + path
	route = paramRegexSingle.ReplaceAllString(route, fsr.Preset.NamedParamReplacement)
	route = paramRegexNested.ReplaceAllString(route, fsr.Preset.WildcardReplacement)

	// apply common replacements to the url
	if strings.HasSuffix(route, "index.html") {
		route = strings.TrimSuffix(route, "index.html")
	} else if strings.HasSuffix(route, ".html") {
		route = strings.TrimSuffix(route, ".html")
	}

	return Route{route, paramNames, path}
}
