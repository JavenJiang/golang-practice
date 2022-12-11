package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"regexp"
)

/*
  1.接收客户端 request，并将 request 中带的 header 写入 response header
  2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
  3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
  4.当访问 localhost/healthz 时，应返回 200
*/
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/request", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "200")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	os.Setenv("VERSION", "v1.0")
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)
	for k, v := range r.Header {
		for _, vv := range v {
			fmt.Println(k, vv)
			w.Header().Add(k, vv)
		}
	}

	clinetIP := getUserId(r)
	log.Printf("IP: %s", clinetIP)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome<h1>"))
}

func getUserId(r *http.Request) string {
	spaceRe, _ := regexp.Compile(`\[|\]`)
	ip := spaceRe.Split(r.RemoteAddr, -1)[1]
	return ip
}
