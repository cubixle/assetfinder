package assetfinder

import (
	"fmt"
	"net/url"
)

type URLScan struct{}

type URLScanResults struct {
	Results []struct {
		Task struct {
			URL string `json:"url"`
		} `json:"task"`

		Page struct {
			URL string `json:"url"`
		} `json:"page"`
	} `json:"results"`
}

func (URLScan) FetchSubDomains(domain string) ([]string, error) {
	resp := URLScanResults{}
	err := fetchJSON(
		fmt.Sprintf("https://urlscan.io/api/v1/search/?q=domain:%s", domain),
		resp,
	)
	if err != nil {
		return []string{}, err
	}
	output := make([]string, 0)

	for _, r := range resp.Results {
		u, err := url.Parse(r.Task.URL)
		if err != nil {
			continue
		}

		output = append(output, u.Hostname())
	}

	for _, r := range resp.Results {
		u, err := url.Parse(r.Page.URL)
		if err != nil {
			continue
		}

		output = append(output, u.Hostname())
	}

	return output, nil
}
