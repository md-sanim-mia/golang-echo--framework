package product

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ProductContller struct {
	Service *ProductService
}

func (c *ProductContller) CreateProduct(ctx echo.Context) error {

	product := new(Product)

	if err := ctx.Bind(product); err != nil {

		return ctx.JSON(http.StatusBadRequest, echo.Map{

			"message": "Invalid request",
		})
	}

	if err := c.Service.CreateProduct(product); err != nil {

		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, echo.Map{
		"success": true,
		"message": "Proudct created success fully!",
		"data":    product,
	})
}
