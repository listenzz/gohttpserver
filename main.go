package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var addr = flag.String("address", ":8080", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()
	os.Setenv("version", "1.0.0")

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/healthz", healthz)

	server := http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	processed := make(chan struct{})

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		fmt.Printf("接收到信号：%s\n", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("服务优雅关闭失败, err: %v\n", err)
		}

		fmt.Println("服务优雅关闭")

		close(processed)
	}()

	fmt.Printf("开始服务监听：%s\n\n", *addr)
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("服务没有优雅终止, err :%v\n", err)
	}

	<-processed

	fmt.Println("服务终止")
}

func home(w http.ResponseWriter, r *http.Request) {
	// 模拟耗时任务
	time.Sleep(5 * time.Second)

	for k, v := range r.Header {
		for _, h := range v {
			writeHeader(w, k, h)
		}
	}

	writeHeader(w, "System-Version", os.Getenv("version"))
	writeHeader(w, "Go-Version", runtime.Version())
	writeStatusCode(w, http.StatusOK)

	fmt.Printf("\n访问者 IP: %s \n", r.RemoteAddr)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("200"))
	if err != nil {
		fmt.Printf("访问 /healthz 时发生异常: %v", err)
		writeStatusCode(w, http.StatusInternalServerError)
		return
	}
	writeStatusCode(w, http.StatusOK)
}

func writeHeader(w http.ResponseWriter, key string, header string) {
	fmt.Printf("响应头: %s, 值: %s \n", key, header)
	w.Header().Add(key, header)
}

func writeStatusCode(w http.ResponseWriter, statusCode int) {
	fmt.Printf("响应码：%d \n", statusCode)
	w.WriteHeader(statusCode)
}
