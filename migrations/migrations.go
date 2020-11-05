package migrations

import (
	//"fmt"
	"github.com/sieusanh/Banking_App/helpers"
	"github.com/sieusanh/Banking_App/interfaces"
	"github.com/sieusanh/Banking_App/database"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func CreateAccounts(users []interfaces.User) {
	
	for i := 0 ; i < len(users); i++ {
		// Correct one way
		generatedPassword := helpers.HashAndSalt([]byte(users[i].Password))
		user := &interfaces.User{Username: users[i].Username, Email: users[i].Email, Password: generatedPassword}
		database.DB.Create(&user)
	
		account := &interfaces.Account{Type: "Daily Account", Name: string(user.Username + "'s" + " account"), Balance: uint(10000 * int(i+1)), UserID: user.ID}
		database.DB.Create(&account)
	}
}

func Migrate() {
	User := &interfaces.User{}
	Account := &interfaces.Account{}
	Transactions := &interfaces.Transaction{}
	database.DB.AutoMigrate(&User, &Account, &Transactions)	
}
