package teapot

import (
	"net/http"

	"github.com/corpix/uarand"
)

const DefaultUserAgent = "teapot"

type RequestInterceptor func(request *http.Request) error
type ResponseInterceptor func(request *http.Response) error

func SetRandomUserAgent(req *http.Request) error {
	req.Header.Set("User-Agent", uarand.GetRandom())
	return nil
}
