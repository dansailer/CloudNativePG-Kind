# EnterpriseDB Operator Workshop

## Install dependencies

Use your package manager to install kind, kubectl, kubectl-cnpg, k9s and sqlelectron.

```sh
brew install kind derailed/k9s/k9s kubectl-cnpg helm kubectl sqlelectron
```

## Setup Kubernetes Cluster with kind

This setup will create a 3 node Kubernetes cluster with `kind`

```sh
kind delete cluster -n pgcluster
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: pgcluster
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 5432
    hostPort: 5432
    listenAddress: "0.0.0.0"
    protocol: TCP
- role: worker
- role: worker
EOF
kubectl cluster-info --context kind-pgcluster
kubectl get nodes -A
```

```
 Context: kind-pgcluster                           <0> all       <a>      Attach     <l>       Logs            <f> Show PortForward       ____  __.________        
 Cluster: kind-pgcluster                           <1> default   <ctrl-d> Delete     <p>       Logs Previous   <t> Transfer              |    |/ _/   __   \______ 
 User:    kind-pgcluster                                         <d>      Describe   <shift-f> Port-Forward    <y> YAML                  |      < \____    /  ___/ 
 K9s Rev: v0.32.4                                                <e>      Edit       <z>       Sanitize                                  |    |  \   /    /\___ \  
 K8s Rev: v1.30.0                                                <?>      Help       <s>       Shell                                     |____|__ \ /____//____  > 
 CPU:     n/a                                                    <ctrl-k> Kill       <o>       Show Node                                         \/            \/  
 MEM:     n/a                                                                                                                                                      
┌───────────────────────────────────────────────────────────────────────── Pods(all)[13] ─────────────────────────────────────────────────────────────────────────┐
│ NAMESPACE↑            NAME                                               PF   READY   STATUS       RESTARTS IP             NODE                        AGE      │
│ kube-system           coredns-7db6d8ff4d-44sq8                           ●    1/1     Running             0 10.244.0.3     pgcluster-control-plane     24s      │
│ kube-system           coredns-7db6d8ff4d-fwghc                           ●    1/1     Running             0 10.244.0.2     pgcluster-control-plane     24s      │
│ kube-system           etcd-pgcluster-control-plane                       ●    1/1     Running             0 172.19.0.4     pgcluster-control-plane     40s      │
│ kube-system           kindnet-dm29d                                      ●    1/1     Running             0 172.19.0.3     pgcluster-worker            21s      │
│ kube-system           kindnet-g6xxh                                      ●    1/1     Running             0 172.19.0.2     pgcluster-worker2           21s      │
│ kube-system           kindnet-wtqjs                                      ●    1/1     Running             0 172.19.0.4     pgcluster-control-plane     25s      │
│ kube-system           kube-apiserver-pgcluster-control-plane             ●    1/1     Running             0 172.19.0.4     pgcluster-control-plane     40s      │
│ kube-system           kube-controller-manager-pgcluster-control-plane    ●    1/1     Running             0 172.19.0.4     pgcluster-control-plane     40s      │
│ kube-system           kube-proxy-67wfq                                   ●    1/1     Running             0 172.19.0.2     pgcluster-worker2           21s      │
│ kube-system           kube-proxy-kf7bx                                   ●    1/1     Running             0 172.19.0.4     pgcluster-control-plane     25s      │
│ kube-system           kube-proxy-mw922                                   ●    1/1     Running             0 172.19.0.3     pgcluster-worker            21s      │
│ kube-system           kube-scheduler-pgcluster-control-plane             ●    1/1     Running             0 172.19.0.4     pgcluster-control-plane     40s      │
│ local-path-storage    local-path-provisioner-988d74bc-d7dhq              ●    1/1     Running             0 10.244.0.4     pgcluster-control-plane     24s      │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```


## Setup NGINX Ingress Controller

To being able to access the deployed services from outside the cluster, a ingress controller is required.
The http services can be configured using different `nip.io` based hostnames on port 443. PostgreSQL is exposed on 5432.

```sh
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
cat <<EOF | kubectl replace -n ingress-nginx -f -
apiVersion: v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/component: controller
    app.kubernetes.io/instance: ingress-nginx
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
    app.kubernetes.io/version: 1.10.1
  name: ingress-nginx-controller
  namespace: ingress-nginx
spec:
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - appProtocol: http
    name: http
    port: 80
    protocol: TCP
    targetPort: http
  - appProtocol: https
    name: https
    port: 443
    protocol: TCP
    targetPort: https
  - name: postgresql
    port: 5432
    protocol: TCP
    targetPort: 5432
  selector:
    app.kubernetes.io/component: controller
    app.kubernetes.io/instance: ingress-nginx
    app.kubernetes.io/name: ingress-nginx
  type: NodePort
EOF
```

Wait until the `ingress-nginx` pods have completed or are running.

