package backend

import (
	"net/http/httputil"
	"net/url"
)

type Backend struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
}
