docker-push:
	docker build -t ghcr.io/kadekchresna/orders-service:latest order-service/ && \
	docker build -t ghcr.io/kadekchresna/product-service:latest product-service/ && \
	docker build -t ghcr.io/kadekchresna/shop-service:latest shop-service/ && \
	docker build -t ghcr.io/kadekchresna/user-service:latest user-service/ && \
	docker build -t ghcr.io/kadekchresna/warehouse-service:latest warehouse-service/ && \
	echo $$G_PAT | docker login ghcr.io -u kadekchresna --password-stdin && \
	docker push ghcr.io/kadekchresna/orders-service:latest && \
	docker push ghcr.io/kadekchresna/product-service:latest && \
	docker push ghcr.io/kadekchresna/shop-service:latest &&  \
	docker push ghcr.io/kadekchresna/user-service:latest && \
	docker push ghcr.io/kadekchresna/warehouse-service:latest

init-bitnami:
	helm repo add bitnami https://charts.bitnami.com/bitnami 
	helm repo update

kafka-up: init-bitnami
	helm upgrade --install kafka bitnami/kafka -f infra/charts/kafka/values.yaml 


ecommerce-up: kafka-up
	helm dependency update infra/charts/warehouse-service
	helm upgrade --install warehouse-service infra/charts/warehouse-service -f infra/charts/warehouse-service/values.yaml
	helm dependency update infra/charts/shop-service
	helm upgrade --install shop-service infra/charts/shop-service -f infra/charts/shop-service/values.yaml
	helm dependency update infra/charts/user-service
	helm upgrade --install user-service infra/charts/user-service -f infra/charts/user-service/values.yaml
	helm dependency update infra/charts/product-service
	helm upgrade --install product-service infra/charts/product-service -f infra/charts/product-service/values.yaml
	helm dependency update infra/charts/order-service
	helm upgrade --install order-service infra/charts/order-service -f infra/charts/order-service/values.yaml

wms-up: kafka-up
	helm dependency update infra/charts/warehouse-service
	helm upgrade --install warehouse-service infra/charts/warehouse-service -f infra/charts/warehouse-service/values.yaml

oms-up: 
	helm dependency update infra/charts/order-service
	helm upgrade --install order-service infra/charts/order-service -f infra/charts/order-service/values.yaml


metrics-up:

	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update

	helm upgrade --install k8s-monitoring prometheus-community/kube-prometheus-stack \
	--reuse-values \
	--set kubelet.enabled=true \
	--set kubelet.serviceMonitor.enabled=true


monitoring-password:
	kubectl --namespace default get secrets k8s-monitoring-grafana -o jsonpath="{.data.admin-password}" | base64 -d ; echo

logging-up:
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo update
	helm install loki-stack grafana/loki-stack \
	--namespace logging --create-namespace \
	--set grafana.enabled=false \
	--set promtail.enabled=true \
	--set loki.persistence.enabled=true \
	--set loki.persistence.size=10Gi


unit-test:
	cd order-service/ && go test ./internal/...
	cd warehouse-service/ && go test ./internal/...

integration-test:
	cd order-service/test/integration && go test -v ./... -tags=integration