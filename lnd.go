package lnd

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

func (c Client) baseURL() string {
	if strings.HasPrefix(c.Host, "https") {
		return c.Host
	}
	return "https://" + c.Host
}

func (c Client) getTlsCert() []byte {
	file, err := ioutil.ReadFile(c.Cert)
	if err != nil {
		panic(err)
	}
	return file
}

func (c Client) getMacaroon() string {
	file, err := ioutil.ReadFile(c.Macaroon)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(file)
}

func (c Client) Call(method string, path string, body map[string]interface{}) (gjson.Result, error) {
	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return gjson.Result{}, err
		}
	}
	req, err := http.NewRequest(method, c.baseURL()+"/"+path, buf)
	if err != nil {
		return gjson.Result{}, err
	}
	req.Header.Set("Grpc-Metadata-macaroon", c.getMacaroon())

	tlscert := x509.NewCertPool()
	tlscert.AppendCertsFromPEM(c.getTlsCert())

	tls_client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: tlscert,
			},
		},
	}
	res, err := tls_client.Do(req)
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
