package lnd

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	Host     string `json:"Host"`
	Cert     string	`json:"Cert"`
	Macaroon string `json:"Macaroon"`
}

func (c Client) BaseURL() string {
	if strings.HasPrefix(c.Host, "https") {
		return c.Host
	}
	return "https://" + c.Host
}

func (c *Client) ConfigFile(path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var client Client
	json.Unmarshal(file, &client)

	c.Host = client.Host
	c.Cert = client.Cert
	c.Macaroon = client.Macaroon
}

func (c Client) GetTlsCert() []byte {
	file, err := ioutil.ReadFile(c.Cert)
	if err != nil {
		panic(err)
	}
	return file
}

func (c Client) GetMacaroon() string {
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
	req, err := http.NewRequest(method, c.BaseURL()+"/"+path, buf)
	if err != nil {
		return gjson.Result{}, err
	}
	req.Header.Set("Grpc-Metadata-macaroon", c.GetMacaroon())

	tlsCert := x509.NewCertPool()
	tlsCert.AppendCertsFromPEM(c.GetTlsCert())

	tlsClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: tlsCert,
			},
		},
	}

	res, err := tlsClient.Do(req)
	if err != nil {
		return gjson.Result{}, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return gjson.Result{}, err
	}
	return gjson.ParseBytes(data), nil
}

func (c Client) CreateHoldInvoice(value int, hash []byte, memo string) (gjson.Result, error) {
	data := map[string]interface{}{
		"value": value,
		"hash": base64.StdEncoding.EncodeToString(hash),
		"memo": memo,
	}
	res, err := c.Call("POST", "v2/invoices/hodl", data)
	if err != nil {
		return gjson.Result{}, err
	}
	return res, nil
}

func (c Client) CreateInvoice(value int, memo string) (gjson.Result, error){
	data := map[string]interface{}{"value": value, "memo": memo}
	res, err := c.Call("POST", "v1/invoices",  data)
	if err != nil {
		return gjson.Result{}, err
	}
	return res, nil
}

func (c Client) CancelInvoice(hash []byte) (gjson.Result, error) {
	data := map[string]interface{}{"payment_hash": base64.StdEncoding.EncodeToString(hash)}
	res, err:= c.Call("POST", "v2/invoices/cancel", data)
	if err != nil {
		return gjson.Result{}, err
	}
	return res, nil
}

func (c Client) SettleInvoice(preimage []byte) (gjson.Result, error) {
	data := map[string]interface{}{"preimage": base64.StdEncoding.EncodeToString(preimage)}
	res, err := c.Call("POST", "v2/invoices/settle", data)
	if err != nil {
		return gjson.Result{}, err
	}
	return res, nil
}

func (c Client) LookupInvoice(hash string) (gjson.Result, error){
	res, err := c.Call("GET", "v2/invoices/subscribe/"+hash, nil)
	if err != nil {
		return gjson.Result{}, err
	}
	return res, nil
}

func (c Client) ListInvoices() (gjson.Result, error){
	res, err := c.Call("GET", "v1/invoices", nil)
	if err != nil {
		return gjson.Result{}, err
	}
	return res, nil
}

func (c Client) PayInvoice(invoice string, timeout int32) (gjson.Result, error) {
	data := map[string]interface{}{"payment_request": invoice, "timeout_seconds": timeout}
	res, err := c.Call("POST", "v2/router/send", data)
	if err != nil {
		return gjson.Result{}, err
	}
	return res, nil
}

func (c Client) DecodeInvoice(invoice string) (gjson.Result, error) {
	res, err := c.Call("GET", "v1/payreq/"+ invoice, nil)
	if err != nil {
		return gjson.Result{}, err
	}
	return res, nil
}
