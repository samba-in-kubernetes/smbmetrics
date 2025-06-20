// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"fmt"
	"net"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// DefaultMetricsPort is the default port used to export prometheus metrics
	DefaultMetricsPort = int(9922)
	// DefaultMetricsPath is the default HTTP path to export prometheus metrics
	DefaultMetricsPath = "/metrics"
)

type smbMetricsExporter struct {
	log           logr.Logger
	reg           *prometheus.Registry
	mux           *http.ServeMux
	port          int
	bindAddresses []net.IP
	profile       bool
}

func newSmbMetricsExporter(
	log logr.Logger, port int, bindAddresses []net.IP, profile bool) *smbMetricsExporter {
	return &smbMetricsExporter{
		log:           log,
		reg:           prometheus.NewRegistry(),
		mux:           http.NewServeMux(),
		port:          port,
		bindAddresses: bindAddresses,
		profile:       profile,
	}
}

func (sme *smbMetricsExporter) init() error {
	sme.log.Info("register collectors")
	return sme.register()
}

func (sme *smbMetricsExporter) serve() error {
	var addr string
	switch len(sme.bindAddresses) {
	case 0:
		addr = fmt.Sprintf(":%d", sme.port)
	case 1:
		addr = fmt.Sprintf("%s:%d", sme.bindAddresses[0], sme.port)
	default:
		return fmt.Errorf("too many listen addresses")
	}
	sme.log.Info("serve metrics", "addr", addr)

	handler := promhttp.HandlerFor(sme.reg, promhttp.HandlerOpts{})
	sme.mux.Handle(DefaultMetricsPath, handler)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		sme.log.Error(err, "failed to listen", "addr", addr)
		return err
	}
	defer listener.Close()

	if err := http.Serve(listener, sme.mux); err != nil {
		sme.log.Error(err, "HTTP server failure", "addr", addr)
		return err
	}
	return nil
}

// RunSmbMetricsExporter executes an HTTP server and exports SMB metrics to
// Prometheus.
func RunSmbMetricsExporter(
	log logr.Logger, port int, bindAddresses []net.IP, profile bool) error {
	if port <= 0 {
		port = DefaultMetricsPort
	}
	sme := newSmbMetricsExporter(log, port, bindAddresses, profile)
	err := sme.init()
	if err != nil {
		return err
	}
	return sme.serve()
}
