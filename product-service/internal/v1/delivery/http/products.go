package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/product-service/helper/jwt"
	"github.com/kadekchresna/ecommerce/product-service/helper/logger"
	"github.com/kadekchresna/ecommerce/product-service/internal/v1/dto"
	usecase_interface "github.com/kadekchresna/ecommerce/product-service/internal/v1/usecase/interface"
	"github.com/labstack/echo/v4"
)

type productsHandler struct {
	ProductsUsecase usecase_interface.IProductsUsecase
}

func NewProductsHandler(
	g *echo.Group,
	ProductsUsecase usecase_interface.IProductsUsecase,
) {
	u := &productsHandler{
		ProductsUsecase: ProductsUsecase,
	}

	v1User := g.Group("/products", jwt.JWTMiddleware)

	v1User.GET("", u.GetProductsPaginate)
	v1User.GET("/:uuid", u.GetProduct)

}

func (h *productsHandler) GetProductsPaginate(c echo.Context) error {

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	req := dto.GetProductsPaginateRequest{}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid input", "request_id": requestID})
	}

	res, err := h.ProductsUsecase.GetProductsPaginate(ctx, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "get products paginate failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "get products paginate success", "request_id": requestID, "data": res})
}

func (h *productsHandler) GetProduct(c echo.Context) error {

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	uuidString := c.Param("uuid")

	productUUID, err := uuid.Parse(uuidString)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid product uuid format", "request_id": requestID})
	}

	res, err := h.ProductsUsecase.GetProduct(ctx, productUUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "get products failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "get products success", "request_id": requestID, "data": res})
}
