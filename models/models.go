package models

import (
	"github.com/jinzhu/gorm"
)

func SetupDB(db_user *string, db_password *string, db_name *string) *gorm.DB {
	database_link := "user=" + *db_user + " password=" + *db_password + " dbname=" + *db_name + " sslmode=disable"

	db, err := gorm.Open("postgres", database_link)
	if err != nil {
		panic(err.Error())
	}

	db.DropTable(&Department{})
	db.CreateTable(&Department{})

	db.DropTable(&Employee{})
	db.CreateTable(&Employee{})

	dep1 := Department{
		Name: "FrontEnd",
		Employee: []Employee{
			{Name: "Sai Chander"},
		},
	}
	dep2 := Department{
		Name: "BackEnd",
		Employee: []Employee{
			{Name: "Mani"},
		},
	}
	db.Save(&dep1)
	db.Save(&dep2)

	return db
}

// Employee schema
type Employee struct {
	gorm.Model
	Name         string
	ManagerID    *int32
	Manager      *Employee
	DepartmentID int
}

// Department schema
type Department struct {
	gorm.Model
	Name     string
	Employee []Employee
}
