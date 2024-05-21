# EnterpriseDB Operator on MiniKube

## Setup Minikube on Mac

```sh
brew install minikube helm kubectl beekeeper-studio
minikube start
kubectl get po -A
minikube dashboard
```

# Setup Bitnami PostgreSQL with Helm

PostgreSQL instance will be accessible within cluster with DNS `pg-minikube-postgresql.default.svc.cluster.local` and to make it accessible from the outside use `kubectl port-forward`.


```sh
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install pg-minikube --set auth.postgresPassword=<your-postgres-password> bitnami/postgresql
kubectl get pods -n default -o wide
kubectl port-forward --namespace default svc/pg-minikube-postgresql 5432:5432 &
```


# Setup PGAdmin with Helm

```sh
helm repo add runix https://helm.runix.net
helm install pgadmin4 runix/pgadmin4  --set env.password=<your-pgadmin-password>  --set env.email=admin@example.com
export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=pgadmin4,app.kubernetes.io/instance=pgadmin4" -o jsonpath="{.items[0].metadata.name}")
kubectl port-forward $POD_NAME 8080:80
open http://localhost:8080
```

Login using the credentials used previously. Configure the connection for `pgAdmin` using these 

```
Host: pg-minikube-postgresql.default.svc.cluster.local
Port: 5432
User: postgres
Password: <your-postgres-password>
Maintenance Database: postgres
```

