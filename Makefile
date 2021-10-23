GO111MODULES=on
APP=gohttpserver
TAG=v1.0.0

.PHONY: build
## build: 构建应用
build: clean
	go build -o ${APP} main.go

.PHONY: run
## run: 在本地运行应用
run:
	go run -race main.go

.PHONY: clean
## clean: 清除构建物
clean:
	go clean

.PHONY: docker-build
## docker-build: 构建 docker 镜像
docker-build:
	docker build -t listenzz/${APP}:${TAG} .

.PHONY: docker-run
## docker-run: 在 docker 容器内运行本应用
docker-run: docker-build
	docker run -p 8000:8000 listenzz/${APP}:${TAG} -- -address :8000

.PHONY: docker-push
## docker-push: 推送 docker 镜像到 docker hub
docker-push: docker-build
	docker push listenzz/${APP}:${TAG}

.PHONY: help
## help: 打印帮助信息
help:
	@echo "使用：\n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'