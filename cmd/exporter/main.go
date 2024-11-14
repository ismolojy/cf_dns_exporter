package main

import (
	"cf_dns_exporter_fork/internal/config"
	"cf_dns_exporter_fork/internal/metrics"
	"cf_dns_exporter_fork/internal/repo"
	"fmt"
	"time"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
	}
	logger := config.SetupLogger(conf.Env)

	go func() {
		for {
			var dnsRecordsResponseAll []repo.DnsRecordsResponse
			dnsZoneIds, err := repo.GetZoneIds(conf.ApiUrl, conf.ApiToken)
			if err != nil {
				logger.Error("Error listing zone ids: ", err)
				continue
			}
			logger.Info("Obtaining dns zone ids is completed.", "The number of received ids is:", len(dnsZoneIds))

			for _, zoneId := range dnsZoneIds {
				var listDnsUrl = fmt.Sprintf("%s/zones/%s/dns_records", conf.ApiUrl, zoneId)
				dnsList, err := repo.ListDnsRecords(listDnsUrl, conf.ApiToken)
				if err != nil {
					logger.Error("Error listing DNS records for zone ID", zoneId, ":", err)
					continue
				}
				dnsRecordsResponseAll = append(dnsRecordsResponseAll, dnsList)
				logger.Debug("Listing DNS records for zone completed", "zoneId", zoneId)
			}
			metrics.GenerateMetrics(dnsRecordsResponseAll)
			logger.Info("Metrics generated successfully")

			time.Sleep(30 * time.Second)
		}
	}()
	err = metrics.Listen(fmt.Sprintf("%s:%s", conf.Address, conf.Port))
	if err != nil {
		logger.Error("Error starting server:", err)
		return
	}

}
