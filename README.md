# EnterpriseDB Operator on MiniKube

## Setup Minikube on Mac

```sh
brew install minikube kubectl beekeeper-studio
minikube start
kubectl get po -A
minikube dashboard
```

## Setup simple PostgreSQL

```sh
kubectl apply -f 1_simple/postgres-deploy.yaml
kubectl get all
minikube service postgres
```

After opening a tunnel to the `postgres` service, open your SQL DB management tool - e.g. Beekeeper Studio or SQLElectron - and connect using the provided port from the URL and the credentials from the deployment yaml.

```
Host: localhost
Port: 51045  (use from minikube service postgres output)
User: root
Password: REDACTED
Default Database: mydb
```

## Setup simple PGAdmin

```sh
kubectl apply -f 1_simple/pgadmin-deploy.yaml
kubectl get all
minikube service pgadmin
```

Login using the credentials from the deployment yaml.

```
User: admin@admin.com
Password: REDACTED
```

Configure the connection for `pgAdmin` to run locally

```
Host: 192.168.49.2
Port: 30432
User: root
Password: REDACTED
Maintenance Database: mydb
```

