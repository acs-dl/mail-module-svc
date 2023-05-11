package data

type Permissions interface {
	New() Permissions

	Upsert(module Permission) error
	Select() ([]ModulePermission, error)
	Get() (*ModulePermission, error)
	Delete(permission Permission) error

	WithModules() Permissions

	FilterByModuleName(name string) Permissions
	FilterByPermissionId(permissionId int64) Permissions
	FilterByStatus(status UserStatus) Permissions

	ResetFilters() Permissions
}

type Permission struct {
	Id       int64      `db:"id" structs:"-"`
	ModuleId int64      `db:"module_id" structs:"module_id"`
	Name     string     `db:"name" structs:"name"`
	Status   UserStatus `db:"status" structs:"status"`
	*Module  `db:"-" structs:",omitempty"`
}
