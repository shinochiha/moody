package helpers

import (
	"encoding/json"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	en "github.com/go-playground/locales/en"
	id "github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	id_translations "github.com/go-playground/validator/v10/translations/id"
	uuid "github.com/google/uuid"
	urn "github.com/leodido/go-urn"

	"github.com/moody/constant"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

type Validations struct {
	validate *validator.Validate
	trans    ut.Translator
}

type Validation struct {
	Lang string
	I    interface{}
}

func Validate(ctx Context, i interface{}) (bool, map[string]interface{}) {
	v := Validation{GetCtx(ctx).UserLang(), i}

	validate = validator.New()
	validate.RegisterTagNameFunc(v.GetJsonTag)
	validate.RegisterValidation("custom", customFunc)
	trans := v.SetTranslator()
	return v.Response(trans)
}

func NewValidation() *Validations {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	// register tag label
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := field.Tag.Get("label")
		return name
	})

	// membuat custom error
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} harus diisi", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	return &Validations{
		validate: validate,
		trans:    trans,
	}
}

func customFunc(fl validator.FieldLevel) bool {
	if fl.Field().String() == "invalid" {
		return false
	}
	return true
}

func (v *Validation) GetJsonTag(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

func (v *Validation) SetTranslator() ut.Translator {
	if v.Lang == "id" {
		return v.BahasaIndonesiaTranslator()
	} else {
		return v.EnglishTranslator()
	}
}

func (v *Validation) BahasaIndonesiaTranslator() ut.Translator {
	id := id.New()
	uni = ut.New(id, id)
	trans, _ := uni.GetTranslator("id")
	id_translations.RegisterDefaultTranslations(validate, trans)
	return trans
}

func (v *Validation) EnglishTranslator() ut.Translator {
	en := en.New()
	uni = ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)
	return trans
}

func (v *Validation) Response(trans ut.Translator) (bool, map[string]interface{}) {
	isValid := true
	res := map[string]interface{}{}
	prefix := GetStructName(v.I) + "."
	err := validate.Struct(v.I)
	if err != nil {
		message := ""
		detail := map[string]interface{}{}
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			msg := e.Translate(trans)
			if message == "" {
				message = msg
			}
			detail[strings.Replace(e.Namespace(), prefix, "", 1)] = map[string]interface{}{e.Tag(): msg}
		}
		isValid = false
		res["error"] = map[string]interface{}{
			"code":    400,
			"message": message,
			"detail":  detail,
		}
	}
	return isValid, res
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func (v *Validations) Struct(s interface{}) interface{} {
	errors := make(map[string]string)

	err := validate.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.StructField()] = e.Translate(v.trans)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (v *Validation) IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func (v *Validation) IsURLEncoded(val string) bool {
	return constant.URLEncodedRegex.MatchString(val)
}

func (v *Validation) IsHTMLEncoded(val string) bool {
	return constant.HTMLEncodedRegex.MatchString(val)
}

func (v *Validation) IsHTML(val string) bool {
	return constant.HTMLRegex.MatchString(val)
}

// IsMAC is the validation function for validating if the field's value is a valid MAC address.
func (v *Validation) IsMAC(val string) bool {
	_, err := net.ParseMAC(val)
	return err == nil
}

// IsCIDRv4 is the validation function for validating if the field's value is a valid v4 CIDR address.
func (v *Validation) IsCIDRv4(val string) bool {
	ip, _, err := net.ParseCIDR(val)
	return err == nil && ip.To4() != nil
}

// IsCIDRv6 is the validation function for validating if the field's value is a valid v6 CIDR address.
func (v *Validation) IsCIDRv6(val string) bool {
	ip, _, err := net.ParseCIDR(val)
	return err == nil && ip.To4() == nil
}

// IsCIDR is the validation function for validating if the field's value is a valid v4 or v6 CIDR address.
func (v *Validation) IsCIDR(val string) bool {
	_, _, err := net.ParseCIDR(val)
	return err == nil
}

// IsIPv4 is the validation function for validating if a value is a valid v4 IP address.
func (v *Validation) IsIPv4(val string) bool {
	ip := net.ParseIP(val)
	return ip != nil && ip.To4() != nil
}

