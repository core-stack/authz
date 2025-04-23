package zrole

import "errors"

var (
	ErrRoleNotFound = errors.New("role not found")
)

type Role struct {
	ID          int
	Name        string
	Permissions int
}

var Roles []Role

func AddRole(role Role) {
	Roles = append(Roles, role)
}
func AddRoles(roles ...Role) {
	Roles = append(Roles, roles...)
}

func GetRole(id int) (Role, error) {
	for _, role := range Roles {
		if role.ID == id {
			return role, nil
		}
	}
	return Role{}, ErrRoleNotFound
}
