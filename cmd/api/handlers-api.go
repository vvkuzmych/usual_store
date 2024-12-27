package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v72"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"usual_store/internal/cards"
	"usual_store/internal/encryption"
	"usual_store/internal/models"
	"usual_store/internal/urlsigner"
)

const contentType = "Content-Type"
const applicationJson = "application/json"

type stripePayload struct {
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Email         string `json:"email"`
	CardBrand     string `json:"card_brand"`
	LastFour      string `json:"last_four"`
	ExpiryMonth   int    `json:"expiry_month"`
	ExpiryYear    int    `json:"expiry_year"`
	ProductID     string `json:"product_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Plan          string `json:"plan"`
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

// GetPaymentIntent get payment intent
func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Println(err)
	}

	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Println(err)
	}

	card := mappingPayloadToCard(app, payload)

	ok := true

	pi, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		app.errorLog.Println(err)
		ok = false
	}

	if ok {
		out, err := json.MarshalIndent(pi, "", "	")
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		w.Header().Set(contentType, applicationJson)
		w.Write(out)
	} else {
		j := jsonResponse{
			OK:      false,
			Message: msg,
			Content: "Invalid amount",
		}

		out, err := json.MarshalIndent(j, "", "  ")
		if err != nil {
			app.errorLog.Println(err)
		}
		w.Header().Set(contentType, applicationJson)
		w.Write(out)
	}
}

// GetWidgetByID get widget by id
func (app *application) GetWidgetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "Widget not found", http.StatusNotFound)
		return
	}
	out, err := json.MarshalIndent(widget, "", "  ")
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set(contentType, applicationJson)
	w.Write(out)
}

// CreateCustomerAndSubscribeToPlan create customer and subscribe to plan
func (app *application) CreateCustomerAndSubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	var data stripePayload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.infoLog.Println(data.LastFour, data.Email, data.PaymentMethod, data.Plan)

	card := mappingPayloadToCard(app, data)

	ok := true
	var subscription *stripe.Subscription
	txnMsg := "Transaction Successful!"

	stripeCustomer, msg, err := card.CreateCustomer(data.PaymentMethod, data.Email)
	if err != nil {
		app.errorLog.Println(err)
		ok = false
		txnMsg = msg
	}
	if ok {
		subscription, err = card.SubscribeToPlan(stripeCustomer, data.Plan, data.Email, data.LastFour, "")
		if err != nil {
			app.errorLog.Println(err)
			ok = false
			txnMsg = "Error Subscribing to Plan"
		}
		app.infoLog.Println("SubscriptionID:", subscription.ID)
	}

	if ok {
		productID, _ := strconv.Atoi(data.ProductID)
		customerID, err := app.SaveCustomer(data.FirstName, data.LastName, data.Email)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		amount, err := strconv.Atoi(data.Amount)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		expiryMonth := data.ExpiryMonth
		expiryYear := data.ExpiryYear
		txn := models.Transaction{
			Amount:              amount,
			Currency:            "EUR",
			LastFour:            data.LastFour,
			ExpiryMonth:         expiryMonth,
			ExpiryYear:          expiryYear,
			TransactionStatusID: 2,
			PaymentIntent:       subscription.ID,
			PaymentMethod:       data.PaymentMethod,
		}
		txnID, err := app.SaveTransaction(txn)
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		order := models.Order{
			TransactionID: txnID,
			CustomerID:    customerID,
			WidgetID:      productID,
			StatusID:      1,
			Quantity:      1,
			Amount:        amount,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		_, err = app.SaveOrder(order)

		if err != nil {
			app.errorLog.Println(err)
			return
		}
	}

	response := jsonResponse{
		OK:      ok,
		Message: txnMsg,
	}
	out, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	w.Header().Set(contentType, "application/json")
	w.Write(out)
}

func mappingPayloadToCard(app *application, data stripePayload) cards.Card {
	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: data.Currency,
	}
	return card
}

// SaveCustomer saves customer and returns id
func (app *application) SaveCustomer(firstName, lastName, email string) (int, error) {
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	err := app.DB.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}

	id, err := app.DB.GetLastInsertedCustomerID()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SaveTransaction saves transaction and returns id
func (app *application) SaveTransaction(txn models.Transaction) (int, error) {

	id, err := app.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// SaveOrder saves order and returns id
func (app *application) SaveOrder(order models.Order) (int, error) {
	id, err := app.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// CreateAuthToken handle creating authenticate token
func (app *application) CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &userInput)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// handle get user by email
	user, err := app.DB.GetUserByEmail(userInput.Email)
	if err != nil {
		fmt.Println("error getting user by email")

		app.invalidCredentials(w)
		return
	}

	validPassword, err := app.passwordMatchers(user.Password, userInput.Password)
	if err != nil {
		fmt.Println("error password")

		app.invalidCredentials(w)
		return
	}

	if !validPassword {
		fmt.Println("error - not valid password")

		app.invalidCredentials(w)
		return
	}

	// Call the TokenService to create and store the token
	token, err := app.tokenService.CreateToken(ctx, user, 24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var payload struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Token   *models.Token `json:"authentication_token"`
	}
	payload.Error = false
	payload.Message = fmt.Sprintf("Token for user %s created.", userInput.Email)
	payload.Token = token
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) CheckAuthentication(w http.ResponseWriter, r *http.Request) {
	user, err := app.authenticateToken(r)
	if err != nil {
		fmt.Println("error authentication")

		app.invalidCredentials(w)
		return
	}

	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Error = false
	payload.Message = fmt.Sprintf("Token for user %s created.", user.Email)
	app.writeJSON(w, http.StatusOK, payload)

}

func (app *application) authenticateToken(r *http.Request) (*models.User, error) {
	ctx := context.Background()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		fmt.Println("no header", authHeader)

		return nil, errors.New("Missing Authorization header")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" || headerParts[1] == "" {
		fmt.Println(" invalid header", authHeader)

		return nil, errors.New("Invalid Authorization header")
	}

	tokenString := headerParts[1]
	if tokenString == "" {
		fmt.Println("token empty", authHeader)

		return nil, errors.New("Invalid Authorization header - no token found")
	}
	if len(tokenString) != 26 {
		fmt.Println("wrong length", authHeader)

		return nil, errors.New("Invalid Authorization header - wrong length")
	}

	user, err := app.tokenService.GetUserForToken(ctx, tokenString)
	if err != nil {
		fmt.Println("token not found", authHeader)

		return nil, errors.New("Invalid Authorization header - matching token not found")
	}

	return user, nil
}

func (app *application) VirtualTerminalPaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	var txnData struct {
		PaymentAmount   int    `json:"amount"`
		PaymentCurrency string `json:"payment_currency"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Email           string `json:"email"`
		PaymentIntent   string `json:"payment_intent"`
		PaymentMethod   string `json:"payment_method"`
		BankReturnCode  string `json:"bank_return_code"`
		ExpiryMonth     int    `json:"expiry_month"`
		ExpiryYear      int    `json:"expiry_year"`
		LastFour        string `json:"last_four"`
	}

	err := app.readJSON(w, r, &txnData)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	payload := stripePayload{}
	card := mappingPayloadToCard(app, payload)

	pi, err := card.RetrievePaymentIntent(txnData.PaymentIntent)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	pm, err := card.GetPaymentMethod(txnData.PaymentMethod)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	txnData.LastFour = pm.Card.Last4
	txnData.ExpiryMonth = int(pm.Card.ExpMonth)
	txnData.ExpiryYear = int(pm.Card.ExpYear)
	txnData.BankReturnCode = pi.Charges.Data[0].ID
	txn := models.Transaction{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         txnData.ExpiryMonth,
		ExpiryYear:          txnData.ExpiryYear,
		BankReturnCode:      pi.Charges.Data[0].ID,
		PaymentMethod:       txnData.PaymentMethod,
		PaymentIntent:       txnData.PaymentIntent,
		TransactionStatusID: 2,
	}
	_, err = app.SaveTransaction(txn)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, txnData)
}

