package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/fatih/color"
	"github.com/senforsce/telemetry"
	"github.com/senforsce/tndr"
	"github.com/senforsce/tndr/cmd/t1/fmtcmd"
	"github.com/senforsce/tndr/cmd/t1/generatecmd"
	"github.com/senforsce/tndr/cmd/t1/infocmd"
	"github.com/senforsce/tndr/cmd/t1/lspcmd"
)

func main() {
	code := run(os.Stdin, os.Stdout, os.Stderr, os.Args)
	if code != 0 {
		os.Exit(code)
	}
}

const usageText = `usage: tndr <command> [<args>...]

tndr - build HTML UIs with Go

See docs at https://senforsce.com/docs/tndr

commands:
  generate   Generates Go code from tndr files
  fmt        Formats tndr files
  lsp        Starts a language server for tndr files
  info       Displays information about the tndr environment
  version    Prints the version
`

func run(stdin io.Reader, stdout, stderr io.Writer, args []string) (code int) {
	if len(args) < 2 {
		_, _ = fmt.Fprint(stderr, usageText)
		return 64 // EX_USAGE
	}
	switch args[1] {
	case "info":
		return infoCmd(stdout, stderr, args[2:])
	case "generate":
		return generateCmd(stdout, stderr, args[2:])
	case "fmt":
		return fmtCmd(stdin, stdout, stderr, args[2:])
	case "lsp":
		return lspCmd(stdin, stdout, stderr, args[2:])
	case "version", "--version":
		_, _ = fmt.Fprintln(stdout, tndr.Version())
		return 0
	case "help", "-help", "--help", "-h":
		_, _ = fmt.Fprint(stdout, usageText)
		return 0
	}
	_, _ = fmt.Fprint(stderr, usageText)
	return 64 // EX_USAGE
}

const infoUsageText = `usage: tndr info [<args>...]

Displays information about the tndr environment.

Args:
  -json
    Output information in JSON format to stdout. (default false)
  -v
    Set log verbosity level to "debug". (default "info")
  -log-level
    Set log verbosity level. (default "info", options: "debug", "info", "warn", "error")
  -help
    Print help and exit.
`

func infoCmd(stdout, stderr io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("diagnose", flag.ExitOnError)
	jsonFlag := cmd.Bool("json", false, "")
	//verboseFlag := cmd.Bool("v", false, "")
	//logLevelFlag := cmd.String("log-level", "info", "")
	helpFlag := cmd.Bool("help", false, "")
	err := cmd.Parse(args)
	if err != nil {
		_, _ = fmt.Fprint(stderr, infoUsageText)
		return 64 // EX_USAGE
	}
	if *helpFlag {
		_, _ = fmt.Fprint(stdout, infoUsageText)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())

	logOptions := telemetry.SenforsceLoggerOptions{
		Context:        ctx,
		LogLevel:       "debug",
		LoggerName:     "tndr",
		ServiceName:    "info",
		ServiceVersion: "0.1.0",
		Verbose:        true,
	}
	tndrLogger, tndrLoggerProvider := telemetry.NewLogger(logOptions)

	defer tndrLoggerProvider.Shutdown(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		_, _ = fmt.Fprintln(stderr, "Stopping...")
		cancel()
	}()

	err = infocmd.Run(ctx, tndrLogger, stdout, infocmd.Arguments{
		JSON: *jsonFlag,
	})
	if err != nil {
		_, _ = color.New(color.FgRed).Fprint(stderr, "(✗) ")
		_, _ = fmt.Fprintln(stderr, "Command failed: "+err.Error())
		return 1
	}
	return 0
}

func generateCmd(stdout, stderr io.Writer, args []string) (code int) {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		_, _ = fmt.Fprintln(stderr, "Stopping...")
		cancel()
	}()

	err := generatecmd.Run(ctx, stdout, stderr, args)
	if err != nil {
		_, _ = color.New(color.FgRed).Fprint(stderr, "(✗) ")
		_, _ = fmt.Fprintln(stderr, "Command failed: "+err.Error())
		exitCode := 1
		if e, ok := err.(ErrorCode); ok {
			exitCode = e.Code()
		}
		return exitCode
	}
	return 0
}

type ErrorCode interface {
	error
	Code() int
}

