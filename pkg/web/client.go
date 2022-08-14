package web

import (
	"context"
	"net/http"
)

type WebClient interface {
	ResponseUrlWithData(ctx context.Context, url string, postdata string, token string, method string) (*http.Response, error)
}
