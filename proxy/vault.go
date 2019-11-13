package proxy

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

func (pc pClient) passToVault(path, body, method string, headers map[string][]string) (*resty.Response, error) {
	var (
		resp     *resty.Response
		err      error
		_headers map[string]string
	)

	_headers = make(map[string]string)

	req := pc.httpClient.R().EnableTrace()

	for h, v := range headers {
		_headers[h] = v[0]
	}

	req.SetHeaders(_headers)
	req.SetBody(body)

	switch method {
	case "GET":
		resp, err = req.Get(path)
	case "PUT":
		resp, err = req.Put(path)
	case "DELETE":
		resp, err = req.Delete(path)
	case "HEAD":
		resp, err = req.Head(path)
	case "OPTIONS":
		resp, err = req.Options(path)
	case "POST":
		resp, err = req.Post(path)
	case "PATCH":
		resp, err = req.Patch(path)
	default:
		return nil, errors.New("Dafuq mate")
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}
