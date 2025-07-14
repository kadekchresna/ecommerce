package standart_error

import "errors"

var (
	ErrorWarehouseInactive = errors.New("warehouse is not active")
	ErrorInsufficientStock = errors.New("insufficient stock")
)
