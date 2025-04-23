package authz

import "context"

type RBACCompareType int

const (
	Every RBACCompareType = iota
	Some
	None
)

func (a *Authz) RBAC(ctx context.Context, compareType RBACCompareType, userPermissions, permission int) bool {
	switch compareType {
	case Every:
		return a.PermissionService.Every(userPermissions, permission)
	case Some:
		return a.PermissionService.Some(userPermissions, permission)
	case None:
		return a.PermissionService.None(userPermissions, permission)
	}
	return false
}
