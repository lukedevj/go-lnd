package lnd

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

type Client struct {
	Host     string
	Cert     string
	Macaroon string
}

func (c Client) Call(method string, path string, body map[string]interface{}) (gjson.Result, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return gjson.Result{}, err
	}

	url := c.Host
	if !strings.HasPrefix(c.Host, "https") {
		url = "https://" + url + "/"
	}
	url += path

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return gjson.Result{}, err
	}

	file, err := ioutil.ReadFile(c.Macaroon)
	if err != nil {
		return gjson.Result{}, err
	}
	req.Header.Set("Grpc-Metadata-macaroon", hex.EncodeToString(file))
	file, err = ioutil.ReadFile(c.Cert)
	if err != nil {
		return gjson.Result{}, err
	}

	tls_cert := x509.NewCertPool()
	tls_cert.AppendCertsFromPEM(file)

	http_client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: tls_cert,
			},
		},
	}
	res, err := http_client.Do(req)
	if err != nil {
		return gjson.Result{}, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))
	if err != nil {
		return gjson.Result{}, err
	}
	return gjson.ParseBytes(b), nil
}
