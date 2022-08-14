package web

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
)

var _ WebClient = &webresponse{}

var (
	ErrMakingRequest = "could not make request successfully"
)

type webresponse struct {
	httpClient *http.Client
}

func (wr *webresponse) ResponseUrlWithData(ctx context.Context, url string, postdata string, token string, method string) (*http.Response, error) {
	request_data, err := http.NewRequest(method, url, strings.NewReader(postdata))
	if err != nil {
		log.Print(err.Error())
		return nil, errors.New(ErrMakingRequest)
	}
	request_data.Header.Set("Content-Type", "application/json")
	request_data.Header.Set("Authorization", "Bearer "+token)
	response, err := wr.httpClient.Do(request_data)
	if err != nil {
		log.Print(err)
	}
	return response, nil
}
