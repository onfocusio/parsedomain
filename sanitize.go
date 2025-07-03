package parsedomain

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

const (
	labelSeparator  = "."
	maxDomainLength = 253
	minLabelLength  = 1
	maxLabelLength  = 63
)

var (
	ipBracketsRegex             = regexp.MustCompile(`^\[|]$`)
	invalidDomainCharatersRegex = regexp.MustCompile(`(?i)[^\da-z-]`)
	lettersAndDashRegex         = regexp.MustCompile(`(?i)[a-z-]`)

	ErrDomainTooLong         = fmt.Errorf("domain is too long. Maximum length is %d, got", maxDomainLength)
	ErrReservedDomain        = fmt.Errorf("domain is reserved")
	ErrLabelInvalidCharacter = fmt.Errorf("invalid character(s) found in label")
	ErrLabelStartsWithDash   = fmt.Errorf("labels cannot start with a dash")
	ErrLabelEndsWithDash     = fmt.Errorf("labels cannot end with a dash")
	ErrLabelTooShort         = fmt.Errorf("label is too short. Minimum length is %d, got", minLabelLength)
	ErrLabelTooLong          = fmt.Errorf("label is too long. Maximum length is %d, got", maxLabelLength)
	ErrLastLabelNumeric      = fmt.Errorf("last label must not be all-numeric. Got")
	ErrUnmanaged             = fmt.Errorf("unmanaged hostname")

	reservedTopDomains = map[string]struct{}{
		"":          {},
		"localhost": {},
		"local":     {},
		"example":   {},
		"invalid":   {},
		"test":      {},
	}
)

type LabelErrors struct {
	Errors []error
}

func (e LabelErrors) Error() string {
	message := "Label error(s):"
	for _, e := range e.Errors {
		message += "\n\t- " + e.Error()
	}
	return message
}

func (e LabelErrors) Is(target error) bool {
	for _, suberror := range e.Errors {
		if errors.Is(suberror, target) {
			return true
		}
	}
	return false
}

func sanitize(hostname string) (HostnameType, []string, error) {
	if hostname == "" {
		return HostnameTypeDomain, []string{}, fmt.Errorf("%w: \"\"", ErrReservedDomain)
	}

	ipString := ipBracketsRegex.ReplaceAllString(hostname, "")
	ip := net.ParseIP(ipString)
	if ip != nil {
		return HostnameTypeIP, nil, nil
	}

	canonical := canonical(hostname)
	if len(canonical) > maxDomainLength {
		return HostnameTypeInvalid, nil, fmt.Errorf("%w \"%s\" (length %d)", ErrDomainTooLong, hostname, len(canonical))
	}

	labels := strings.Split(canonical, labelSeparator)
	labelErrors := validateLabelsStrict(labels)
	if labelErrors != nil {
		return HostnameTypeInvalid, nil, LabelErrors{labelErrors}
	}

	lastLabel := labels[len(labels)-1]
	if _, ok := reservedTopDomains[lastLabel]; ok {
		return HostnameTypeDomain, labels, fmt.Errorf("%w: \"%s\"", ErrReservedDomain, lastLabel)
	}

	return HostnameTypeDomain, labels, nil
}

func canonical(hostname string) string {
	lastChar := hostname[len(hostname)-1:]
	if lastChar == labelSeparator {
		return hostname[:len(hostname)-1]
	}

	return hostname
}

func validateLabelsStrict(labels []string) []error {
	errors := []error{}

	for _, label := range labels {
		if invalidDomainCharatersRegex.MatchString(label) {
			errors = append(errors, fmt.Errorf("%w \"%s\"", ErrLabelInvalidCharacter, label))
		}

		if strings.HasPrefix(label, "-") {
			errors = append(errors, fmt.Errorf("%w (\"%s\")", ErrLabelStartsWithDash, label))
		}
		if strings.HasSuffix(label, "-") {
			errors = append(errors, fmt.Errorf("%w (\"%s\")", ErrLabelEndsWithDash, label))
		}

		if len(label) < minLabelLength {
			errors = append(errors, fmt.Errorf("%w \"%s\" (length %d)", ErrLabelTooShort, label, len(label)))
		} else if len(label) > maxLabelLength {
			errors = append(errors, fmt.Errorf("%w \"%s\" (length %d)", ErrLabelTooLong, label, len(label)))
		}
	}

	if len(labels) > 0 {
		lastLabel := labels[len(labels)-1]
		if !lettersAndDashRegex.MatchString(lastLabel) {
			errors = append(errors, fmt.Errorf("%w \"%s\"", ErrLastLabelNumeric, lastLabel))
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}
