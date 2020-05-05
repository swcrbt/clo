package db

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var MysqlDB *gorm.DB

// NullTime is an alias for mysql.NullTime data type
type NullTime struct {
	mysql.NullTime
}

// MarshalJSON for NullTime
func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format("2006-01-02 15:04:05"))
	return []byte(val), nil
}

type Info struct {
	Id         int      `gorm:"column:id;primary_key" json:"id"`
	Name       string   `gorm:"column:name" json:"name"`
	Company    string   `gorm:"column:company" json:"company"`
	Position   string   `gorm:"column:position" json:"position"`
	Mobile     string   `gorm:"column:mobile" json:"mobile"`
	IsSign     bool     `gorm:"column:issign" json:"issign"`
	Status     int      `gorm:"column:status" json:"status"`
	SignTime   NullTime `gorm:"column:signtime" "default: null" json:"signtime"`
	AgreeTime  NullTime `gorm:"column:agreetime" "default: null" json:"agreetime"`
	CreateTime NullTime `gorm:"column:createtime" "default: null" json:"createtime"`
}
