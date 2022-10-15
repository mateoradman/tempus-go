// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateAbsence(ctx context.Context, arg CreateAbsenceParams) (Absence, error)
	CreateCompany(ctx context.Context, name string) (Company, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreatePermission(ctx context.Context, name string) (Permission, error)
	CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateTeam(ctx context.Context, arg CreateTeamParams) (Team, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAbsence(ctx context.Context, id int64) (Absence, error)
	DeleteCompany(ctx context.Context, id int64) (Company, error)
	DeleteEntry(ctx context.Context, id int64) (Entry, error)
	DeletePermission(ctx context.Context, id int64) (Permission, error)
	DeleteRole(ctx context.Context, id int64) (Role, error)
	DeleteTeam(ctx context.Context, id int64) (Team, error)
	DeleteUser(ctx context.Context, id int64) (User, error)
	GetAbsence(ctx context.Context, id int64) (Absence, error)
	GetCompany(ctx context.Context, id int64) (Company, error)
	GetEntry(ctx context.Context, id int64) (Entry, error)
	GetPermission(ctx context.Context, id int64) (Permission, error)
	GetRole(ctx context.Context, id int64) (Role, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetTeam(ctx context.Context, id int64) (Team, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	ListAbsences(ctx context.Context, arg ListAbsencesParams) ([]Absence, error)
	ListCompanies(ctx context.Context, arg ListCompaniesParams) ([]Company, error)
	ListCompanyEmployees(ctx context.Context, arg ListCompanyEmployeesParams) ([]User, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
	ListPermissions(ctx context.Context, arg ListPermissionsParams) ([]Permission, error)
	ListRoles(ctx context.Context, arg ListRolesParams) ([]Role, error)
	ListTeamMembers(ctx context.Context, arg ListTeamMembersParams) ([]User, error)
	ListTeams(ctx context.Context, arg ListTeamsParams) ([]Team, error)
	ListUserAbsences(ctx context.Context, arg ListUserAbsencesParams) ([]Absence, error)
	ListUserEntries(ctx context.Context, arg ListUserEntriesParams) ([]Entry, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateAbsence(ctx context.Context, arg UpdateAbsenceParams) (Absence, error)
	UpdateCompany(ctx context.Context, arg UpdateCompanyParams) (Company, error)
	UpdateEntry(ctx context.Context, arg UpdateEntryParams) (Entry, error)
	UpdatePermission(ctx context.Context, arg UpdatePermissionParams) (Permission, error)
	UpdateRole(ctx context.Context, arg UpdateRoleParams) (Role, error)
	UpdateTeam(ctx context.Context, arg UpdateTeamParams) (Team, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
