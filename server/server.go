package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	pb "grpc_exercise.com/exercise/empmgmt"
	"grpc_exercise.com/exercise/models"
)

// declaring the port number
const (
	port = "localhost:8080"
)

// connecting to the unimplemented methods via a struct
type empServer struct {
	pb.UnimplementedEmployeeDatabaseCrudServer
	db *gorm.DB //linking to db via this struct
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Implementing the CreateEmployee Function
func (s *empServer) CreateEmployee(ctx context.Context, in *pb.NewEmp) (*pb.Emp, error) {
	log.Printf("Creating a Employee in the Database")
	//First Getting the manager details
	id := in.GetManagerId()
	manager := models.Employee{}
	s.db.First(&manager, id)
	newEmp := models.Employee{
		Name:         in.GetName(),
		ManagerID:    &id,
		Manager:      &manager, //populating the manager
		DepartmentID: int(in.GetDepartmentId()),
	}

	s.db.Save(&newEmp)
	return &pb.Emp{Id: int32(newEmp.ID)}, nil
}

// Implementing the GetEmployees Function
func (s *empServer) GetEmployees(ctx context.Context, in *pb.Request) (*pb.Emps, error) {
	log.Printf("Getting All the Employee Details")
	Employees := []models.Employee{}
	EmployeesData := []*pb.Emp{}
	s.db.Find(&Employees)
	for _, employee := range Employees {
		EmployeesData = append(EmployeesData, &pb.Emp{Id: int32(employee.ID), Name: employee.Name, DepartmentId: int32(employee.DepartmentID)})
	}
	return &pb.Emps{Emps: EmployeesData}, nil
}

// Implementing the UpdateEmployee Function
func (s *empServer) UpdateEmployee(ctx context.Context, in *pb.Emp) (*pb.Emp, error) {
	log.Printf("Updating the Employee Details")
	updated_emp := models.Employee{}
	s.db.Model(&models.Employee{}).Where("id=?", int32(in.Id)).Update(in)
	s.db.Find(&updated_emp, int32(in.Id))

	// If the given Employee details are not present in the database
	if updated_emp.ID == 0 {
		return &pb.Emp{}, errors.New("Employee With the given Id is not present in the database")
	} else {
		return &pb.Emp{Id: int32(updated_emp.ID)}, nil
	}
}

// Implementing the DeleteEmployee Function
func (s *empServer) DeleteEmployee(ctx context.Context, in *pb.Emp) (*pb.Response, error) {
	log.Printf("Deleting Employee Details")
	employee := models.Employee{}
	s.db.Find(&employee, int32(in.Id))

	//If the given Employee details are not present in the database
	if employee.ID == 0 {
		return &pb.Response{Response: "Delete Operation Failed"}, errors.New("Employee Details are not present in the database could not delete the employee")
	}
	s.db.Delete(&models.Employee{}, in.GetId())
	return &pb.Response{Response: "Deleted the Employee Successfully"}, nil
}

// Main Function
func main() {

	//Using command line flags to connect with the database
	db_user := flag.String("user", "postgres", "database user")
	db_password := flag.String("password", "root", "database password")
	db_name := flag.String("name", "grpc_exercise", "database name")
	db := models.SetupDB(db_user, db_password, db_name)

	listener, err := net.Listen("tcp", port)
	checkErr(err)

	defer db.Close()

	//Creating a New Server
	s := grpc.NewServer()
	log.Printf("server listening at %v", listener.Addr())

	pb.RegisterEmployeeDatabaseCrudServer(s, &empServer{
		db: db,
	})

	//if connection fails
	err = s.Serve(listener)
	checkErr(err)

}
