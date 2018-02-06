// +build !unittest

package urllib

import (
	"net/http"
)

// Response executes request client gets response mannually.
func (b *HttpRequest) Response() (*http.Response, error) {
	return b.getResponse()
}
