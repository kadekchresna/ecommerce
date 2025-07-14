package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/shop-service/helper/jwt"
	"github.com/kadekchresna/ecommerce/shop-service/helper/logger"
	usecase_interface "github.com/kadekchresna/ecommerce/shop-service/internal/v1/usecase/interface"
	"github.com/labstack/echo/v4"
)

type shopsHandler struct {
	ShopsUsecase usecase_interface.IShopsUsecase
}

func NewShopsHandler(
	g *echo.Group,
	ShopsUsecase usecase_interface.IShopsUsecase,
) {
	u := &shopsHandler{
		ShopsUsecase: ShopsUsecase,
	}

	v1User := g.Group("/shops", jwt.JWTMiddleware)

	v1User.GET("/:uuid", u.GetShops)

}

func (h *shopsHandler) GetShops(c echo.Context) error {

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	shopUUIDString := c.Param("uuid")

	shopUUID, err := uuid.Parse(shopUUIDString)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid shop uuid format", "request_id": requestID})
	}

	res, err := h.ShopsUsecase.GetShops(ctx, shopUUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "get shop data failed", "error": err.Error(), "request_id": requestID})
	}

	if res == nil {
		return c.JSON(http.StatusOK, echo.Map{"message": "shop is not found", "request_id": requestID, "data": nil})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "get shop data successfully", "request_id": requestID, "data": res})
}
