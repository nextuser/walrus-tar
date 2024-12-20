package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"unsafe"
)

var wc sync.WaitGroup

func getOutFile(out string) *os.File {
	var outFile *os.File = nil
	if len(out) == 0 {
		outFile = os.Stdout
	} else {
		var file, openErr = os.OpenFile(out, os.O_WRONLY|os.O_CREATE, 0664)
		if openErr != nil {
			debug(openErr)
			os.Exit(1)
		}
		outFile = file
	}
	return outFile
}

/*
*
read file from the walrus,process read blob data

func extractFile(fr io.Reader, pathInTar string, out string)
*/
func process(body io.Reader, action string, pathInTar string, out string) {

	if action == "read" && len(pathInTar) == 0 {
		_, write_err := io.Copy(getOutFile(out), body)
		log.Fatalln(write_err)
	} else if action == "read" {
		extractFile(body, pathInTar, out)
	} else if action == "list" {
		listFilesInTar(body, getOutFile(out))
	}
}
func read(c *http.Client, url string, action string, pathInTar string, out string) {
	info("url=", url)
	defer wc.Done()

	var response, err = c.Get(url)
	if err != nil {
		log.Fatalln(err)
	} else {
		if response.StatusCode == 200 {
			process(response.Body, action, pathInTar, out)
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

func info(args ...any) {
	if *verbose >= 1 {
		fmt.Println(args...)
	}
}

func store(c *http.Client, url string, from string) {
	info("url=", url)
	defer wc.Done()

	if c == nil {
		log.Fatalln("client is nil")
		return
	}

	debug("post:", url)

	// from file or dir  => pipe => tar =>  http.request.body
	var buf bytes.Buffer
	TarWrite(from, &buf)

	debug("begin new request")
	req, err := http.NewRequest("PUT", url, &buf)

	if err == nil {
		if rep, err := c.Do(req); err == nil {
			content, _ := io.ReadAll(rep.Body)
			info("get response: \n" + string(content))
			var response StoreResponse
			json.Unmarshal(content, &response)
			// fmt.Println("unmarshal response:")
			// fmt.Println("txDigest:", response.Certified.EventOrObject.Event.TxDigest)
			// fmt.Println("blobId:", response.Certified.BlobId)
			// fmt.Println("endEpoch:", response.Certified.EndEpoch)
			rep.Body.Close()
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}

}

var blobId = flag.String("b", "", "{blob}, required when action=read or list ")
var pathInTar = flag.String("p", "", "path in tar file,,used when action=read")
var out = flag.String("o", "", "output file,,used when action=read")
var from_path = flag.String("f", "", "source file or directory ,  required when action=store")
var action = flag.String("c", "read", "read|store|list")
var epochs = flag.Int("e", 1, "epoch number,used when action=store")
var verbose = flag.Int("v", 1, "verbose 0|1|2")

/*
*
go run main/main.go -c=store -f=go.mod -e=3
go run main/main.go -c=read -b=gTZQ1xeTlgY9NG7QSLDWra5uaXIcV5NCDRJcPpQTkFY -o a.tar -p  /d/a.txt
*/
func main() {
	//store
	var parseFail bool = false
	flag.Parse()
	if *verbose == 2 {
		enableDebug(true)
	}
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
		} else if *action == "read" || *action == "list" {

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

		fmt.Fprintf(os.Stderr, "使用方式:", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s  -c list -b {blob id} -o {output file} \n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s  -c read -b {blob id} -p {path in tar} -o {output file/dir} \n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s  -c store -f  {file list} \n", os.Args[0])
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
	if *action == "read" || *action == "list" {
		var url = AGGREGATOR + "/v1/" + *blobId
		read(client, url, *action, *pathInTar, *out)
	} else if *action == "store" {
		var url = PUBLISHER + "/v1/store?epochs=" + strconv.Itoa(*epochs)
		store(client, url, *from_path)
	}

	wc.Wait()
}
