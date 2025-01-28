// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	collectorsNamespace = "smb"
)

func (sme *smbMetricsExporter) register() error {
	cols := []prometheus.Collector{
		sme.newSMBVersionsCollector(),
		sme.newSMBStatusCollector(),
		sme.newSMBProfileCollector(),
	}
	for _, c := range cols {
		if err := sme.reg.Register(c); err != nil {
			sme.log.Error(err, "failed to register collector")
			return err
		}
	}
	return nil
}

type smbCollector struct {
	// nolint:structcheck
	sme *smbMetricsExporter
	dsc []*prometheus.Desc
}

func (col *smbCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, d := range col.dsc {
		ch <- d
	}
}

type smbVersionsCollector struct {
	smbCollector
	clnt *kclient
}

func (col *smbVersionsCollector) Collect(ch chan<- prometheus.Metric) {
	status := 0
	vers, err := ResolveVersions(col.clnt)
	if err != nil {
		status = 1
	}
	ch <- prometheus.MustNewConstMetric(
		col.dsc[0],
		prometheus.GaugeValue,
		float64(status),
		vers.Version,
		vers.CommitID,
		vers.SambaImage,
		vers.SambaVersion,
		vers.CtdbVersion,
	)
}

func (sme *smbMetricsExporter) newSMBVersionsCollector() prometheus.Collector {
	col := &smbVersionsCollector{}
	col.sme = sme
	col.clnt, _ = newKClient()
	col.dsc = []*prometheus.Desc{
		prometheus.NewDesc(
			collectorName("metrics", "status"),
			"Current metrics-collector status versions",
			[]string{
				"version",
				"commitid",
				"sambaimage",
				"sambavers",
				"ctdbvers",
			}, nil),
	}
	return col
}

type smbStatusCollector struct {
	smbCollector
}

func (col *smbStatusCollector) Collect(ch chan<- prometheus.Metric) {
	smbInfo, err := NewUpdatedSMBInfo()
	if err != nil {
		return
	}
	ch <- prometheus.MustNewConstMetric(col.dsc[0],
		prometheus.GaugeValue, float64(smbInfo.TotalSessions()))

	ch <- prometheus.MustNewConstMetric(col.dsc[1],
		prometheus.GaugeValue, float64(smbInfo.TotalTreeCons()))

	ch <- prometheus.MustNewConstMetric(col.dsc[2],
		prometheus.GaugeValue, float64(smbInfo.TotalConnectedUsers()))

	serviceToMachine := smbInfo.MapServiceToMachines()
	for service, machines := range serviceToMachine {
		ch <- prometheus.MustNewConstMetric(col.dsc[3],
			prometheus.GaugeValue,
			float64(len(machines)),
			service)
	}
	machineToServices := smbInfo.MapMachineToServies()
	for machine, services := range machineToServices {
		ch <- prometheus.MustNewConstMetric(col.dsc[4],
			prometheus.GaugeValue,
			float64(len(services)),
			machine)
	}
}

func (sme *smbMetricsExporter) newSMBStatusCollector() prometheus.Collector {
	col := &smbStatusCollector{}
	col.sme = sme
	col.dsc = []*prometheus.Desc{
		prometheus.NewDesc(
			collectorName("sessions", "total"),
			"Number of currently active SMB sessions",
			[]string{}, nil),

		prometheus.NewDesc(
			collectorName("tcon", "total"),
			"Number of currently active SMB tree-connections",
			[]string{}, nil),

		prometheus.NewDesc(
			collectorName("users", "total"),
			"Number of currently active SMB users",
			[]string{}, nil),

		prometheus.NewDesc(
			collectorName("share", "activity"),
			"Number of remote machines currently using a share",
			[]string{"service"}, nil),

		prometheus.NewDesc(
			collectorName("share", "byremote"),
			"Number of shares served for remote machine",
			[]string{"machine"}, nil),
	}
	return col
}

type smbProfileCollector struct {
	smbCollector
}

