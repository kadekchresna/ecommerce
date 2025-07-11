package handler

import (
	"net/http"

	"github.com/kadekchresna/go-boilerplate/helper/logger"
	usecase_interface "github.com/kadekchresna/go-boilerplate/internal/v1/usecase/interface"
	"github.com/labstack/echo/v4"
)

type usersHandler struct {
	UsersUsecase usecase_interface.IUsersUsecase
}

func NewUsersHandler(
	g *echo.Group,
	UsersUsecase usecase_interface.IUsersUsecase,
) {
	u := &usersHandler{
		UsersUsecase: UsersUsecase,
	}

	v1User := g.Group("/user")

	v1User.GET("", u.Get)

}

func (h *usersHandler) Get(c echo.Context) error {

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	res, err := h.UsersUsecase.Get()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "get user data successfully failed", "error": err.Error(), "request_id": requestID})
	}

	if res == nil {
		return c.JSON(http.StatusOK, echo.Map{"message": "user is not found", "request_id": requestID, "data": nil})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "get user data successfully", "request_id": requestID, "data": res})
}
