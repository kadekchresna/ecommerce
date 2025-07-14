package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/repository/dao"
	repository_interface "github.com/kadekchresna/ecommerce/order-service/internal/v1/repository/interface"
)

type productsRepository struct {
	baseURL     string
	staticToken string
}

func NewProductRepository(
	baseURL string,
	staticToken string,
) repository_interface.IProductRepository {
	return &productsRepository{
		baseURL:     baseURL,
		staticToken: staticToken,
	}
}

func (r *productsRepository) GetProduct(ctx context.Context, productUUID uuid.UUID) (*model.Products, error) {

	endpoint := fmt.Sprintf("%s/products/%s", r.baseURL, productUUID.String())

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

	var result dao.ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode > http.StatusOK {
		return nil, fmt.Errorf("http request failed: %s", result.Message)
	}

	if result.Data.Title == "" {
		return nil, nil
	}

	return &model.Products{
		UUID:  result.Data.UUID,
		Title: result.Data.Title,
		Desc:  result.Data.Desc,
		Price: result.Data.Price,
	}, nil

}
