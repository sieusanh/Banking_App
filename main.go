package main

import (
	
	"github.com/sieusanh/Banking_App/migrations"
	"github.com/sieusanh/Banking_App/interfaces"
	
	"github.com/sieusanh/Banking_App/api"
	"github.com/sieusanh/Banking_App/database"
)

func main() {

	
	// Initialize Database
	database.InitDatabase()

	// Auto Migrate Tables
	migrations.Migrate()

	users := []interfaces.User{
		{Username: "Martin", Email: "martin@martin.com", Password: "martin123"},
		{Username: "Michael", Email: "michael@michael.com", Password: "michael123"},
	}
	
	migrations.CreateAccounts(users)
	
	// Start Rest API Server	
	api.StartApi()

	//api.Shutdown()
}