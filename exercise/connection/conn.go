package connection

import (
	"github.com/jinzhu/gorm"
	"grpc_exercise.com/exercise/models"
)

func SetupDB() {
	db, err := gorm.Open("postgres", "user=postgres password=root dbname=grpc_exercise sslmode=disable")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	db.DropTable(&models.Department{})
	db.CreateTable(&models.Department{})

	db.DropTable(&models.Employee{})
	db.CreateTable(&models.Employee{})

	dep1 := models.Department{
		DeptName: "FrontEnd",
		Employee: []models.Employee{
			{Name: "Sai Chander", DepartmentID: 1},
		},
	}
	dep2 := models.Department{
		DeptName: "BackEnd",
		Employee: []models.Employee{
			{Name: "Mani", DepartmentID: 2},
		},
	}
	db.Save(&dep1)
	db.Save(&dep2)

}
