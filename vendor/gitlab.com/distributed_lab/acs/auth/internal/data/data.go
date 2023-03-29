package data

const (
	ModuleName     = "auth"
	TokenRegExpStr = `/^([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_\-\+\/=]*)/gm`
)

type GenerateTokens struct {
	User              User
	AccessLife        int64
	RefreshLife       int64
	Secret            string
	PermissionsString string
}

type JwtClaims struct {
	ExpiresAt        int64  `json:"exp"`
	CreatedAtNano    int64  `json:"iat_nano"`
	OwnerId          int64  `json:"owner_id"`
	Email            string `json:"email"`
	ModulePermission string `json:"module.permission"`
}

type ModulePayload struct {
	RequestId string `json:"request_id"`
	Action    string `json:"action"`

	//other fields that are required for module
	ModulePermissions ModulePermissions `json:"module_permissions,omitempty"`
	ModuleName        string            `json:"module_name"`
}

type StatusPermission map[string]string
type ModulePermissions map[string]StatusPermission
