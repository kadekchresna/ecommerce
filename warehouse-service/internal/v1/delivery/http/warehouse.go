package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/warehouse-service/helper/jwt"
	"github.com/kadekchresna/ecommerce/warehouse-service/helper/logger"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/helper/dto"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/model"
	usecase_interface "github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/usecase/interface"
	"github.com/labstack/echo/v4"
)

type warehouseHandler struct {
	WarehouseUsecase usecase_interface.IWarehouseUsecase
}

func NewWarehouseHandler(
	g *echo.Group,
	WarehouseUsecase usecase_interface.IWarehouseUsecase,
) {
	u := &warehouseHandler{
		WarehouseUsecase: WarehouseUsecase,
	}

	v1Warehouse := g.Group("/warehouse")

	v1Warehouse.GET("/:uuid", u.GetWarehouse)
	v1Warehouse.GET("/stock/:uuid", u.GetProductStock)
	v1Warehouse.POST("/transfer-product", u.TransferProduct)
	v1Warehouse.PATCH("/status", u.UpdateStatusWarehouse)
	v1Warehouse.POST("/run-inbox", u.ProcessInbox)
	v1Warehouse.POST("/run-outbox", u.ProcessOutbox)

}

func (h *warehouseHandler) GetWarehouse(c echo.Context) error {

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	warehouseUUIDString := c.Param("uuid")

	warehouseUUID, err := uuid.Parse(warehouseUUIDString)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid warehouse uuid format", "request_id": requestID})
	}

	res, err := h.WarehouseUsecase.GetWarehouse(ctx, &dto.GetWarehouseRequest{
		UUID: warehouseUUID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "get warehouse data failed", "error": err.Error(), "request_id": requestID})
	}

	if res == nil {
		return c.JSON(http.StatusOK, echo.Map{"message": "warehouse is not found", "request_id": requestID, "data": nil})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "get warehouse data successfully", "request_id": requestID, "data": res})
}

type TransferProductRequest struct {
	ProductUUID         string `json:"product_uuid"`
	TargetWarehouseUUID string `json:"target_warehouse_uuid"`
}

func (h *warehouseHandler) TransferProduct(c echo.Context) error {

	req := TransferProductRequest{}
	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid input", "request_id": requestID})
	}

	productUUID, err := uuid.Parse(req.ProductUUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid product uuid format", "request_id": requestID})
	}

	targetWarehouseUUID, err := uuid.Parse(req.TargetWarehouseUUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid target warehouse uuid format", "request_id": requestID})
	}

	if err := h.WarehouseUsecase.TransferProduct(ctx, &dto.TransferProductRequest{
		ProductUUID:         productUUID,
		TargetWarehouseUUID: targetWarehouseUUID,
	}); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "transfer product failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "transfer product successfully executed", "request_id": requestID})
}

type UpdateStatusWarehouseRequest struct {
	WarehouseUUID string                    `json:"warehouse_uuid"`
	Status        model.WarehouseStatusType `json:"status"`
}

func (h *warehouseHandler) UpdateStatusWarehouse(c echo.Context) error {

	req := UpdateStatusWarehouseRequest{}
	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid input", "request_id": requestID})
	}

	userUUID, _ := c.Get(jwt.USER_UUID_KEY).(uuid.UUID)

	warehouseUUID, err := uuid.Parse(req.WarehouseUUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid warehouse uuid format", "request_id": requestID})
	}

	if err := h.WarehouseUsecase.UpdateStatusWarehouse(ctx, &dto.UpdateStatusWarehouseRequest{
		WarehouseUUID: warehouseUUID,
		Status:        req.Status,
		UserUUID:      userUUID,
	}); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "update warehouse status failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "update warehouse status successfully executed", "request_id": requestID})
}

func (h *warehouseHandler) GetProductStock(c echo.Context) error {

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	productUUIDString := c.Param("uuid")

	productUUID, err := uuid.Parse(productUUIDString)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid product uuid format", "request_id": requestID})
	}

	res, err := h.WarehouseUsecase.GetProductStock(ctx, &dto.GetProductStockRequest{
		ProductUUID: productUUID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "get product stock data failed", "error": err.Error(), "request_id": requestID})
	}

	if res == nil {
		return c.JSON(http.StatusOK, echo.Map{"message": "product stock is not found", "request_id": requestID, "data": nil})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "get product stock successfully", "request_id": requestID, "data": res})
}

func (h *warehouseHandler) ProcessInbox(c echo.Context) error {

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	err := h.WarehouseUsecase.ProcessInbox(ctx, &dto.ProcessInboxRequest{
		Statuses: []model.InboxStatusType{model.InboxStatusCreated, model.InboxStatusFailed},
		Limit:    10,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "process inbox failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "process inbox successfully", "request_id": requestID})
}

func (h *warehouseHandler) ProcessOutbox(c echo.Context) error {

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	err := h.WarehouseUsecase.ProcessOutbox(ctx, &dto.ProcessOutboxRequest{
		Statuses: []model.OutboxStatusType{model.OutboxStatusCreated, model.OutboxStatusFailed},
		Limit:    10,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "process outbox failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "process outbox successfully", "request_id": requestID})
}