```
 Context: kind-pgcluster                           <0> all       <a>      Attach     <l>       Logs            <f> Show PortForward       ____  __.________        
 Cluster: kind-pgcluster                           <1> default   <ctrl-d> Delete     <p>       Logs Previous   <t> Transfer              |    |/ _/   __   \______ 
 User:    kind-pgcluster                                         <d>      Describe   <shift-f> Port-Forward    <y> YAML                  |      < \____    /  ___/ 
 K9s Rev: v0.32.4                                                <e>      Edit       <z>       Sanitize                                  |    |  \   /    /\___ \  
 K8s Rev: v1.30.0                                                <?>      Help       <s>       Shell                                     |____|__ \ /____//____  > 
 CPU:     n/a                                                    <ctrl-k> Kill       <o>       Show Node                                         \/            \/  
 MEM:     n/a                                                                                                                                                      
┌───────────────────────────────────────────────────────────────────────── Pods(all)[16] ─────────────────────────────────────────────────────────────────────────┐
│ NAMESPACE↑           NAME                                               PF   READY   STATUS         RESTARTS IP            NODE                       AGE       │
│ ingress-nginx        ingress-nginx-admission-create-pwpfk               ●    0/1     Completed             0 10.244.2.2    pgcluster-worker           86s       │
│ ingress-nginx        ingress-nginx-admission-patch-jkkf7                ●    0/1     Completed             2 10.244.1.2    pgcluster-worker2          86s       │
│ ingress-nginx        ingress-nginx-controller-56cbc5d9d4-4jklx          ●    1/1     Running               0 10.244.0.5    pgcluster-control-plane    86s       │
│ kube-system          coredns-7db6d8ff4d-44sq8                           ●    1/1     Running               0 10.244.0.3    pgcluster-control-plane    3m50s     │
│ kube-system          coredns-7db6d8ff4d-fwghc                           ●    1/1     Running               0 10.244.0.2    pgcluster-control-plane    3m50s     │
│ kube-system          etcd-pgcluster-control-plane                       ●    1/1     Running               0 172.19.0.4    pgcluster-control-plane    4m6s      │
│ kube-system          kindnet-dm29d                                      ●    1/1     Running               0 172.19.0.3    pgcluster-worker           3m47s     │
│ kube-system          kindnet-g6xxh                                      ●    1/1     Running               0 172.19.0.2    pgcluster-worker2          3m47s     │
│ kube-system          kindnet-wtqjs                                      ●    1/1     Running               0 172.19.0.4    pgcluster-control-plane    3m51s     │
│ kube-system          kube-apiserver-pgcluster-control-plane             ●    1/1     Running               0 172.19.0.4    pgcluster-control-plane    4m6s      │
│ kube-system          kube-controller-manager-pgcluster-control-plane    ●    1/1     Running               0 172.19.0.4    pgcluster-control-plane    4m6s      │
│ kube-system          kube-proxy-67wfq                                   ●    1/1     Running               0 172.19.0.2    pgcluster-worker2          3m47s     │
│ kube-system          kube-proxy-kf7bx                                   ●    1/1     Running               0 172.19.0.4    pgcluster-control-plane    3m51s     │
│ kube-system          kube-proxy-mw922                                   ●    1/1     Running               0 172.19.0.3    pgcluster-worker           3m47s     │
│ kube-system          kube-scheduler-pgcluster-control-plane             ●    1/1     Running               0 172.19.0.4    pgcluster-control-plane    4m6s      │
│ local-path-storage   local-path-provisioner-988d74bc-d7dhq              ●    1/1     Running               0 10.244.0.4    pgcluster-control-plane    3m50s     │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
│                                                                                                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```


## Setup Minio

A S3 compatible object storage is required for backup/restore with cloudnativePG. That is why Min.io is setup.

```sh
helm repo add minio https://charts.min.io/
helm delete -n minio minio
helm install -n minio --create-namespace --set "rootUser=rootuser,rootPassword=rootpass123,replicas=1,resources.requests.memory=512Mi,mode=standalone,consoleIngress.enabled=true,consoleIngress.hosts[0]=minio.127.0.0.1.nip.io" minio minio/minio
open https://minio.127.0.0.1.nip.io/
```


kubectl apply -f https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/release-1.23/releases/cnpg-1.23.1.yaml

kubectl get deploy -n cnpg-system cnpg-controller-manager



kubectl get deploy -n cnpg-system cnpg-controller-manager -o wide

docker run -p 9000:9000 -p 9001:9001 \
           -e MINIO_ROOT_USER=admin \
           -e MINIO_ROOT_PASSWORD=password \
           --rm \
           minio/minio server /data \
           --console-address ":9001" 
kubectl create secret generic minio-creds \
  --from-literal=MINIO_ACCESS_KEY=admin \
  --from-literal=MINIO_SECRET_KEY=password



kubectl-cnpg status cluster-example

kubectl exec -i cluster-example-1 -- psql < sqltest.sql


kubectl port-forward --namespace default svc/cluster-example-rw 5432:5432 &





kubectl-cnpg promote cluster-example cluster-example-2


kubectl apply -f pgcluster-minorupdate-backup.yaml
