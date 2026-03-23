package collector

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/disk"
)

type DiskCollector struct {
	desc *prometheus.Desc
}

func NewDiskCollector() *DiskCollector {
	return &DiskCollector{
		desc: prometheus.NewDesc(
			"disk_usage_bytes",
			"Disk usage for key mountpoints and total",
			[]string{"mount", "type"}, // total, used, free, used_percent
			nil,
		),
	}
}

func (c *DiskCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *DiskCollector) Collect(ch chan<- prometheus.Metric) {
	var mounts []string

	switch runtime.GOOS {
	case "windows":
		mounts = []string{"C:\\", "D:\\"}
	default:
		mounts = []string{"/", "/tmp"}
	}

	for _, m := range mounts {
		usage, err := disk.Usage(m)
		if err != nil {
			continue
		}

		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(usage.Total), m, "total")
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(usage.Used), m, "used")
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(usage.Free), m, "free")
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, usage.UsedPercent, m, "used_percent")
	}
}
