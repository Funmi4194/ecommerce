package enum

type Role string

func (r Role) String() string {
	return string(r)
}

const (
	Admin Role = "ADMIN"

	User Role = "USER"
)
