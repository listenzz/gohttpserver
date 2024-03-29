package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/golang/glog"
	"golang.org/x/sync/errgroup"
)

var addr = flag.String("address", ":8080", "The address to listen on for HTTP requests.")

func main() {
	defer glog.Flush()

	flag.Parse()
	os.Setenv("version", "1.0.0")

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/healthz", healthz)

	server := http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	group, ctx := errgroup.WithContext(context.Background())

	group.Go(func() error {
		glog.Infof("开始服务监听：%s\n\n", *addr)
		return server.ListenAndServe()
	})

	group.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			glog.Infof("接收到信号：%s\n", sig.String())
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			glog.Infoln("关闭服务...")
			return server.Shutdown(ctx)
		}
	})

	err := group.Wait()
	if err != nil && err != http.ErrServerClosed {
		glog.Errorf("服务异常终止, err :%s\n", err.Error())
	} else {
		glog.Infoln("服务正常终止")
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	// 模拟耗时任务
	glog.Infoln("接收到请求，1秒后返回结果")
	time.Sleep(1 * time.Second)

	for k, v := range r.Header {
		for _, h := range v {
			writeHeader(w, k, h)
		}
	}

	fmt.Fprintln(w, "Hello, World!")

	fmt.Fprintln(w, "版本：v1.0.5")

	_, err := fmt.Fprintf(w, "GOARCH: %s \n响应时间：%s \n", runtime.GOARCH, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		glog.Errorf("访问 / 时发生异常: %s \n", err.Error())
		writeStatusCode(w, http.StatusInternalServerError)
		return
	}
	writeHeader(w, "System-Version", os.Getenv("version"))
	writeHeader(w, "Go-Version", runtime.Version())

	glog.Infof("\n访问者 IP: %s \n", r.RemoteAddr)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	data := time.Now().String()
	glog.Infof("healthz 探测 %s \n", data)
	_, err := w.Write([]byte(data))
	if err != nil {
		glog.Errorf("访问 /healthz 时发生异常: %s", err.Error())
		writeStatusCode(w, http.StatusInternalServerError)
		return
	}
}

func writeHeader(w http.ResponseWriter, key string, header string) {
	glog.Infof("响应头: %s, 值: %s \n", key, header)
	w.Header().Add(key, header)
}

func writeStatusCode(w http.ResponseWriter, statusCode int) {
	glog.Infof("响应码：%d \n", statusCode)
	w.WriteHeader(statusCode)
}
