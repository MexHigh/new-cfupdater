# New CFUpdater

Docker image to keep your A and AAAA records in Cloudflare up to date. Just supply an API Token and the hosts you want to keep track of. Supports multiple zones.

This is the successor of https://git.leon.wtf/leon/cfupdater written in Go (was Python).

## Third-party services

The CFUpdater uses [Cloudflares trace service](https://www.cloudflare.com/cdn-cgi/trace) to determine the external IPv4 and IPv6.

**IPv4:**
- `https://1.1.1.1/cdn-cgi/trace`
- `https://1.0.0.1/cdn-cgi/trace` (fallback on timeout)

**IPv6:**
- `https://[2606:4700:4700::1111]/cdn-cgi/trace`
- `https://[2606:4700:4700::1001]/cdn-cgi/trace` (fallback on timeout)

Only the `ip=` field will be extracted from the response.