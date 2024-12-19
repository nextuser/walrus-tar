package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

func Test_Tar(t *testing.T) {
	// 修改日志格式，显示出错代码的所在行，方便调试，实际项目中一般不记录这个。

	var src = "data/apt"
	var dst = fmt.Sprintf("%s.tar", src)

	// 将步骤写入了一个函数中，这样处理错误方便一些
	if err := TarFile(src, dst); err != nil {
		log.Fatalln(err)
	}
}

// func extractFile(fr io.Reader, pathInTar string, out string) {
func Test_Extract(t *testing.T) {
	log.Println("extract file")
	f, e := os.OpenFile("./data/apt.tar", os.O_RDONLY, os.ModePerm)
	ErrPrintln(e)
	r := bufio.NewReader(f)

	extractFile(r, "apt/a", "data/other")

}

func Test_List(t *testing.T) {
	f, e := os.OpenFile("./data/apt.tar", os.O_RDONLY, os.ModePerm)
	ErrPrintln(e)
	var reader = bufio.NewReader(f)

	listFilesInTar(reader, os.Stdout)
}
