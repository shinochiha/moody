package routes

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/moody/controllers"
	"github.com/moody/middlewares"
)

func Routes(db *gorm.DB) *echo.Echo {

	e := echo.New()
	e.HTTPErrorHandler = middlewares.ErrorHandler
	e.Use(middlewares.TransactionHandler(db))

	e.POST("/api/v1/register", controllers.Register)
	e.POST("/api/v1/login", controllers.Login)
	e.POST("/api/v1/forgot_password", controllers.ForgotPassword)
	e.POST("/api/v1/reset_password", controllers.ResetPassword)
	e.GET("/api/v1/reset_password/:id", controllers.ResetPassword)
	e.GET("/api/v1/emailver/:id", controllers.EmailVer)
	e.GET("/api/v1/version", controllers.Version)

	// e.GET("/api/v1/users/:id", controllers.UserGetByIdApiHandle, middlewares.UserRoleValidation)
	// e.GET("/api/v1/users", controllers.UserGetPaginatedApiHandle, middlewares.UserRoleValidation)
	// e.PUT("/api/v1/users/:id", controllers.UserUpdateByIdApiHandle, middlewares.UserRoleValidation)
	// e.PATCH("/api/v1/users/:id", controllers.UserPartialUpdateByIdApiHandle, middlewares.UserRoleValidation)
	// e.DELETE("/api/v1/users/:id", controllers.UserDeleteByIdApiHandle, middlewares.UserRoleValidation)

	g := e.Group("/api/v1")
	g.Use(middlewares.JWTMiddleware)
	e.GET("/api/v1/logout", controllers.Logout)

	g.GET("/users/:id", controllers.UserGetById)
	g.GET("/users", controllers.UserGetPaginated)
	g.PUT("/users/:id", controllers.UserUpdateById)
	g.DELETE("/users/:id", controllers.UserDeleteById)

	return e
}
