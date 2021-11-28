# 编写 Kubernetes 部署脚本将 httpserver 部署到 kubernetes 集群

## 优雅启动、探活

通过[配置存活、就绪和启动探测器](https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) 来实现

## 优雅终止

结合 `errorgroup` 和 `signal` 来实现

为了方便测试，将 replicas 设置为 1

通过以下命令持续查看日志

```sh
kubectl logs pod/gohttpserver-7f84fcd76b-69kq9 -f
```

在一终端输入如下命令，发起请求

```sh
curl -i  -H 'Accept: text/html' -H 'Accept: application/xml' 192.168.32.33:31533
```

在另一终端，输入如下命令，删除 pod

```sh
kubectl delete pod/gohttpserver-7f84fcd76b-69kq9
```

可以看到如下日志输出

```
I1128 02:13:55.600998       1 main.go:67] 接收到请求，5秒后返回结果
I1128 02:13:59.571209       1 main.go:49] 接收到信号：terminated
I1128 02:13:59.571246       1 main.go:52] 关闭服务...
I1128 02:14:00.601695       1 main.go:101] 响应头: Accept, 值: text/html
I1128 02:14:00.601743       1 main.go:101] 响应头: Accept, 值: application/xml
I1128 02:14:00.601749       1 main.go:101] 响应头: User-Agent, 值: curl/7.77.0
I1128 02:14:00.601787       1 main.go:101] 响应头: System-Version, 值: 1.0.0
I1128 02:14:00.601791       1 main.go:101] 响应头: Go-Version, 值: go1.17.1
I1128 02:14:00.601797       1 main.go:86]
访问者 IP: 10.211.55.5:16252
I1128 02:14:00.647916       1 main.go:61] 服务正常终止
```

服务在接收到请求后，还没来得及返回响应，此时收到 k8s 发过来的 terminated 信号，进入优雅终止流程，不再接收新请求，待已有请求完成响应后，服务终止退出。

## 资源需求和 QoS 保证

通过配置 resources 和 limit 来实现

```yaml
resources:
  limits:
    memory: "128Mi"
    cpu: "500m"
```

通过 kubectl describe pod gohttpserver-7f84fcd76b-h988n 命令查看，可看到 QoS Class 为 Guaranteed

```
QoS Class:                   Guaranteed
```

参考

- [为容器和 Pod 分配内存资源](https://kubernetes.io/zh/docs/tasks/configure-pod-container/assign-memory-resource/)
- [为容器和 Pods 分配 CPU 资源](https://kubernetes.io/zh/docs/tasks/configure-pod-container/assign-cpu-resource/)
- [配置 Pod 的服务质量](https://kubernetes.io/zh/docs/tasks/configure-pod-container/quality-service-pod/)

## 日常运维需求，日志等级

使用 glog 记录日志，合理使用 Info、Error 等日志等级，运行应用时，传递 `-logtostderr` 参数，以将日志打印到标准输出

## 配置和代码分离

- 使用 [ConfigMap](https://kubernetes.io/zh/docs/concepts/configuration/configmap/) 来配置，未考虑

- 通过 spec#containers#args 来传递参数，如

```yaml
spec:
  containers:
    - name: gohttpserver
      args:
        - "-address"
        - ":8080"
        - "-logtostderr"
```

## Service

使用 Service 对外对内暴露服务，可以将 Service 的类型设置为 NodePort，方便通过节点访问

## Ingress

使用 ingress-nginx 作为 Ingress 的实现，并配置了证书

结合之前搭建的[基础集群](https://todoit.tech/k8s/#%E5%9F%BA%E7%A1%80%E9%9B%86%E7%BE%A4)，该集群部署了 ingress-nginx，使用 cert-manager 来自动签发和管理证书。

由于集群通过虚拟机部署在本地，这里使用 SwitchHosts 作域名映射

![deploy-2021-11-28-10-32-02](https://todoit.oss-cn-shanghai.aliyuncs.com/todoit/deploy-2021-11-28-10-32-02.png)

通过浏览器访问 gohttpserver.todoit.tech

![deploy-2021-11-28-10-30-38](https://todoit.oss-cn-shanghai.aliyuncs.com/todoit/deploy-2021-11-28-10-30-38.png)
