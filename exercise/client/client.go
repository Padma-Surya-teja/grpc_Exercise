package main

import (
	"context"
	"fmt"

	// "fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "grpc_exercise.com/exercise/empmgmt"
)

// declaring port for the client to run on
const (
	address = "localhost:8080"
)

func CreateEmployee(c pb.EmployeeDatabaseCrudClient, ctx context.Context) {
	new_emp, err := c.CreateEmp(ctx, &pb.NewEmp{Name: "Sandeep", ManagerId: 1, DepartmentId: 1})
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Employee Name: %v, Manager Id: %v, Department Id: %v", new_emp.GetName(), new_emp.GetManagerId(), new_emp.GetDepartmentId())
}

func GetAllEmployees(c pb.EmployeeDatabaseCrudClient, ctx context.Context) {
	// total emp from client
	total_emp, err := c.ReadEmp(ctx, &pb.Request{})
	if err != nil {
		log.Printf("error getting all employees")
	}
	//printing emp details from the total_emp
	for _, employee := range total_emp.Emps {
		fmt.Println(employee.GetId(), employee.GetName(), employee.GetDepartmentId())
	}
}

func UpdateEmployees(c pb.EmployeeDatabaseCrudClient, ctx context.Context) {
	//updating the manager of an emp
	upd_res, err := c.UpdateEmp(ctx, &pb.Emp{Id: 3, Name: "Pyata Sandeep"})
	if err != nil {
		log.Fatalf("\nError Updating employe %v", err)
	}
	fmt.Printf("\nUpdated Employee : %v", upd_res)
}

func DeleteEmployee(c pb.EmployeeDatabaseCrudClient, ctx context.Context) {
	//deleting an instance of an emp
	deleted_emp, err := c.DeleteEmp(ctx, &pb.Emp{Name: "Sandeep"})
	if err != nil {
		fmt.Println(deleted_emp)
		panic(err.Error())
	}
}

func main() {

	connection, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
	log.Printf("Client Started")
	if err != nil {
		log.Fatal("Connection Failed", err.Error())
	}
	defer connection.Close()
	c := pb.NewEmployeeDatabaseCrudClient(connection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//create emp from client
	CreateEmployee(c, ctx)

	GetAllEmployees(c, ctx)

	UpdateEmployees(c, ctx)

	DeleteEmployee(c, ctx)
}
