package courier

type Courier struct {
	Id          int64 `gorm:"primary_key;uniqueIndex"`
	Name        string
	PhoneNumber string
	Available   bool
	AddedBy     int64
}
