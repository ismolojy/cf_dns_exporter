package metrics

import (
	"cf_dns_exporter_fork/internal/config"
	"cf_dns_exporter_fork/internal/repo"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	CfDnsModifiedTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cf_dns_record_modified_time",
			Help: "DNS record modification time",
		},
		[]string{"name", "type"},
	)
	CfDnsDomainCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cf_dns_domain_counter",
			Help: "Number of DNS records in the domain",
		},
		[]string{"domain"},
	)
	env    = "prod"
	logger = config.SetupLogger(env)
)

func Listen(address string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	logger.Info("Server started ", `to: `, address)
	return http.ListenAndServe(address, mux)
}

func GenerateMetrics(responseAll []repo.DnsRecordsResponse) {
	CfDnsModifiedTime.Reset()
	CfDnsDomainCounter.Reset()
	for _, dnsList := range responseAll {
		if len(dnsList.Result) == 0 {
			continue
		}
		for _, record := range dnsList.Result {
			modifiedOn, err := time.Parse(time.RFC3339, record.ModifiedOn)
			if err != nil {
				fmt.Println(err)
			}
			CfDnsModifiedTime.WithLabelValues(record.Name, record.Type).Set(float64(modifiedOn.Unix()))
		}
		CfDnsDomainCounter.WithLabelValues(dnsList.Result[0].ZoneName).Set(float64(dnsList.ResultInfo.Count))
	}
}

func init() {
	prometheus.MustRegister(CfDnsModifiedTime)
	prometheus.MustRegister(CfDnsDomainCounter)
}
