# FSRouter: Generic File System Router for Go

FSRouter is a file system router library for Go, designed to work seamlessly with popular web frameworks like Fiber and Go-Chi. This library simplifies routing by allowing you to structure your routes in a directory hierarchy, making it easy to organize and manage your web application's endpoints.

## Features

- File System-Based Routing

- Compatible with Fiber and Go-Chi

- Dynamic Route Parameters: FSRouter uses the main NextJS conventions and allows you to define dynamic route parameters using placeholders like `[param]` and `[...param]` that gets mapped to the framework syntax (e.g. `:param` and `*` for Fiber)

## Getting Started

To start using FSRouter in your Go project, follow these simple steps:

```bash shell
go get -v -u github.com/aziis98/go-fsrouter
```

and import the package with

```go
import "github.com/aziis98/go-fsrouter"
```

## Usage (Fiber)

Create a file system router and then use it to load all the routes.

```go
app := fiber.New()

fsr := fsrouter.New("./pages", fsrouter.FiberPreset)
engine := fsrouter.NewTemplateCache(true)

routes, err := fsr.LoadRoutes()
if err != nil {
    log.Fatal(err)
}

for _, r := range routes {
    app.Get(r.Name, func(ctx *fiber.Ctx) error {
        params := map[string]any{}
        for _, p := range r.ParamNames {
            params[p.Name] = ctx.Params(p.Name)
        }

        return engine.Render(path.Join(fsr.Root, r.Name), ctx, params)
    })
}
```

## Example Directory Structure

Here's an example directory structure for your routes:

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

## Customization

You can customize the way route parameters are replaced using the `Preset` structure. The default configuration is set to `FiberPreset`, which replaces route parameters with `:paramName` for named parameters and `*` for wildcard parameters.

```go
customPreset := fsrouter.Preset{
    NamedParamReplacement: ":$1",
    WildcardReplacement:   "*",
}

fsr := fsrouter.New("./path/to/your/pages", customPreset)
```
