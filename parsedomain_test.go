package parsedomain

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testParseDomain(t *testing.T, hostname string, expected *Host, expectedErr error) {
	d, err := Parse(hostname)
	assert.Equal(t, expected, d)
	assert.ErrorIs(t, err, expectedErr)
}

func TestParseDomain(t *testing.T) {
	testParseDomain(t, "amazon.co.uk", &Host{Hostname: "amazon.co.uk", Domain: "amazon", Management: ManagementICANN, Subdomains: []string{}, TopLevelDomains: []string{"co", "uk"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "books.amazon.co.uk", &Host{Hostname: "books.amazon.co.uk", Domain: "amazon", Management: ManagementICANN, Subdomains: []string{"books"}, TopLevelDomains: []string{"co", "uk"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "www.books.amazon.co.uk", &Host{Hostname: "www.books.amazon.co.uk", Domain: "amazon", Management: ManagementICANN, Subdomains: []string{"www", "books"}, TopLevelDomains: []string{"co", "uk"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "amazon.com", &Host{Hostname: "amazon.com", Domain: "amazon", Management: ManagementICANN, Subdomains: []string{}, TopLevelDomains: []string{"com"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "example0.debian.net", &Host{Hostname: "example0.debian.net", Domain: "example0", Management: ManagementPrivate, Subdomains: []string{}, TopLevelDomains: []string{"debian", "net"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "example1.debian.org", &Host{Hostname: "example1.debian.org", Domain: "debian", Management: ManagementICANN, Subdomains: []string{"example1"}, TopLevelDomains: []string{"org"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "golang.dev", &Host{Hostname: "golang.dev", Domain: "golang", Management: ManagementICANN, Subdomains: []string{}, TopLevelDomains: []string{"dev"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "golang.net", &Host{Hostname: "golang.net", Domain: "golang", Management: ManagementICANN, Subdomains: []string{}, TopLevelDomains: []string{"net"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "play.golang.org", &Host{Hostname: "play.golang.org", Domain: "golang", Management: ManagementICANN, Subdomains: []string{"play"}, TopLevelDomains: []string{"org"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "gophers.in.space.museum", &Host{Hostname: "gophers.in.space.museum", Domain: "in", Management: ManagementICANN, Subdomains: []string{"gophers"}, TopLevelDomains: []string{"space", "museum"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "0emm.com", &Host{Hostname: "0emm.com", Domain: "0emm", Management: ManagementICANN, Subdomains: []string{}, TopLevelDomains: []string{"com"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "a.0emm.com", &Host{Hostname: "a.0emm.com", Domain: "0emm", Management: ManagementPrivate, Subdomains: []string{"a"}, TopLevelDomains: []string{"com"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "b.c.d.0emm.com", &Host{Hostname: "b.c.d.0emm.com", Domain: "c", Management: ManagementPrivate, Subdomains: []string{"b"}, TopLevelDomains: []string{"d", "0emm", "com"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "there.is.no.such-tld", nil, ErrUnmanaged)
	testParseDomain(t, "foo.org", &Host{Hostname: "foo.org", Domain: "foo", Management: ManagementICANN, Subdomains: []string{}, TopLevelDomains: []string{"org"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "foo.co.uk", &Host{Hostname: "foo.co.uk", Domain: "foo", Management: ManagementICANN, Subdomains: []string{}, TopLevelDomains: []string{"co", "uk"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "foo.dyndns.org", &Host{Hostname: "foo.dyndns.org", Domain: "foo", Management: ManagementPrivate, Subdomains: []string{}, TopLevelDomains: []string{"dyndns", "org"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "foo.blogspot.co.uk", &Host{Hostname: "foo.blogspot.co.uk", Domain: "foo", Management: ManagementPrivate, Subdomains: []string{}, TopLevelDomains: []string{"blogspot", "co", "uk"}, Type: HostnameTypeDomain}, nil)
	testParseDomain(t, "cromulent", nil, ErrUnmanaged)
	testParseDomain(t, "127.0.0.1", &Host{Hostname: "127.0.0.1", Domain: "", Management: "", Subdomains: nil, TopLevelDomains: nil, Type: HostnameTypeIP}, nil)
	testParseDomain(t, "0000:0000:0000:0000:0000:0000:0000:0001", &Host{Hostname: "0000:0000:0000:0000:0000:0000:0000:0001", Domain: "", Management: "", Subdomains: nil, TopLevelDomains: nil, Type: HostnameTypeIP}, nil)
	testParseDomain(t, "::1", &Host{Hostname: "::1", Domain: "", Management: "", Subdomains: nil, TopLevelDomains: nil, Type: HostnameTypeIP}, nil)

	testParseDomain(t, "", nil, ErrReservedDomain)
	testParseDomain(t, "localhost", nil, ErrReservedDomain)
	testParseDomain(t, "foo.example.com.", nil, ErrUnmanaged)
	testParseDomain(t, strings.Repeat("verylongdomain", 20), nil, ErrDomainTooLong)
	testParseDomain(t, "invalidæcharacter.example.org", nil, ErrLabelInvalidCharacter)
	testParseDomain(t, "-label.example.org", nil, ErrLabelStartsWithDash)
	testParseDomain(t, "label-.example.org", nil, ErrLabelEndsWithDash)
	testParseDomain(t, ".example.org", nil, ErrLabelTooShort)
	testParseDomain(t, "extremelylongsubdomainiswaytoolonglongerthanthemaximumlabellength.example.org", nil, ErrLabelTooLong)
	testParseDomain(t, "label.example.123", nil, ErrLastLabelNumeric)
	testParseDomain(t, "label.example.123a", nil, ErrUnmanaged)
}

func TestLabelErrors(t *testing.T) {
	err := &LabelErrors{Errors: []error{
		fmt.Errorf("%w \"%s\" (length %d)", ErrLabelTooLong, "label", 64),
		fmt.Errorf("%w (\"%s\")", ErrLabelEndsWithDash, "label-"),
	}}

	assert.ErrorIs(t, err, ErrLabelTooLong)
	assert.ErrorIs(t, err, ErrLabelEndsWithDash)
	assert.NotErrorIs(t, err, ErrLastLabelNumeric)
	assert.Equal(t, "Label error(s):\n\t- Label is too long. Maximum length is 63, got \"label\" (length 64)\n\t- Labels cannot end with a dash (\"label-\")", err.Error())
}
