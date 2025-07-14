package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kadekchresna/ecommerce/warehouse-service/helper/logger"
	helper_messaging "github.com/kadekchresna/ecommerce/warehouse-service/helper/messaging"
	"github.com/kadekchresna/ecommerce/warehouse-service/infrastructure/messaging"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/model"
	usecase_interface "github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/usecase/interface"
)

type consumerHandler struct {
	WarehouseUsecase usecase_interface.IWarehouseUsecase
}

func NewConsumerHandler(
	WarehouseUsecase usecase_interface.IWarehouseUsecase,
) *consumerHandler {
	return &consumerHandler{
		WarehouseUsecase: WarehouseUsecase,
	}
}

func (h *consumerHandler) HandleMessage(ctx context.Context, msg messaging.Message) error {

	switch msg.Topic {
	case helper_messaging.TOPIC_ORDER_EVENTS:

		return h.OrderStockRequestEvent(ctx, msg)

	default:
		logger.LogWithContext(ctx).Info(fmt.Sprintf("Unknown topic %s", msg.Topic))

	}

	return nil

}

func (h *consumerHandler) OrderStockRequestEvent(ctx context.Context, msg messaging.Message) error {

	payload := model.OutboxOrderCreatedMetaRequest{}

	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		err = fmt.Errorf("error Marshal :: warehouseHandler.OrderStockEvent(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}

	err := h.WarehouseUsecase.StoreToInbox(ctx, model.Inbox{
		Metadata:  string(msg.Value),
		Type:      helper_messaging.TOPIC_ORDER_EVENTS,
		Reference: payload.Order.UUID.String(),
		Status:    model.InboxStatusCreated,
		Action:    payload.Action,
		Response:  "{}",
	})
	if err != nil {
		err = fmt.Errorf("error query create :: WarehouseUsecase.StoreToInbox(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}
	return nil
}
