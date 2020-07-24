package assetfinder

import (
	"encoding/json"
	"net/http"
	"time"
)

var LoadedSources = []Source{
	BufferOverrun{},
	CertSpotter{},
	CrtSh{},
	URLScan{},
}

// Source is an interface that defines the methods required by a source.
type Source interface {
	FetchSubDomains(domain string) ([]string, error)
}

// Result is the basic type that should be return from the sources.Fetch method.
type Result struct {
	Domain     string
	StatusCode int
	Error      error
}

func fetchJSON(url string, i interface{}) error {
	client := httpClient()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(i)
}

func fetchStatus(url string) Result {
	res := Result{Domain: url}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		res.Error = err
		return res
	}

	client := httpClient()
	resp, err := client.Do(request)
	if err != nil {
		res.Error = err
		return res
	}
	res.StatusCode = resp.StatusCode
	return res
}

func httpClient() http.Client {
	return http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
}
