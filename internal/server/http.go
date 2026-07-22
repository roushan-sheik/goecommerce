package server

import (
	"goecommerce/internal/config"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Start(db *gorm.DB, cfg *config.Config) {
	db.AutoMigrate()

	e := echo.New()

	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	e.Start(":" + cfg.Port)
}
