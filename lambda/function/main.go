package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	h "github.com/liftplan/liftplan/serve/handler"
)

var (
	badRequest = events.APIGatewayProxyResponse{Body: http.StatusText(http.StatusBadRequest), StatusCode: http.StatusBadRequest}
	rootRoute  = h.Root()
	planRoute  = h.Plan()
)

func eventToRequest(request events.APIGatewayProxyRequest) (*http.Request, error) {
	v := url.Values(request.MultiValueQueryStringParameters)
	u := url.URL{
		Path:     request.Path,
		RawQuery: v.Encode(),
	}
	r, err := http.NewRequest(request.HTTPMethod, u.String(), strings.NewReader(request.Body))
	r.Header = request.MultiValueHeaders
	return r, err
}

func newProxyResponse() proxyResponse {
	return proxyResponse{
		PR: &events.APIGatewayProxyResponse{
			MultiValueHeaders: make(map[string][]string),
		},
	}
}

type proxyResponse struct {
	PR *events.APIGatewayProxyResponse
}

func (p proxyResponse) Header() http.Header {
	return p.PR.MultiValueHeaders
}

func (p proxyResponse) Write(d []byte) (int, error) {
	var b strings.Builder
	b.WriteString(p.PR.Body)
	i, err := b.Write(d)
	if err == nil {
		p.PR.Body = b.String()
	}
	return i, err
}

func (p proxyResponse) WriteHeader(statusCode int) {
	p.PR.StatusCode = statusCode
}

func newRequestResponse(request events.APIGatewayProxyRequest) (proxyResponse, *http.Request, error) {
	w := newProxyResponse()
	r, err := eventToRequest(request)
	return w, r, err
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	w, r, err := newRequestResponse(request)
	if err != nil {
		return badRequest, err
	}
	switch request.Path {
	case "/":
		rootRoute(w, r)
	case "/plan":
		planRoute(w, r)
	}
	return *w.PR, nil
}

func main() {
	lambda.Start(handler)
}
