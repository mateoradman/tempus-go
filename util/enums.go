package	util

type Gender int64

const (
	Undefined Gender = iota
	Male
	Female
	Other
	Unknown
)

func (g Gender) String() string {
	switch g {
	case Male:
		return "male"
	case Female:
		return "female"
	case Other:
		return "other"
	case Unknown:
		return "unknown"
	}
	return ""
}

