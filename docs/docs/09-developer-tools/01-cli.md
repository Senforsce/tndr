# CLI

`tndr` provides a command line interface. Most users will only need to run the `t1 generate` command to generate Go code from `*.t1` files.

```
usage: t1 <command> [<args>...]

tndr - build HTML UIs with Go

See docs at https://templ.guide

commands:
  generate   Generates Go code from t1 files
  fmt        Formats t1 files
  lsp        Starts a language server for t1 files
  info       Displays information about the t1 environment
  version    Prints the version
```

## Generating Go code from tndr files

The `t1 generate` command generates Go code from `*.t1` files in the current directory tree.

The command provides additional options:

```
usage: t1 generate [<args>...]

Generates Go code from t1 files.

Args:
  -path <path>
    Generates code for all files in path. (default .)
  -f <file>
    Optionally generates code for a single file, e.g. -f header.t1
  -source-map-visualisations
    Set to true to generate HTML files to visualise the tndr code and its corresponding Go code.
  -include-version
    Set to false to skip inclusion of the tndr version in the generated code. (default true)
  -include-timestamp
    Set to true to include the current time in the generated code.
  -watch
    Set to true to watch the path for changes and regenerate code.
  -cmd <cmd>
    Set the command to run after generating code.
  -proxy
    Set the URL to proxy after generating code and executing the command.
  -proxyport
    The port the proxy will listen on. (default 7331)
  -proxybind
    The address the proxy will listen on. (default 127.0.0.1)
  -w
    Number of workers to use when generating code. (default runtime.NumCPUs)
  -lazy
    Only generate .go files if the source .t1 file is newer.	
  -pprof
    Port to run the pprof server on.
  -keep-orphaned-files
    Keeps orphaned generated tndr files. (default false)
  -v
    Set log verbosity level to "debug". (default "info")
  -log-level
    Set log verbosity level. (default "info", options: "debug", "info", "warn", "error")
  -help
    Print help and exit.
```

For example, to generate code for a single file:

```
t1 generate -f header.templ
```

## Formatting tndr files

The `t1 fmt` command formats template files. You can use this command in different ways:

1. Format all template files in the current directory and subdirectories:

```
t1 fmt .
```

2. Format input from stdin and output to stdout:

```
t1 fmt
```

Alternatively, you can run `fmt` in CI to ensure that invalidly formatted templatess do not pass CI. This will cause the command
to exit with unix error-code `1` if any templates needed to be modified.

```
t1 fmt -fail .
```

If `prettierd`, `prettier` or `npx` is found in your `PATH`, `t1 fmt` will use prettier to format `script` and `style` elements in files.

## Language Server for IDE integration

`t1 lsp` provides a Language Server Protocol (LSP) implementation to support IDE integrations.

This command isn't intended to be used directly by users, but is used by IDE integrations such as the VSCode extension and by Neovim support.

By default, `t1 lsp` starts its own instance of gopls. However, gopls supports a [shared daemon mode](https://github.com/golang/tools/blob/master/gopls/doc/daemon.md), allowing multiple clients to connect to a single, long-lived instance. You can enable this mode using the `-gopls-remote` flag, which will either connect to an existing shared gopls instance or create one if none is running. This can improve performance and reduce resource usage.

A number of additional options are provided to enable runtime logging and profiling tools.

```
  -goplsLog string
        The file to log gopls output, or leave empty to disable logging.
  -goplsRPCTrace
        Set gopls to log input and output messages.
  -gopls-remote
        Specify remote gopls instance to connect to.
  -help
        Print help and exit.
  -http string
        Enable http debug server by setting a listen address (e.g. localhost:7474)
  -log string
        The file to log t1 LSP output to, or leave empty to disable logging.
  -pprof
        Enable pprof web server (default address is localhost:9999)
```
