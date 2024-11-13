package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DnsRecord struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	ModifiedOn string `json:"modified_on"`
	ZoneName   string `json:"zone_name"`
}

type DnsRecordsResponse struct {
	Success    bool        `json:"success"`
	Result     []DnsRecord `json:"result"`
	ResultInfo struct {
		Count int `json:"count"`
	} `json:"result_info"`
}

type ZoneId struct {
	Id string `json:"id"`
}
type ZoneIdResponse struct {
	Result []ZoneId `json:"result"`
}

func GetZoneIds(apiUrl, apiToken string) ([]string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/zones", apiUrl), nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+apiToken)
	req.Header.Add("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Printf("HTTP Status Code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var r ZoneIdResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println(err)
	}

	var ids []string
	for _, zone := range r.Result {
		ids = append(ids, zone.Id)
	}

	return ids, err

}

func ListDnsRecords(dnsUrl, apiToken string) (DnsRecordsResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", dnsUrl, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+apiToken)
	req.Header.Add("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Printf("HTTP Status Code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var r DnsRecordsResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println(err)
	}

	return r, err
}
