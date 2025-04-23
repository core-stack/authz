package authz

const (
	PermissionRead = 1 << iota
	PermissionWrite
	PermissionDelete
	PermissionAdmin
)

type PermissionService struct{}

func (s *PermissionService) Every(userPermissions, permissions int) bool {
	return userPermissions&permissions == permissions
}

func (s *PermissionService) Some(userPermissions, permissions int) bool {
	return userPermissions&permissions != 0
}

func (s *PermissionService) None(userPermissions, permissions int) bool {
	return userPermissions&permissions == 0
}

func (s *PermissionService) AddPermission(userPermissions, permission int) int {
	return userPermissions | permission
}

func (s *PermissionService) RemovePermission(userPermissions, permission int) int {
	return userPermissions &^ permission
}
func NewPermissionService() *PermissionService {
	return &PermissionService{}
}
