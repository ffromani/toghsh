# togsh - extracts shell commands from github actions workflows

`toghsh` is a helper tool to translate github action workflows into equivalent shell scripts.
`toghsh` cannot, and *will not* translate github actions into shell-script equivalent.

## When I should use toghsh

In the narrow but relevant use case on which you run end-to-end (e2e) tests as github workflows,
relying on shell commands to setup and run the suite, then `toghsh` can help you.

You can keep the master source of the e2e tests as github workflows, and use `toghsh` to generate
a runnable shell script for your local/custom worker.

You are advised and strongly encouraged to look at [act](https://github.com/nektos/act) for a more
comprehensive solution.

## Usage

1. fetch the `toghsh` source tree
2. run `make`
3. copy the binary (`_out/toghsh` by default) in your `$PATH`

## Example

```bash
$ ./_out/toghsh --help
Usage of ./_out/toghsh:
  -J, --job-id string   process job
  -L, --list            list available jobs and exit
pflag: help requested

$ curl -L https://raw.githubusercontent.com/k8stopologyawareschedwg/resource-topology-exporter/master/.github/workflows/e2e.yml -o example.yaml
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  3518  100  3518    0     0  14242      0 --:--:-- --:--:-- --:--:-- 14242

$ ./_out/toghsh --list example.yaml 
job: "e2e-ic"
job: "e2e-ip"

$ ./_out/toghsh --job-id e2e-ic example.yaml 
### setup environment as per job "e2e-ic"
export RTE_CONTAINER_IMAGE=quay.io/k8stopologyawarewg/resource-topology-exporter:ci

### order=0 - ID="" - name="checkout sources"
### nothing to run

### order=1 - ID="go" - name="setup golang"
### nothing to run

### order=2 - ID="" - name="build test binary"
make build-e2e

### order=3 - ID="" - name="build image"
RTE_CONTAINER_IMAGE=${RTE_CONTAINER_IMAGE} RUNTIME=docker make image

### order=4 - ID="" - name="generate manifests"
RTE_CONTAINER_IMAGE=${RTE_CONTAINER_IMAGE} RTE_POLL_INTERVAL=10s make gen-manifests | tee rte-e2e.yaml

### order=5 - ID="" - name="create K8S kind cluster"
# kind is part of 20.04 image, see: https://github.com/actions/virtual-environments/blob/main/images/linux/Ubuntu2004-README.md
kind create cluster --config=hack/kind-config-e2e.yaml --image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6
kind load docker-image ${RTE_CONTAINER_IMAGE}

### order=6 - ID="" - name="deploy RTE"
# TODO: what about the other workers (if any)?
kubectl label node kind-worker node-role.kubernetes.io/worker=''
kubectl create -f rte-e2e.yaml

### order=7 - ID="" - name="cluster info"
kubectl get nodes
kubectl get pods -A
kubectl describe pod -l name=resource-topology || :
kubectl logs -l name=resource-topology -c resource-topology-exporter-container || :

### order=8 - ID="" - name="cluster ready"
hack/check-ds.sh
kubectl logs -l name=resource-topology -c resource-topology-exporter-container || :
kubectl get noderesourcetopologies.topology.node.k8s.io -A -o yaml

### order=9 - ID="" - name="run E2E tests"
export KUBECONFIG=${HOME}/.kube/config 
_out/rte-e2e.test -ginkgo.focus='\[(RTE|TopologyUpdater)\].*\[(Local|InfraConsuming)\]'
```
