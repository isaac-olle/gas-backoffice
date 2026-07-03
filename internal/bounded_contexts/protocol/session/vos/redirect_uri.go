package vos

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

type RedirectURI struct{ value vo.Url }

func NewRedirectURI(v string) (RedirectURI, error) {
	url, err := vo.NewUrl(v, vo.NewFragmentNotAllowedOption())
	if err != nil {
		return RedirectURI{}, err
	}
	return RedirectURI{value: url}, nil
}

func (r RedirectURI) Value() string  { return r.value.String() }
func (r RedirectURI) String() string { return r.value.String() }
