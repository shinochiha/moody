package constant

import "regexp"

const (
	AlphaRegexString                 = "^[a-zA-Z]+$"
	AlphaNumericRegexString          = "^[a-zA-Z0-9]+$"
	AlphaUnicodeRegexString          = "^[\\p{L}]+$"
	AlphaUnicodeNumericRegexString   = "^[\\p{L}\\p{N}]+$"
	NumericRegexString               = "^[-+]?[0-9]+(?:\\.[0-9]+)?$"
	NumberRegexString                = "^[0-9]+$"
	HexadecimalRegexString           = "^[0-9a-fA-F]+$"
	HexcolorRegexString              = "^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{6})$"
	RgbRegexString                   = "^rgb\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*\\)$"
	RgbaRegexString                  = "^rgba\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*,\\s*(?:(?:0.[1-9]*)|[01])\\s*\\)$"
	HslRegexString                   = "^hsl\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*\\)$"
	HslaRegexString                  = "^hsla\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0.[1-9]*)|[01])\\s*\\)$"
	EmailRegexString                 = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	E164RegexString                  = "^\\+[1-9]?[0-9]{7,14}$"
	Base64RegexString                = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	Base64URLRegexString             = "^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3}=|[A-Za-z0-9-_]{4})$"
	ISBN10RegexString                = "^(?:[0-9]{9}X|[0-9]{10})$"
	ISBN13RegexString                = "^(?:(?:97(?:8|9))[0-9]{10})$"
	UUID3RegexString                 = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	UUID4RegexString                 = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	UUID5RegexString                 = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	UUIDRegexString                  = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	UUID3RFC4122RegexString          = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-3[0-9a-fA-F]{3}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
	UUID4RFC4122RegexString          = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$"
	UUID5RFC4122RegexString          = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-5[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$"
	UUIDRFC4122RegexString           = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
	ASCIIRegexString                 = "^[\x00-\x7F]*$"
	PrintableASCIIRegexString        = "^[\x20-\x7E]*$"
	MultibyteRegexString             = "[^\x00-\x7F]"
	DataURIRegexString               = `^data:((?:\w+\/(?:([^;]|;[^;]).)+)?)`
	LatitudeRegexString              = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"
	LongitudeRegexString             = "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"
	SSNRegexString                   = `^[0-9]{3}[ -]?(0[1-9]|[1-9][0-9])[ -]?([1-9][0-9]{3}|[0-9][1-9][0-9]{2}|[0-9]{2}[1-9][0-9]|[0-9]{3}[1-9])$`
	HostnameRegexStringRFC952        = `^[a-zA-Z][a-zA-Z0-9\-\.]+[a-zA-Z0-9]$`                                            // https://tools.ietf.org/html/rfc952
	HostnameRegexStringRFC1123       = `^([a-zA-Z0-9]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*?$` // accepts hostname starting with a digit https://tools.ietf.org/html/rfc1123
	BtcAddressRegexString            = `^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$`                                                // bitcoin address
	BtcAddressUpperRegexStringBech32 = `^BC1[02-9AC-HJ-NP-Z]{7,76}$`                                                      // bitcoin bech32 address https://en.bitcoin.it/wiki/Bech32
	BtcAddressLowerRegexStringBech32 = `^bc1[02-9ac-hj-np-z]{7,76}$`                                                      // bitcoin bech32 address https://en.bitcoin.it/wiki/Bech32
	EthAddressRegexString            = `^0x[0-9a-fA-F]{40}$`
	EthAddressUpperRegexString       = `^0x[0-9A-F]{40}$`
	EthAddressLowerRegexString       = `^0x[0-9a-f]{40}$`
	URLEncodedRegexString            = `(%[A-Fa-f0-9]{2})`
	HTMLEncodedRegexString           = `&#[x]?([0-9a-fA-F]{2})|(&gt)|(&lt)|(&quot)|(&amp)+[;]?`
	HTMLRegexString                  = `<[/]?([a-zA-Z]+).*?>`
	SplitParamsRegexString           = `'[^']*'|\S+`
)

var (
	AlphaRegex                 = regexp.MustCompile(AlphaRegexString)
	AlphaNumericRegex          = regexp.MustCompile(AlphaNumericRegexString)
	AlphaUnicodeRegex          = regexp.MustCompile(AlphaUnicodeRegexString)
	AlphaUnicodeNumericRegex   = regexp.MustCompile(AlphaUnicodeNumericRegexString)
	NumericRegex               = regexp.MustCompile(NumericRegexString)
	NumberRegex                = regexp.MustCompile(NumberRegexString)
	HexadecimalRegex           = regexp.MustCompile(HexadecimalRegexString)
	HexcolorRegex              = regexp.MustCompile(HexcolorRegexString)
	RgbRegex                   = regexp.MustCompile(RgbRegexString)
	RgbaRegex                  = regexp.MustCompile(RgbaRegexString)
	HslRegex                   = regexp.MustCompile(HslRegexString)
	HslaRegex                  = regexp.MustCompile(HslaRegexString)
	E164Regex                  = regexp.MustCompile(E164RegexString)
	EmailRegex                 = regexp.MustCompile(EmailRegexString)
	Base64Regex                = regexp.MustCompile(Base64RegexString)
	Base64URLRegex             = regexp.MustCompile(Base64URLRegexString)
	ISBN10Regex                = regexp.MustCompile(ISBN10RegexString)
	ISBN13Regex                = regexp.MustCompile(ISBN13RegexString)
	UUID3Regex                 = regexp.MustCompile(UUID3RegexString)
	UUID4Regex                 = regexp.MustCompile(UUID4RegexString)
	UUID5Regex                 = regexp.MustCompile(UUID5RegexString)
	UUIDRegex                  = regexp.MustCompile(UUIDRegexString)
	UUID3RFC4122Regex          = regexp.MustCompile(UUID3RFC4122RegexString)
	UUID4RFC4122Regex          = regexp.MustCompile(UUID4RFC4122RegexString)
	UUID5RFC4122Regex          = regexp.MustCompile(UUID5RFC4122RegexString)
	UUIDRFC4122Regex           = regexp.MustCompile(UUIDRFC4122RegexString)
	ASCIIRegex                 = regexp.MustCompile(ASCIIRegexString)
	PrintableASCIIRegex        = regexp.MustCompile(PrintableASCIIRegexString)
	MultibyteRegex             = regexp.MustCompile(MultibyteRegexString)
	DataURIRegex               = regexp.MustCompile(DataURIRegexString)
	LatitudeRegex              = regexp.MustCompile(LatitudeRegexString)
	LongitudeRegex             = regexp.MustCompile(LongitudeRegexString)
	SSNRegex                   = regexp.MustCompile(SSNRegexString)
	HostnameRegexRFC952        = regexp.MustCompile(HostnameRegexStringRFC952)
	HostnameRegexRFC1123       = regexp.MustCompile(HostnameRegexStringRFC1123)
	BtcAddressRegex            = regexp.MustCompile(BtcAddressRegexString)
	BtcUpperAddressRegexBech32 = regexp.MustCompile(BtcAddressUpperRegexStringBech32)
	BtcLowerAddressRegexBech32 = regexp.MustCompile(BtcAddressLowerRegexStringBech32)
	EthAddressRegex            = regexp.MustCompile(EthAddressRegexString)
	EthaddressRegexUpper       = regexp.MustCompile(EthAddressUpperRegexString)
	EthAddressRegexLower       = regexp.MustCompile(EthAddressLowerRegexString)
	URLEncodedRegex            = regexp.MustCompile(URLEncodedRegexString)
	HTMLEncodedRegex           = regexp.MustCompile(HTMLEncodedRegexString)
	HTMLRegex                  = regexp.MustCompile(HTMLRegexString)
	SplitParamsRegex           = regexp.MustCompile(SplitParamsRegexString)
)
