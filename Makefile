docker-push:
	docker build -t ghcr.io/kadekchresna/order-service:v1 order-service/ && \
	docker build -t ghcr.io/kadekchresna/product-service:v1 product-service/ && \
	docker build -t ghcr.io/kadekchresna/shop-service:v1 shop-service/ && \
	docker build -t ghcr.io/kadekchresna/user-service:v1 user-service/ && \
	docker build -t ghcr.io/kadekchresna/warehouse-service:v1 warehouse-service/ && \
	echo $$G_PAT | docker login ghcr.io -u kadekchresna --password-stdin && \
	docker push ghcr.io/kadekchresna/order-service:v1 && \
	docker push ghcr.io/kadekchresna/product-service:v1 && \
	docker push ghcr.io/kadekchresna/shop-service:v1 &&  \
	docker push ghcr.io/kadekchresna/user-service:v1 && \
	docker push ghcr.io/kadekchresna/warehouse-service:v1 
