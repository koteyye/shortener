package handlers

import (
	"net"
	"net/http"
)

const IPHeader = "X-Real-IP"

func (h *Handlers) Subnet(next http.Handler) http.Handler {
	subnetFn := func(res http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get(IPHeader)
		if ip == "" || !h.subnet.Contains(net.ParseIP(ip)) {
			res.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(res, r)
	}
	return http.HandlerFunc(subnetFn)
}
