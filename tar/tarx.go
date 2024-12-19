package tar

import (
	"archive/tar"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func join(dir string, name string) string {
	if len(dir) == 0 {
		return name
	} else {
		return dir + "/" + name
	}

}
func writeEntry(fi fs.FileInfo, tr io.Reader, path string, out string) {
	target := join(out, path)
	debug("wirte path,target:", path, target)
	if fi.IsDir() {
		e := os.MkdirAll(target, fi.Mode().Perm())
		if e == nil {
			log.Fatal("ERROR CREATE", target, fi.Mode().Perm())
		}
		debug("mkdir target,perm:", target, fi.Mode().Perm())
		return
	}

	var dir = filepath.Dir(target)
	debug("dir=", dir)
	e2 := os.MkdirAll(dir, os.ModePerm)
	ErrPrintln(e2)
	// 创建一个空文件，用来写入解包后的数据
	fw, err := os.Create(target)
	defer fw.Close()
	ErrPrintln(err)

	// 将 tr 写入到 fw
	n, err := io.Copy(fw, tr)
	ErrPrintln(err)
	log.Printf("解包： 到 %s ，共处理了 %d 个字符的数据。", target, n)
	// 设置文件权限，这样可以保证和原始文件权限相同，如果不设置，会根据当前系统的 umask 来设置。
	////os.Chmod(fi.Name(), fi.Mode().Perm())
}

func extractFile(fr io.Reader, pathInTar string, out string) {
	out = strings.Trim(out, " ")
	if len(out) != 0 {
		os.Mkdir(out, os.ModePerm)
	}
	// 通过 fr 创建一个 tar.*Reader 结构，然后将 tr 遍历，并将数据保存到磁盘中
	tr := tar.NewReader(fr)

	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		// 处理 err ！= nil 的情况
		ErrPrintln(err)
		// 获取文件信息
		fi := hdr.FileInfo()
		debug("extract fi,hdr:", fi.Name(), hdr.Name)
		if hdr.Name == pathInTar {
			writeEntry(fi, tr, hdr.Name, out)
		}
	}
}
