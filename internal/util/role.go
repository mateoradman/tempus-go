package util

// AccessRole represents user role type
type AccessRole int32

const (
	// SuperUserRole has all permissions
	SuperUserRole AccessRole = iota + 1

	// AdminRole has admin permissions
	AdminRole

	// CompanyAdmin has company specific admin permissions
	CompanyAdminRole

	// TeamManagerRole has team specific admin permissions
	TeamManagerRole

	// DefaultRole is a default standard permission
	DefaultRole
)
