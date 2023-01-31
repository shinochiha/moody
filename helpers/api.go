package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/moody/config"
)

type Map = map[string]interface{}

func CallAPI(method, url string, data interface{}, header map[string]string) (int, interface{}) {
	reqBody, _ := json.Marshal(data)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	for k, v := range header {
		req.Header.Add(k, v)
	}

	if header["Content-Type"] == "" {
		req.Header.Add("Content-Type", "application/json")
	}

	// TODO: change this with custom client and transport.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("failed request to %s: %s", url, err.Error())
		return http.StatusInternalServerError, err
	}

	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)

	if config.Get("API_IS_DEBUG").Bool() {
		fmt.Println("--------------------")
		fmt.Println(method, url)
		fmt.Println(string(reqBody))
		fmt.Println("--------------------")
		fmt.Println(string(resBody))
	}

	if res.StatusCode < http.StatusOK && res.StatusCode >= http.StatusBadRequest {
		fmt.Printf("%d: failed request to %s: %s", res.StatusCode, url, string(resBody))
	}

	if res.StatusCode == http.StatusOK &&
		res.Header.Get("Content-Type") != "application/json" {
		// Special case For Woo-WA, since it does not return json format,
		// just single string status in the response body.
		return res.StatusCode, string(resBody)
	}

	// TODO: change the way we parsing response, let the caller decide what it expects from calling this API.
	// the response struct defined by the caller of this func.
	// one parameter possibly added to this func, something like `response interface{}`
	var d interface{}
	json.Unmarshal(resBody, &d)
	return res.StatusCode, d
}

func IsErrorHttpCode(code int) bool {
	if code < 200 || code >= 400 {
		return true
	}
	return false
}
