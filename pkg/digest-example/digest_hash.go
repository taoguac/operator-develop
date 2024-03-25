package main

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

var fileName = "/Users/gutao/gutaodev/gocode/hellogolang/file/config.toml"

func digest_file() {
	file, _ := os.Open(fileName)
	defer file.Close()
	m5 := md5.New()
	h2 := sha256.New()

	io.Copy(m5, file)
	m := m5.Sum(nil)
	fmt.Printf("md5: %x\n", m)

	io.Copy(h2, file)
	h := h2.Sum(nil)
	fmt.Printf("sha256: %x\n", h)
}

func digest_md5() {
	t := time.Now()
	fmt.Println("md5 begin, time=", t)

	file, _ := os.Open(fileName)
	defer file.Close()
	m5 := md5.New()
	io.Copy(m5, file)
	m := m5.Sum(nil)
	fmt.Printf("md5: %x\n", m)

	fmt.Println("md5 end, time=", time.Now(), "cost:", time.Now().Sub(t))
}

func digest_sha256() {
	t := time.Now()
	fmt.Println("sha256 begin, time=", t)

	file, _ := os.Open(fileName)
	defer file.Close()
	h2 := sha256.New()
	io.Copy(h2, file)
	h := h2.Sum(nil)
	fmt.Printf("sha256: %x\n", h)

	fmt.Println("sha256 end, time=", time.Now(), "cost:", time.Now().Sub(t))
}

func main() {
	digest_md5()
	digest_sha256()
}
