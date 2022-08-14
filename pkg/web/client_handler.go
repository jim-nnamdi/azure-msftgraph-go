package web

import (
	"context"
	"net/http"
)

var _ WebClient = &webresponse{}

type webresponse struct {
}

func (wr *webresponse) ResponseUrlWithData(ctx context.Context, url string, postdata string, token string, method string) (*http.Response, error) {
}
