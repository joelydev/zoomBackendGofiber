package main

import (
	"log"
	"net/http"
	"strings"

	"go-fiber-auth/utilities"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	proxy.KeepDestinationHeaders = true
	proxy.KeepHeader = false

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			proxyAuthToken := r.Header.Get("Proxy-Authorization")

			if proxyAuthToken == "" {
				return r, goproxy.NewResponse(r,
					goproxy.ContentTypeText, http.StatusForbidden,
					"Invalid proxy authorization token!")
			}

			trimmedToken := strings.TrimSpace(proxyAuthToken)

			if trimmedToken == "" {
				return r, goproxy.NewResponse(r,
					goproxy.ContentTypeText, http.StatusForbidden,
					"Invalid proxy authorization token!")
			}

			claims, parsingError := utilities.ParseClaims(trimmedToken)
			if parsingError != nil {
				return r, goproxy.NewResponse(r,
					goproxy.ContentTypeText, http.StatusForbidden,
					"Proxy token is not authorized!")
			}

			log.Println(claims.UserId)

			return r, nil
		})

	http.ListenAndServe(":8080", proxy)
}
