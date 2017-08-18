package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/koding/websocketproxy"
)

func isWebsocket(req *http.Request) bool {
	connHdr := ""
	connHdrs := req.Header["Connection"]
	if len(connHdrs) > 0 {
		connHdr = connHdrs[0]
	}

	upgradeWebsocket := false
	if strings.ToLower(connHdr) == "upgrade" {
		upgradeHdrs := req.Header["Upgrade"]
		if len(upgradeHdrs) > 0 {
			upgradeWebsocket = (strings.ToLower(upgradeHdrs[0]) == "websocket")
		}
	}

	return upgradeWebsocket
}

func newWebsocketHandler(target *url.URL, alternative http.Handler) http.Handler {
	u := *target
	if target.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	ws := websocketproxy.NewProxy(&u)

	websocketproxy.DefaultUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isWebsocket(r) {
			ws.ServeHTTP(w, r)
		} else {
			alternative.ServeHTTP(w, r)
		}
	})
}
