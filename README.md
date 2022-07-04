# parsedomain

[![Version](https://img.shields.io/github/v/release/onfocusio/parsedomain?include_prereleases)](https://github.com/onfocusio/parsedomain/releases)
[![Build Status](https://github.com/onfocusio/parsedomain/workflows/Test/badge.svg)](https://github.com/onfocusio/parsedomain/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/onfocusio/parsedomain.svg)](https://pkg.go.dev/github.com/onfocusio/parsedomain)

Hostname parsing library inspired by the JS library [parse-domain](https://github.com/peerigon/parse-domain).

It is not always easy to split a hostname into subdomains, domains and top-level domains because domain name registrar organize their namespaces in different ways. This library uses [golang.org/x/net/publicsuffix](https://pkg.go.dev/golang.org/x/net/publicsuffix) along with additional validation checks to parse a given hostname and give you the following information:
- The domain
- A `[]string` of subdomains
- A `[]string` of **effective** top-level domains
- The type of Management*
    - **Unmanaged** (the hostname is valid but is not managed privately nor by ICANN, typically a custom entry in the hostfile of your computer)
    - **[ICANN](https://www.icann.org/) managed**
    - **Privately managed**

The library supports IP as hostnames and won't return an error if the IP is valid. You should always check the returned host type.

**Example**
```go
import (
	"fmt"

	"github.com/onfocusio/parsedomain"
)

func main() {
    host, err := parsedomain.Parse("books.amazon.co.uk")
    if err != nil {
        fmt.Printf("Invalid hostname \"%s\". Error: %s\n", hostname, err.Error())
        return
    }
    fmt.Println(host.Hostname)        // "books.amazon.co.uk"
    fmt.Println(host.Subdomains)      // ["books"]
    fmt.Println(host.Domain)          // "amazon"
    fmt.Println(host.TopLevelDomains) // ["co", "uk"]
    fmt.Println(host.Type)            // 2 (parsedomain.HostnameTypeDomain)
    fmt.Println(host.Management)      // "ICANN Managed" (parsedomain.ManagementICANN)
}
```

<details>
    <summary><b>More examples</b></summary>
  
```go
package main

import (
    "fmt"

    "github.com/onfocusio/parsedomain"
)

func main() {
    hosts := []string{
        "amazon.co.uk",
        "books.amazon.co.uk",
        "www.books.amazon.co.uk",
        "amazon.com",
        "example0.debian.net",
        "example1.debian.org",
        "golang.dev",
        "golang.net",
        "play.golang.org",
        "gophers.in.space.museum",
        "0emm.com",
        "a.0emm.com",
        "b.c.d.0emm.com",
        "there.is.no.such-tld",
        "foo.org",
        "foo.co.uk",
        "foo.dyndns.org",
        "foo.blogspot.co.uk",
        "cromulent",
        "foo.example.com.",
        "127.0.0.1",
        "0000:0000:0000:0000:0000:0000:0000:0001",
        "localhost",
        "",
    }

    for _, hostname := range hosts {
        host, err := parsedomain.Parse(hostname)
        if err != nil {
            fmt.Printf("Invalid hostname \"%s\". Error: %s\n", hostname, err.Error())
            continue
        }
        fmt.Printf("%#v\n", host)
    }
}
```
```
&parsedomain.Host{Hostname:"amazon.co.uk", Domain:"amazon", Management:"ICANN Managed", Subdomains:[]string{}, TopLevelDomains:[]string{"co", "uk"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"books.amazon.co.uk", Domain:"amazon", Management:"ICANN Managed", Subdomains:[]string{"books"}, TopLevelDomains:[]string{"co", "uk"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"www.books.amazon.co.uk", Domain:"amazon", Management:"ICANN Managed", Subdomains:[]string{"www", "books"}, TopLevelDomains:[]string{"co", "uk"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"amazon.com", Domain:"amazon", Management:"ICANN Managed", Subdomains:[]string{}, TopLevelDomains:[]string{"com"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"example0.debian.net", Domain:"example0", Management:"Privately Managed", Subdomains:[]string{}, TopLevelDomains:[]string{"debian", "net"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"example1.debian.org", Domain:"debian", Management:"ICANN Managed", Subdomains:[]string{"example1"}, TopLevelDomains:[]string{"org"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"golang.dev", Domain:"golang", Management:"ICANN Managed", Subdomains:[]string{}, TopLevelDomains:[]string{"dev"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"golang.net", Domain:"golang", Management:"ICANN Managed", Subdomains:[]string{}, TopLevelDomains:[]string{"net"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"play.golang.org", Domain:"golang", Management:"ICANN Managed", Subdomains:[]string{"play"}, TopLevelDomains:[]string{"org"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"gophers.in.space.museum", Domain:"in", Management:"ICANN Managed", Subdomains:[]string{"gophers"}, TopLevelDomains:[]string{"space", "museum"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"0emm.com", Domain:"0emm", Management:"ICANN Managed", Subdomains:[]string{}, TopLevelDomains:[]string{"com"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"a.0emm.com", Domain:"0emm", Management:"Privately Managed", Subdomains:[]string{"a"}, TopLevelDomains:[]string{"com"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"b.c.d.0emm.com", Domain:"c", Management:"Privately Managed", Subdomains:[]string{"b"}, TopLevelDomains:[]string{"d", "0emm", "com"}, Type:parsedomain.HostnameTypeDomain}
Invalid hostname "there.is.no.such-tld". Error: Unmanaged hostname: "there.is.no.such-tld"
&parsedomain.Host{Hostname:"foo.org", Domain:"foo", Management:"ICANN Managed", Subdomains:[]string{}, TopLevelDomains:[]string{"org"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"foo.co.uk", Domain:"foo", Management:"ICANN Managed", Subdomains:[]string{}, TopLevelDomains:[]string{"co", "uk"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"foo.dyndns.org", Domain:"foo", Management:"Privately Managed", Subdomains:[]string{}, TopLevelDomains:[]string{"dyndns", "org"}, Type:parsedomain.HostnameTypeDomain}
&parsedomain.Host{Hostname:"foo.blogspot.co.uk", Domain:"foo", Management:"Privately Managed", Subdomains:[]string{}, TopLevelDomains:[]string{"blogspot", "co", "uk"}, Type:parsedomain.HostnameTypeDomain}
Invalid hostname "cromulent". Error: Unmanaged hostname: "cromulent"
Invalid hostname "foo.example.com.". Error: Unmanaged hostname: "foo.example.com."
&parsedomain.Host{Hostname:"127.0.0.1", Domain:"", Management:"", Subdomains:[]string(nil), TopLevelDomains:[]string(nil), Type:parsedomain.HostnameTypeIP}
&parsedomain.Host{Hostname:"0000:0000:0000:0000:0000:0000:0000:0001", Domain:"", Management:"", Subdomains:[]string(nil), TopLevelDomains:[]string(nil), Type:parsedomain.HostnameTypeIP}
Invalid hostname "localhost". Error: Domain is reserved: "localhost"
Invalid hostname "". Error: Domain is reserved: ""
```
</details>
