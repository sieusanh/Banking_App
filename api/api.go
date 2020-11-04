package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"context"
	"os/signal"
	"github.com/sieusanh/Banking_App/helpers"
	"github.com/sieusanh/Banking_App/users"
	"github.com/sieusanh/Banking_App/useraccounts"
	"github.com/sieusanh/Banking_App/transactions"
	"github.com/gorilla/mux"
)

type Login struct {
	Username string
	Password string
}

type Register struct {
	Username string
	Email string
	Password string
}

// Interface that will be responsible for the body of the transaction call
type TransactionBody struct {
	UserId uint
	From uint
	To uint
	Amount int
}

func login(w http.ResponseWriter, r *http.Request) {
	// Read the body of API request
	body := readBody(r)

	// Handle login
	var formattedBody Login
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)

	login, expirationTime := users.Login(formattedBody.Username, formattedBody.Password)

	// Set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: login["jwt"].(string), // convert interface{} to string
		Expires: expirationTime, // convert interface{} to time.Time
	})

	// Prepare response
	apiResponse(login, w)
}

func register(w http.ResponseWriter, r *http.Request) {
	// Read the body of API request
	body := readBody(r)

	// Handle registration
	var formattedBody Register
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)

	register, expirationTime := users.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password)

	// Set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: register["jwt"].(string), // convert interface{} to string
		Expires: expirationTime, // convert interface{} to time.Time
	})

	// Prepare response
	apiResponse(register, w)
}

func getUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userId := vars["id"]

	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	auth := tknStr
	
	user := users.GetUser(userId, auth)
	apiResponse(user, w)
}

func transaction(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	
	// Get token string from Cookie
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	auth := c.Value
	var formattedBody TransactionBody
	err = json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)

	transaction := useraccounts.Transaction(formattedBody.UserId, formattedBody.From, formattedBody.To, formattedBody.Amount, auth)
	apiResponse(transaction, w)
}

func getMyTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userID"]

	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	auth := tknStr
	
	transactions := transactions.GetMyTransactions(userId, auth)
	apiResponse(transactions, w)
}

// Create Router, handle our API endpoints
func StartApi() {
	router := mux.NewRouter()
	// Add panic handler middleware
	router.Use(helpers.PanicHandler)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/transaction", transaction).Methods("POST")
	router.HandleFunc("/transaction/{userID}", getMyTransactions).Methods("GET")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	fmt.Println("App is working on port: 8888")

	// Create Http Server
	server := &http.Server {
		Addr: ":8888",
		Handler: router,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout: 1*time.Second,
	}

	//log.Fatal(http.ListenAndServe(":8888", router))

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	
	// Graceful Shutdown
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	fmt.Println("Received terminate, graceful shutdown", sig) // Ctrl + C

	timeout_context, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// Actually important when we need to update the server
	server.Shutdown(timeout_context)
}

//
func readBody(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)

	return body
}

//
func apiResponse(call map[string]interface{}, w http.ResponseWriter) {
	// Prepare response
	// Check if the login message is equal to "all is fine"
	if call["message"] == "all is fine" {
		resp := call
		json.NewEncoder(w).Encode(resp)

		// Handle error in else
	} else {
		//resp := interfaces.ErrResponse{Message: "Wrong username or password"}

		// we should return the whole call variable instead of "data"
		resp := call
		json.NewEncoder(w).Encode(resp)
	}
}
