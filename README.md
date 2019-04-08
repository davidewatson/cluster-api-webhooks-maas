# cluster-api-webhook-maas

An experimental webhook for allocating machines using MAAS for use with the Cluster API.

# Development

## Rebuild changed artifacts

```
export IMG=quay.io/samsung_cnct/cluster-api-webhook-maas
make docker-build
make docker-push
```

## Deploy webhook

```
kind create cluster
export KUBECONFIG="$(kind get kubeconfig-path --name="1")"
make deploy
```

## Install and verify `maas` command

```
docker run -it --rm ubuntu:bionic-20190307 /bin/bash
export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get install maas-cli jq -yq
```

```
root@907cbf6a59e2:/# maas login smithtower http://192.168.2.24:5240/MAAS/ <REDACTED-TOKEN>

You are now logged in to the MAAS server at
http://192.168.2.24:5240/MAAS/api/2.0/ with the profile name
'smithtower'.

For help with the available commands, try:

  maas smithtower --help

root@907cbf6a59e2:/# maas smithtower nodes read | jq '.[].ip_addresses' | head -n2
[
  "192.168.2.1"
root@907cbf6a59e2:/#
```

## Recreate repository

```
cd "${GOPATH}" 
mkdir -p src/github.com/davidewatson/cluster-api-webhooks-maas
cd src/github.com/davidewatson/cluster-api-webhooks-maas

git init
```

```
# Respond  "y" when prompted "Run `dep ensure` to fetch dependencies (Recommended) [y/n]?""
kubebuilder init --domain cluster.k8s.io --license apache2 --owner "The Kubernetes Authors"

git add .
git commit -m "Generate scaffolding."
```

```
kubebuilder alpha webhook --group generic --version v1alpha1 --kind Machine --type=mutatng --operations=create,delete --make=false

git add .
git commit -m "Generate webhook for Machine resource."
```

```
# Edit imports of github.com/davidewatson/cluster-api-webhooks-maas/pkg/apis/generic/v1alpha1
# to be github.com/cattlek8s/cluster-api-provider-generic/pkg/apis/generic/v1alpha1
vi pkg/webhook/default_server/machine/mutatng/create_delete_webhook.go
vi pkg/webhook/default_server/machine/mutatng/machine_create_delete_handler.go 
dep ensure

git add -u
git commit -m "Fixup imports to pull from the upstream generic provider."
```


