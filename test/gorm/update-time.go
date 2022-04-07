package main

import (
	"fmt"
	"time"

	// "github.com/tidwall/pretty"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	db *gorm.DB
)

type Config struct {
	User     string
	Password string
	Endpoint string
	Port     int
	Database string
}

func initDB() {
	cfg := Config{
		User:     "root",
		Password: "root",
		Endpoint: "localhost",
		Port:     3306,
		Database: "gorm_test",
	}
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Endpoint, cfg.Port, cfg.Database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&RelatedTime{}, &DataTime{})
}
func main() {

	initDB()
	m := RelatedTime{
		ContentId: "125",
		MatchName: "match_name6",
		DataTimes: []DataTime{
			{Name: "hello world6", ContentId: "125"},
		},
	}
	res := db.Clauses(clause.OnConflict{
		UpdateAll: true,
		Columns:   []clause.Column{{Name: "ContentId"}},
	}).Create(&m)
	// res := db.Session(&gorm.Session{FullSaveAssociations: true}).Where("content_id = ?", m.ContentId).Updates(&m)
	fmt.Println("res:", res)
	// if res.RowsAffected == 0 {
	// 	db.Clauses(clause.OnConflict{
	// 		UpdateAll: true,
	// 	}).Create(&m.DataTimes)
	// }
	fmt.Printf("%+v\n", m)
	m = RelatedTime{}
	db.Where("content_id = ?", "125").Find(&m)
	fmt.Println(m)
	readData()
}
func readData() {
	m := RelatedTime{
		ContentId: "123",
		DataTimes: []DataTime{},
	}

	db.Model(&m).Association("DataTimes").Find(&m.DataTimes)
	fmt.Printf("%v\n", m)
}

type RelatedTime struct {
	ID        int        `json:"id"`
	MatchName string     `gorm:"column:match_name;type:varchar(128)" json:"match_name,omitempty"`
	ContentId string     `gorm:"column:content_id;type:varchar(128);uniqueIndex:by_content_id" json:"content_id,omitempty"`
	DataTimes []DataTime `gorm:"foreignKey:ContentId;references:ContentId;constraint:OnUpdate:CASCADE" json:"data_times,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time  `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

type DataTime struct {
	ID        int       `json:"id"`
	Name      string    `gorm:"column:name;type:varchar(128)" json:"name,omitempty"`
	ContentId string    `gorm:"column:content_id;type:varchar(128);uniqueIndex:by_content_id" json:"content_id,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}
