// SPDX-License-Identifier: Apache-2.0

package metrics

import (
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
	sme         *smbMetricsExporter
	dsc         []*prometheus.Desc
	netbiosName string
}

func (col *smbCollector) Refresh() {
	if col.netbiosName != "" {
		return
	}
	netbiosName, err := resolveNetbiosName()
	if err != nil {
		return
	}
	col.netbiosName = netbiosName
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
	col.Refresh()
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
		col.netbiosName,
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
				"netbiosname",
			}, nil),
	}
	return col
}

type smbStatusCollector struct {
	smbCollector
}

func (col *smbStatusCollector) Collect(ch chan<- prometheus.Metric) {
	smbInfo, err := NewUpdatedSMBInfo(col.sme.log)
	if err != nil {
		return
	}
	col.Refresh()
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
	if !col.sme.profile {
		return
	}
	smbProfileInfo, err := NewUpdatedSMBProfileInfo(col.sme.log)
	if err != nil {
		return
	}
	col.Refresh()
	smb2Calls := smbProfileInfo.profileStatus.SMB2Calls
	if smb2Calls != nil {
		col.collectSMB2CallsMetrics(ch, smb2Calls, "", "")
	}
	sysCalls := smbProfileInfo.profileStatus.SystemCalls
	if sysCalls != nil {
		col.collectSysCallsMetrics(ch, sysCalls, "", "")
	}
	for key, extended := range smbProfileInfo.profileStatus.Extended {
		sharename, client := ParseExtendedProfileKey(key)
		if sharename == "" || client == "" {
			continue
		}
		smb2Calls = extended.SMB2Calls
		if smb2Calls != nil {
			col.collectSMB2CallsMetrics(ch, smb2Calls, sharename, client)
		}
		sysCalls = extended.SystemCalls
		if sysCalls != nil {
			col.collectSysCallsMetrics(ch, sysCalls, sharename, client)
		}
	}
}

func (col *smbProfileCollector) collectSMB2CallsMetrics(
	ch chan<- prometheus.Metric, smb2Calls *SMBProfileSMB2Calls,
	sharename, client string) {
	operationToProfileCallEntry := map[string]*SMBProfileCallEntry{
		"negprot":   &smb2Calls.NegProt,
		"sesssetup": &smb2Calls.SessSetup,
		"logoff":    &smb2Calls.LogOff,
		"tcon":      &smb2Calls.Tcon,
		"tdis":      &smb2Calls.Tdis,
		"create":    &smb2Calls.Create,
		"close":     &smb2Calls.Close,
		"flush":     &smb2Calls.Flush,
		"read":      &smb2Calls.Read,
		"write":     &smb2Calls.Write,
		"lock":      &smb2Calls.Lock,
		"ioctl":     &smb2Calls.Ioctl,
		"cancel":    &smb2Calls.Cancel,
		"keepalive": &smb2Calls.KeepAlive,
		"find":      &smb2Calls.Find,
		"notify":    &smb2Calls.Notify,
		"getinfo":   &smb2Calls.GetInfo,
		"setinfo":   &smb2Calls.SetInfo,
		"break":     &smb2Calls.Break,
	}
	for op, pce := range operationToProfileCallEntry {
		ch <- col.smb2RequestTotalMetric(sharename, client, op, pce)
		ch <- col.smb2RequestInbytesMetric(sharename, client, op, pce)
		ch <- col.smb2RequestOutbytesMetric(sharename, client, op, pce)
		ch <- col.smb2RequestDurationMetric(sharename, client, op, pce)
	}
}

func (col *smbProfileCollector) collectSysCallsMetrics(
	ch chan<- prometheus.Metric, sysCalls *SMBProfileSyscalls,
	sharename, client string) {
	operationToProfileIOEntry := map[string]*SMBProfileIOEntry{
		"pread":       &sysCalls.PRead,
		"asys_pread":  &sysCalls.AsysPRead,
		"pwrite":      &sysCalls.PWrite,
		"asys_pwrite": &sysCalls.AsysPWrite,
		"asys_fsync":  &sysCalls.AsysFSync,
	}
	operationToProfileEntry := map[string]*SMBProfileEntry{
		"opendir":    &sysCalls.Opendir,
		"fdopendir":  &sysCalls.FDOpendir,
		"readdir":    &sysCalls.Readdir,
		"rewinddir":  &sysCalls.Rewinddir,
		"mkdirat":    &sysCalls.Mkdirat,
		"closedir":   &sysCalls.Closedir,
		"open":       &sysCalls.Open,
		"openat":     &sysCalls.OpenAt,
		"createfile": &sysCalls.CreateFile,
		"close":      &sysCalls.Close,
		"lseek":      &sysCalls.Lseek,
		"renameat":   &sysCalls.RenameAt,
		"stat":       &sysCalls.Stat,
		"fstat":      &sysCalls.FStat,
		"lstat":      &sysCalls.LStat,
		"fstatat":    &sysCalls.FStatAt,
		"unlinkat":   &sysCalls.UnlinkAt,
		"chmod":      &sysCalls.Chmod,
		"fchmod":     &sysCalls.FChmod,
		"fchown":     &sysCalls.FChown,
		"lchown":     &sysCalls.LChown,
		"chdir":      &sysCalls.Chdir,
		"getwd":      &sysCalls.GetWD,
		"fntimes":    &sysCalls.Fntimes,
		"ftruncate":  &sysCalls.FTruncate,
		"fallocate":  &sysCalls.FAllocate,
		"readlinkat": &sysCalls.ReadLinkAt,
		"symlinkat":  &sysCalls.SymLinkAt,
		"linkat":     &sysCalls.LinkAt,
		"mknodat":    &sysCalls.MknodAt,
		"realpath":   &sysCalls.RealPath,
	}
	for op, pioe := range operationToProfileIOEntry {
		ch <- col.vfsIOTotalMetric(sharename, client, op, pioe)
		ch <- col.vfsIOBytesMetric(sharename, client, op, pioe)
		ch <- col.vfsIODurationMetric(sharename, client, op, pioe)
	}
	for op, pe := range operationToProfileEntry {
		ch <- col.vfsTotalMetric(sharename, client, op, pe)
		ch <- col.vfsDurationMetric(sharename, client, op, pe)
	}
}

