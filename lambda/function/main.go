package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	h "github.com/liftplan/liftplan/serve/handler"
)

var (
	badRequest = events.APIGatewayV2HTTPResponse{Body: http.StatusText(http.StatusBadRequest), StatusCode: http.StatusBadRequest}
	rootRoute  = h.Root()
	planRoute  = h.Plan()
)

func eventToRequest(request events.APIGatewayV2HTTPRequest) (*http.Request, error) {
	u := url.URL{
		Path:     request.RequestContext.HTTP.Path,
		RawQuery: request.RawQueryString,
	}

	method := request.RequestContext.HTTP.Method

	r, err := http.NewRequest(method, u.String(), strings.NewReader(request.Body))
	for k, v := range request.Headers {
		r.Header.Add(k, v)
	}
	return r, err
}

func newResponse() response {
	return response{
		PR: &events.APIGatewayV2HTTPResponse{
			MultiValueHeaders: make(map[string][]string),
		},
	}
}

type response struct {
	PR *events.APIGatewayV2HTTPResponse
}

func (r response) Header() http.Header {
	return r.PR.MultiValueHeaders
}

func (r response) Write(d []byte) (int, error) {
	var b strings.Builder
	b.WriteString(r.PR.Body)
	i, err := b.Write(d)
	if err == nil {
		r.PR.Body = b.String()
	}
	return i, err
}

func (r response) WriteHeader(statusCode int) {
	r.PR.StatusCode = statusCode
}

func newRequestResponse(request events.APIGatewayV2HTTPRequest) (response, *http.Request, error) {
	w := newResponse()
	r, err := eventToRequest(request)
	return w, r, err
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	w, r, err := newRequestResponse(request)
	if err != nil {
		log.Println(err)
		return badRequest, err
	}
	switch request.RequestContext.HTTP.Path {
	case "/":
		rootRoute(w, r)
	case "/plan":
		planRoute(w, r)
	}
	log.Println(w.PR.StatusCode)
	return events.APIGatewayV2HTTPResponse{
		StatusCode:        200,
		Body:              w.PR.Body,
		MultiValueHeaders: w.PR.MultiValueHeaders,
	}, nil
	// return *w.PR, nil
}

func main() {
	lambda.Start(handler)
}
