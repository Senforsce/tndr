package main

const usageText = `usage: t1 <command> [<args>...]

Tndr - build Semantic Web With OWL & HTMX UIs with Go

See docs at https://senforsce.com/t1/guide and https://senforsce.com/tndr/guide (coming soon)
Based off https://templ.guide

commands:
  generate   Generates Go code from t1 or t1 files
  fmt        Formats t1 or t1 files
  lsp        Starts a language server for t1 or t1 files
  version    Prints the version
`

const generateUsageText = `usage: t1 generate [<args>...]

Generates Go code from t1 or t1 files.

Args:
  -path <path>
    Generates code for all files in path. (default .)
  -f <file>
    Optionally generates code for a single file, e.g. -f header.t1
  -sourceMapVisualisations
    Set to true to generate HTML files to visualise the t1 or t1 code and its corresponding Go code.
  -include-version
    Set to false to skip inclusion of the t1 version in the generated code. (default true)
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
  -w
    Number of workers to use when generating code. (default runtime.NumCPUs)
  -pprof
    Port to run the pprof server on.
  -keep-orphaned-files
    Keeps orphaned generated t1 files. (default false)
  -v
    Set log verbosity level to "debug". (default "info")
  -log-level
    Set log verbosity level. (default "info", options: "debug", "info", "warn", "error")
  -help
    Print help and exit.

Examples:

  Generate code for all files in the current directory and subdirectories:

    t1 generate

  Generate code for a single file:

    t1 generate -f header.t1

  Watch the current directory and subdirectories for changes and regenerate code:

    t1 generate -watch
`

const fmtUsageText = `usage: t1 fmt [<args> ...]

Format all files in directory:

  t1 fmt .

Format stdin to stdout:

  t1 fmt < header.t1

Format file or directory to stdout:

  t1 fmt -stdout FILE

Args:
  -stdout
    Prints to stdout instead of in-place format
  -v
    Set log verbosity level to "debug". (default "info")
  -log-level
    Set log verbosity level. (default "info", options: "debug", "info", "warn", "error")
  -w
    Number of workers to use when formatting code. (default runtime.NumCPUs).
  -help
    Print help and exit.
`

const lspUsageText = `usage: t1 lsp [<args> ...]

Starts a language server for t1.

Args:
  -log string
    The file to log t1 LSP output to, or leave empty to disable logging.
  -goplsLog string
    The file to log gopls output, or leave empty to disable logging.
  -goplsRPCTrace
    Set gopls to log input and output messages.
  -help
    Print help and exit.
  -pprof
    Enable pprof web server (default address is localhost:9999)
  -http string
    Enable http debug server by setting a listen address (e.g. localhost:7474)
`
