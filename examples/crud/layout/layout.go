package layout

import (
	"net/http"

	"github.com/senforsce/tndr"
)

func Handler(content tndr.Component) http.Handler {
	return tndr.Handler(Page(content))
}