func (col *smbProfileCollector) Collect(ch chan<- prometheus.Metric) {
	if col.sme.profile {
		smbProfileInfo, err := NewUpdatedSMBProfileInfo()
		if err == nil {
			smb2Calls := smbProfileInfo.profileStatus.SMB2Calls
			ch <- col.smb2RequestMetric(&smb2Calls.NegProt, "negprot")
			ch <- col.smb2RequestMetric(&smb2Calls.SessSetup, "sesssetup")
			ch <- col.smb2RequestMetric(&smb2Calls.LogOff, "logoff")
			ch <- col.smb2RequestMetric(&smb2Calls.Tcon, "tcon")
			ch <- col.smb2RequestMetric(&smb2Calls.Tdis, "tdis")
			ch <- col.smb2RequestMetric(&smb2Calls.Create, "create")
			ch <- col.smb2RequestMetric(&smb2Calls.Close, "close")
			ch <- col.smb2RequestMetric(&smb2Calls.Flush, "flush")
			ch <- col.smb2RequestMetric(&smb2Calls.Read, "read")
			ch <- col.smb2RequestMetric(&smb2Calls.Write, "write")
			ch <- col.smb2RequestMetric(&smb2Calls.Lock, "lock")
			ch <- col.smb2RequestMetric(&smb2Calls.Ioctl, "ioctl")
			ch <- col.smb2RequestMetric(&smb2Calls.Cancel, "cancel")
			ch <- col.smb2RequestMetric(&smb2Calls.KeepAlive, "keepalive")
			ch <- col.smb2RequestMetric(&smb2Calls.Find, "find")
			ch <- col.smb2RequestMetric(&smb2Calls.Notify, "notify")
			ch <- col.smb2RequestMetric(&smb2Calls.GetInfo, "getinfo")
			ch <- col.smb2RequestMetric(&smb2Calls.SetInfo, "setinfo")
			ch <- col.smb2RequestMetric(&smb2Calls.Break, "break")

			sysCalls := smbProfileInfo.profileStatus.SystemCalls
			ch <- col.vfsIORequestMetric(&sysCalls.PRead, "pread")
			ch <- col.vfsIORequestMetric(&sysCalls.AsysPRead, "asys_pread")
			ch <- col.vfsIORequestMetric(&sysCalls.PWrite, "pwrite")
			ch <- col.vfsIORequestMetric(&sysCalls.AsysPWrite, "asys_pwrite")
			ch <- col.vfsIORequestMetric(&sysCalls.AsysFSync, "asys_fsync")

			ch <- col.vfsRequestMetric(&sysCalls.Opendir, "opendir")
			ch <- col.vfsRequestMetric(&sysCalls.FDOpendir, "fdopendir")
			ch <- col.vfsRequestMetric(&sysCalls.Readdir, "readdir")
			ch <- col.vfsRequestMetric(&sysCalls.Rewinddir, "rewinddir")
			ch <- col.vfsRequestMetric(&sysCalls.Mkdirat, "mkdirat")
			ch <- col.vfsRequestMetric(&sysCalls.Closedir, "closedir")
			ch <- col.vfsRequestMetric(&sysCalls.Open, "open")
			ch <- col.vfsRequestMetric(&sysCalls.OpenAt, "openat")
			ch <- col.vfsRequestMetric(&sysCalls.CreateFile, "createfile")
			ch <- col.vfsRequestMetric(&sysCalls.Close, "close")
			ch <- col.vfsRequestMetric(&sysCalls.Lseek, "lseek")
			ch <- col.vfsRequestMetric(&sysCalls.RenameAt, "renameat")
			ch <- col.vfsRequestMetric(&sysCalls.Stat, "stat")
			ch <- col.vfsRequestMetric(&sysCalls.FStat, "fstat")
			ch <- col.vfsRequestMetric(&sysCalls.LStat, "lstat")
			ch <- col.vfsRequestMetric(&sysCalls.FStatAt, "fstatat")
			ch <- col.vfsRequestMetric(&sysCalls.UnlinkAt, "unlinkat")
			ch <- col.vfsRequestMetric(&sysCalls.Chmod, "chmod")
			ch <- col.vfsRequestMetric(&sysCalls.FChmod, "fchmod")
			ch <- col.vfsRequestMetric(&sysCalls.FChown, "fchown")
			ch <- col.vfsRequestMetric(&sysCalls.LChown, "lchown")
			ch <- col.vfsRequestMetric(&sysCalls.Chdir, "chdir")
			ch <- col.vfsRequestMetric(&sysCalls.GetWD, "getwd")
			ch <- col.vfsRequestMetric(&sysCalls.Fntimes, "fntimes")
			ch <- col.vfsRequestMetric(&sysCalls.FTruncate, "ftruncate")
			ch <- col.vfsRequestMetric(&sysCalls.FAllocate, "fallocate")
			ch <- col.vfsRequestMetric(&sysCalls.ReadLinkAt, "readlinkat")
			ch <- col.vfsRequestMetric(&sysCalls.SymLinkAt, "symlinkat")
			ch <- col.vfsRequestMetric(&sysCalls.LinkAt, "linkat")
			ch <- col.vfsRequestMetric(&sysCalls.MknodAt, "mknodat")
			ch <- col.vfsRequestMetric(&sysCalls.RealPath, "realpath")
		}
	}
}

func (col *smbProfileCollector) smb2RequestMetric(pce *SMBProfileCallEntry,
	operation string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		col.dsc[0],
		prometheus.GaugeValue,
		float64(pce.Count),
		strconv.Itoa(pce.Time),
		strconv.Itoa(pce.Idle),
		strconv.Itoa(pce.Inbytes),
		strconv.Itoa(pce.Outbytes),
		operation)
}

func (col *smbProfileCollector) vfsIORequestMetric(pioe *SMBProfileIOEntry,
	operation string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		col.dsc[1],
		prometheus.GaugeValue,
		float64(pioe.Count),
		strconv.Itoa(pioe.Time),
		strconv.Itoa(pioe.Idle),
		strconv.Itoa(pioe.Bytes),
		operation)
}

func (col *smbProfileCollector) vfsRequestMetric(pe *SMBProfileEntry,
	operation string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		col.dsc[2],
		prometheus.GaugeValue,
		float64(pe.Count),
		strconv.Itoa(pe.Time),
		operation)
}

func (sme *smbMetricsExporter) newSMBProfileCollector() prometheus.Collector {
	col := &smbProfileCollector{}
	col.sme = sme
	col.dsc = []*prometheus.Desc{
		prometheus.NewDesc(
			collectorName("smb2", "request_total"),
			"Total number of SMB2 requests",
			[]string{"time", "idle", "inbytes", "outbytes", "operation"}, nil),
		prometheus.NewDesc(
			collectorName("vfs_io", "call_total"),
			"Total number of I/O calls to underlying VFS layer",
			[]string{"time", "idle", "bytes", "operation"}, nil),
		prometheus.NewDesc(
			collectorName("vfs", "call_total"),
			"Total number of calls to underlying VFS layer",
			[]string{"time", "operation"}, nil),
	}

	return col
}

func collectorName(subsystem, name string) string {
	return prometheus.BuildFQName(collectorsNamespace, subsystem, name)
}
