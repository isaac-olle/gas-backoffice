package response_types

import (
	"net/url"
	"strings"
)

type AuthorizeResponse struct {
	baseURL        string
	queryParams    map[string]string
	fragmentParams map[string]string
}

func NewAuthorizeResponse(baseURL string) *AuthorizeResponse {
	return &AuthorizeResponse{
		baseURL:        baseURL,
		queryParams:    make(map[string]string),
		fragmentParams: make(map[string]string),
	}
}

func (r *AuthorizeResponse) AddQueryParameter(key string, value string) *AuthorizeResponse {
	r.queryParams[key] = value
	return r
}

func (r *AuthorizeResponse) AddFragmentParameter(key string, value string) *AuthorizeResponse {
	r.fragmentParams[key] = value
	return r
}

func (r *AuthorizeResponse) AddState(state string) *AuthorizeResponse {
	if len(r.fragmentParams) > 0 {
		r.fragmentParams["state"] = state
	} else {
		r.queryParams["state"] = state
	}
	return r
}

func (r *AuthorizeResponse) Build() string {
	result := r.baseURL

	if len(r.queryParams) > 0 {
		separator := "?"
		if strings.Contains(r.baseURL, "?") {
			separator = "&"
		}
		result += separator + r.buildParamsString(r.queryParams)
	}

	if len(r.fragmentParams) > 0 {
		result += "#" + r.buildParamsString(r.fragmentParams)
	}

	return result
}

// buildParamsString percent-encodes every pair. Concatenar-los crus corrompia
// qualsevol valor amb caràcters reservats — en particular un state en base64
// estàndard, on el '+' arribava al destí convertit en espai.
func (r *AuthorizeResponse) buildParamsString(params map[string]string) string {
	values := url.Values{}
	for key, value := range params {
		if value != "" {
			values.Set(key, value)
		}
	}
	return values.Encode()
}
