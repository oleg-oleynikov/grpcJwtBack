package service

type Permissions []string

type Role struct {
	Id          string
	Title       string
	Permissions Permissions
}

func (r *Role) HasPermission(permission string) bool {
	for _, role := range r.Permissions {
		if role == permission {
			return true
		}
	}

	return false
}
