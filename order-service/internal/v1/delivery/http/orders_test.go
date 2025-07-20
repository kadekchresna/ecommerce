package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/order-service/helper/jwt"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/dto"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/usecase/interface/mocks"
	"github.com/labstack/echo/v4"
)

func Test_ordersHandler_Checkout(t *testing.T) {

	userUUID := uuid.New()
	productUUID := uuid.New()
	req := dto.CheckoutRequest{
		Order: dto.Order{
			OrderDetails: []dto.OrderDetails{
				{
					ProductUUID: productUUID,
					Quantity:    1,
				},
			},
		}}

	type args struct {
		body dto.CheckoutRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(u *mocks.MockIOrdersUsecase)
	}{
		{
			name: "Success Checkout",
			args: args{
				body: req,
			},
			wantErr: false,
			beforeFunc: func(u *mocks.MockIOrdersUsecase) {

				reqUsecase := req
				reqUsecase.UserUUID = userUUID
				u.EXPECT().Checkout(context.Background(), &reqUsecase).Return(nil).Once()
			},
		},
		{
			name: "Failed Checkout",
			args: args{
				body: req,
			},
			wantErr: true,
			beforeFunc: func(u *mocks.MockIOrdersUsecase) {

				reqUsecase := req
				reqUsecase.UserUUID = userUUID
				u.EXPECT().Checkout(context.Background(), &reqUsecase).Return(errors.New("FATAL ERROR")).Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockUsecase := mocks.NewMockIOrdersUsecase(t)
			e := echo.New()
			v1 := e.Group("v1")

			jsonBody, _ := json.Marshal(tt.args.body)
			req := httptest.NewRequest(http.MethodPost, "/checkout", bytes.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			tt.beforeFunc(mockUsecase)

			c.Set(jwt.USER_UUID_KEY, userUUID)

			h := NewOrdersHandler(v1, mockUsecase)
			if err := h.Checkout(c); (err != nil) != tt.wantErr {
				t.Errorf("ordersHandler.Checkout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
