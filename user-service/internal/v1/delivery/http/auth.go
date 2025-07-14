package handler

import (
	"net/http"

	"github.com/kadekchresna/ecommerce/user-service/helper/logger"
	"github.com/kadekchresna/ecommerce/user-service/internal/v1/helper/dto"
	usecase_interface "github.com/kadekchresna/ecommerce/user-service/internal/v1/usecase/interface"
	"github.com/labstack/echo/v4"
)

type authHandler struct {
	AuthUsecase usecase_interface.IAuthUsecase
}

func NewAuthHandler(
	g *echo.Group,
	AuthUsecase usecase_interface.IAuthUsecase,
) {
	u := &authHandler{
		AuthUsecase: AuthUsecase,
	}

	v1User := g.Group("/auth")

	v1User.POST("/login", u.Login)

}

func (h *authHandler) Login(c echo.Context) error {

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	req := dto.LoginUserRequest{}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid input", "request_id": requestID})
	}

	res, err := h.AuthUsecase.Login(ctx, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "login failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "login success", "request_id": requestID, "data": res})
}
