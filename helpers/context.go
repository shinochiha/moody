package helpers

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type Context echo.Context

type Ctx struct {
	C echo.Context
}

func CheckCookie(ctx echo.Context) string {
	cookie, _ := ctx.Cookie("token")
	tokenString := cookie.Value
	var id string
	if tokenString != "" {
		jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			u := t.Claims.(jwt.MapClaims)
			id = u["ID"].(string)
			// params["created.user.id"] = []string{id}
			return id, nil
		})
	}
	return id
}

func SetContext(c echo.Context) Context {
	return Context(c)
}

func NewCtx(ctx echo.Context) *Ctx {
	return &Ctx{C: SetContext(ctx)}
}

func GetCtx(ctx Context) *Ctx {
	return &Ctx{C: ctx}
}

func (ctx *Ctx) DB() *gorm.DB {
	return ctx.C.Get("db").(*gorm.DB)
}

func (ctx *Ctx) Tx() *gorm.DB {
	return ctx.C.Get("tx").(*gorm.DB)
}

func (ctx *Ctx) UserLang() string {
	lang := ctx.C.Request().Header.Get("Content-Language")
	if lang == "" {
		lang = "en"
	}
	return lang
}

func (ctx *Ctx) ErrorMessage() map[string]interface{} {
	ErrorMessage := map[string]interface{}{}
	CtxErrorMessage := ctx.C.Get("error_message")
	if CtxErrorMessage != nil {
		ErrorMessage = CtxErrorMessage.(map[string]interface{})
	}
	return ErrorMessage
}

func (ctx *Ctx) SetErrorMessage(ErrorMessage map[string]interface{}) {
	ctx.C.Set("error_message", ErrorMessage)
}
