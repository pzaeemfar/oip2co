package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

const (
	geoliteURL = "https://github.com/PrxyHunter/GeoLite2/releases/latest/download/GeoLite2-Country.mmdb"
	targetPath = "/tmp/GeoLite2-Country.mmdb"
)

func main() {
	silent := flag.Bool("s", true, "silent mode - suppress all output (default: true)")
	flag.Parse()

	// Check if file already exists
	if _, err := os.Stat(targetPath); err != nil {
		// Create the file
		out, err := os.Create(targetPath)
		if err != nil {
			if !*silent {
				fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
			}
			return
		}
		defer out.Close()

		// Download the file
		resp, err := http.Get(geoliteURL)
		if err != nil {
			if !*silent {
				fmt.Fprintf(os.Stderr, "Error downloading file: %v\n", err)
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if !*silent {
				fmt.Fprintf(os.Stderr, "Bad status: %s\n", resp.Status)
			}
			return
		}

		// Copy the response body to the file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			if !*silent {
				fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			}
			return
		}

		if !*silent {
			fmt.Fprintln(os.Stderr, "Successfully downloaded GeoLite2-Country.mmdb to /tmp")
		}
	}

	// Open the GeoLite2 database
	db, err := geoip2.Open(targetPath)
	if err != nil {
		if !*silent {
			fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		}
		return
	}
	defer db.Close()

	// Read IPs from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ipStr := strings.TrimSpace(scanner.Text())
		if ipStr == "" {
			continue
		}

		ip := net.ParseIP(ipStr)
		if ip == nil {
			if !*silent {
				fmt.Fprintf(os.Stderr, "Invalid IP address: %s\n", ipStr)
			}
			continue
		}

		record, err := db.Country(ip)
		if err != nil {
			if !*silent {
				fmt.Fprintf(os.Stderr, "Error looking up IP %s: %v\n", ipStr, err)
			}
			continue
		}

		if !*silent {
			fmt.Fprintf(os.Stderr, "Debug - Record: %+v\n", record)
		}

		countryCode := record.Country.IsoCode
		if countryCode == "" {
			countryCode = "Unknown"
		}

		fmt.Printf("%s - %s\n", ipStr, countryCode)
	}

	if err := scanner.Err(); err != nil {
		if !*silent {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		}
	}
} 