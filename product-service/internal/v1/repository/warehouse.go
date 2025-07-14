package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/product-service/internal/v1/model"
	"github.com/kadekchresna/ecommerce/product-service/internal/v1/repository/dao"
	repository_interface "github.com/kadekchresna/ecommerce/product-service/internal/v1/repository/interface"
)

type warehouseRepository struct {
	baseURL     string
	staticToken string
}

func NewWarehouseRepository(
	baseURL string,
	staticToken string,
) repository_interface.IWarehouseRepository {
	return &warehouseRepository{
		baseURL:     baseURL,
		staticToken: staticToken,
	}
}

func (r *warehouseRepository) GetProductStock(ctx context.Context, productUUID uuid.UUID) (*model.ProductStock, error) {

	endpoint := fmt.Sprintf("%s/warehouse/stock/%s", r.baseURL, productUUID.String())

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("X-App-Token", "Bearer "+r.staticToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var result dao.ProductStockResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode > http.StatusOK {
		return nil, fmt.Errorf("http request failed: %s", result.Message)
	}

	if result.Data.WarehouseName == "" {
		return nil, nil
	}

	return &model.ProductStock{
		ProductUUID:     result.Data.ProductUUID,
		WarehouseUUID:   result.Data.WarehouseUUID,
		WarehouseName:   result.Data.WarehouseName,
		ShopUUID:        result.Data.ShopUUID,
		Status:          result.Data.Status,
		ReserveQuantity: result.Data.ReserveQuantity,
		StartQuantity:   result.Data.StartQuantity,
	}, nil

}
