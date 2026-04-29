# Template generation

To generate Go code from `*.t1` files, use the `t1` command line tool.

```
t1 generate
```

The `t1 generate` recurses into subdirectories and generates Go code for each `*.t1` file it finds.

The command outputs warnings, and a summary of updates.

```
(!) void element <input> should not have child content [ from=12:2 to=12:7 ]
(✓) Complete [ updates=62 duration=144.677334ms ]
```

## Advanced options

The `t1 generate` command has a `--help` option that prints advanced options.

These include the ability to generate code for a single file and to choose the number of parallel workers that `t1 generate` uses to create Go files.

By default `t1 generate` uses the number of CPUs that your machine has installed.

```
t1 generate --help
```

```
usage: t1 generate [<args>...]

Generates Go code from tndr files.

Args:
  -path <path>
    Generates code for all files in path. (default .)
  -f <file>
    Optionally generates code for a single file, e.g. -f header.t1
  -stdout
    Prints to stdout instead of writing generated files to the filesystem.
    Only applicable when -f is used.
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
  -notify-proxy
    If present, the command will issue a reload event to the proxy 127.0.0.1:7331, or use proxyport and proxybind to specify a different address.
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

Examples:

  Generate code for all files in the current directory and subdirectories:

    t1 generate

  Generate code for a single file:

    t1 generate -f header.t1

  Watch the current directory and subdirectories for changes and regenerate code:

    t1 generate -watch
```

:::tip
The `t1 generate --watch` option watches files for changes and runs t1 generate when required.

However, the code generated in this mode is not optimised for production use.
:::