// IsIPv6 is the validation function for validating if the field's value is a valid v6 IP address.
func (v *Validation) IsIPv6(val string) bool {
	ip := net.ParseIP(val)
	return ip != nil && ip.To4() == nil
}

// IsIP is the validation function for validating if the field's value is a valid v4 or v6 IP address.
func (v *Validation) IsIP(val string) bool {
	ip := net.ParseIP(val)
	return ip != nil
}

// IsSSN is the validation function for validating if the field's value is a valid SSN.
func (v *Validation) IsSSN(val string) bool {
	if len(val) != 11 {
		return false
	}
	return constant.SSNRegex.MatchString(val)
}

// IsLongitude is the validation function for validating if the field's value is a valid longitude coordinate.
func (v *Validation) IsLongitude(val string) bool {
	return constant.LongitudeRegex.MatchString(val)
}

// IsLatitude is the validation function for validating if the field's value is a valid latitude coordinate.
func (v *Validation) IsLatitude(val string) bool {
	return constant.LatitudeRegex.MatchString(val)
}

// IsDataURI is the validation function for validating if the field's value is a valid data URI.
func (v *Validation) IsDataURI(val string) bool {
	uri := strings.SplitN(val, ",", 2)
	if len(uri) != 2 {
		return false
	}
	if !constant.DataURIRegex.MatchString(uri[0]) {
		return false
	}
	return constant.Base64Regex.MatchString(uri[1])
}

// HasMultiByteCharacter is the validation function for validating if the field's value has a multi byte character.
func (v *Validation) HasMultiByteCharacter(val string) bool {
	return constant.MultibyteRegex.MatchString(val)
}

// IsPrintableASCII is the validation function for validating if the field's value is a valid printable ASCII character.
func (v *Validation) IsPrintableASCII(val string) bool {
	return constant.PrintableASCIIRegex.MatchString(val)
}

// IsASCII is the validation function for validating if the field's value is a valid ASCII character.
func (v *Validation) IsASCII(val string) bool {
	return constant.ASCIIRegex.MatchString(val)
}

// IsUUID5 is the validation function for validating if the field's value is a valid v5 UUID.
func (v *Validation) IsUUID5(val string) bool {
	return constant.UUID5Regex.MatchString(val)
}

// IsUUID4 is the validation function for validating if the field's value is a valid v4 UUID.
func (v *Validation) IsUUID4(val string) bool {
	return constant.UUID4Regex.MatchString(val)
}

// IsUUID3 is the validation function for validating if the field's value is a valid v3 UUID.
func (v *Validation) IsUUID3(val string) bool {
	return constant.UUID3Regex.MatchString(val)
}

// IsUUID is the validation function for validating if the field's value is a valid UUID of any version.
func (v *Validation) IsUUID(val string) bool {
	return constant.UUIDRegex.MatchString(val)
}

// IsUUID5RFC4122 is the validation function for validating if the field's value is a valid RFC4122 v5 UUID.
func (v *Validation) IsUUID5RFC4122(val string) bool {
	return constant.UUID5RFC4122Regex.MatchString(val)
}

// IsUUID4RFC4122 is the validation function for validating if the field's value is a valid RFC4122 v4 UUID.
func (v *Validation) IsUUID4RFC4122(val string) bool {
	return constant.UUID4RFC4122Regex.MatchString(val)
}

// IsUUID3RFC4122 is the validation function for validating if the field's value is a valid RFC4122 v3 UUID.
func (v *Validation) IsUUID3RFC4122(val string) bool {
	return constant.UUID3RFC4122Regex.MatchString(val)
}

// IsUUIDRFC4122 is the validation function for validating if the field's value is a valid RFC4122 UUID of any version.
func (v *Validation) IsUUIDRFC4122(val string) bool {
	return constant.UUIDRFC4122Regex.MatchString(val)
}

// IsISBN is the validation function for validating if the field's value is a valid v10 or v13 ISBN.
func (v *Validation) IsISBN(val string) bool {
	return v.IsISBN10(val) || v.IsISBN13(val)
}

// IsISBN13 is the validation function for validating if the field's value is a valid v13 ISBN.
func (v *Validation) IsISBN13(val string) bool {

	s := strings.Replace(strings.Replace(val, "-", "", 4), " ", "", 4)

	if !constant.ISBN13Regex.MatchString(s) {
		return false
	}

	var checksum int32
	var i int32

	factor := []int32{1, 3}

	for i = 0; i < 12; i++ {
		checksum += factor[i%2] * int32(s[i]-'0')
	}

	return (int32(s[12]-'0'))-((10-(checksum%10))%10) == 0
}

