# oip2co

A simple command-line tool that converts IP addresses to country codes. It uses the MaxMind GeoLite2 database for accurate IP-to-country lookups.

## Features

- Fast IP-to-country lookups
- Reads IPs from stdin
- Clean, simple output format
- Automatic database download
- Silent mode by default

## Installation

```bash
go install github.com/pzaeemfar/oip2co@latest
```

### Options

- `-s=false`: Show debug messages (default is silent mode)

## Examples

```bash
# Look up Google's DNS
echo "8.8.8.8" | oip2co
# Output: 8.8.8.8 - US

# Look up a private IP
echo "192.168.1.1" | oip2co
# Output: 192.168.1.1 - Unknown

# Look up an IPv6 address
echo "2001:4860:4860::8888" | oip2co
# Output: 2001:4860:4860::8888 - US
```

## Notes

- The program automatically downloads and updates the GeoLite2 database
- Invalid IP addresses are skipped
- Private IP ranges return "Unknown" as country code
- The database is stored in `/tmp/GeoLite2-Country.mmdb`