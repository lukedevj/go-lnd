package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
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

func (c Client) parseUrl(path string) string {
	if !strings.HasPrefix(c.Host, "https") {
		return "https://" + c.Host + "/" + path
	}
	return c.Host + "/" + path
}

func (c Client) tlsCertRead() []byte {
	file, err := ioutil.ReadFile(c.Cert)
	if err != nil {
		panic(err)
	}
	return file
}

func (c Client) macaroonRead() string {
	file, err := ioutil.ReadFile(c.Macaroon)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(file)
}

func (c Client) Call(method string, path string, body map[string]interface{}) (gjson.Result, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		panic(err)
	}

	url := c.parseUrl(path)
	req, err := http.NewRequest(method, url, buf)
	req.Header.Set("Grpc-Metadata-macaroon", c.macaroonRead())
	if err != nil {
		panic(err)
	}

	cert_pool := x509.NewCertPool()
	cert_pool.AppendCertsFromPEM(c.tlsCertRead())
	http_client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: cert_pool,
			},
		},
	}
	res, err := http_client.Do(req)
	if err != nil {
		return gjson.Result{}, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return gjson.Result{}, err
	}
	return gjson.ParseBytes(b), nil
}
