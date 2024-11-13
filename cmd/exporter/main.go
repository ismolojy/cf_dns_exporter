package main

import (
	"cf_dns_exporter_fork/internal/config"
	"cf_dns_exporter_fork/internal/metrics"
	"cf_dns_exporter_fork/internal/repo"
	"flag"
	"fmt"
	"time"
)

var ConfigFile = flag.String("c", "./Config.yaml", "Config path")

func main() {
	flag.Parse()
	loadConfig, err := config.LoadConfig(ConfigFile)
	if err != nil {
		fmt.Println("Error loading loadConfig:", err)
	}

	logger := config.SetupLogger(loadConfig.Env)

	go func() {
		for {
			var dnsRecordsResponseAll []repo.DnsRecordsResponse
			dnsZoneIds, err := repo.GetZoneIds(loadConfig.ApiUrl, loadConfig.ApiToken)
			if err != nil {
				logger.Error("Error listing zone ids: ", err)
				continue
			}
			logger.Info("Obtaining dns zone ids is completed.", "The number of received ids is:", len(dnsZoneIds))

			for _, zoneId := range dnsZoneIds {
				var listDnsUrl = fmt.Sprintf("%s/zones/%s/dns_records", loadConfig.ApiUrl, zoneId)
				dnsList, err := repo.ListDnsRecords(listDnsUrl, loadConfig.ApiToken)
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
	err = metrics.Listen(fmt.Sprintf(":%s", loadConfig.Port))
	logger.Info("Server started on port ", loadConfig.Port)
	if err != nil {
		logger.Error("Error starting server:", err)
		return
	}

}
