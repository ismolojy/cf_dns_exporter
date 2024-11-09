package app

import (
	"cf_dns_exporter_fork/internal/config"
	"cf_dns_exporter_fork/internal/metrics"
	"fmt"
	"time"
)

func StartApp() {
	configFile, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	go func() {
		for {
			var dnsRecordsResponseAll []metrics.DnsRecordsResponse
			dnsZoneIds, err := metrics.ListZoneIds(configFile.ApiUrl, configFile.ApiToken)
			if err != nil {
				fmt.Println("Error listing zone ids: ", err)
				continue
			}

			for _, zoneId := range dnsZoneIds {
				var listDnsUrl = fmt.Sprintf("%s/zones/%s/dns_records", configFile.ApiUrl, zoneId)
				dnsList, err := metrics.ListDnsRecords(listDnsUrl, configFile.ApiToken)
				if err != nil {
					fmt.Println("Error listing DNS records for zone ID", zoneId, ":", err)
					continue
				}
				dnsRecordsResponseAll = append(dnsRecordsResponseAll, dnsList)
			}
			metrics.GenerateMetrics(dnsRecordsResponseAll)

			time.Sleep(30 * time.Second)
		}
	}()
	err = metrics.Listen(fmt.Sprintf(":%s", configFile.Port))
	if err != nil {
		return
	}

}
