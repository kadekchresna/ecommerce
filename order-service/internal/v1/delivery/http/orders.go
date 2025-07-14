package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/order-service/helper/jwt"
	"github.com/kadekchresna/ecommerce/order-service/helper/logger"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/dto"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
	usecase_interface "github.com/kadekchresna/ecommerce/order-service/internal/v1/usecase/interface"
	"github.com/labstack/echo/v4"
)

type ordersHandler struct {
	OrdersUsecase usecase_interface.IOrdersUsecase
}

func NewOrdersHandler(
	g *echo.Group,
	OrdersUsecase usecase_interface.IOrdersUsecase,
) {
	u := &ordersHandler{
		OrdersUsecase: OrdersUsecase,
	}

	v1Order := g.Group("/order", jwt.JWTMiddleware)

	v1Order.POST("/checkout", u.Checkout)
	v1Order.POST("/run-outbox", u.ProcessOutbox)
	v1Order.POST("/run-inbox", u.ProcessInbox)
	v1Order.POST("/set-completed/:uuid", u.UpdateOrderCompleted)
	v1Order.POST("/set-expired", u.UpdateOrderExpired)

}

func (h *ordersHandler) Checkout(c echo.Context) error {
	req := dto.CheckoutRequest{}
	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)
	userUUID, _ := c.Get(jwt.USER_UUID_KEY).(uuid.UUID)

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid input", "request_id": requestID})
	}

	req.UserUUID = userUUID
	if err := h.OrdersUsecase.Checkout(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "checkout failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "checkout success", "request_id": requestID})
}

func (h *ordersHandler) ProcessOutbox(c echo.Context) error {

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	err := h.OrdersUsecase.ProcessOutbox(ctx, &dto.ProcessOutboxRequest{
		Statuses: []model.OutboxStatusType{model.OutboxStatusCreated, model.OutboxStatusFailed},
		Limit:    10,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "process outbox failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "process outbox successfully", "request_id": requestID})
}

func (h *ordersHandler) ProcessInbox(c echo.Context) error {

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	err := h.OrdersUsecase.ProcessInbox(ctx, &dto.ProcessInboxRequest{
		Statuses: []model.InboxStatusType{model.InboxStatusCreated, model.InboxStatusFailed},
		Limit:    10,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "process inbox failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "process inbox successfully", "request_id": requestID})
}

func (h *ordersHandler) UpdateOrderCompleted(c echo.Context) error {

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	orderUUIDString := c.Param("uuid")

	orderUUID, err := uuid.Parse(orderUUIDString)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid order uuid format", "request_id": requestID})
	}

	err = h.OrdersUsecase.UpdateOrderStatusCompleted(ctx, orderUUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "update order failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "update order successfully", "request_id": requestID})
}

func (h *ordersHandler) UpdateOrderExpired(c echo.Context) error {

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	err := h.OrdersUsecase.UpdateOrderStatusExpired(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "update order failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "update order successfully", "request_id": requestID})
}
