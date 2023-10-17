package types

// Constants for all supported gender choices
const (
	Male    = "male"
	Female  = "female"
	Other   = "other"
	Unknown = "unknown"
)

// IsSupportedGender returns true if the provided gender is supported
func IsSupportedGender(gender string) bool {
	switch gender {
	case Male, Female, Other, Unknown:
		return true
	}
	return false
}
