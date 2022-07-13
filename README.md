# Kubernetes Operator for Cascade Job starting

## Create manifest

`make manifests`

## Create CRD

`make generate`

## Image build

`docker build -t ghcr.io/randsw/cascademanualoperator .`

## Image push

`docker login`

`docker push ghcr.io/randsw/cascademanualoperator`

## Image deploy 

`make deploy IMG=ghcr.io/randsw/cascademanualoperator`

## Delete operator

`make undeploy`

## Deploy using helm

`make install-helm helm-namespace=<your-namespace>`
