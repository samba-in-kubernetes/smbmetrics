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

func collectorName(subsystem, name string) string {
	return prometheus.BuildFQName(collectorsNamespace, subsystem, name)
}
