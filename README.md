Echo http header

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
