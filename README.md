# FSRouter

FSRouter is a simple **file system router library** for Go, designed to easily integrate with most http router libraries. This uses the "NextJS" convention to retrive routes directly as a directory-file hierarchy. Example directory structure for your routes:

```bash shell
pages/
├── dashboard/[...all].html # => /dashboard/*   (useful for SPAs)
└── user/
    ├── [name].html         # => /user/:name
    └── [name]/
        └── posts/
            └── [post].html # => /user/:name/posts/:post
```

In this structure, `[name]` and `[post]` are dynamic route parameters.

## Features

- **File System Routing**

    FSRouter uses the main NextJS conventions and allows you to define dynamic route parameters using placeholders like `[param]` and `[...param]` that gets mapped to the framework syntax (e.g. `:param` and `*` for Fiber)

- **Simple format**

    This library just reads all `**/*.html` (can be changed using `FSRouter.IncludePattern`) files in a directory and parses route names using the NextJS convention into a `[]Route` slice.

    ```go
    type RouteParam struct {
        Name   string
        Nested bool
    }

    type Route struct {
        Name       string
        ParamNames []RouteParam

        Path string
    }
    ```

- **Presets** 

    There are already presets for [Fiber](https://github.com/gofiber/fiber) and [Chi](https://github.com/go-chi/chi)

## Usage

To start using FSRouter in your Go project:

```bash shell
go get -v -u github.com/aziis98/go-fsrouter
```

and import the package with

```go
import "github.com/aziis98/go-fsrouter"
```

### With Chi

Create an `FSRouter` and then use it to load all the routes.

```go
// ExtractChiParams retrives all params needed by this route from the current context
func ExtractChiParams(r *http.Request, route fsrouter.Route) map[string]string {
    return route.ExtractMap(func(key string) string { return chi.URLParam(r, key) })
}

```

```go
r := chi.NewRouter()

fsr := fsrouter.New("./pages", fsrouter.ChiPreset)
engine := fsrouter.NewTemplateCache(true)

routes, err := fsr.LoadRoutes()
if err != nil {
    log.Fatal(err)
}

for _, route := range routes {
    route := route

    r.Get(route.Name, func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        if err := engine.Render(w, 
            path.Join(fsr.Root, route.Path), 
            ExtractChiParams(r, route),
        ); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    })
}
```

### With Fiber

Create an `FSRouter` and then use it to load all the routes.

```go
// ExtractFiberParams retrives all params needed by this route from the current context
func ExtractFiberParams(c *fiber.Ctx, route fsrouter.Route) map[string]string {
    return route.ExtractMap(func(key string) string { return c.Params(key) })
}
```

```go
app := fiber.New()

fsr := fsrouter.New("./pages", fsrouter.FiberPreset)
engine := fsrouter.NewTemplateCache(true)

routes, err := fsr.LoadRoutes()
if err != nil {
    log.Fatal(err)
}

for _, route := range routes {
    route := route

    app.Get(r.Name, func(c *fiber.Ctx) error {
        c.Type(path.Ext(route.Path))
        return engine.Render(ctx, 
            path.Join(fsr.Root, route.Path), 
            ExtractFiberParams(c, route),
        )
    })
}
```

### Custom Preset

You can customize the way route parameters are replaced using the `Preset` structure, for example `FiberPreset` uses the following

```go
fsr := fsrouter.New("./path/to/your/pages", fsrouter.Preset{
    NamedParamReplacement: ":$1",
    WildcardReplacement:   "*",
})
```
