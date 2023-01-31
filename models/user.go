package models

import (
	"encoding/json"
	"html/template"
	"time"

	"github.com/moody/config"
	"github.com/moody/helpers"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type User struct {
	ID              string    `json:"id,omitempty" gorm:"primaryKey;type:char(36)"`
	Username        string    `json:"username,omitempty"`
	Email           string    `json:"email,omitempty" validate:"required|unique" gorm:"type:varchar(100)"`
	Password        string    `json:"password,omitempty" validate:"required" gorm:"type:varchar(255)"`
	ConfirmPassword string    `json:"confirm_password,omitempty" gorm:"type:varchar(255)"`
	PhoneNumber     string    `json:"phone_number,omitempty" gorm:"type:varchar(20)"`
	IsActive        *bool     `json:"is_active,omitempty" gorm:"default:false"`
	IsAdmin         *bool     `json:"is_admin,omitempty" gorm:"default:true"`
	ExpiredAt       time.Time `json:"expired_at,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty" example:"2020-03-16T13:55:09.598136+07:00"`
	UpdatedAt       time.Time `json:"updated_at,omitempty" example:"2020-03-16T13:55:09.598136+07:00"`
}

func (User) TableName() string {
	return "users"
}

func (o *User) Schema() map[string]interface{} {
	return map[string]interface{}{
		"table": map[string]string{"name": "users", "as": "us"},
		"fields": map[string]map[string]string{
			"id":           {"name": "us.id", "as": "id", "type": "string"},
			"username":     {"name": "us.username", "as": "username", "type": "string"},
			"email":        {"name": "us.email", "as": "email", "type": "string"},
			"phone_number": {"name": "us.phone_number", "as": "phone_number", "type": "string"},
			"is_active":    {"name": "us.is_active", "as": "is_active", "type": "string"},
			"expired_at":   {"name": "us.expired_at", "as": "expired_at", "type": "string"},
			"created_at":   {"name": "us.created_at", "as": "created_at"},
			"updated_at":   {"name": "us.updated_at", "as": "updated_at"},
		},
	}
}

func (o *User) GetById(ctx helpers.Context, id string, params map[string][]string) map[string]interface{} {
	return helpers.GetById(ctx, "users", "id", id, params, o.Schema(), map[string]interface{}{})
}

func (o *User) GetPaginated(ctx helpers.Context, params map[string][]string) map[string]interface{} {
	return helpers.GetPaginated(ctx, params, o.Schema(), map[string]interface{}{})
}

func (o *User) UpdateById(ctx helpers.Context) map[string]interface{} {
	helpers.GetDB(ctx).Model(User{}).Where("id = ?", o.ID).Updates(o)
	return o.GetById(ctx, helpers.Convert(o.ID).String(), map[string][]string{})
}

func (o *User) DeleteById(ctx helpers.Context) map[string]interface{} {
	id := helpers.Convert(o.ID).String()
	helpers.GetDB(ctx).Model(User{}).Where("id = ?", o.ID).Delete(&User{})
	return helpers.DeletedMessage("users", "id", id)
}

func (o *User) EmailVerification(ctx helpers.Context, id string, params map[string][]string) error {
	err := helpers.GetDB(ctx).Table("users").Where("id = ?", id).Update("is_active", true).Error
	if err != nil {
		return err
	}

	response := helpers.Map{
		"message": "Verification email success!",
	}

	// fmt.Println(filepath)
	temp, err := template.ParseFiles("views/index.html")
	if err != nil {
		return err
	}

	err = temp.Execute(ctx.Response(), response)
	if err != nil {
		return err
	}
	return err
}

var ts User

func (o *User) ResetPassword(ctx helpers.Context) (res helpers.Map, err error) {
	helpers.GetDB(ctx).Model(User{}).Where("id = ?", ctx.Param("id")).First(&ts)
	if ctx.Request().Method == "GET" {
		temp, err := template.ParseFiles("views/reset_password.html")
		if err != nil {
			panic(err)
		}
		temp.Execute(ctx.Response(), nil)
	}
	if ctx.Request().Method == "POST" {
		ctx.Request().ParseForm()
		var password User
		var data = make(map[string]interface{})

		password.Email = ctx.Request().Form.Get("email")
		if password.Email == "" {
			data["required_email"] = "email is a required field "
			temp, _ := template.ParseFiles("views/reset_password.html")
			temp.Execute(ctx.Response(), data)
			return res, err
		}
		password.Password = ctx.Request().Form.Get("password")
		if password.Password == "" {
			data["required_password"] = "password is a required field "
			temp, _ := template.ParseFiles("views/reset_password.html")
			temp.Execute(ctx.Response(), data)
			return res, err
		}
		password.ConfirmPassword = ctx.Request().Form.Get("confirm_password")
		if password.ConfirmPassword == "" {
			data["required_confirm_password"] = "confirm password is a required field "
			temp, _ := template.ParseFiles("views/reset_password.html")
			temp.Execute(ctx.Response(), data)
			return res, err
		}
		if ts.Email != password.Email {
			data["invalid_email"] = "email with " + "'" + password.Email + "'" + " is invalid"
			temp, _ := template.ParseFiles("views/reset_password.html")
			temp.Execute(ctx.Response(), data)
			return res, err
		}
		err := helpers.GetDB(ctx).Model(User{}).Where("email = ?", password.Email).First(&o).Error
		if err != nil {
			data["invalid_email"] = "email with " + "'" + password.Email + "'" + " not found"
			temp, _ := template.ParseFiles("views/reset_password.html")
			temp.Execute(ctx.Response(), data)
			return res, err
		}
		if password.Password != password.ConfirmPassword {
			data["invalid"] = "Your password and confirmation password do not match."
			temp, _ := template.ParseFiles("views/reset_password.html")
			temp.Execute(ctx.Response(), data)
			return res, err
		}

		if ts.ID != "" {
			data["message"] = "password changed successfully"
			decoder := json.NewDecoder(ctx.Request().Body)
			o.Password = password.Password
			decoder.Decode(o)
			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(o.Password), bcrypt.DefaultCost)
			o.Password = string(hashPassword)
			helpers.GetDB(ctx).Model(User{}).Where("id = ?", ts.ID).Updates(User{Password: o.Password, ConfirmPassword: o.Password})
			temp, _ := template.ParseFiles("views/reset_password_success.html")
			temp.Execute(ctx.Response(), data)
		}
	}
	return res, err
}

