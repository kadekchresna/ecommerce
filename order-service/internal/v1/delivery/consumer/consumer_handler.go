package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kadekchresna/ecommerce/order-service/helper/logger"
	helper_messaging "github.com/kadekchresna/ecommerce/order-service/helper/messaging"
	"github.com/kadekchresna/ecommerce/order-service/infrastructure/messaging"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
	usecase_interface "github.com/kadekchresna/ecommerce/order-service/internal/v1/usecase/interface"
)

type consumerHandler struct {
	OrdersUsecase usecase_interface.IOrdersUsecase
}

func NewConsumerHandler(
	OrdersUsecase usecase_interface.IOrdersUsecase,
) *consumerHandler {
	return &consumerHandler{
		OrdersUsecase: OrdersUsecase,
	}
}

func (h *consumerHandler) HandleMessage(ctx context.Context, msg messaging.Message) {

	switch msg.Topic {
	case helper_messaging.TOPIC_WAREHOUSE_EVENTS:

		h.OrderStockResponseEvent(ctx, msg)

	default:
		logger.LogWithContext(ctx).Info(fmt.Sprintf("Unknown topic %s", msg.Topic))

	}

}

func (h *consumerHandler) OrderStockResponseEvent(ctx context.Context, msg messaging.Message) error {

	payload := model.OutboxOrderCreatedMetaRequest{}

	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		err = fmt.Errorf("error Marshal :: consumerHandler.OrderStockEvent(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}

	err := h.OrdersUsecase.StoreToInbox(ctx, model.Inbox{
		Metadata:  string(msg.Value),
		Type:      helper_messaging.TOPIC_WAREHOUSE_EVENTS,
		Reference: payload.Order.UUID.String(),
		Status:    model.InboxStatusCreated,
		Action:    payload.Action,
		Response:  "{}",
	})
	if err != nil {
		err = fmt.Errorf("error query create :: OrdersUsecase.StoreToInbox(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}
	return nil
}
