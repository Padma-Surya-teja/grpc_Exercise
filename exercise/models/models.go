package models

import (
	"github.com/jinzhu/gorm"
)

type Employee struct {
	gorm.Model
	Name         string
	ManagerID    *int32
	Manager      *Employee
	DepartmentID int
}

type Department struct {
	gorm.Model
	DeptName string
	Employee []Employee
}
