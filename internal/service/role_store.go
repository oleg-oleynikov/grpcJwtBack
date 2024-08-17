package service

type RoleStore interface {
	Save(*Role) error
	FindByTitle(string) (*Role, error)
}
