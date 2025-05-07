package auth

type User struct {
	Id       int64 `gorm:"primary_key;uniqueIndex"`
	Username string
	Password string
}
