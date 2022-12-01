package main

import (
	"context"
	"log"
	"net"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"grpc_exercise.com/exercise/connection"
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

// server side implementation of creating emp
func (s *empServer) CreateEmp(ctx context.Context, in *pb.NewEmp) (*pb.Emp, error) {
	log.Printf("createEmp method called from server side")
	id := in.GetManagerId()
	manager := models.Employee{}
	s.db.First(&manager, id)
	newEmp := models.Employee{
		Name:         in.GetName(),
		ManagerID:    &id,
		Manager:      &manager,
		DepartmentID: int(in.GetDepartmentId()),
	}
	s.db.Save(&newEmp)
	return &pb.Emp{Id: int32(newEmp.ID), Name: in.GetName(), ManagerId: in.GetManagerId(), DepartmentId: in.GetDepartmentId()}, nil
}

// server side implementation of reading emp
func (s *empServer) ReadEmp(ctx context.Context, in *pb.Request) (*pb.Emps, error) {
	log.Printf("readEmp method called from server side")
	Employees := []models.Employee{}
	EmployeesData := []*pb.Emp{}
	s.db.Find(&Employees)
	for _, employee := range Employees {
		EmployeesData = append(EmployeesData, &pb.Emp{Id: int32(employee.ID), Name: employee.Name, DepartmentId: int32(employee.DepartmentID)})
	}
	return &pb.Emps{Emps: EmployeesData}, nil
}

// server side implementation of update emp
// updates the manager of an emp
func (s *empServer) UpdateEmp(ctx context.Context, in *pb.Emp) (*pb.Emp, error) {
	log.Printf("updatedEmp method called from server side")
	var res pb.Emp
	log.Printf("id : %v", in.Id)
	s.db.Model(&models.Employee{}).Where("id=?", int32(in.Id)).Update(in)
	s.db.Find(&res, int32(in.Id))

	return &res, nil
}

// server side implementation of delete emp
func (s *empServer) DeleteEmp(ctx context.Context, in *pb.Emp) (*pb.Response, error) {
	log.Printf("deleteEmp method called from server side")
	s.db.Where(&models.Employee{Name: in.GetName()}).Delete(&models.Employee{})
	return &pb.Response{}, nil
}

func main() {

	//creating gorm instance
	connection.SetupDB()
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err.Error())
	}

	//connecting gorm with grpc
	connection, err := gorm.Open("postgres", "user=postgres password=root dbname=grpc_exercise sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	defer connection.Close()

	//creating a new server
	s := grpc.NewServer()
	log.Printf("server listening at %v", listener.Addr())

	pb.RegisterEmployeeDatabaseCrudServer(s, &empServer{
		db: connection,
	})

	//if connection fails
	if err := s.Serve(listener); err != nil {
		log.Fatal(err.Error())
	}

}