// IsISBN10 is the validation function for validating if the field's value is a valid v10 ISBN.
func (v *Validation) IsISBN10(val string) bool {

	s := strings.Replace(strings.Replace(val, "-", "", 3), " ", "", 3)

	if !constant.ISBN10Regex.MatchString(s) {
		return false
	}

	var checksum int32
	var i int32

	for i = 0; i < 9; i++ {
		checksum += (i + 1) * int32(s[i]-'0')
	}

	if s[9] == 'X' {
		checksum += 10 * 10
	} else {
		checksum += 10 * int32(s[9]-'0')
	}

	return checksum%11 == 0
}

// ContainsRune is the validation function for validating that the field's value contains the rune specified within the param.
func (v *Validation) ContainsRune(val, param string) bool {
	r, _ := utf8.DecodeRuneInString(param)
	return strings.ContainsRune(val, r)
}

// ContainsAny is the validation function for validating that the field's value contains any of the characters specified within the param.
func (v *Validation) ContainsAny(val, param string) bool {
	return strings.ContainsAny(val, param)
}

// Contains is the validation function for validating that the field's value contains the text specified within the param.
func (v *Validation) Contains(val, param string) bool {
	return strings.Contains(val, param)
}

// StartsWith is the validation function for validating that the field's value starts with the text specified within the param.
func (v *Validation) StartsWith(val, param string) bool {
	return strings.HasPrefix(val, param)
}

// EndsWith is the validation function for validating that the field's value ends with the text specified within the param.
func (v *Validation) EndsWith(val, param string) bool {
	return strings.HasSuffix(val, param)
}

// IsBase64 is the validation function for validating if the current field's value is a valid base 64.
func (v *Validation) IsBase64(val string) bool {
	return constant.Base64Regex.MatchString(val)
}

// IsBase64URL is the validation function for validating if the current field's value is a valid base64 URL safe string.
func (v *Validation) IsBase64URL(val string) bool {
	return constant.Base64URLRegex.MatchString(val)
}

// IsURI is the validation function for validating if the current field's value is a valid URI.
func (v *Validation) IsURI(s string) bool {
	// checks needed as of Go 1.6 because of change https://github.com/golang/go/commit/617c93ce740c3c3cc28cdd1a0d712be183d0b328#diff-6c2d018290e298803c0c9419d8739885L195
	// emulate browser and strip the '#' suffix prior to validation. see issue-#237
	if i := strings.Index(s, "#"); i > -1 {
		s = s[:i]
	}
	if len(s) == 0 {
		return false
	}
	_, err := url.ParseRequestURI(s)
	return err == nil
}

// IsURL is the validation function for validating if the current field's value is a valid URL.
func (v *Validation) IsURL(s string) bool {
	var i int
	// checks needed as of Go 1.6 because of change https://github.com/golang/go/commit/617c93ce740c3c3cc28cdd1a0d712be183d0b328#diff-6c2d018290e298803c0c9419d8739885L195
	// emulate browser and strip the '#' suffix prior to validation. see issue-#237
	if i = strings.Index(s, "#"); i > -1 {
		s = s[:i]
	}
	if len(s) == 0 {
		return false
	}
	url, err := url.ParseRequestURI(s)
	if err != nil || url.Scheme == "" {
		return false
	}
	return true
}

// IsUrnRFC2141 is the validation function for validating if the current field's value is a valid URN as per RFC 2141.
func (v *Validation) IsUrnRFC2141(val string) bool {
	_, match := urn.Parse([]byte(val))
	return match
}

