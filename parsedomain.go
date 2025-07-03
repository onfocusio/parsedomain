package parsedomain

import (
	"fmt"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type HostnameType uint

type Management string

const (
	HostnameTypeInvalid HostnameType = iota
	HostnameTypeIP
	HostnameTypeDomain

	ManagementUnmanaged Management = "Unmanaged"
	ManagementICANN     Management = "ICANN Managed"
	ManagementPrivate   Management = "Privately Managed"
)

type Host struct {
	Hostname        string
	Domain          string
	Management      Management
	Subdomains      []string
	TopLevelDomains []string
	Type            HostnameType
}

func Parse(hostname string) (*Host, error) {
	hostnameType, labels, err := sanitize(hostname)
	if err != nil {
		return nil, err
	}

	if hostnameType == HostnameTypeIP {
		// Input "hostname" can be safely parsed by net.ParseIP.
		// To know if it's an IPv4, use `ip.To4() != nil`
		return &Host{
			Hostname: hostname,
			Type:     hostnameType,
		}, nil
	}

	eTLD, icann := publicsuffix.PublicSuffix(hostname)

	management := ManagementUnmanaged
	if icann {
		management = ManagementICANN
	} else if strings.IndexByte(eTLD, '.') >= 0 {
		management = ManagementPrivate
	}

	if management == ManagementUnmanaged {
		return nil, fmt.Errorf("%w: \"%s\"", ErrUnmanaged, hostname)
	}

	topLevelDomains := strings.Split(eTLD, labelSeparator)
	i := len(labels) - len(topLevelDomains) - 1

	if i == -1 {
		i = 1
	}

	return &Host{
		Hostname:        hostname,
		Domain:          labels[i],
		Management:      management,
		Subdomains:      labels[:i],
		TopLevelDomains: labels[i+1:],
		Type:            hostnameType,
	}, nil
}
