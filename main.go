package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
)

var addr = flag.String("address", "127.0.0.1:8000", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()

	os.Setenv("version", "1.0.0")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			for _, h := range v {
				writeHeader(w, k, h)
			}
		}

		writeHeader(w, "System-Version", os.Getenv("version"))
		writeHeader(w, "Go-Version", runtime.Version())
		writeStatusCode(w, http.StatusOK)

		fmt.Printf("\n访问者 IP: %s \n", r.RemoteAddr)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("200"))
		if err != nil {
			fmt.Printf("访问 /healthz 时发生异常: %v", err)
			writeStatusCode(w, http.StatusInternalServerError)
			return
		}
		writeStatusCode(w, http.StatusOK)
	})

	fmt.Printf("开始服务监听：%s\n\n", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Printf("服务异常退出: %v", err)
	}
}

func writeHeader(w http.ResponseWriter, key string, header string) {
	fmt.Printf("响应头: %s, 值: %s \n", key, header)
	w.Header().Add(key, header)
}

func writeStatusCode(w http.ResponseWriter, statusCode int) {
	fmt.Printf("响应码：%d \n", statusCode)
	w.WriteHeader(statusCode)
}
