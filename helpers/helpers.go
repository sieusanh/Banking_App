package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"regexp"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	//"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/sieusanh/Banking_App/interfaces"
)

// Error Handling Function
func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// Password Encrypt Function
func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	HandleErr(err)

	return string(hashed)
}	

// USer Data Validation Function
func Validation(values []interfaces.Validation) bool {
	username := regexp.MustCompile(`^([A-Za-z0-9]{5,})+$`)
	email := regexp.MustCompile(`^[A-Za-z0-9]+[@]+[A-Za-z0-9]+[.]+[A-Za-z]+$`)

	for i := 0; i < len(values); i++ {
		switch values[i].Valid {
			case "username":
				if !username.MatchString(values[i].Value) {
					return false
				}
			case "email":
				if !email.MatchString(values[i].Value){
					return false
				}
			case "password":
				if len(values[i].Value) < 5 {
					return false
				}
		}
	}
	return true
}

// *Middleware*

// Server Internal Error Handler
func PanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			error := recover()
			if error != nil {
				log.Println(error)

				resp := interfaces.ErrResponse{Message: "Internal server error"}
				json.NewEncoder(w).Encode(resp)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
/* * */

// JWT Validation
func ValidateToken(id, jwtToken string) bool {

	// Remove "Bearer" from Token string
	cleanJWT := strings.Replace(jwtToken, "Bearer ", "", -1)

	// Verify the JWT token
	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cleanJWT, tokenData, func (token *jwt.Token) (interface{}, error) {
		return []byte("TokenPassword"), nil
	})
	HandleErr(err)

	// And verify if the id from the authenticated token is the same as id that sent to API
	var userId, _ = strconv.ParseFloat(id, 8)
	if token.Valid && tokenData["user_id"] == userId {
		return true
	} 
	
	return false
}