package response_types

import "strings"

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

func (r *AuthorizeResponse) buildParamsString(params map[string]string) string {
	pairs := make([]string, 0, len(params))
	for key, value := range params {
		if value != "" {
			pairs = append(pairs, key+"="+value)
		}
	}
	return strings.Join(pairs, "&")
}
