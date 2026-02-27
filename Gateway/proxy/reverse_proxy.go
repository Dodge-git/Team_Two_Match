package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewReverseProxy(target string) (*httputil.ReverseProxy, error) {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api")

		if req.URL.Path == "" {
			req.URL.Path = "/"
		}
	}

	return proxy, nil
}











// package proxy

// import (
// 	"net/http/httputil"
// 	"net/url"
// )

// func NewReverseProxy(target string) (*httputil.ReverseProxy,error){
// 	parsedURL,err := url.Parse(target)
// 	if err != nil {
// 		return nil,err
// 	}

// 	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

// 	return proxy,nil
// }