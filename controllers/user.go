package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/moody/config"
	"github.com/moody/helpers"
	"github.com/moody/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserGetById(c echo.Context) error {
	m := models.User{}
	res := m.GetById(helpers.SetContext(c), c.Param("id"), c.QueryParams())
	return helpers.Response(c, 200, res)
}

func UserGetPaginated(c echo.Context) error {
	m := models.User{}
	res := m.GetPaginated(helpers.SetContext(c), c.QueryParams())
	return helpers.Response(c, 200, res)
}

func UserUpdateById(c echo.Context) error {
	o := new(models.User)
	if err := c.Bind(o); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	o.ID = helpers.Convert(c.Param("id")).String()
	res := o.UpdateById(helpers.SetContext(c))
	return helpers.Response(c, 200, res)
}

func UserPartialUpdateById(c echo.Context) error {
	return UserUpdateById(c)
}

func UserDeleteById(c echo.Context) error {
	o := new(models.User)
	if err := c.Bind(o); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	o.ID = helpers.Convert(c.Param("id")).String()
	res := o.DeleteById(helpers.SetContext(c))
	return helpers.Response(c, 200, res)
}

func Login(c echo.Context) error {
	// mengambil inputan json
	var userInput models.User
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := helpers.Map{
			"error": helpers.Map{
				"message": err.Error(),
			},
		}
		helpers.Response(c, http.StatusBadRequest, response)
	}

	// ambil data user berdasarkan Email
	var user models.User

	if err := helpers.GetDB(c).Table("users").Where("email = ?", userInput.Email).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := helpers.Map{
				"error": helpers.Map{
					"message": "email atau password salah",
				},
			}
			return helpers.Response(c, http.StatusBadRequest, response)

		default:
			response := helpers.Map{
				"error": helpers.Map{
					"message": "email atau password salah",
				},
			}
			return helpers.Response(c, http.StatusBadRequest, response)
		}
	}

	// cek apakah password valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		response := helpers.Map{"message": "Email atau password salah"}
		return helpers.Response(c, http.StatusBadRequest, response)
	}

	if err := helpers.GetDB(c).Table("users").Where("email = ?", userInput.Email).Where("is_active = ?", true).First(&user).Error; err != nil {
		if err != nil {
			response := helpers.Map{
				"error": helpers.Map{
					"message":        "The account has not been activated, please check your email",
					"user_is_active": false,
				},
			}
			return helpers.Response(c, http.StatusBadRequest, response)
		}
	}

	// proses pembuatan token jwt
	expTime := time.Now().Add(time.Hour * 24)
	claims := &config.JWTClaim{
		ID:       user.ID,
		Email:    user.Email,
		UserName: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// medeklarasikan algoritma yang akan digunakan untuk signing
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// signed token
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		response := helpers.Map{
			"error": helpers.Map{
				"message": err.Error(),
			},
		}
		return helpers.Response(c, http.StatusBadRequest, response)
	}

	// set token yang ke cookie
	http.SetCookie(c.Response().Writer, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	response := helpers.Map{
		"code":             200,
		"message":          "Login Success, the token has been set in the cookies",
		"expired_at_token": expTime,
		"token":            token,
		"user": helpers.Map{
			"is_admin": user.IsAdmin,
		},
	}
	return helpers.Response(c, http.StatusOK, response)

}

func Register(c echo.Context) error {
	o := new(models.User)
	if err := c.Bind(o); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	res := o.Register(helpers.SetContext(c))
	return helpers.Response(c, 201, res)
}

func Logout(c echo.Context) error {
	http.SetCookie(c.Response().Writer, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})
	response := helpers.Map{"message": "Logout Success"}
	return helpers.Response(c, http.StatusOK, response)
}

func EmailVer(c echo.Context) error {
	m := models.User{}
	err := m.EmailVerification(helpers.SetContext(c), c.Param("id"), c.QueryParams())
	if err != nil {
		return err
	}
	return err
}

func ForgotPassword(c echo.Context) error {
	o := new(models.User)
	if err := c.Bind(o); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	res := o.ForgotPassword(helpers.SetContext(c))
	return helpers.Response(c, 201, res)
}

func ResetPassword(c echo.Context) error {
	o := new(models.User)
	if err := c.Bind(o); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	res, _ := o.ResetPassword(helpers.SetContext(c))
	return helpers.Response(c, 201, res)
}
