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
	debug := flag.Bool("debug", false, "enable debug output")
	flag.Parse()

	if _, err := os.Stat(targetPath); err != nil {
		out, err := os.Create(targetPath)
		if err != nil {
			if *debug {
				fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
			}
			return
		}
		defer out.Close()

		resp, err := http.Get(geoliteURL)
		if err != nil {
			if *debug {
				fmt.Fprintf(os.Stderr, "Error downloading file: %v\n", err)
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			if *debug {
				fmt.Fprintf(os.Stderr, "Bad status: %s\n", resp.Status)
			}
			return
		}

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			if *debug {
				fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			}
			return
		}

		if *debug {
			fmt.Fprintln(os.Stderr, "Successfully downloaded GeoLite2-Country.mmdb to /tmp")
		}
	}

	db, err := geoip2.Open(targetPath)
	if err != nil {
		if *debug {
			fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		}
		return
	}
	defer db.Close()

	stat, _ := os.Stdin.Stat()
	var ips []string

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ip := strings.TrimSpace(scanner.Text())
			if ip != "" {
				ips = append(ips, ip)
			}
		}
		if err := scanner.Err(); err != nil && *debug {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		}
	} else {
		ips = flag.Args()
		if len(ips) == 0 {
			if *debug {
				fmt.Fprintln(os.Stderr, "No IP addresses provided via stdin or arguments")
			}
			return
		}
	}

	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			if *debug {
				fmt.Fprintf(os.Stderr, "Invalid IP address: %s\n", ipStr)
			}
			continue
		}

		record, err := db.Country(ip)
		if err != nil {
			if *debug {
				fmt.Fprintf(os.Stderr, "Error looking up IP %s: %v\n", ipStr, err)
			}
			continue
		}

		if *debug {
			fmt.Fprintf(os.Stderr, "Debug - Record: %+v\n", record)
		}

		countryCode := record.Country.IsoCode
		if countryCode == "" {
			countryCode = "Unknown"
		}

		fmt.Printf("%s - %s\n", ipStr, countryCode)
	}
}
