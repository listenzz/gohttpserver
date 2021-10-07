一个用 GO 写的返回请求头的 http 服务

## 运行

```sh
go run main.go
```

## 访问

```sh
curl -i  -H 'Accept: text/html' -H 'Accept: application/xml'   localhost:8000
```
