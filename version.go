package t1

import _ "embed"

//go:embed .version
var version string

func Version() string {
	return "v" + version
}
