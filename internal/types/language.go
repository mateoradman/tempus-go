package types

// Constants for all supported gender choices
const (
	English  = "en"
	German   = "de"
	Croatian = "hr"
)

// IsSupportedLanguage function checks whether a language choice is supported
func IsSupportedLanguage(gender string) bool {
	switch gender {
	case English, German, Croatian:
		return true
	}
	return false
}
