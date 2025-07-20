package repository

import (
	"context"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	helper_messaging "github.com/kadekchresna/ecommerce/order-service/helper/messaging"
	helper_time "github.com/kadekchresna/ecommerce/order-service/helper/time"
	helper_uuid "github.com/kadekchresna/ecommerce/order-service/helper/uuid"
	helper_db "github.com/kadekchresna/ecommerce/order-service/infrastructure/db/helper"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/dto"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/repository/dao"
)

func Test_ordersRepository_Checkout(t *testing.T) {

	uuid := uuid.New()
	o := model.Order{
		UUID:   uuid,
		Status: model.OrderStatusCreated,
	}
	od := model.OrderDetail{
		OrderUUID: uuid,
	}

	ods := []model.OrderDetail{
		od,
	}

	meta := model.OutboxOrderCreatedMetaRequest{
		MessageID: uuid.String(),
	}

	metaB, _ := json.Marshal(meta)
	outbox := dao.OutboxDAO{
		UUID:      uuid,
		Status:    model.OutboxStatusCreated,
		Action:    helper_messaging.ACTION_ORDER_CREATED,
		Type:      helper_messaging.TOPIC_ORDER_EVENTS,
		Metadata:  string(metaB),
		Reference: uuid.String(),
		Response:  "{}",
	}

	now := time.Now()
	type args struct {
		ctx     context.Context
		request *dto.CreateCheckoutRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "orders" ("code","metadata","user_uuid","total_amount","expired_at","status","created_by","updated_by","uuid") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "uuid","created_at","updated_at"`)).WithArgs(o.Code, o.Metadata, o.UserUUID, o.TotalAmount, o.ExpiredAt, o.Status, o.CreatedBy, o.UpdatedBy, o.UUID).WillReturnRows(sqlmock.NewRows([]string{"uuid", "created_at", "updated_at"}).AddRow(uuid, now, now))

				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "orders_detail" ("product_uuid","product_title","product_price","quantity","sub_total","order_uuid") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "uuid"`)).WithArgs(od.ProductUUID, od.ProductTitle, od.ProductPrice, od.Quantity, od.SubTotal, od.OrderUUID).WillReturnRows(sqlmock.NewRows([]string{"uuid"}).AddRow(uuid))

				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "outbox" ("uuid","metadata","response","status","type","action","reference","retry_count") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "created_at","updated_at"`)).WithArgs(outbox.UUID, outbox.Metadata, outbox.Response, outbox.Status, outbox.Type, outbox.Action, outbox.Reference, outbox.RetryCount).WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).AddRow(now, now))

				mockDB.ExpectCommit()
			},
			args: args{
				ctx: context.Background(),
				request: &dto.CreateCheckoutRequest{
					Order:        o,
					OrderDetails: ods,
					EventType:    helper_messaging.TOPIC_ORDER_EVENTS,
					EventPayload: model.OutboxOrderCreatedMetaRequest{},
					UserUUID:     uuid,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)

			timeHelper := helper_time.NewTime(&now)
			uuidHeler := helper_uuid.NewUUID(uuid)

			r := NewOrdersRepository(db, timeHelper, uuidHeler, nil)
			if err := r.Checkout(tt.args.ctx, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("ordersRepository.Checkout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
