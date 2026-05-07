.PHONY: cluster-up cluster-label cluster-setup cluster-down cluster-status

CLUSTER_NAME=sre-bootcamp
NODES=3
DRIVER=docker

cluster-up:
	minikube start --profile $(CLUSTER_NAME) --nodes $(NODES) --driver $(DRIVER)

cluster-label:
	minikube kubectl --profile $(CLUSTER_NAME) -- label node $(CLUSTER_NAME)       type=application
	minikube kubectl --profile $(CLUSTER_NAME) -- label node $(CLUSTER_NAME)-m02   type=database
	minikube kubectl --profile $(CLUSTER_NAME) -- label node $(CLUSTER_NAME)-m03   type=dependent_services

cluster-setup: cluster-up cluster-label

cluster-down:
	minikube delete --profile $(CLUSTER_NAME)

cluster-status:
	minikube kubectl --profile $(CLUSTER_NAME) -- get nodes
