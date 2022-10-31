package rbac

import (
	"github.com/gin-gonic/gin"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
	"github.com/mateoradman/tempus/internal/token"
	"github.com/mateoradman/tempus/internal/util"
)

// Creates new RBAC service
func NewRBACService(store db.Store) *RBACService {
	return &RBACService{store: store}
}

// RBAC application service
type RBACService struct {
	store db.Store
}

// Get the User object from the bearer token payload.
// AuthPayloadKey is set in the authMiddleware.
func (s *RBACService) getUser(ctx *gin.Context) (db.User, error) {
	payload := ctx.MustGet(util.AuthPayloadKey).(*token.Payload)
	return s.store.GetUserByUsername(ctx, payload.Username)
}

// User is an admin if the role is either AdminRole or SuperUserRole
func (s *RBACService) isAdmin(user db.User) bool {
	return user.Role <= int32(util.AdminRole)
}

// User is a company admin if the role is <= CompanyAdminRole and belongs to the company
func (s *RBACService) isCompanyAdmin(user db.User, companyID *int64) bool {
	if s.isAdmin(user) {
		return true
	}
	return (user.Role <= int32(util.CompanyAdminRole)) && (user.CompanyID == companyID)
}

// EnforceRole checks that the user role is equal or stronger than the required AccessRole.
func (s *RBACService) EnforceRole(ctx *gin.Context, r util.AccessRole) bool {
	user, err := s.getUser(ctx)
	if err != nil {
		return false
	}
	return !(user.Role > int32(r))
}

// EnforceUser checks whether the request to change user data is done by the admin,
// user's company admin, user's manager or the user themselves.
func (s *RBACService) EnforceUser(ctx *gin.Context, userID int64) bool {
	user, err := s.getUser(ctx)
	if err != nil {
		return false
	}
	if user.ID == userID {
		// if userID matches the user's ID from the context, it is the same user
		return true
	}

	userSearchedFor, err := s.store.GetUser(ctx, userID)
	if err != nil {
		return false
	}

	// Check if the user making the request is a company admin of the user we're looking for.
	return s.isCompanyAdmin(user, userSearchedFor.CompanyID)

	// TODO check if the user making the request is a team manager of the user we're looking for.
}

// EnforceCompany checks whether the request to apply change to company data
// is done by the user belonging to the that company and that the user has role CompanyAdmin.
// If user has admin role, the check for company doesnt need to pass.
func (s *RBACService) EnforceCompany(ctx *gin.Context, companyID *int64) bool {
	user, err := s.getUser(ctx)
	if err != nil {
		return false
	}
	return s.isCompanyAdmin(user, companyID)
}
