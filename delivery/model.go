package delivery

type Delivery struct {
	Id        int64 `gorm:"primary_key;uniqueIndex"`
	CourierId int64
	OrderId   int64
	Status    string
	AddedBy   int64
}
