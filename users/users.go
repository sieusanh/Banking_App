package users

import (
	"time"
	"github.com/sieusanh/Banking_App/helpers"
	"github.com/sieusanh/Banking_App/interfaces"
	"github.com/sieusanh/Banking_App/database"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// User Login Function
func Login(username, pass string) (map[string]interface{}, time.Time) {

	// Input Data Validation
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: pass, Valid: "password"},
		})

	if valid {

		user := &interfaces.User{}

		// Check if this username exis
		if database.DB.Where("username = ?", username).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "User not found"}, time.Time{}
		}

		// Verify password
		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{"message": "Wrong password"}, time.Time{}
		}

		// Find accounts for the user
		// Return accounts if user.ID match with UserID of the Account
		accounts := []interfaces.ResponseAccount{}
		database.DB.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

		response, expirationTime := prepareResponse(user, accounts, true) //
		
		return response, expirationTime
	} 
	
	return map[string]interface{}{"message": "not valid values"}, time.Time{}
}

// prepareToken function will process token logic matters
func prepareToken(user *interfaces.User) (string, time.Time) {
	expirationTime := time.Now().Add(time.Minute * 60)
	tokenContent := jwt.MapClaims {
		"user_id" : user.ID, 
		"expiry" : expirationTime.Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	return token, expirationTime
}

// prepareResponse
func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount, withToken bool) (map[string]interface{}, time.Time) {
	responseUser := &interfaces.ResponseUser {
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		Accounts: accounts,
	}

	var response = map[string]interface{}{"message": "all is fine"}
	var expirationTime time.Time
	var token string
	// Add withToken feature to prepare response
	if withToken {
		token, expirationTime = prepareToken(user)
		response["jwt"] = token
	}
	response["data"] = responseUser

	return response, expirationTime
}

// User Registration Function
func Register(username, email, pass string) (map[string]interface{}, time.Time) {
	// Add validation to registration
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
		})	

	if valid {
		user := &interfaces.User{}
		
		if database.DB.Where("username = ?", username).First(&user).RecordNotFound() {
			// Correct one way
			generatedPassword := helpers.HashAndSalt([]byte(pass))
			user := &interfaces.User{Username: username, Email: email, Password: generatedPassword}
			database.DB.Create(&user)
			
			account := &interfaces.Account{Type: "Daily Account", Name: string(username + "'s" + " account"), Balance: uint(10000 * int(user.ID+1)), UserID: user.ID}
			database.DB.Create(&account)
			
			accounts := []interfaces.ResponseAccount{}
			respAccount := interfaces.ResponseAccount{ID: account.ID, Name: account.Name, Balance: int(account.Balance)}
			accounts = append(accounts, respAccount)
			response, expirationTime := prepareResponse(user, accounts, true)

			return response, expirationTime

		}
		return map[string]interface{}{"message": "Username has alrealdy existed"}, time.Time{}
	} 
	
	return map[string]interface{}{"message": "not valid values"}, time.Time{}
}

//
func GetUser(id string, jwt string) map[string]interface{} {
	isValid := helpers.ValidateToken(id, jwt)
	// Find and return user
	if isValid {
		
		user := &interfaces.User{}

		// Check if this username exis
		if database.DB.Where("id = ?", id).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "User not found"}
		}

		// Find accounts for the user
		// Return accounts if user.ID match with UserID of the Account
		accounts := []interfaces.ResponseAccount{}
		database.DB.Table("accounts").Select("id, name, balance").Where("user_id = ? ", user.ID).Scan(&accounts)

		var response, _ = prepareResponse(user, accounts, false) //

		return response
	} 

	return map[string]interface{}{"message": "Not valid token"}
}

