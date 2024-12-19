package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

var enableDebug = false

func debug(a ...any) {
	if enableDebug {
		fmt.Println(a...)
	}

}

var wc sync.WaitGroup

func getHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 600,
	}
}
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

func process(body io.Reader, out string) {

	var written int64 = 0
	var write_err error = nil
	if len(out) == 0 {
		written, write_err = io.Copy(os.Stdout, body)
	} else {
		var file, openErr = os.OpenFile(out, os.O_WRONLY|os.O_CREATE, 0664)
		if openErr != nil {
			debug(openErr)
			return
		}

		written, write_err = io.Copy(file, body)
	}

	if write_err != nil {
		fmt.Println("Write error:", write_err)
	} else {
		fmt.Printf("write in file %s, %d bytes", out, written)
	}
}
func read(c *http.Client, url string, out string) {
	fmt.Println("url=", url)
	defer wc.Done()

	var response, err = c.Get(url)
	if err != nil {
		debug(err)
	} else {
		if response.StatusCode == 200 {
			process(response.Body, out)
		}
		debug(response)
	}
}

type StoreEvent struct {
	TxDigest string `json:"txDigest"`
	EventSeq string `json:"eventSeq"`
}

type EventOrObject struct {
	Event StoreEvent `json:"Event"`
}
type Certified struct {
	BlobId        string        `json:"blobId"`
	EventOrObject EventOrObject `json:"eventOrObject"`
	EndEpoch      int           `json:"endEpoch"`
}
type StoreResponse struct {
	Certified Certified `json:"alreadyCertified"`
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func store(c *http.Client, url string, from string) {
	fmt.Println("url=", url)
	defer wc.Done()

	if c == nil {
		log.Fatalln("client is nil")
		return
	}

	var data []byte
	if _, err := os.Lstat(from); err == nil {
		file, _ := os.Open(from)
		defer file.Close()

		data, _ = io.ReadAll(file)
	} else {
		log.Fatal("file not exist")
		return
	}

	debug("post:", url)

	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	var data_len = len(data)
	debug("data length:", data_len)

	if err == nil {
		if rep, err := c.Do(req); err == nil {
			content, _ := io.ReadAll(rep.Body)
			debug("get response: " + string(content))
			var response StoreResponse
			json.Unmarshal(content, &response)
			fmt.Println("unmarshal response:")
			fmt.Println("txDigest:", response.Certified.EventOrObject.Event.TxDigest)
			fmt.Println("blobId:", response.Certified.BlobId)
			fmt.Println("endEpoch:", response.Certified.EndEpoch)
			rep.Body.Close()
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}

}

var blobId = flag.String("blob-id", "", "{blob-id}")
var pathInTar = flag.String("path", "", "path in tar file")
var out = flag.String("out", "", "output file")
var from_path = flag.String("from", "", "from dir")
var action = flag.String("action", "read", "read|store")
var epochs = flag.Int("epochs", 1, "epoch number")

/*
*
go run main/main.go -action=store -from=go.mod -epochs=3
go run main/main.go -action=read -blob-id=gTZQ1xeTlgY9NG7QSLDWra5uaXIcV5NCDRJcPpQTkFY -out a.tar -path  /d/a.txt
*/
func main() {
	//store
	var parseFail bool = false
	flag.Parse()

	debug("blobId=", *blobId)
	debug("action=", *action)
	debug("from=", *from_path)
	debug("epochs=", *epochs)
	debug("out=", *out)
	debug("path in tar=", *pathInTar)

	if len(os.Args) > 2 {
		debug("action ,from :", *action, *from_path)
		if *action == "store" {

			if len(*from_path) == 0 {
				parseFail = true
			}
		} else if *action == "read" {

			if *blobId == "" {
				parseFail = true
			}
		} else {
			parseFail = true
		}
	} else {
		parseFail = true
	}

	// 自定义帮助信息
	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "使用方式: %s  store| read [选项]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s  -action=read -blob-id {blobid} -path {path in tar} -out {output file} \n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s  -action=store -from  {file list} \n", os.Args[0])
		fmt.Fprintf(os.Stderr, "选项:\n")

		flag.PrintDefaults()
	}
	// 如果没有提供任何命令行参数，则打印帮助信息
	if parseFail {
		flag.Usage()
		return
	}

	const AGGREGATOR = "https://aggregator.walrus-testnet.walrus.space"
	//const AGGREGATOR = "http://walrus.krates.ai:9000"
	const PUBLISHER = "https://publisher.walrus-testnet.walrus.space"
	//const PUBLISHER = "http://walrus.krates.ai:9001"
	wc.Add(1)

	client := getHttpsClient()
	if *action == "read" {
		var url = AGGREGATOR + "/v1/" + *blobId
		read(client, url, *out)
	} else if *action == "store" {
		var url = PUBLISHER + "/v1/store?epochs=" + strconv.Itoa(*epochs)
		store(client, url, *from_path)
	}

	wc.Wait()
}
