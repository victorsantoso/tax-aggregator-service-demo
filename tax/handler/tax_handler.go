package handler

import (
	"log"
	"net/http"
	"tax-aggregator-service-demo/tax/domain"

	"github.com/labstack/echo/v4"
)

type taxHandler struct {
	taxUsecase domain.TaxUsecase
}

func NewTaxHandler(taxUsecase domain.TaxUsecase) domain.TaxHandler {
	return &taxHandler{
		taxUsecase: taxUsecase,
	}
}

func (th *taxHandler) Routes(echo *echo.Echo) {
	echo.GET("/tax", th.GetTax)
}

func (th *taxHandler) GetTax(ctx echo.Context) error {
	taxDate := &domain.TaxDate{}
	err := echo.QueryParamsBinder(ctx).
		MustInt64("start_date", &taxDate.StartDate).
		MustInt("amount_of_days", &taxDate.AmountOfDays).
		BindError()
	if err != nil {
		log.Println("[TaxHandler.GetTax]:: error bind query params")
		return ctx.JSON(http.StatusBadRequest, &domain.Response{
			Code:    http.StatusInternalServerError,
			Message: "bad request",
		})
	}
	tax, err := th.taxUsecase.GetTax(taxDate)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &domain.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, &domain.Response{
		Code:    http.StatusOK,
		Message: "success get tax",
		Data:    tax,
	})
}
