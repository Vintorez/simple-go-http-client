package client

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type IHttpClient interface {
	GET(urlStr string, v interface{}) error
	PUT(urlStr string, body io.Reader, v interface{}) error
	POST(urlStr string, body io.Reader, v interface{}) error
	DELETE(urlStr string, body io.Reader) error
}

type ILogger interface {
	Print(v ...interface{})
}

func NewHttpClient(host, baseUrlPath, username, password string, log ILogger, insecure bool, localCertFile, userAgent string, httpClient *http.Client) (IHttpClient, error) {
	if httpClient == nil {
		tlsConfig := &tls.Config{InsecureSkipVerify: insecure}

		if !insecure {
			rootCAs, _ := x509.SystemCertPool()
			if rootCAs == nil {
				rootCAs = x509.NewCertPool()
			}

			if len(localCertFile) != 0 {
				certs, err := ioutil.ReadFile(localCertFile)
				if err != nil {
					return nil, err
				}

				if ok := rootCAs.AppendCertsFromPEM(certs); !ok && log != nil {
					log.Print("No certs appended, using system certs only")
				}
			}

			tlsConfig.RootCAs = rootCAs
		}

		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}
	}

	hostUrl, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	return &client{
		log:         log,
		httpClient:  httpClient,
		hostUrl:     hostUrl,
		baseUrlPath: baseUrlPath,
		userAgent:   userAgent,
		username:    username,
		password:    password,
	}, nil
}
