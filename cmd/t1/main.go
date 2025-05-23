package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/senforsce/o/tr"
	t1 "github.com/senforsce/tndr"
	"github.com/senforsce/tndr/cmd/t1/fmtcmd"
	"github.com/senforsce/tndr/cmd/t1/generatecmd"
	"github.com/senforsce/tndr/cmd/t1/lspcmd"
)

func init() {
	env := os.Getenv("T1_ENV")
	if env == "" {
		env = "development"
	}

	// godotenv.Load(".env." + env + ".local")
	// if env != "test" {
	// 	godotenv.Load(".env.local")
	// }
	// godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env
	_, ok := os.LookupEnv("O8ROOT")

	if !ok {
		tr.Ln("O8ROOT is not present, please make sure O8ROOT is exported before running the command")
	} else {
		//tr.Ln("OntologyRootFolder: %s\n", driver)
	}

	_, ok = os.LookupEnv("O8_META_NAMESPACE")

	if !ok {
		tr.Ln("O8NAMESPACE is not present, please make sure O8NAMESPACE is exported before running the command")
	} else {
		//tr.Ln("OntologyRootNamespace: %s\n", ns)
	}

	_, ok = os.LookupEnv("T1_O8_TERM_PREFIX")

	if !ok {
		tr.Ln("T1_O8_TERM_PREFIX is not present, please make sure T1_O8_TERM_PREFIX is exported before running the command")
		os.Setenv("T1_O8_TERM_PREFIX", "§")
		//tr.Ln("using default T1_O8_TERM_PREFIX = §")

	} else {
		//tr.Ln("T1_O8_TERM_PREFIX: %s\n", nstp)
	}

	_, ok = os.LookupEnv("T1_O8_TERM_SUFFIX")

	if !ok {
		tr.Ln("T1_O8_TERM_SUFFIX is not present, please make sure T1_O8_TERM_SUFFIX is exported before running the command")
		os.Setenv("T1_O8_TERM_SUFFIX", "!")
		tr.Ln("using default T1_O8_TERM_SUFFIX = !")
	} else {
		//tr.Ln("T1_O8_TERM_SUFFIX: %s\n", nssf)
	}

}

func main() {
	code := run(os.Stdout, os.Args)
	if code != 0 {
		os.Exit(code)
	}
}

func run(w io.Writer, args []string) (code int) {
	if len(args) < 2 {
		fmt.Fprint(w, usageText)
		return 0
	}
	switch args[1] {
	case "generate":
		return generateCmd(w, args[2:])
	case "fmt":
		return fmtCmd(w, args[2:])
	case "lsp":
		return lspCmd(w, args[2:])
	case "version":
		fmt.Fprintln(w, t1.Version())
		return 0
	case "--version":
		fmt.Fprintln(w, t1.Version())
		return 0
	}
	fmt.Fprint(w, usageText)
	return 0
}

func generateCmd(w io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("generate", flag.ExitOnError)
	cmd.SetOutput(w)
	fileNameFlag := cmd.String("f", "", "")
	pathFlag := cmd.String("path", ".", "")
	sourceMapVisualisationsFlag := cmd.Bool("source-map-visualisations", false, "")
	includeVersionFlag := cmd.Bool("include-version", true, "")
	includeTimestampFlag := cmd.Bool("include-timestamp", false, "")
	watchFlag := cmd.Bool("watch", false, "")
	openBrowserFlag := cmd.Bool("open-browser", true, "")
	cmdFlag := cmd.String("cmd", "", "")
	proxyFlag := cmd.String("proxy", "", "")
	proxyPortFlag := cmd.Int("proxyport", 7331, "")
	workerCountFlag := cmd.Int("w", runtime.NumCPU(), "")
	pprofPortFlag := cmd.Int("pprof", 0, "")
	keepOrphanedFilesFlag := cmd.Bool("keep-orphaned-files", false, "")
	verboseFlag := cmd.Bool("v", false, "")
	logLevelFlag := cmd.String("log-level", "info", "")
	helpFlag := cmd.Bool("help", false, "")
	err := cmd.Parse(args)
	if err != nil || *helpFlag {
		fmt.Fprint(w, generateUsageText)
		return
	}

	logLevel := *logLevelFlag
	if *verboseFlag {
		logLevel = "debug"
	}

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Fprintln(w, "Stopping...")
		cancel()
	}()
	err = generatecmd.Run(ctx, w, generatecmd.Arguments{
		FileName:                        *fileNameFlag,
		Path:                            *pathFlag,
		Watch:                           *watchFlag,
		OpenBrowser:                     *openBrowserFlag,
		Command:                         *cmdFlag,
		Proxy:                           *proxyFlag,
		ProxyPort:                       *proxyPortFlag,
		WorkerCount:                     *workerCountFlag,
		GenerateSourceMapVisualisations: *sourceMapVisualisationsFlag,
		IncludeVersion:                  *includeVersionFlag,
		IncludeTimestamp:                *includeTimestampFlag,
		LogLevel:                        logLevel,
		PPROFPort:                       *pprofPortFlag,
		KeepOrphanedFiles:               *keepOrphanedFilesFlag,
	})
	if err != nil {
		color.New(color.FgRed).Fprint(w, "(✗) ")
		fmt.Fprintln(w, err.Error())
		return 1
	}
	return 0
}

func fmtCmd(w io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("fmt", flag.ExitOnError)
	cmd.SetOutput(w)
	cmd.Usage = func() {
		fmt.Fprint(w, fmtUsageText)
	}
	helpFlag := cmd.Bool("help", false, "")
	workerCountFlag := cmd.Int("w", runtime.NumCPU(), "")
	verboseFlag := cmd.Bool("v", false, "")
	logLevelFlag := cmd.String("log-level", "info", "")
	stdout := cmd.Bool("stdout", false, "")

	err := cmd.Parse(args)
	if err != nil || *helpFlag {
		cmd.Usage()
		return
	}

	logLevel := *logLevelFlag
	if *verboseFlag {
		logLevel = "debug"
	}

	err = fmtcmd.Run(w, fmtcmd.Arguments{
		ToStdout:    *stdout,
		Files:       cmd.Args(),
		LogLevel:    logLevel,
		WorkerCount: *workerCountFlag,
	})
	if err != nil {
		return 1
	}
	return 0
}

func lspCmd(w io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("lsp", flag.ExitOnError)
	cmd.SetOutput(w)
	log := cmd.String("log", "", "")
	goplsLog := cmd.String("goplsLog", "", "")
	goplsRPCTrace := cmd.Bool("goplsRPCTrace", false, "")
	helpFlag := cmd.Bool("help", false, "")
	pprofFlag := cmd.Bool("pprof", false, "")
	httpDebugFlag := cmd.String("http", "", "")
	err := cmd.Parse(args)
	if err != nil || *helpFlag {
		fmt.Fprint(w, lspUsageText)
		return
	}
	err = lspcmd.Run(w, lspcmd.Arguments{
		Log:           *log,
		GoplsLog:      *goplsLog,
		GoplsRPCTrace: *goplsRPCTrace,
		PPROF:         *pprofFlag,
		HTTPDebug:     *httpDebugFlag,
	})
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return 1
	}
	return 0
}
