package helper_messaging

const (
	TOPIC_ORDER_EVENTS     = "order-events"
	TOPIC_WAREHOUSE_EVENTS = "warehouse-events"

	ACTION_RESERVE_SUCCESS                   = "reserve-success"
	ACTION_ORDER_EXPIRED                     = "order-expired"
	ACTION_RESERVE_FAILED_INSUFFICIENT_STOCK = "reserve-failed-insufficient-stock"
	ACTION_RESERVE_FAILED_WAREHOUSE_INACTIVE = "reserve-failed-warehouse-inactive"

	ORDER_CONSUMER_GROUP     = "order-service-consumer-group"
	WAREHOUSE_CONSUMER_GROUP = "warehouse-service-consumer-group"
)
