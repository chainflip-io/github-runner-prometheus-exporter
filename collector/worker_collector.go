package collector

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/thineshsubramani/github-runner-prometheus-exporter/internal/parser"
)

type WorkerCollector struct {
	logPath         string
	workflowStart   *prometheus.GaugeVec
	workflowEnd     *prometheus.GaugeVec
	workflowRuntime *prometheus.GaugeVec
}

func NewWorkerCollector(path string) *WorkerCollector {
	labelKeys := []string{
		"run_id",
		"repository",
		"repository_owner",
		"workflow",
	}

	return &WorkerCollector{
		logPath: path,
		workflowStart: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "github_workflow_start_timestamp_seconds",
			Help: "Start time of GitHub workflow run",
		}, labelKeys),
		workflowEnd: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "github_workflow_end_timestamp_seconds",
			Help: "End time of GitHub workflow run",
		}, labelKeys),
		workflowRuntime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "github_workflow_duration_seconds",
			Help: "Duration of GitHub workflow run",
		}, labelKeys),
	}
}

func (c *WorkerCollector) Describe(ch chan<- *prometheus.Desc) {
	c.workflowStart.Describe(ch)
	c.workflowEnd.Describe(ch)
	c.workflowRuntime.Describe(ch)
}

func (c *WorkerCollector) Collect(ch chan<- prometheus.Metric) {
	worker, err := parser.ParseLatestWorkerLog(c.logPath)
	if err != nil || worker == nil || worker.RunID == "" {
		return
	}

	labels := []string{
		defaultIfEmpty(worker.RunID),
		defaultIfEmpty(worker.Repo),
		defaultIfEmpty(worker.Owner),
		defaultIfEmpty(worker.Workflow),
	}

	c.workflowStart.WithLabelValues(labels...).Set(float64(worker.StartTime.Unix()))
	c.workflowEnd.WithLabelValues(labels...).Set(float64(worker.EndTime.Unix()))
	c.workflowRuntime.WithLabelValues(labels...).Set(worker.TotalRuntime.Seconds())

	c.workflowStart.Collect(ch)
	c.workflowEnd.Collect(ch)
	c.workflowRuntime.Collect(ch)
}

func defaultIfEmpty(s string) string {
	if strings.TrimSpace(s) == "" {
		return "unknown"
	}
	return strings.Trim(s, "{}\" ")
}