func (app *application) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// verify that email exists
	_, err = app.DB.GetUserByEmail(payload.Email)
	if err != nil {
		var resp struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		resp.Error = true
		resp.Message = "No matching email found in DB"
		app.writeJSON(w, http.StatusAccepted, resp)
		return
	}

	link := fmt.Sprintf("%s/reset-password?email=%s", app.config.frontend, payload.Email)

	sign := urlsigner.Signer{
		Secret: []byte(app.config.secretkey),
	}
	signedLink := sign.GenerateTokenFromString(link)

	var data struct {
		Link string
	}
	data.Link = signedLink

	//send email
	err = app.SendEmail("info@usual_store.com", payload.Email, "Password Reset Request", "password-reset", data)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	resp.Error = false
	app.writeJSON(w, http.StatusCreated, resp)
}

func (app *application) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	encyptor := encryption.Encryption{
		Key: []byte(app.config.secretkey),
	}

	realEmail, err := encyptor.Decrypt(payload.Email)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.DB.GetUserByEmail(realEmail)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	err = app.DB.UpdatePasswordForUser(user, string(newHash))
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	resp.Error = false
	resp.Message = "password changed"

	app.writeJSON(w, http.StatusCreated, resp)
}

func (app *application) AllSales(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		PageSize    int `json:"page_size"`
		CurrentPage int `json:"current_page"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	allSales, lastPage, totalRecords, err := app.DB.GetAllOrders(payload.PageSize, payload.CurrentPage)
	if err != nil {
		fmt.Println("error getting all sales", err)
		app.badRequest(w, r, err)
		return
	}

	var response struct {
		CurrentPage  int             `json:"current_page"`
		PageSize     int             `json:"page_size"`
		LastPage     int             `json:"last_page"`
		TotalRecords int             `json:"total_records"`
		Orders       []*models.Order `json:"orders"`
	}

	response.CurrentPage = payload.CurrentPage
	response.PageSize = payload.PageSize
	response.LastPage = lastPage
	response.TotalRecords = totalRecords
	response.Orders = allSales

	app.writeJSON(w, http.StatusOK, response)
}

func (app *application) AllSubscriptions(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		PageSize    int `json:"page_size"`
		CurrentPage int `json:"current_page"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	allSubscriptions, lastPage, totalRecords, err := app.DB.GetAllSubscriptions(payload.PageSize, payload.CurrentPage)
	if err != nil {
		fmt.Println("error getting all sales", err)

		app.badRequest(w, r, err)
		return
	}

	var response struct {
		CurrentPage   int             `json:"current_page"`
		PageSize      int             `json:"page_size"`
		LastPage      int             `json:"last_page"`
		TotalRecords  int             `json:"total_records"`
		Subscriptions []*models.Order `json:"subscriptions"`
	}

	response.CurrentPage = payload.CurrentPage
	response.PageSize = payload.PageSize
	response.LastPage = lastPage
	response.TotalRecords = totalRecords
	response.Subscriptions = allSubscriptions

	app.writeJSON(w, http.StatusOK, response)
}

