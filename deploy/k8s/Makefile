.PHONY: cluster-up cluster-label cluster-setup cluster-down cluster-status \
        vault-configure pull-secret deploy teardown

CLUSTER_NAME=sre-bootcamp
NODES=3
DRIVER=docker

-include .env

cluster-up:
	minikube start --profile $(CLUSTER_NAME) --nodes $(NODES) --driver $(DRIVER)

cluster-label:
	minikube kubectl --profile $(CLUSTER_NAME) -- label node $(CLUSTER_NAME)     type=application
	minikube kubectl --profile $(CLUSTER_NAME) -- label node $(CLUSTER_NAME)-m02 type=database
	minikube kubectl --profile $(CLUSTER_NAME) -- label node $(CLUSTER_NAME)-m03 type=dependent_services

cluster-setup: cluster-up cluster-label

cluster-down:
	minikube delete --profile $(CLUSTER_NAME)

cluster-status:
	minikube kubectl --profile $(CLUSTER_NAME) -- get nodes

vault-configure:
	kubectl wait pod/vault-0 -n vault --for=condition=Ready --timeout=120s
	kubectl exec -n vault vault-0 -- vault kv put secret/student-api/db \
		username=pguser password=pgpass

pull-secret:
	kubectl create namespace student-api --dry-run=client -o yaml | kubectl apply -f -
	kubectl create secret docker-registry ghcr-pull-secret \
		--docker-server=ghcr.io \
		--docker-username=$(GHCR_USER) \
		--docker-password=$(GHCR_TOKEN) \
		--namespace=student-api \
		--dry-run=client -o yaml | kubectl apply -f -

deploy:
	helmfile apply
	kubectl wait pod/vault-0 -n vault --for=condition=Ready --timeout=120s
	$(MAKE) vault-configure
	kubectl wait --for=condition=Established --timeout=120s crd/clustersecretstores.external-secrets.io
	kubectl wait deployment/external-secrets-webhook -n external-secrets --for=condition=Available --timeout=120s
	kubectl apply -f manifests/eso-config.yaml
	$(MAKE) pull-secret
	kubectl apply -f manifests/api.yaml

teardown:
	kubectl delete namespace student-api --ignore-not-found
	kubectl delete clustersecretstore vault-backend --ignore-not-found 2>/dev/null || true
	kubectl delete secret vault-token -n external-secrets --ignore-not-found
	helmfile destroy
