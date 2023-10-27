package types

const (
	AdminRole    = "admin"
	ManagerRole  = "manager"
	EmployeeRole = "employee"
)

// IsValidRole returns true if the provided role is supported
func IsValidRole(role string) bool {
	switch role {
	case AdminRole, ManagerRole, EmployeeRole:
		return true
	}
	return false
}
