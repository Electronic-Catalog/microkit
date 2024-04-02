package locale

type Language string

const (
	English Language = "en"
	Arabic  Language = "ar"
	Turkish Language = "tr"
	Persian Language = "fa"
)

type Locale interface {
	Message(lang Language, key string, params ...interface{}) string
}
