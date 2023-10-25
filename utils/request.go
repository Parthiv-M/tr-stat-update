package utils

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type CustomRequest struct {
	method  string
	url     string
	params  map[string]string
	headers map[string]string
	cookie  http.Cookie
}

type CustomResponse struct {
	resBody   []byte
	resHeader http.Header
	status    int
	isError   bool
}

func (cres CustomResponse) GetResponseBody() []byte {
	return cres.resBody
}

func (cres CustomResponse) GetResponseHeaders() http.Header {
	return cres.resHeader
}

func (cres CustomResponse) GetResponseStatusCode() int {
	return cres.status
}

func CreateCustomRequest(url string, method string, params map[string]string, headers map[string]string, cookie *http.Cookie) CustomRequest {
	return CustomRequest{
		url:     url,
		method:  method,
		params:  params,
		headers: headers,
		cookie:  *cookie,
	}
}

func (creq CustomRequest) MakeRequest() CustomResponse {
	// create a new HTTP client and a new request
	jar, err := cookiejar.New(nil)
	client := http.Client{Jar: jar}
	data := url.Values{}

	// add query params to request
	for key, value := range creq.params {
		data.Add(key, value)
	}

	req, err := http.NewRequest(creq.method, creq.url, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err.Error())
	}

	// add headers to request
	for key, value := range creq.headers {
		http.Header.Add(req.Header, key, value)
	}

	// add cookie to request
	req.AddCookie(&creq.cookie)

	// make the API call
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()

	// convert the response body to []byte
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	return CustomResponse{
		resBody:   body,
		resHeader: res.Header,
		status:    res.StatusCode,
		isError:   res.StatusCode == 200,
	}
}
