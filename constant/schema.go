package constant

var Operator = map[string]string{
	"$eq":      "=",
	"$ne":      "<>",
	"$gt":      ">",
	"$gte":     ">=",
	"$lte":     "<=",
	"$lt":      "<",
	"$like":    "like",
	"$ilike":   "like",
	"$nlike":   "not like",
	"$nilike":  "not like",
	"$in":      "in",
	"$nin":     "not in",
	"$regexp":  "regexp",
	"$nregexp": "not regexp",
}

var Keywords []string = []string{
	"type",
	"key",
	"value",
}

var Wrappings = map[string]string{
	"mysql":       "`",
	"firebirdsql": "\"",
	"postgres":    "\"",
}

func IsDatabaseKeywords(field string) bool {
	f := false
	for _, v := range Keywords {
		if field == v {
			f = true
		}
	}
	return f
}
