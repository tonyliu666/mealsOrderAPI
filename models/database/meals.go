package database

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/beego/beego/orm"
	"github.com/joho/godotenv"
)

var o orm.Ormer
var connectionInfo string

type Ingredients struct {
	ID           int     `json:"id" orm:"auto"`
	Carolie      float64 `json:"carolie"`
	Protein      float64 `json:"protein"`
	Fat          float64 `json:"fat"`
	Carbohydrate float64 `json:"carbohydrate"`
}

type Diets struct {
	ID       int    `json:"id" orm:"auto"`
	Name     string `json:"food_name"`
	Location string `json:"where_eaten"`
	// the time the meal was eaten
	Date      string `json:"date_eaten"`
	Time      string `json:"time_eaten"`
	Periods   string `json:"periods"`
	TimeSlots string `json:"time_slots"`
	// foreign key
	Ingredients *Ingredients `orm:"rel(fk)"`
}

func init() {
	// read the connection information
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	connectionInfo = "user=" + os.Getenv("DBuser") + " password=" + os.Getenv("DBpassword") + " port=" + os.Getenv("DBport") + " dbname=" + os.Getenv("DBname") + " sslmode=" + os.Getenv("DBmode")
	log.Println("Database connection info: ", connectionInfo)
	orm.RegisterModel(new(Diets), new(Ingredients), new(Client))
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", connectionInfo)
	log.Println("Database connection established")
	o = orm.NewOrm()
}
func (f *Ingredients) Save() error {
	_, err := o.Insert(f)
	return err
}
func (f *Ingredients) Read() error {
	err := o.Read(f)
	return err
}

func (f *Diets) Save() error {
	_, err := o.Insert(f.Ingredients)
	if err != nil {
		return err
	}
	_, err = o.Insert(f)
	if err != nil {
		return err
	}
	return nil
}
func (f *Diets) Read() error {
	err := o.Read(f)
	return err
}
func QueryDates(timeslots string, period string) ([]Diets, error) {
	var tmpResults []Diets
	var diets []Diets
	qs := o.QueryTable("diets")
	_, err := qs.Filter("time_slots", timeslots).RelatedSel().All(&tmpResults)

	if err != nil {
		return nil, err
	}
	// filter the results by the period, I need to find the meals prior to the current date within the period.
	periodInt, err := strconv.Atoi(period)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	//substract the period from the current date
	prevDate := now.AddDate(0, 0, -periodInt)

	// filter the queries that are within the prevDate and now
	for _, diet := range tmpResults {
		if diet.Date >= prevDate.Format("2006-01-02") && diet.Date <= now.Format("2006-01-02") {
			diets = append(diets, diet)
		}
	}

	return diets, nil
}
