package main

import (
	"flag"
	"log"

	"github.com/cubixle/assetfinder"
)

func main() {
	// subsOnly := flag.Bool("subs-only", false, "Only search for subdomains.")
	disableStatusCheck := flag.Bool("status", false, "Enable/Disable checking the status code for each sub domain found.")
	checkHTTPS := flag.Bool("https", true, "Enable/Disable the checking of status of https on a domain.")
	outfile := flag.String("output", "", "Where to save results to.")
	verbose := flag.Bool("verbose", false, "Turns on verbose logging. Mainly used for debugging development.")

	flag.Parse()

	domain := flag.Arg(0)
	if domain == "" {
		log.Fatal("no domain given")
	}
	if err := assetfinder.Scanner(*verbose, *checkHTTPS, *disableStatusCheck, *outfile, domain); err != nil {
		log.Fatal(err)
	}
}
