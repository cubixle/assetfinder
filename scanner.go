package assetfinder

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
)

func Scanner(verboseLogging, checkHTTPS, disableStatusCheck bool, outputFile, domain string) error {
	results := []string{}
	for _, source := range LoadedSources {
		subdomains, err := source.FetchSubDomains(domain)
		if err != nil {
			continue
		}
		results = append(results, subdomains...)
	}
	sortedDomains := sortDomains(results)

	ress := []Result{}
	if !disableStatusCheck {
		for _, r := range sortedDomains {
			// http
			u, err := url.Parse("http://" + r)
			if err != nil {
				if verboseLogging {
					log.Println(err)
				}
				continue
			}

			re := fetchStatus(u.String())
			ress = append(ress, re)

			if checkHTTPS && re.StatusCode != 0 {
				u.Scheme = "https"
				re := fetchStatus(u.String())
				ress = append(ress, re)
			}
		}
	} else {
		ress = convertToResults(results)
	}

	if outputFile == "" {
		for _, r := range ress {
			if r.StatusCode == 0 || r.Error != nil {
				continue
			}
			log.Printf("| %d | %s\n", r.StatusCode, r.Domain)
		}
		return nil
	}

	return writeFile(outputFile, ress, disableStatusCheck)
}

func writeFile(outputFile string, res []Result, disableStatusCheck bool) error {
	resBytes := []byte("domain, status code")
	for _, r := range res {
		if r.StatusCode == 0 || r.Error != nil {
			continue
		}
		if disableStatusCheck {
			resBytes = append(resBytes, []byte(fmt.Sprintf("%s\n", r.Domain))...)
		} else {
			resBytes = append(resBytes, []byte(fmt.Sprintf("%s, %d\n", r.Domain, r.StatusCode))...)
		}
	}

	err := ioutil.WriteFile(outputFile, resBytes, 0775)
	if err != nil {
		return err
	}
	return nil
}

func sortDomains(domains []string) []string {
	encountered := map[string]string{}
	fixedList := []string{}
	for _, d := range domains {
		if strings.Contains(d, "\n") {
			for _, d := range strings.Split(d, "\n") {
				d = cleanDomain(d)
				fixedList = append(fixedList, d)
			}
		} else {
			d = cleanDomain(d)
			fixedList = append(fixedList, d)
		}
	}

	// dedup
	newList := []string{}
	for _, d := range fixedList {
		_, ok := encountered[d]
		if !ok {
			encountered[d] = d
			newList = append(newList, d)
		}
	}
	return newList
}

func convertToResults(strRes []string) []Result {
	ress := []Result{}
	for _, r := range strRes {
		ress = append(ress, Result{
			Domain: r,
		})
	}
	return ress
}

func cleanDomain(d string) string {
	d = strings.ToLower(d)
	if len(d) < 2 {
		return d
	}
	if d[0] == '*' || d[0] == '%' {
		d = d[1:]
	}

	if d[0] == '.' {
		d = d[1:]
	}
	return d
}
