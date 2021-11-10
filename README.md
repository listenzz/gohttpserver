Echo http header

> 基于 errgroup 实现一个 http server 的启动和关闭 ，
> 以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

## 查看帮助信息

```sh
make help
```

## 运行

```sh
make run
```

## 访问

```sh
curl -i  -H 'Accept: text/html' -H 'Accept: application/xml' localhost:8080
```

## 构建镜像

```sh
make docker-build
```

## 启动镜像

```sh
make docker-run
```

## 发布镜像

```sh
docker login
make docker-push
```

[镜像地址](https://hub.docker.com/repository/docker/listenzz/gohttpserver)
