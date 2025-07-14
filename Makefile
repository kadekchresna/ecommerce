docker-push:
	docker build -t ghcr.io/kadekchresna/order-service:latest order-service/ && \
	docker build -t ghcr.io/kadekchresna/product-service:latest product-service/ && \
	docker build -t ghcr.io/kadekchresna/shop-service:latest shop-service/ && \
	docker build -t ghcr.io/kadekchresna/user-service:latest user-service/ && \
	docker build -t ghcr.io/kadekchresna/warehouse-service:latest warehouse-service/ && \
	echo $$G_PAT | docker login ghcr.io -u kadekchresna --password-stdin && \
	docker push ghcr.io/kadekchresna/order-service:latest && \
	docker push ghcr.io/kadekchresna/product-service:latest && \
	docker push ghcr.io/kadekchresna/shop-service:latest &&  \
	docker push ghcr.io/kadekchresna/user-service:latest && \
	docker push ghcr.io/kadekchresna/warehouse-service:latest


init-bitnami:
	helm repo add bitnami https://charts.bitnami.com/bitnami 
	helm repo update

kafka-up: init-bitnami
	helm install kafka bitnami/kafka -f infra/charts/kafka/values.yaml 


ecommerce-up: kafka-up redis-up
	helm dependency update infra/charts/warehouse-service
	helm install warehouse-service infra/charts/warehouse-service -f infra/charts/warehouse-service/values.yaml
	helm dependency update infra/charts/shop-service
	helm install shop-service infra/charts/shop-service -f infra/charts/shop-service/values.yaml
	helm dependency update infra/charts/user-service
	helm install user-service infra/charts/user-service -f infra/charts/user-service/values.yaml
	helm dependency update infra/charts/product-service
	helm install product-service infra/charts/product-service -f infra/charts/product-service/values.yaml
	helm dependency update infra/charts/order-service
	helm install order-service infra/charts/order-service -f infra/charts/order-service/values.yaml
