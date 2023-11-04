```shell
$ make deploy IMG=lengdanran/hello-operator:v0.0.1
# Warning: 'patchesStrategicMerge' is deprecated. Please use 'patches' instead. Run 'kustomize edit fix' to update your Kustomization automatically.
namespace/hellooperator-system created
customresourcedefinition.apiextensions.k8s.io/helloapps.apps.zsh.io unchanged
serviceaccount/hellooperator-controller-manager created
role.rbac.authorization.k8s.io/hellooperator-leader-election-role created
clusterrole.rbac.authorization.k8s.io/hellooperator-manager-role unchanged
clusterrole.rbac.authorization.k8s.io/hellooperator-metrics-reader unchanged
clusterrole.rbac.authorization.k8s.io/hellooperator-proxy-role unchanged
rolebinding.rbac.authorization.k8s.io/hellooperator-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/hellooperator-manager-rolebinding unchanged
clusterrolebinding.rbac.authorization.k8s.io/hellooperator-proxy-rolebinding unchanged
service/hellooperator-controller-manager-metrics-service created
deployment.apps/hellooperator-controller-manager created

$ kubectl get crd 
NAME                                                     CREATED AT
helloapps.apps.zsh.io                                    2023-11-03T10:07:28Z
loadbalancerresources.networking.tke.cloud.tencent.com   2023-11-03T09:53:57Z
tkeserviceconfigs.cloud.tencent.com                      2023-11-03T09:53:52Z
volumesnapshotclasses.snapshot.storage.k8s.io            2023-11-03T09:55:49Z
volumesnapshotcontents.snapshot.storage.k8s.io           2023-11-03T09:55:49Z
volumesnapshots.snapshot.storage.k8s.io                  2023-11-03T09:55:49Z

```

config/samples/apps_v1_helloapp.yaml

```yaml
apiVersion: apps.zsh.io/v1
kind: HelloApp
metadata:
  labels:
    label: helloapp-v1
    app.kubernetes.io/name: helloapp
    app.kubernetes.io/instance: helloapp-v1
    app.kubernetes.io/part-of: hellooperator
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: hellooperator
  name: helloapp-v1
spec:
  # 新增下述内容
  replicas: 1 # 设置副本数
  label: helloapp-v1
  template: # 按照 pod 模板设置 pod 规格
    spec:
      containers:
        - image:  nginx
          name:  app

```

```shell
$ kubectl apply -f config/samples/apps_v1_helloapp.yaml
helloapp.apps.zsh.io/helloapp-v1 created
```



get pods

```shell
$ kubectl get pods                                     
NAME                                READY   STATUS    RESTARTS   AGE
helloapp-v1-nxn2m                   1/1     Running   0          115s
kubernetes-proxy-78474bcd56-bkc9b   1/1     Running   0          16h
kubernetes-proxy-78474bcd56-qdqqj   1/1     Running   0          16h

```

modify replicas to 3

```yaml
apiVersion: apps.zsh.io/v1
kind: HelloApp
metadata:
  labels:
    label: helloapp-v1
    app.kubernetes.io/name: helloapp
    app.kubernetes.io/instance: helloapp-v1
    app.kubernetes.io/part-of: hellooperator
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: hellooperator
  name: helloapp-v1
spec:
  # 新增下述内容
  replicas: 3 # 设置副本数
  label: helloapp-v1
  template: # 按照 pod 模板设置 pod 规格
    spec:
      containers:
        - image:  nginx
          name:  app
```

```shell
$ kubectl apply -f config/samples/apps_v1_helloapp.yaml
helloapp.apps.zsh.io/helloapp-v1 configured


$ kubectl get pods                                     
NAME                                READY   STATUS    RESTARTS   AGE
helloapp-v1-4g8g5                   1/1     Running   0          5s
helloapp-v1-ctch5                   1/1     Running   0          5s
helloapp-v1-nxn2m                   1/1     Running   0          2m29s
kubernetes-proxy-78474bcd56-bkc9b   1/1     Running   0          16h
kubernetes-proxy-78474bcd56-qdqqj   1/1     Running   0          16h

```

```shell
$ kubectl delete HelloApp helloapp-v1                  
helloapp.apps.zsh.io "helloapp-v1" deleted
$ kubectl get pods                   
NAME                                READY   STATUS    RESTARTS   AGE
kubernetes-proxy-78474bcd56-bkc9b   1/1     Running   0          16h
kubernetes-proxy-78474bcd56-qdqqj   1/1     Running   0          16h
```


operator-log

