一个用 GO 写的返回请求头的 http 服务

## 运行

```sh
go run main.go
```

## 访问

```sh
curl -i  -H 'Accept: text/html' -H 'Accept: application/xml'   localhost:8000
```

## 构建镜像

```
docker build . -t listenzz/gohttpserver:v1.0.0
```

## 启动镜像

```
docker run -p 8000:8000 listenzz/gohttpserver:v1.0.0 --
```

## 发布镜像

```
docker login
docker push listenzz/gohttpserver:v1.0.0
```

[镜像地址](https://hub.docker.com/repository/docker/listenzz/gohttpserver)
