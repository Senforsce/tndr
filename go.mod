module github.com/senforsce/tndr

go 1.25.0

replace github.com/senforsce/level0/ => ../level0/

replace github.com/senforsce/o/ => ../o/

replace github.com/senforsce/toolbelt/ => ../toolbelt/

replace github.com/senforsce/parse/ => ../parse/

require (
	github.com/PuerkitoBio/goquery v1.8.1
	github.com/a-h/htmlformat v0.0.0-20231108124658-5bd994fe268e
	github.com/a-h/parse v0.0.0-20240121214402-3caf7543159a
	github.com/a-h/protocol v0.0.0-20230224160810-b4eec67c1c22
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/cli/browser v1.3.0
	github.com/deiu/rdf2go v0.0.0-20230522072227-c39eb296d51c
	github.com/fatih/color v1.16.0
	github.com/fsnotify/fsnotify v1.7.0
	github.com/google/go-cmp v0.6.0
	github.com/natefinch/atomic v1.0.1
	go.lsp.dev/jsonrpc2 v0.10.0
	go.lsp.dev/uri v0.3.0
	go.uber.org/zap v1.27.0
	golang.org/x/mod v0.17.0
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d
)

require github.com/a-h/lexical v0.0.53

require (
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	github.com/samber/lo v1.51.0 // indirect
	github.com/samber/oops v1.19.0 // indirect
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)

require (
	github.com/andybalholm/brotli v1.1.0
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/deiu/gon3 v0.0.0-20230411081920-f0f8f879f597 // indirect
	github.com/linkeddata/gojsonld v0.0.0-20170418210642-4f5db6791326 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rychipman/easylex v0.0.0-20160129204217-49ee7767142f // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/segmentio/encoding v0.4.0 // indirect
	go.lsp.dev/pkg v0.0.0-20210717090340-384b27a52fb2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
)
