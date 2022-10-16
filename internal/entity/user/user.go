package user

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Role struct {
	ID         int    `json:"id"`
	Prefix     string `json:"prefix"`
	Permission string `json:"permission"`
}

type UserAccess struct {
	ID    int    `json:"id"`
	User  User   `json:"user"`
	Roles []Role `json:"roles"`
}

type UserDetails struct {
	User  User   `json:"user"`
	Roles []Role `json:"roles"`
}

const (
	RoleAdmin     string = "admin"
	RoleLead      string = "lead"
	RoleUser      string = "user"
	RoleSuperUser string = "superuser"
)