func (o *User) ForgotPassword(ctx helpers.Context) (res map[string]interface{}) {
	if o.Email == "" {
		res = helpers.Map{
			"error": helpers.Map{
				"code": 400,
				"detail": helpers.Map{
					"email": helpers.Map{
						"required": "email is a required field",
					},
				},
				"message": "email is a required field",
			},
		}
		return res
	}

	pr := User{}
	err := helpers.GetDB(ctx).Model(User{}).Where("email = ?", o.Email).First(&pr).Error
	if err != nil {
		res = helpers.Map{
			"error": helpers.Map{
				"code":    400,
				"message": "email is not found",
			},
		}
		return res
	}

	var appVer string
	env := config.Get("APP_ENV").String()
	if env == "local" {
		appVer = "http://localhost:4000/api/v1/reset_password/" + pr.ID
	}

	body := `
		<html>
		<h1>Reset Password</h1>
			<table">
			<tr>
				<td>Reset password request</td>
			</tr>
			</table>
			<a href="` + appVer + `">Reset Password</a>
		</html>
		`

	err = o.SendEmail(body, "Your Password Reset Link")
	if err != nil {
		return helpers.Map{
			"code":    500,
			"message": err.Error(),
		}
	}
	res = helpers.Map{
		"code":    200,
		"message": "Your password link has been sent to your email",
	}
	return res
}

func (o *User) Register(ctx helpers.Context) map[string]interface{} {
	isValid, msg := helpers.Validate(ctx, o)
	if !isValid {
		return msg
	}
	params, err := o.SetDefaultValue(ctx)
	if err != nil {
		return params
	}

	var appVer string
	env := config.Get("APP_ENV").String()
	if env == "local" {
		appVer = "http://localhost:4000/api/v1/emailver/" + o.ID
	}

	emailExist, _ := helpers.IsExists(ctx, o.TableName(), "email", o.Email)
	if emailExist != nil {
		return emailExist
	}

	userNameExist, _ := helpers.IsExists(ctx, o.TableName(), "username", o.Email)
	if userNameExist != nil {
		return userNameExist
	}

	phoneNumberExist, _ := helpers.IsExists(ctx, o.TableName(), "phone_number", o.Email)
	if phoneNumberExist != nil {
		return phoneNumberExist
	}

	body := `
		<html>
		<h1>Activation Email</h1>
			<table">
			<tr>
				<td>You are just one step one</td>
			</tr>
			</table>
			<a href="` + appVer + `">Confirm Email</a>
		</html>
		`

	err = o.SendEmail(body, "Activation Email")
	if err != nil {
		return helpers.Map{
			"code":    500,
			"message": err.Error(),
		}
	}

	helpers.GetDB(ctx).Create(o)
	res := helpers.Map{
		"code":    201,
		"message": "Register Success, Please check your email to activate your account",
		"user": helpers.Map{
			"id":         o.ID,
			"username":   o.Username,
			"email":      o.Email,
			"is_active":  o.IsActive,
			"expired_at": o.ExpiredAt,
		},
	}
	return res
}

func (o *User) SetDefaultValue(ctx helpers.Context) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	o.ID = helpers.NewUUID()
	o.ExpiredAt = time.Now().Add(time.Hour * 1)
	decoder := json.NewDecoder(ctx.Request().Body)
	decoder.Decode(o)
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(o.Password), bcrypt.DefaultCost)
	o.Password = string(hashPassword)
	o.ConfirmPassword = o.Password
	return params, nil
}

func (o *User) SendEmail(body string, params string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", "test@gmail.com")
	m.SetHeader("To", o.Email)
	m.SetHeader("Subject", params)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(string(config.Get("SMTP")), 587, string(config.Get("SMTP_USER")), string(config.Get("SMTP_PASSWORD")))

	if err = d.DialAndSend(m); err != nil {
		return err
	}
	return err
}
