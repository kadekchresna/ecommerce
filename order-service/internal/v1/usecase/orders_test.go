package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	helper_messaging "github.com/kadekchresna/ecommerce/order-service/helper/messaging"
	helper_time "github.com/kadekchresna/ecommerce/order-service/helper/time"
	helper_uuid "github.com/kadekchresna/ecommerce/order-service/helper/uuid"
	mock_lock "github.com/kadekchresna/ecommerce/order-service/infrastructure/lock/mocks"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/dto"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
	mock_repository_interface "github.com/kadekchresna/ecommerce/order-service/internal/v1/repository/interface/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_ordersUsecase_Checkout(t *testing.T) {
	uuid := uuid.New()
	now := time.Now()

	p := model.Products{
		UUID:  uuid,
		Title: "title",
		Price: 10000,
	}

	order := model.Order{
		UUID:        uuid,
		Code:        fmt.Sprintf("ORD-%s", uuid.String()[:8]),
		UserUUID:    uuid,
		TotalAmount: p.Price * 1,
		ExpiredAt:   now.Add(15 * time.Minute),
		Status:      model.OrderStatusCreated,
		CreatedBy:   uuid,
		UpdatedBy:   uuid,
		Metadata:    "{}",
	}

	orderDetails := []model.OrderDetail{
		{
			ProductUUID:  p.UUID,
			ProductTitle: p.Title,
			ProductPrice: p.Price,
			SubTotal:     p.Price * 1,
			Quantity:     1,
		},
	}

	eventPayload := model.OutboxOrderCreatedMetaRequest{
		Order:       order,
		OrderDetail: orderDetails,
		Action:      helper_messaging.ACTION_ORDER_CREATED,
	}

	type args struct {
		ctx     context.Context
		request *dto.CheckoutRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(mockProductRepo *mock_repository_interface.MockIProductRepository, mockOrderRepo *mock_repository_interface.MockIOrdersRepository, mockDistributedLock *mock_lock.MockDistributedLock)
	}{
		{
			name: "Success Checkout",
			beforeFunc: func(mockProductRepo *mock_repository_interface.MockIProductRepository, mockOrderRepo *mock_repository_interface.MockIOrdersRepository, mockDistributedLock *mock_lock.MockDistributedLock) {
				mockProductRepo.EXPECT().GetProduct(mock.Anything, uuid).Return(&p, nil).Once()

				req := &dto.CreateCheckoutRequest{
					EventType:    helper_messaging.TOPIC_ORDER_EVENTS,
					Order:        order,
					OrderDetails: orderDetails,
					EventPayload: eventPayload,
					UserUUID:     uuid,
				}
				mockOrderRepo.EXPECT().Checkout(mock.Anything, req).Return(nil).Once()
			},
			wantErr: false,
			args: args{
				ctx: context.Background(),
				request: &dto.CheckoutRequest{
					Order: dto.Order{
						OrderDetails: []dto.OrderDetails{
							{
								ProductUUID: uuid,
								Quantity:    1,
							},
						},
					},
					UserUUID: uuid,
				},
			},
		},
		{
			name: "Failed Checkout - Error Create Checkout",
			beforeFunc: func(mockProductRepo *mock_repository_interface.MockIProductRepository, mockOrderRepo *mock_repository_interface.MockIOrdersRepository, mockDistributedLock *mock_lock.MockDistributedLock) {
				mockProductRepo.EXPECT().GetProduct(mock.Anything, uuid).Return(&p, nil).Once()

				req := &dto.CreateCheckoutRequest{
					EventType:    helper_messaging.TOPIC_ORDER_EVENTS,
					Order:        order,
					OrderDetails: orderDetails,
					EventPayload: eventPayload,
					UserUUID:     uuid,
				}
				mockOrderRepo.EXPECT().Checkout(mock.Anything, req).Return(errors.New("FATAL ERROR")).Once()
			},
			wantErr: true,
			args: args{
				ctx: context.Background(),
				request: &dto.CheckoutRequest{
					Order: dto.Order{
						OrderDetails: []dto.OrderDetails{
							{
								ProductUUID: uuid,
								Quantity:    1,
							},
						},
					},
					UserUUID: uuid,
				},
			},
		},
		{
			name: "Failed Checkout - Error Product Not Found",
			beforeFunc: func(mockProductRepo *mock_repository_interface.MockIProductRepository, mockOrderRepo *mock_repository_interface.MockIOrdersRepository, mockDistributedLock *mock_lock.MockDistributedLock) {
				mockProductRepo.EXPECT().GetProduct(mock.Anything, uuid).Return(nil, nil).Once()

			},
			wantErr: true,
			args: args{
				ctx: context.Background(),
				request: &dto.CheckoutRequest{
					Order: dto.Order{
						OrderDetails: []dto.OrderDetails{
							{
								ProductUUID: uuid,
								Quantity:    1,
							},
						},
					},
					UserUUID: uuid,
				},
			},
		},
		{
			name: "Failed Checkout - Error Fetch Product",
			beforeFunc: func(mockProductRepo *mock_repository_interface.MockIProductRepository, mockOrderRepo *mock_repository_interface.MockIOrdersRepository, mockDistributedLock *mock_lock.MockDistributedLock) {
				mockProductRepo.EXPECT().GetProduct(mock.Anything, uuid).Return(nil, errors.New("FATAL ERROR")).Once()

			},
			wantErr: true,
			args: args{
				ctx: context.Background(),
				request: &dto.CheckoutRequest{
					Order: dto.Order{
						OrderDetails: []dto.OrderDetails{
							{
								ProductUUID: uuid,
								Quantity:    1,
							},
						},
					},
					UserUUID: uuid,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockOrderRepo := mock_repository_interface.NewMockIOrdersRepository(t)

			mockProductRepo := mock_repository_interface.NewMockIProductRepository(t)

			timeHelper := helper_time.NewTime(&now)
			uuidHeler := helper_uuid.NewUUID(uuid)
			mockLock := mock_lock.NewMockDistributedLock(t)

			tt.beforeFunc(mockProductRepo, mockOrderRepo, mockLock)

			u := NewOrdersUsecase(timeHelper, uuidHeler, mockLock, mockOrderRepo, mockProductRepo)

			if err := u.Checkout(tt.args.ctx, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("ordersUsecase.Checkout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
