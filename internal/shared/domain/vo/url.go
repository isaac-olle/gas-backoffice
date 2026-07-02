package vo

import (
	"fmt"
	"net/url"
	"strings"
)

type Url string

func NewUrl(raw string, opts ...UrlOption) (Url, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid url: %w", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("invalid url, scheme must be http or https")
	}

	if parsed.Host == "" {
		return "", fmt.Errorf("invalid url, host cannot be empty")
	}

	for _, opt := range opts {
		if err = opt.IsAllowed(parsed); err != nil {
			return "", err
		}
	}

	return Url(raw), nil
}

func (u Url) String() string {
	return string(u)
}

type UrlOption interface {
	IsAllowed(url *url.URL) error
}

type FragmentNotAllowedOption struct{}

func NewFragmentNotAllowedOption() *FragmentNotAllowedOption {
	return &FragmentNotAllowedOption{}
}

func (f *FragmentNotAllowedOption) IsAllowed(u *url.URL) error {
	if u.Fragment != "" {
		return fmt.Errorf("url cannot contain fragment")
	}
	return nil
}

type LocalHostOption struct{}

func NewLocalHostOption() *LocalHostOption {
	return &LocalHostOption{}
}

func (LocalHostOption) IsAllowed(u *url.URL) error {
	if strings.Contains(u.Host, "localhost") {
		return fmt.Errorf("url contains localhost")
	}
	return nil
}

type HttpsOnlyOption struct{}

func NewHttpsOnlyOption() *HttpsOnlyOption {
	return &HttpsOnlyOption{}
}

func (h *HttpsOnlyOption) IsAllowed(u *url.URL) error {
	if u.Scheme != "https" {
		return fmt.Errorf("url scheme must be https")
	}
	return nil
}
