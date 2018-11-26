package adapters

import (
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/parnurzeal/gorequest"
)

//dump all http request and response to logger? defaults to true
var dump bool

func init() {
	dump = beego.AppConfig.DefaultBool("http.dump", true)
}

//IHTTPClient HTTP client adapter
type IHTTPClient interface {
	SetHeader(name string, value string)
	PostJSON(url string, json interface{}) (gorequest.Response, string, []error)
	PostFormData(url string, formdata string) (gorequest.Response, string, []error)
}

//HTTPClient struct extends gorequest.SuperAgent
type httpClient struct {
	*gorequest.SuperAgent
}

func (client *httpClient) SetHeader(name string, value string) {
	client.Set(name, value)
}

//PostJSON from a struct
func (client *httpClient) PostJSON(url string, json interface{}) (gorequest.Response, string, []error) {
	res, body, err := client.
		Post(url).
		Type(gorequest.TypeJSON).
		Send(json).
		End()
	return res, body, err
}

//PostFormData as url encoded string
func (client *httpClient) PostFormData(url string, formdata string) (gorequest.Response, string, []error) {
	res, body, err := client.
		Post(url).
		Type(gorequest.TypeUrlencoded).
		Send(formdata).
		End()
	return res, body, err
}

//NewHTTPClient factory method to initialize default HTTP service
func NewHTTPClient(timeoutms time.Duration, retry int) IHTTPClient {
	req := gorequest.New().
		SetDebug(dump).
		Timeout(timeoutms*time.Millisecond).
		Retry(retry, time.Second, http.StatusInternalServerError) //retry attempt, 1 second after failure

	return &httpClient{
		SuperAgent: req,
	}
}
