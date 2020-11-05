package transactions

import (
	"github.com/sieusanh/Banking_App/helpers"
	"github.com/sieusanh/Banking_App/interfaces"
	"github.com/sieusanh/Banking_App/database"
	"github.com/sieusanh/Banking_App/users"
)

// Save info about a transfer in transation history
func CreateTransaction(From, To uint, Amount int) {
	transaction := &interfaces.Transaction{From: From, To: To, Amount: Amount}
	
	database.DB.Create(&transaction)
}

func GetTransactionByAccount(id uint) []interfaces.ResponseTransaction{
	transactions := []interfaces.ResponseTransaction{}
	database.DB.Table("transactions").Select("id, transactions.from, transactions.to, amount").Where(interfaces.Transaction{From: id}).Or(interfaces.Transaction{To: id}).Scan(&transactions)
	return transactions
}

// Get Transaction History List with User's ID
func GetMyTransactions(id string, jwt string) map[string]interface{} {
	isValid := helpers.ValidateToken(id, jwt)
	if isValid {
		accounts := []interfaces.ResponseAccount{}
		database.DB.Table("accounts").Select("id, name, balance").Where("user_id = ? ", id).Scan(&accounts)

		transactions := []interfaces.ResponseTransaction{}
		for i := 0; i < len(accounts); i++ {
			accTransactions := GetTransactionByAccount(accounts[i].ID)
			transactions = append(transactions, accTransactions...)
		}

		var response = map[string]interface{}{"message": "all is fine"}
		response["data"] = transactions
		return response
	} 
		
	return users.FailResponse("Not valid token")
	
}