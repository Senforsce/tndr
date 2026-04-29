module github.com/senforsce/tndr

go 1.25.0

replace github.com/senforsce/tndr/ => ./

replace github.com/senforsce/level0/ => ../level0/

replace github.com/senforsce/o/ => ../o/

replace github.com/senforsce/toolbelt/ => ../toolbelt/

replace github.com/senforsce/parsers/ => ../parsers/

replace github.com/senforsce/t1parsers/ => ../t1parsers/

replace github.com/senforsce/generator/ => ../generator/



require (
	github.com/a-h/parse v0.0.0-20240121214402-3caf7543159a
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/cli/browser v1.3.0
	github.com/fatih/color v1.16.0
	github.com/fsnotify/fsnotify v1.7.0
	github.com/google/go-cmp v0.6.0
	github.com/natefinch/atomic v1.0.1
	golang.org/x/mod v0.17.0
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d
)

require (
	github.com/rs/cors v1.11.1
	github.com/stretchr/testify v1.10.0
	golang.org/x/sync v0.11.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/andybalholm/brotli v1.1.0
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/net v0.25.0
	golang.org/x/sys v0.20.0 // indirect
)
