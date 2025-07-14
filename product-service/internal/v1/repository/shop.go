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

type shopRepository struct {
	baseURL     string
	staticToken string
}

func NewShopRepository(
	baseURL string,
	staticToken string,
) repository_interface.IShopRepository {
	return &shopRepository{
		baseURL:     baseURL,
		staticToken: staticToken,
	}
}

func (r *shopRepository) GetShop(ctx context.Context, shopUUID uuid.UUID) (*model.Shop, error) {

	endpoint := fmt.Sprintf("%s/shops/%s", r.baseURL, shopUUID.String())

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
	var result dao.ShopResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode > http.StatusOK {
		return nil, fmt.Errorf("http request failed: %s", result.Message)
	}

	if result.Data.Name == "" {
		return nil, nil
	}

	return &model.Shop{
		UUID: result.Data.UUID,
		Code: result.Data.Code,
		Name: result.Data.Name,
		Desc: result.Data.Desc,
	}, nil

}
