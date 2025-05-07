package order

import (
	"time"

	"gorm.io/datatypes"
)

type Order struct {
	Id        int64 `gorm:"primary_key;uniqueIndex"`
	SenderId  int64
	Items     datatypes.JSON
	Address   string
	Status    string
	CreatedAt time.Time
}
