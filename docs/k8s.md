# k8s usage 


```sh
#build 镜像

docker build -t golang-test:v2 . 


```




```sh
概念
Nodeport
ClusterIP

Kind:Service


kubectl version 
# 查看k8s集群信息
kubectl cluster-info


# 获取pod
kubectl get pods 

# 获取service
kubectl  get services

# 查看services=kubernetes-bootcamp
kubectl describe services/kubernetes-bootcamp

# 获取service= kubernetes-bootcamp port
NODE_PORT=$(kubectl get services/kubernetes-bootcamp -o go-template='{{(index .spec.ports 0).nodePort}}')

# 查看deployment的信息
kubectl describe deployment


# 暴露服务，外部可以访问
kubectl expose deployment/kubernetes-bootcamp --type="NodePort" --port 8080
kubectl expose deployment/firstapp --type="NodePort" --port 8080

kubectl create deployment first-k8s-deploy \
  --image="laxman/nodejs/express-app" \
  -o yaml \
  --dry-run
# 会生成一个yaml文件，修改yaml文件之后，使用

kubectl apply -f a.yaml

# 使用标签
kubectl get pods -l app=kubernetes-bootcamp

# 查看标签等于特定值的服务
kubectl get services -l app=kubernetes-bootcamp

# 删除service,删除之后，不可从外部访问，只能从容器内部访问

kubectl delete service -l app=kubernetes-bootcamp

```

## 部署应用

```sh
# 创建一个deployment

# firstapp 是名字，--image=golang-test 指定了app image
kubectl create deployment firstapp  --image=golang-test
kubectl get deployments

kubectl expose deployment/firstapp --type="NodePort" --port 8080


```

## 部署服务到k8s

```sh
# generate a deployment yaml files  
kubectl create deployment first-app --image="golang-test:v3" -o yaml --dry-run=client > first-app.yaml

# the reason generate deployment yaml is that ,we will modify imagePullPolicy in yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: first-app
  name: first-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: first-app
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: first-app
    spec:
      containers:
      - image: golang-test:v3
        name: golang-test
        imagePullPolicy: Never ## Never: use local docker image other than pull docker image from docker-hub 
        resources: {}
status: {}

# apply deployment
kubectl apply -f first-app.yaml

# 查看pod
kubectl get pods
# 查看部署
kubectl get deployment
# now container is running
docker ps
# output like this:
# CONTAINER ID   IMAGE               COMMAND      CREATED STATUS     PORTS                NAMES
# 0f3afe30a460   5a4d02edae48        "/app/main"  About a minute ago Up About a minute                                                       k8s_golang-test_first-app-6cdc85f58-hk8pv_default_8932e5bd-0d33-494f-b6b0-42e772186b89_0

# expose service outside k8s cluster
kubectl expose deployment/first-app --type="NodePort" --port 9090
# - type:expose type 
# - port: is the server listen port in container

# 查看service
kubectl get services
# output like this
# NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
# first-app    NodePort    10.105.179.76   <none>        9090:32158/TCP   4s

# now you can access server in your local machine 
curl localhost:32158/
# server output for you request
# welcome 1 guest
```



## 服务升级

```shell 
# 升级镜像

kubectl set image deployments/firstapp  kubernetes-bootcamp=golang-test:v3

kubectl set image deployment/first-app first-app=glang-test:v5
kubectl set image deployment/first-app golang-test=golang-tes:v5 --record


# update the images
# describe in kubectl apply yaml files where images name 
# kubectl set image deployment/deployment-name imagename=new_imagename --record 
kubectl set image deployment/first-app golang-test=golang-test:v5 --record 



# rollingup
kubectl rollout status deployments/first-app

kubectl set image deployment/first-app golang-test=golang-test:v5 --record 
kubectl rollout undo deployments/first-app






```

## ambassador 

```sh


kubectl create namespace ambassador
helm install --devel edge-stack --namespace ambassador datawire/edge-stack 
kubectl -n ambassador wait --for condition=available --timeout=90s deploy -lproduct=aes


kubectl apply -f - <<EOF
---
apiVersion: x.getambassador.io/v3alpha1
kind: AmbassadorListener
metadata:
  name: edge-stack-listener-8080
  namespace: ambassador
spec:
  port: 8080
  protocol: HTTP
  securityModel: XFP
  hostBinding:
    namespace:
      from: ALL
---
apiVersion: x.getambassador.io/v3alpha1
kind: AmbassadorListener
metadata:
  name: edge-stack-listener-8443
  namespace: ambassador
spec:
  port: 8443
  protocol: HTTPS
  securityModel: XFP
  hostBinding:
    namespace:
      from: ALL
EOF


# create deployment 
kubectl apply -n ambassador -f https://app.getambassador.io/yaml/v2-docs/latest/quickstart/qotm.yaml


# see the status
kubectl get services,deployments quote --namespace ambassador


## quote-backend.yaml file content
---
apiVersion: x.getambassador.io/v3alpha1
kind: AmbassadorMapping
metadata:
  name: quote-backend
  namespace: ambassador
spec:
  hostname: "*"
  prefix: /backend/
  service: quote



```