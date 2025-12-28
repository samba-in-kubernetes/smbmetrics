// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"net"
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
	var bindAddress net.IP
	pflag.IPVar(&bindAddress, "address", bindAddress,
		"Prometheus metrics-exporter bind address")
	var noProfile bool
	pflag.BoolVar(&noProfile, "no-profile", false,
		"Run without collecting profile information")
	var showVersions bool
	pflag.BoolVar(&showVersions, "show-versions", false,
		"Show versions info and exit")
	pflag.Parse()

	if showVersions {
		showVersionsAndExit()
	}

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

	var bindAddrs []net.IP
	if len(bindAddress) > 0 {
		bindAddrs = append(bindAddrs, bindAddress)
		log.Info("User supplied bind addresses", "bindAddrs", bindAddrs)
	}
	err = metrics.RunSmbMetricsExporter(log, port, bindAddrs, !noProfile)
	if err != nil {
		os.Exit(1)
	}
}

func showVersionsAndExit() {
	vers, _ := metrics.ResolveVersions(nil)
	fmt.Println("Progname:", os.Args[0])
	fmt.Println("Version:", vers.Version)
	fmt.Println("CommitID:", vers.CommitID)
	fmt.Println("GoVersion:", goruntime.Version())
	fmt.Println("Arch:", goruntime.GOARCH)
	fmt.Println("SambaVersion:", vers.SambaVersion)
	os.Exit(0)
}
