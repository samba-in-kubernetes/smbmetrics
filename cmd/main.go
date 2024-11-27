// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	goruntime "runtime"

	"github.com/spf13/pflag"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/samba-in-kubernetes/smbmetrics/internal/metrics"
)

var (
	// Version of the software at compile time.
	Version = "(unset)"
	// CommitID of the revision used to compile the software.
	CommitID = "(unset)"
)

func init() {
	metrics.UpdateDefaultVersions(Version, CommitID)
}

func main() {
	var port int
	pflag.IntVar(&port, "port", metrics.DefaultMetricsPort,
		"Prometheus metrics-exporter port number")
	pflag.Parse()

	log := zap.New(zap.UseDevMode(true))
	log.Info("Initializing smbmetrics",
		"ProgramName", os.Args[0],
		"GoVersion", goruntime.Version())

	vers, _ := metrics.ResolveVersions(nil)
	log.Info("Versions", "Versions", vers)

	podid := metrics.GetSelfPodID()
	if len(podid.Name) > 0 {
		log.Info("Self", "PodID", podid)
	}

	loc, err := metrics.LocateSMBStatus()
	if err != nil {
		log.Error(err, "Failed to locate smbstatus")
		os.Exit(1)
	}
	ver, err := metrics.RunSMBStatusVersion()
	if err != nil {
		log.Error(err, "Failed to run smbstatus")
		os.Exit(1)
	}
	log.Info("Located smbstatus", "path", loc, "version", ver)

	err = metrics.RunSmbMetricsExporter(log, port)
	if err != nil {
		os.Exit(1)
	}
}