```log
2023-11-04T10:25:01+08:00       INFO    setup   starting manager
2023-11-04T10:25:01+08:00       INFO    controller-runtime.metrics      Starting metrics server
2023-11-04T10:25:01+08:00       INFO    starting server {"kind": "health probe", "addr": "[::]:8081"}
2023-11-04T10:25:01+08:00       INFO    controller-runtime.metrics      Serving metrics server  {"bindAddress": ":8080", "secure": false}
2023-11-04T10:25:01+08:00       INFO    Starting EventSource    {"controller": "helloapp", "controllerGroup": "apps.zsh.io", "controllerKind": "HelloApp", "source": "kind source: *v1.HelloApp"}
2023-11-04T10:25:01+08:00       INFO    Starting Controller     {"controller": "helloapp", "controllerGroup": "apps.zsh.io", "controllerKind": "HelloApp"}
2023-11-04T10:25:02+08:00       INFO    Starting workers        {"controller": "helloapp", "controllerGroup": "apps.zsh.io", "controllerKind": "HelloApp", "worker count": 1}
2023-11-04T10:25:03+08:00       INFO    Updating pod Replicas.....      {"controller": "helloapp", "controllerGroup": "apps.zsh.io", "controllerKind": "HelloApp", "HelloApp": {"name":"helloapp-v1","namespace":"default"}, "namespace": "default", "name": "helloapp-v1", "reconcileID": "1b7059aa-0c7c-4d86-99f5-5915a691075e"}
2023-11-04T10:25:03+08:00       INFO    Current pod cnt = 0     {"controller": "helloapp", "controllerGroup": "apps.zsh.io", "controllerKind": "HelloApp", "HelloApp": {"name":"helloapp-v1","namespace":"default"}, "namespace": "default", "name": "helloapp-v1", "reconcileID": "1b7059aa-0c7c-4d86-99f5-5915a691075e"}
2023-11-04T10:25:03+08:00       INFO    Less than desired replicas 1    {"controller": "helloapp", "controllerGroup": "apps.zsh.io", "controllerKind": "HelloApp", "HelloApp": {"name":"helloapp-v1","namespace":"default"}, "namespace": "default", "name": "helloapp-v1", "reconcileID": "1b7059aa-0c7c-4d86-99f5-5915a691075e"}
2023-11-04T10:25:03+08:00       INFO    Current pod cnt matches with desired cnt 1      {"controller": "helloapp", "controllerGroup": "apps.zsh.io", "controllerKind": "HelloApp", "HelloApp": {"name":"helloapp-v1","namespace":"default"}, "namespace": "default", "name": "helloapp-v1", "reconcileID": "1b7059aa-0c7c-4d86-99f5-5915a691075e"}
2023-11-04T10:27:27+08:00       INFO    Updating pod Replicas.....      {"controller": "helloapp", "controllerGroup": "apps.zsh.io", "controllerKind": "HelloApp", "HelloApp": {"name":"helloapp-v1","namespace":"default"}, "namespace": "default", "name": "helloapp-v1", "reconcileID": "ee6e9000-a9de-4176-aaea-de2348070533"}
2023-11-04T10:27:27+08:00       INFO    Current pod cnt = 1     {"controller": "helloapp", "controllerGroup": "apps.zsh.io", "controllerKind": "HelloApp", "HelloApp": {"name":"helloapp-v1","namespace":"default"}, "namespace": "default", "name": "helloapp-v1", "reconcileID": "ee6e9000-a9de-4176-aaea-de2348070533"}
2023-11-04T10:27:27+08:00       INFO    Less than desired replicas 3    {"controller": "helloapp", "controllerGroup": "apps.zsh.io", "controllerKind": "HelloApp", "HelloApp": {"name":"helloapp-v1","namespace":"default"}, "namespace": "default", "name": "helloapp-v1", "reconcileID": "ee6e9000-a9de-4176-aaea-de2348070533"}
2023-11-04T10:27:28+08:00       INFO    Current pod cnt matches with desired cnt 3      {"controller": "helloapp", "controllerGroup": "apps.zsh.io", "controllerKind": "HelloApp", "HelloApp": {"name":"helloapp-v1","namespace":"default"}, "namespace": "default", "name": "helloapp-v1", "reconcileID": "ee6e9000-a9de-4176-aaea-de2348070533"}
2023-11-04T10:29:32+08:00       INFO    HelloApp CRD is deleted, delete the associated Pods     {"controller": "helloapp", "controllerGroup": "apps.zsh.io", "controllerKind": "HelloApp", "HelloApp": {"name":"helloapp-v1","namespace":"default"}, "namespace": "default", "name": "helloapp-v1", "reconcileID": "14502f78-b7fd-4e0c-98ba-ecfee873ec16"}

```