func (col *smbProfileCollector) smb2RequestTotalMetric(
	sharename, client, operation string, pce *SMBProfileCallEntry) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		col.dsc[0],
		prometheus.GaugeValue,
		float64(pce.Count),
		col.netbiosName,
		sharename,
		client,
		operation)
}

func (col *smbProfileCollector) smb2RequestInbytesMetric(
	sharename, client, operation string, pce *SMBProfileCallEntry) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		col.dsc[1],
		prometheus.GaugeValue,
		float64(pce.Inbytes),
		col.netbiosName,
		sharename,
		client,
		operation)
}

func (col *smbProfileCollector) smb2RequestOutbytesMetric(
	sharename, client, operation string, pce *SMBProfileCallEntry) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		col.dsc[2],
		prometheus.GaugeValue,
		float64(pce.Outbytes),
		col.netbiosName,
		sharename,
		client,
		operation)
}

func (col *smbProfileCollector) smb2RequestDurationMetric(
	sharename, client, operation string, pce *SMBProfileCallEntry) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		col.dsc[3],
		prometheus.GaugeValue,
		float64(pce.Time),
		col.netbiosName,
		sharename,
		client,
		operation)
}

func (col *smbProfileCollector) vfsIOTotalMetric(
	sharename, client, operation string, pioe *SMBProfileIOEntry) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		col.dsc[4],
		prometheus.GaugeValue,
		float64(pioe.Count),
		col.netbiosName,
		sharename,
		client,
		operation)
}

func (col *smbProfileCollector) vfsIOBytesMetric(
	sharename, client, operation string, pioe *SMBProfileIOEntry) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		col.dsc[5],
		prometheus.GaugeValue,
		float64(pioe.Bytes),
		col.netbiosName,
		sharename,
		client,
		operation)
}

func (col *smbProfileCollector) vfsIODurationMetric(
	sharename, client, operation string, pioe *SMBProfileIOEntry) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		col.dsc[6],
		prometheus.GaugeValue,
		float64(pioe.Time),
		col.netbiosName,
		sharename,
		client,
		operation)
}

func (col *smbProfileCollector) vfsTotalMetric(
	sharename, client, operation string, pe *SMBProfileEntry) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		col.dsc[7],
		prometheus.GaugeValue,
		float64(pe.Count),
		col.netbiosName,
		sharename,
		client,
		operation)
}

func (col *smbProfileCollector) vfsDurationMetric(
	sharename, client, operation string, pe *SMBProfileEntry) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		col.dsc[8],
		prometheus.GaugeValue,
		float64(pe.Time),
		col.netbiosName,
		sharename,
		client,
		operation)
}

func (sme *smbMetricsExporter) newSMBProfileCollector() prometheus.Collector {
	variableLabels := []string{"netbiosname", "share", "client", "operation"}
	col := &smbProfileCollector{}
	col.sme = sme
	col.dsc = []*prometheus.Desc{
		prometheus.NewDesc(
			collectorName("smb2", "request_total"),
			"Total number of SMB2 requests",
			variableLabels, nil),
		prometheus.NewDesc(
			collectorName("smb2", "request_inbytes"),
			"Bytes received for SMB2 requests",
			variableLabels, nil),
		prometheus.NewDesc(
			collectorName("smb2", "request_outbytes"),
			"Bytes replied for SMB2 requests",
			variableLabels, nil),
		prometheus.NewDesc(
			collectorName("smb2", "request_duration_microseconds_sum"),
			"Execution time in microseconds of SMB2 requests",
			variableLabels, nil),
		prometheus.NewDesc(
			collectorName("vfs_io", "total"),
			"Total number of I/O calls to underlying VFS layer",
			variableLabels, nil),
		prometheus.NewDesc(
			collectorName("vfs_io", "bytes"),
			"Number of bytes transferred via underlying VFS I/O layer",
			variableLabels, nil),
		prometheus.NewDesc(
			collectorName("vfs_io", "duration_microseconds_sum"),
			"Execution time in microseconds of VFS I/O requests",
			variableLabels, nil),
		prometheus.NewDesc(
			collectorName("vfs", "total"),
			"Total number of calls to underlying VFS layer",
			variableLabels, nil),
		prometheus.NewDesc(
			collectorName("vfs", "duration_microseconds_sum"),
			"Execution time in microseconds of VFS requests",
			variableLabels, nil),
	}

	return col
}

func collectorName(subsystem, name string) string {
	return prometheus.BuildFQName(collectorsNamespace, subsystem, name)
}
