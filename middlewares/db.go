package middlewares

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo/v4"
)

func TransactionHandler(db *gorm.DB) echo.MiddlewareFunc {
  return func(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
      c.Set("db", db)
      tx := db.Begin()
      c.Set("tx", tx)
      n := next(c)
      if c.Response().Status >= 200 && c.Response().Status < 400 {
        tx.Commit()
      } else {
        tx.Rollback()
      }
      return n
    }
  }
}
