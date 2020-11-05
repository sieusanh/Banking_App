package useraccounts

import (
	"github.com/sieusanh/Banking_App/helpers"
	"github.com/sieusanh/Banking_App/interfaces"
	"github.com/sieusanh/Banking_App/transactions"
	"github.com/sieusanh/Banking_App/database"
	"github.com/sieusanh/Banking_App/users"
	"fmt"
)

// Define Transaction Type
var WITHDRAW int = -1
var TRANSFER int = 0
var DEPOSIT int = 1

// Update an account by id and amount
func updateAccount(id uint, amount int) interfaces.ResponseAccount{

	account := interfaces.Account{}
	responseAcc := interfaces.ResponseAccount{}

	database.DB.Where("id = ?", id).First(&account)
	account.Balance = uint(amount)
	database.DB.Save(&account)

	responseAcc.ID = account.ID
	responseAcc.Name = account.Name
	responseAcc.Balance = int(account.Balance)

	return responseAcc
}

// Find an account by id
func getAccount(id uint) *interfaces.Account {
	
	account := &interfaces.Account{}

	// Check if this account exis
	if database.DB.Where("id = ?", id).First(&account).RecordNotFound() {
		return nil
	}

	return account
}

// // Bank transfer between two different account's ID
func Transaction(userId, from, to uint, amount int, jwt string, trans_type int) map[string]interface{} {
	// Convert userId from uint to string
	userIdString := fmt.Sprint(userId)

	// Validate JWT
	isValid := helpers.ValidateToken(userIdString, jwt)
	if isValid {
		var updatedAccount interfaces.ResponseAccount
		if from != to && trans_type == TRANSFER{ // from != to
			// Sender's account variable
			fromAccount := getAccount(from)
			// Receiver's account variable
			toAccount := getAccount(to)

			// Handle errors
			// Verify if both accounts exist
			if fromAccount == nil || toAccount == nil {
				return users.FailResponse("Account not found")
			} else if fromAccount.UserID != userId {
				return users.FailResponse("You are not owner of the account")
			} else if int(fromAccount.Balance) < amount {
				return users.FailResponse("Sender Account Balance is too small")
			}

			// Update accounts
			updatedAccount = updateAccount(from, int(fromAccount.Balance) - amount)
			updateAccount(to, int(toAccount.Balance) + amount)

		} else if from == to { // from == to
			if trans_type == WITHDRAW {
				fromAccount := getAccount(from)

				// Handle errors
				// Verify if both accounts exist
				if fromAccount == nil {
					return users.FailResponse("Account not found")
				} else if fromAccount.UserID != userId {
					return users.FailResponse("You are not owner of the account")
				} else if int(fromAccount.Balance) < amount {
					return users.FailResponse("Account Balance is too small for withdrawing")
				}
				// Update account
				updatedAccount = updateAccount(from, int(fromAccount.Balance) - amount)

			} else if trans_type == DEPOSIT {
				// Sender's account variable
				fromAccount := getAccount(from)

				// Handle errors
				// Verify if both accounts exist
				if fromAccount == nil {
					return users.FailResponse("Account not found")
				} else if fromAccount.UserID != userId {
					return users.FailResponse("You are not owner of the account")
				}

				// Update accounts
				updatedAccount = updateAccount(from, int(fromAccount.Balance) + amount)
			}
		}
			
		// To have the info about a transfer in our history
		transactions.CreateTransaction(from, to, amount)

		// Prepare Response
		var response = map[string]interface{}{"message": "all is fine"}
		response["data"] = updatedAccount
		return response
	}
	return users.FailResponse("Not valid token")
}
