package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func init() {
	cfg := Config{
		User:     "root",
		Password: "root",
		Endpoint: "localhost",
		Port:     3306,
		Database: "gorm_test",
	}
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Endpoint, cfg.Port, cfg.Database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Matches{}, &Series{}, &Ground{}, &Country{}, &Team{}, &Town{})
}
func main() {
	m := Matches{
		MatchName: "match name",
		SeriesId:  100,
		Series: Series{
			SeriesId: 100,
			Name:     "series",
		},
		GroundId: 101,
		Ground: Ground{
			GroundId: 101,
			Name:     "ground",
			TownId:   102,
			Town: Town{
				TownId: 102,
				Name:   "town",
			},
			// Country: Country{
			// 	Name: "country",
			// },
		},
	}
	// db.Create(&m)
	g := Ground{
		GroundId: 101,
		Name:     "ground",
		TownId:   102,
		Town: Town{
			TownId: 102,
			Name:   "town",
		},
	}
	db.Create(&g)
	// db.Find(&m)
	fmt.Printf("%+v", m)
}

type Matches struct {
	ID        int    `json:"id"`
	MatchName string `gorm:"column:match_name;type:varchar(128)" json:"match_name,omitempty"`
	Stage     string `gorm:"column:stage;type:varchar(128)" json:"stage,omitempty"`
	State     string `gorm:"column:state;type:varchar(128)" json:"state,omitempty"`
	// StartDate        time.Time `gorm:"column:start_date;type:datetime" json:"start_date,omitempty"`
	// EndDate          time.Time `gorm:"column:end_date;type:datetime" json:"end_date,omitempty"`
	// StartTime        time.Time `gorm:"column:start_time;type:datetime" json:"start_time,omitempty"`
	IsCancelled      bool      `gorm:"column:is_cancelled;type:TINYINT(1)" json:"is_cancelled,omitempty"`
	Status           string    `gorm:"column:status;type:varchar(128)" json:"status,omitempty"`
	StatusText       string    `gorm:"column:status_text;type:varchar(128)" json:"status_text,omitempty"`
	TossWinnerTeamId int       `gorm:"column:toss_winner_team_id;type:int(11)" json:"toss_winner_team_id,omitempty"`
	ResultStatus     int       `gorm:"column:result_status;type:int(11)" json:"result_status,omitempty"`
	LiveInning       int       `gorm:"column:live_inning;type:int(11)" json:"live_inning,omitempty"`
	CurrentSeriesId  int       `gorm:"column:current_series_id;type:int(11)" json:"current_series_id,omitempty"`
	SeriesId         int       `gorm:"column:series_id;type:int(11)" json:"series_id,omitempty"`
	Series           Series    `gorm:"references:SeriesId" json:"series,omitempty"`
	CurrentMatchId   int       `gorm:"column:current_match_id;type:int(11)" json:"current_match_id,omitempty"`
	TeamsIdList      string    `gorm:"column:teams_id_list;type:varchar(128)" json:"teams_id_list,omitempty"`
	CurrentGroundId  int       `gorm:"column:current_ground_id;type:int(11)" json:"current_ground_id,omitempty"`
	GroundId         int       `gorm:"column:ground_id;type:int(11)" json:"ground_id,omitempty"`
	Ground           Ground    `gorm:"references:GroundId" json:"ground,omitempty"`
	CmsMatchId       string    `gorm:"column:cms_match_id;type:varchar(128)" json:"cms_match_id,omitempty"`
	ContentId        string    `gorm:"column:content_id;type:varchar(128)" json:"content_id,omitempty"`
	CreatedAt        time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP(3)" json:"created_at,omitempty"`
}

type Series struct {
	Id            int    `json:"ID"`
	SeriesId      int    `gorm:"uniqueIndex:by_series_id" json:"series_id,omitempty"`
	Name          string `json:"name,omitempty"`
	LongName      string `json:"long_name,omitempty"`
	AlternateName string `json:"alternate_name,omitempty"`
	Year          int    `json:"year,omitempty"`
	Description   string `json:"description,omitempty"`
	Season        string `json:"season,omitempty"`
	// StartDate     time.Time `json:"start_date,omitempty"`
	// EndDate       time.Time `json:"end_date,omitempty"`
}

type Ground struct {
	ID        int     `json:"id"`
	GroundId  int     `gorm:"uniqueIndex:by_ground_id" json:"ground_id,omitempty"`
	Name      string  `json:"name,omitempty"`
	SmallName string  `json:"small_name,omitempty"`
	LongName  string  `json:"long_name,omitempty"`
	Location  string  `json:"location,omitempty"`
	TownId    int     `gorm:"column:town_id;type:int(11)" json:"town_id,omitempty"`
	Town      Town    `gorm:"references:TownId" json:"town,omitempty"`
	CountryId int     `gorm:"column:country_id;type:int(11)"  json:"country_id,omitempty"`
	Country   Country `gorm:"references:CountryId" json:"country,omitempty"`
}

type Town struct {
	ID       int    `json:"id"`
	TownId   int    `gorm:"column:town_id;type:int(11);uniqueIndex:by_town_id" json:"town_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Area     string `json:"area,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

type Country struct {
	ID        int    `json:"id"`
	CountryId int    `gorm:"uniqueIndex:by_country_id" json:"country_id,omitempty"`
	Name      string `json:"name,omitempty"`
	ShortName string `json:"short_name,omitempty"`
}

type Team struct {
	ID           int    `json:"id"`
	TeamId       int    `gorm:"column:team_id;type:int(11)" json:"team_id,omitempty"`
	Slug         string `json:"slug,omitempty"`
	Name         string `json:"name,omitempty"`
	LongName     string `json:"long_name,omitempty"`
	Abbreviation string `json:"abbreviation,omitempty"`
	IsCountry    bool   `json:"is_country,omitempty"`
}
