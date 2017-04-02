//go:generate goversioninfo -icon=icon.ico -manifest luzhibo.manifest

package main

import (
	"runtime"
	"flag"
	"strconv"
	"fmt"
	"time"
	"path/filepath"
	"os"
)

const ver = 2017040200
const p = "录直播"

var port = 12216

var nhta *bool

func main() {
	p := flag.Int("port", port, "WebUI监听端口")
	nopen := flag.Bool("nopenui", false, "不自动打开WebUI")
	nhta = flag.Bool("nhta", false, "禁用hta(仅Windows有效)")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		time.Sleep(time.Second * 5)
		d, f := filepath.Split(os.Args[0])
		tp := filepath.Join(d, "."+f+".old")
		os.Remove(tp)
	}()
	port = *p
	s := ":" + strconv.Itoa(port)
	fmt.Printf("正在\"%s\"处监听WebUI...\n", s)
	go startServer(s)
	if !*nopen {
		openWebUI(!*nhta)
	}
	cmd()
}
