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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Creating a new Employee
func CreateEmployee(c pb.EmployeeDatabaseCrudClient, ctx context.Context) {
	new_emp, err := c.CreateEmployee(ctx, &pb.NewEmp{Name: "Sandeep", ManagerId: 1, DepartmentId: 1})
	checkErr(err)

	log.Printf("Employee Id : %v Added to Database Successfully\n", new_emp.GetId())
}

// Getting all Employees in the database
func GetEmployees(c pb.EmployeeDatabaseCrudClient, ctx context.Context) {
	Employees, err := c.GetEmployees(ctx, &pb.Request{Request: "Get All Employee Details"})
	checkErr(err)

	//printing emp details from the Employees
	log.Printf("Employee Details are as follows:")
	for _, employee := range Employees.Emps {
		fmt.Println(employee.GetId(), employee.GetName(), employee.GetDepartmentId())
	}
}

// Updating the details of Employee
func UpdateEmployees(c pb.EmployeeDatabaseCrudClient, ctx context.Context) {
	updated_data, err := c.UpdateEmployee(ctx, &pb.Emp{Id: 3, Name: "Pyata Sandeep"})
	checkErr(err)

	fmt.Printf("Updated Employee details with Id : %v\n", updated_data.Id)
}

// Deleting the Employee
func DeleteEmployee(c pb.EmployeeDatabaseCrudClient, ctx context.Context) {
	deleted_emp, err := c.DeleteEmployee(ctx, &pb.Emp{Id: 3})
	checkErr(err)

	fmt.Printf(deleted_emp.Response)
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

	//Calling the CreateEmployee Function
	CreateEmployee(c, ctx)

	//Calling the GetEmployees Function
	GetEmployees(c, ctx)

	//Calling the UpdateEmployees Function
	UpdateEmployees(c, ctx)

	//Calling the DeleteEmployee Function
	DeleteEmployee(c, ctx)
}
