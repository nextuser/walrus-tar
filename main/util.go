package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var debug_enabled = false

func debug(a ...any) {
	if debug_enabled {
		fmt.Println(a...)
	}
}

func debugf(f string, a ...any) {
	if debug_enabled {
		fmt.Printf(f, a...)
	}
}

func enableDebug(val bool) {
	debug_enabled = val
}

// 定义一个用来打印的函数，少写点代码，因为要处理很多次的 err
// 后面其他示例还会继续使用这个函数，就不单独再写，望看到此函数了解
func PrintError(err error) {
	if err != nil {
		log.Fatalln("error:", err)
	}
}

func getHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 600,
	}
}

/*
*
https client
*/
func getHttpsClient() *http.Client {
	var tlsConfig *tls.Config

	switch os.Getenv("AUTH_TYPE") {
	case "mtls":
		debug("AUTH_TYPE set as mtls")
		cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
		if err != nil {
			panic(err)
		}
		caCert, err := os.ReadFile("cacert.pem")
		if err != nil {
			panic(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}

	case "tls":
		debug("AUTH_TYPE set as tls")
		caCert, err := os.ReadFile("cacert.pem")
		if err != nil {
			panic(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig = &tls.Config{
			RootCAs: caCertPool,
		}

	default:
		debug("Insecure communication selected, skipping server verification")
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		MaxIdleConns:    10,
		IdleConnTimeout: 600 * time.Second,
	}
	return &http.Client{Transport: transport}
}
