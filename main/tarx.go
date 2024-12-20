package main

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func join(dir string, name string) string {
	if len(dir) == 0 {
		return name
	} else if len(name) == 0 {
		return dir
	} else {
		return dir + "/" + name
	}

}
func writeEntry(fi fs.FileInfo, tr io.Reader, path string, target string) {
	debug("writeEntry:wirte path,target:", path, target)
	if fi.IsDir() {
		if len(target) == 0 {
			return
		}
		e := os.MkdirAll(target, fi.Mode().Perm())
		PrintError(e)
		debug("mkdir target,perm:", target, fi.Mode().Perm())
		return
	}

	// var dir = filepath.Dir(target)
	// debug("writeEntry: mkdir=", dir)
	// e2 := os.MkdirAll(dir, fi.Mode().Perm())
	// PrintError(e2)

	var file *os.File = os.Stdout
	if len(target) > 0 {
		// 创建一个空文件，用来写入解包后的数据
		fw, err := os.Create(target)
		PrintError(err)
		defer fw.Close()
		file = fw
	}

	// 将 tr 写入到 fw
	n, err := io.Copy(file, tr)
	PrintError(err)
	if *verbose >= 1 {
		fmt.Printf("--------解包： 到 %s ，共处理了 %d 个字符的数据。--------\n", target, n)
	}
	// 设置文件权限，这样可以保证和原始文件权限相同，如果不设置，会根据当前系统的 umask 来设置。
	////os.Chmod(fi.Name(), fi.Mode().Perm())
}

func extractFile(fr io.Reader, pathInTar string, out string) {
	out = strings.Trim(out, " ")
	// if len(out) != 0 {
	// 	os.Mkdir(out, os.ModePerm)
	// }

	// 通过 fr 创建一个 tar.*Reader 结构，然后将 tr 遍历，并将数据保存到磁盘中
	tr := tar.NewReader(fr)

	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		// 处理 err ！= nil 的情况
		PrintError(err)
		// 获取文件信息
		fi := hdr.FileInfo()
		debug("extract fi,hdr:", fi.Name(), hdr.Name)
		if hdr.Name == pathInTar && !fi.IsDir() {
			writeEntry(fi, tr, hdr.Name, out)
			return
		}
		if strings.HasPrefix(hdr.Name, pathInTar) {
			tailName, _ := strings.CutPrefix(hdr.Name, pathInTar)
			//remove /
			tailName = strings.TrimPrefix(tailName, string(filepath.Separator))
			writeEntry(fi, tr, hdr.Name, join(out, tailName))
		}
	}
}

func listFilesInTar(fr io.Reader, file *os.File) {

	// 通过 fr 创建一个 tar.*Reader 结构，然后将 tr 遍历，并将数据保存到磁盘中
	tr := tar.NewReader(fr)
	fmt.Println("------------------list file path in tar---------------")
	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		// 处理 err ！= nil 的情况
		PrintError(err)
		// 获取文件信息
		fi := hdr.FileInfo()
		var t = "f"
		if fi.IsDir() {
			t = "d"
		}
		fmt.Fprintf(file, "%s : size(%d), mtime(%s) type(%s)\n", hdr.Name, fi.Size(), fi.ModTime().Format(time.UnixDate), t)
	}
}