const fmtUsageText = `usage: tndr fmt [<args> ...]

Format all files in directory:

  tndr fmt .

Format stdin to stdout:

  tndr fmt < header.tndr

Format file or directory to stdout:

  tndr fmt -stdout FILE

Args:
  -stdout
    Prints to stdout instead of in-place format
  -stdin-filepath
    Provides the formatter with filepath context when using -stdout.
    Required for organising imports.
  -v
    Set log verbosity level to "debug". (default "info")
  -log-level
    Set log verbosity level. (default "info", options: "debug", "info", "warn", "error")
  -w
    Number of workers to use when formatting code. (default runtime.NumCPUs).
  -prettier-command
    Set the command to use for formatting HTML, CSS, and JS blocks. Default is "prettier --stdin-filepath $TEMPL_PRETTIER_FILENAME".
  -prettier-required
    Set to true to return an error the prettier command is not available. Default is false.
  -fail
    Fails with exit code 1 if files are changed. (e.g. in CI)
  -help
    Print help and exit.
`

func fmtCmd(stdin io.Reader, stdout, stderr io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("fmt", flag.ExitOnError)
	helpFlag := cmd.Bool("help", false, "")
	workerCountFlag := cmd.Int("w", runtime.NumCPU(), "")
	//verboseFlag := cmd.Bool("v", false, "")
	//logLevelFlag := cmd.String("log-level", "info", "")
	failIfChanged := cmd.Bool("fail", false, "")
	prettierCommand := cmd.String("prettier-command", "", "")
	prettierRequired := cmd.Bool("prettier-required", false, "")
	stdoutFlag := cmd.Bool("stdout", false, "")
	stdinFilepath := cmd.String("stdin-filepath", "", "")
	err := cmd.Parse(args)
	if err != nil {
		_, _ = fmt.Fprint(stderr, fmtUsageText)
		return 64 // EX_USAGE
	}
	if *helpFlag {
		_, _ = fmt.Fprint(stdout, fmtUsageText)
		return
	}

	tndrContext := context.Background()
	logOptions := telemetry.SenforsceLoggerOptions{
		Context:        tndrContext,
		LogLevel:       "debug",
		LoggerName:     "tndr",
		ServiceName:    "fmt",
		ServiceVersion: "0.1.0",
		Verbose:        true,
	}
	tndrLogger, tndrLoggerProvider := telemetry.NewLogger(logOptions)

	defer tndrLoggerProvider.Shutdown(tndrContext)

	err = fmtcmd.Run(tndrLogger, stdin, stdout, fmtcmd.Arguments{
		ToStdout:         *stdoutFlag,
		Files:            cmd.Args(),
		WorkerCount:      *workerCountFlag,
		StdinFilepath:    *stdinFilepath,
		FailIfChanged:    *failIfChanged,
		PrettierCommand:  *prettierCommand,
		PrettierRequired: *prettierRequired,
	})
	if err != nil {
		return 1
	}
	return 0
}

const lspUsageText = `usage: tndr lsp [<args> ...]

Starts a language server for tndr.

Args:
  -log string
    The file to log tndr LSP output to, or leave empty to disable logging.
  -goplsLog string
    The file to log gopls output, or leave empty to disable logging.
  -goplsRPCTrace
    Set gopls to log input and output messages.
  -gopls-remote
    Specify remote gopls instance to connect to.
  -help
    Print help and exit.
  -pprof
    Enable pprof web server (default address is localhost:9999)
  -http string
    Enable http debug server by setting a listen address (e.g. localhost:7474)
  -no-preload
    Disable preloading of tndr files on server startup and use custom GOPACKAGESDRIVER for lazy loading (useful for large monorepos). GOPACKAGESDRIVER environment variable must be set.
`

func lspCmd(stdin io.Reader, stdout, stderr io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("lsp", flag.ExitOnError)
	logFlag := cmd.String("log", "", "")
	goplsLog := cmd.String("goplsLog", "", "")
	goplsRPCTrace := cmd.Bool("goplsRPCTrace", false, "")
	goplsRemote := cmd.String("gopls-remote", "", "")
	helpFlag := cmd.Bool("help", false, "")
	pprofFlag := cmd.Bool("pprof", false, "")
	httpDebugFlag := cmd.String("http", "", "")
	noPreloadFlag := cmd.Bool("no-preload", false, "")
	err := cmd.Parse(args)
	if err != nil {
		_, _ = fmt.Fprint(stderr, lspUsageText)
		return 64 // EX_USAGE
	}
	if *helpFlag {
		_, _ = fmt.Fprint(stdout, lspUsageText)
		return
	}

	err = lspcmd.Run(stdin, stdout, stderr, lspcmd.Arguments{
		Log:           *logFlag,
		GoplsLog:      *goplsLog,
		GoplsRPCTrace: *goplsRPCTrace,
		GoplsRemote:   *goplsRemote,
		PPROF:         *pprofFlag,
		HTTPDebug:     *httpDebugFlag,
		NoPreload:     *noPreloadFlag && os.Getenv("GOPACKAGESDRIVER") != "",
	})
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err.Error())
		return 1
	}
	return 0
}
