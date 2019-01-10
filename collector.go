package main

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	namespace = "awsses_exporter"
)

// Exporter struct
type Exporter struct {
	accessKey       string
	secretAccessKey string

	bounces          *prometheus.Desc
	complaints       *prometheus.Desc
	deliveryAttempts *prometheus.Desc
	rejects          *prometheus.Desc
	timestamp        *prometheus.Desc
}

// NewExporter returns an initialized exporter.
func NewExporter() *Exporter {
	return &Exporter{
		bounces: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "bounces"),
			"Bounces per region",
			[]string{"aws_region"},
			nil,
		),
		complaints: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "complaints"),
			"Complaints per region",
			[]string{"aws_region"},
			nil,
		),
		deliveryAttempts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "deliveryAttempts"),
			"Delivery attempts per region",
			[]string{"aws_region"},
			nil,
		),
		rejects: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "rejects"),
			"Rejects per region",
			[]string{"aws_region"},
			nil,
		),
		timestamp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "timestamp"),
			"Timestamp per region",
			[]string{"aws_region"},
			nil,
		),
	}
}

// Describe describes all the metrics exported by the memcached exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.bounces
	ch <- e.complaints
	ch <- e.deliveryAttempts
	ch <- e.rejects
	ch <- e.timestamp
}

// Collect fetches the statistics from the configured memcached server, and
// delivers them as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	awsSesRegions := []string{"us-east-1", "us-west-2"} //, "eu-west-1"}
	latest := make(map[string]*ses.SendDataPoint)
	for _, regionName := range awsSesRegions {
		svc := ses.New(session.New(), aws.NewConfig().WithRegion(regionName))
		input := &ses.GetSendStatisticsInput{}

		result, err := svc.GetSendStatistics(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}
		sort.Slice(result.SendDataPoints, func(i, j int) bool {
			return result.SendDataPoints[i].Timestamp.Unix() > result.SendDataPoints[j].Timestamp.Unix()
		})
		if len(result.SendDataPoints) > 0 {
			latest[regionName] = result.SendDataPoints[0]
		}
		//for _, value := range latest {
		ch <- prometheus.MustNewConstMetric(e.bounces, prometheus.GaugeValue, float64(*latest[regionName].Bounces), regionName)
		ch <- prometheus.MustNewConstMetric(e.complaints, prometheus.GaugeValue, float64(*latest[regionName].Complaints), regionName)
		ch <- prometheus.MustNewConstMetric(e.deliveryAttempts, prometheus.GaugeValue, float64(*latest[regionName].DeliveryAttempts), regionName)
		ch <- prometheus.MustNewConstMetric(e.rejects, prometheus.GaugeValue, float64(*latest[regionName].Rejects), regionName)
		//ch <- prometheus.MustNewConstMetric(e.timestamp, prometheus.GaugeValue, float64(value.Timestamp.Unix()), regionName)
		//}

		//fmt.Printf("%v", latest)
	}

}

func main() {
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9199").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	)
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("awsses_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting awsses_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	prometheus.MustRegister(NewExporter())
	fmt.Println("RESULTSSSSSS")

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>AWS SES Exporter</title></head>
             <body>
             <h1>AWS SES Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Infoln("Starting HTTP server on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
