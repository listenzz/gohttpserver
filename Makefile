GO111MODULES=on
APP=gohttpserver
TAG=v1.0.7

.PHONY: build
## build: 构建应用
build: clean
	go build -o ${APP} .

.PHONY: run
## run: 在本地运行应用
run:
	go run -race main.go -logtostderr

.PHONY: clean
## clean: 清除构建物
clean:
	go clean

.PHONY: docker-build
## docker-build: 构建 docker 镜像
docker-build:
	docker build -t listenzz/${APP}:${TAG} -t listenzz/${APP}:latest .

.PHONY: docker-run
## docker-run: 在 docker 容器内运行本应用
docker-run: docker-build
	docker run -p 8080:8080 listenzz/${APP}:${TAG} -address :8080 -logtostderr

.PHONY: docker-push
## docker-push: 推送 docker 镜像到 docker hub
docker-push:
	docker buildx build --platform linux/amd64,linux/arm64 -t listenzz/${APP}:${TAG} -t listenzz/${APP}:latest -o type=registry .

.PHONY: help
## help: 打印帮助信息
help:
	@echo "使用：\n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'