func (app *application) GetSale(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	orderID, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	order, err := app.DB.GetOrderByID(orderID)
	if err != nil {
		fmt.Println("error getting sale", err)

		app.badRequest(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, order)
}

func (app *application) RefundCharge(w http.ResponseWriter, r *http.Request) {
	var chargeToRefund struct {
		ID            int    `json:"id"`
		PaymentIntent string `json:"pi"`
		Amount        int    `json:"amount"`
		Currency      string `json:"currency"`
	}

	err := app.readJSON(w, r, &chargeToRefund)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: chargeToRefund.Currency,
	}

	err = card.Refund(chargeToRefund.PaymentIntent, chargeToRefund.Amount)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	//update status in DB
	err = app.DB.UpdateOrderStatus(chargeToRefund.ID, 2)
	if err != nil {
		app.badRequest(w, r, errors.New("charge was refunded, but error happens while updating order in DB"))
		return
	}

	// response message with error
	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = "refunded successfully"
	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	var subscriptionToCancel struct {
		ID            int    `json:"id"`
		PaymentIntent string `json:"pi"`
		Currency      string `json:"currency"`
	}

	err := app.readJSON(w, r, &subscriptionToCancel)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: subscriptionToCancel.Currency,
	}

	err = card.CancelSubscription(subscriptionToCancel.PaymentIntent)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	//update status in DB
	err = app.DB.UpdateOrderStatus(subscriptionToCancel.ID, 3)
	if err != nil {
		app.badRequest(w, r, errors.New("subscription was canceled, but error happens while updating order in DB"))
		return
	}

	// response message with error
	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = "successfully cancelled subscription"
	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := app.DB.GetAllUsers()
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, allUsers)
}