// IsFile is the validation function for validating if the current field's value is a valid file path.
func (v *Validation) IsFile(val string) bool {
	fileInfo, err := os.Stat(val)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

// IsE164 is the validation function for validating if the current field's value is a valid e.164 formatted phone number.
func (v *Validation) IsE164(val string) bool {
	return constant.E164Regex.MatchString(val)
}

// IsEmail is the validation function for validating if the current field's value is a valid email address.
func (v *Validation) IsEmail(val string) bool {
	return constant.EmailRegex.MatchString(val)
}

// IsHSLA is the validation function for validating if the current field's value is a valid HSLA color.
func (v *Validation) IsHSLA(val string) bool {
	return constant.HslaRegex.MatchString(val)
}

// IsHSL is the validation function for validating if the current field's value is a valid HSL color.
func (v *Validation) IsHSL(val string) bool {
	return constant.HslRegex.MatchString(val)
}

// IsRGBA is the validation function for validating if the current field's value is a valid RGBA color.
func (v *Validation) IsRGBA(val string) bool {
	return constant.RgbaRegex.MatchString(val)
}

// IsRGB is the validation function for validating if the current field's value is a valid RGB color.
func (v *Validation) IsRGB(val string) bool {
	return constant.RgbRegex.MatchString(val)
}

// IsHEXColor is the validation function for validating if the current field's value is a valid HEX color.
func (v *Validation) IsHEXColor(val string) bool {
	return constant.HexcolorRegex.MatchString(val)
}

// IsHexadecimal is the validation function for validating if the current field's value is a valid hexadecimal.
func (v *Validation) IsHexadecimal(val string) bool {
	return constant.HexadecimalRegex.MatchString(val)
}

// IsNumber is the validation function for validating if the current field's value is a valid number.
func (v *Validation) IsNumber(val string) bool {
	// switch reflect.ValueOf(val).Kind() {
	// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
	//   return true
	// default:
	// }
	return constant.NumberRegex.MatchString(val)
}

// IsNumeric is the validation function for validating if the current field's value is a valid numeric value.
func (v *Validation) IsNumeric(val string) bool {
	// switch reflect.ValueOf(val).Kind() {
	// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
	//   return true
	// default:
	// }
	return constant.NumericRegex.MatchString(val)
}

// IsAlphanum is the validation function for validating if the current field's value is a valid alphanumeric value.
func (v *Validation) IsAlphanum(val string) bool {
	return constant.AlphaNumericRegex.MatchString(val)
}

// IsAlpha is the validation function for validating if the current field's value is a valid alpha value.
func (v *Validation) IsAlpha(val string) bool {
	return constant.AlphaRegex.MatchString(val)
}

// IsAlphanumUnicode is the validation function for validating if the current field's value is a valid alphanumeric unicode value.
func (v *Validation) IsAlphanumUnicode(val string) bool {
	return constant.AlphaUnicodeNumericRegex.MatchString(val)
}

// IsAlphaUnicode is the validation function for validating if the current field's value is a valid alpha unicode value.
func (v *Validation) IsAlphaUnicode(val string) bool {
	return constant.AlphaUnicodeRegex.MatchString(val)
}

// IsTCP4AddrResolvable is the validation function for validating if the field's value is a resolvable tcp4 address.
func (v *Validation) IsTCP4AddrResolvable(val string) bool {
	_, err := net.ResolveTCPAddr("tcp4", val)
	return err == nil
}

// IsTCP6AddrResolvable is the validation function for validating if the field's value is a resolvable tcp6 address.
func (v *Validation) IsTCP6AddrResolvable(val string) bool {
	_, err := net.ResolveTCPAddr("tcp6", val)
	return err == nil
}

// IsTCPAddrResolvable is the validation function for validating if the field's value is a resolvable tcp address.
func (v *Validation) IsTCPAddrResolvable(val string) bool {
	_, err := net.ResolveTCPAddr("tcp", val)
	return err == nil
}

// IsUDP4AddrResolvable is the validation function for validating if the field's value is a resolvable udp4 address.
func (v *Validation) IsUDP4AddrResolvable(val string) bool {
	_, err := net.ResolveUDPAddr("udp4", val)
	return err == nil
}

// IsUDP6AddrResolvable is the validation function for validating if the field's value is a resolvable udp6 address.
func (v *Validation) IsUDP6AddrResolvable(val string) bool {
	_, err := net.ResolveUDPAddr("udp6", val)
	return err == nil
}

// IsUDPAddrResolvable is the validation function for validating if the field's value is a resolvable udp address.
func (v *Validation) IsUDPAddrResolvable(val string) bool {
	_, err := net.ResolveUDPAddr("udp", val)
	return err == nil
}

// IsIP4AddrResolvable is the validation function for validating if the field's value is a resolvable ip4 address.
func (v *Validation) IsIP4AddrResolvable(val string) bool {
	_, err := net.ResolveIPAddr("ip4", val)
	return err == nil
}

// IsIP6AddrResolvable is the validation function for validating if the field's value is a resolvable ip6 address.
func (v *Validation) IsIP6AddrResolvable(val string) bool {
	_, err := net.ResolveIPAddr("ip6", val)
	return err == nil
}

// IsIPAddrResolvable is the validation function for validating if the field's value is a resolvable ip address.
func (v *Validation) IsIPAddrResolvable(val string) bool {
	_, err := net.ResolveIPAddr("ip", val)
	return err == nil
}

// IsUnixAddrResolvable is the validation function for validating if the field's value is a resolvable unix address.
func (v *Validation) IsUnixAddrResolvable(val string) bool {
	_, err := net.ResolveUnixAddr("unix", val)
	return err == nil
}

func (v *Validation) IsIP4Addr(val string) bool {
	if idx := strings.LastIndex(val, ":"); idx != -1 {
		val = val[0:idx]
	}
	ip := net.ParseIP(val)
	return ip != nil && ip.To4() != nil
}

func (v *Validation) IsIP6Addr(val string) bool {
	if idx := strings.LastIndex(val, ":"); idx != -1 {
		if idx != 0 && val[idx-1:idx] == "]" {
			val = val[1 : idx-1]
		}
	}
	ip := net.ParseIP(val)
	return ip != nil && ip.To4() == nil
}

func (v *Validation) IsHostnameRFC952(val string) bool {
	return constant.HostnameRegexRFC952.MatchString(val)
}

func (v *Validation) IsHostnameRFC1123(val string) bool {
	return constant.HostnameRegexRFC1123.MatchString(val)
}

func (v *Validation) IsFQDN(val string) bool {
	if val == "" {
		return false
	}
	if val[len(val)-1] == '.' {
		val = val[0 : len(val)-1]
	}
	return strings.ContainsAny(val, ".") && constant.HostnameRegexRFC952.MatchString(val)
}

// IsDir is the validation function for validating if the current field's value is a valid directory.
func (v *Validation) IsDir(val string) bool {
	fileInfo, err := os.Stat(val)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// IsJSON is the validation function for validating if the current field's value is a valid json string.
func (v *Validation) IsJSON(val string) bool {
	return json.Valid([]byte(val))
}

// IsHostnamePort validates a <dns>:<port> combination for fields typically used for socket address.
func (v *Validation) IsHostnamePort(val string) bool {
	host, port, err := net.SplitHostPort(val)
	if err != nil {
		return false
	}
	// Port must be a iny <= 65535.
	if portNum, err := strconv.ParseInt(port, 10, 32); err != nil || portNum > 65535 || portNum < 1 {
		return false
	}
	// If host is specified, it should match a DNS name
	if host != "" {
		return constant.HostnameRegexRFC1123.MatchString(host)
	}
	return true
}

// IsLowercase is the validation function for validating if the current field's value is a lowercase string.
func (v *Validation) IsLowercase(val string) bool {
	if val == "" {
		return false
	}
	return val == strings.ToLower(val)
}

// IsUppercase is the validation function for validating if the current field's value is an uppercase string.
func (v *Validation) IsUppercase(val string) bool {
	if val == "" {
		return false
	}
	return val == strings.ToUpper(val)
}

// IsDateTime is the validation function for validating if the current field's value is a valid datetime string.
func (v *Validation) IsDateTime(val string) bool {
	_, err := time.Parse("2006-01-02T15:04:05Z", val)
	if err != nil {
		_, err = time.Parse("RFC3339", val)
		if err != nil {
			_, err = time.Parse("2006-01-02T15:04:05.999999+07:00", val)
			if err != nil {
				_, err = time.Parse("RFC3339Nano", val)
				if err != nil {
					return false
				}
			}
		}
	}
	return true
}

// IsDate is the validation function for validating if the current field's value is a valid date string.
func (v *Validation) IsDate(val string) bool {
	_, err := time.Parse("2006-01-02", val)
	if err != nil {
		return false
	}
	return true
}

// IsTime is the validation function for validating if the current field's value is a valid time string.
func (v *Validation) IsTime(val string) bool {
	_, err := time.Parse("15:04:05", val)
	if err != nil {
		return false
	}
	return true
}
