package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"usual_store/internal/cards"
	"usual_store/internal/messaging"
	"usual_store/internal/models"
	"usual_store/internal/urlsigner"
	"usual_store/internal/validator"

	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v72"
	"golang.org/x/crypto/bcrypt"
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

// Invoice describes payload that sends to microservice
type Invoice struct {
	ID        int       `json:"id"`
	WidgetID  int       `json:"widget_id"`
	Amount    int       `json:"amount"`
	Quantity  int       `json:"quantity"`
	Product   string    `json:"product"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
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
		_, err = w.Write(out)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
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
		_, err = w.Write(out)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
	}
}

// GetAllWidgets get all widgets/products
func (app *application) GetAllWidgets(w http.ResponseWriter, r *http.Request) {
	widgets, err := app.DB.GetAllWidgets()
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}

	out, err := json.MarshalIndent(widgets, "", "  ")
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, applicationJson)
	_, err = w.Write(out)
	if err != nil {
		app.errorLog.Println(err)
		return
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
	_, err = w.Write(out)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

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

	v := validator.New()
	v.Check(len(data.FirstName) > 2, "first_name", "must be at least 3 characters")
	v.Check(len(data.LastName) > 2, "last_name", "must be at least 3 characters")
	if !v.Valid() {
		app.failedValidation(w, r, v.Errors)
		return
	}

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
		subscription, err = card.SubscribeToPlan(stripeCustomer, data.Plan, data.LastFour, "")
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

		invoice := Invoice{
			ID:        order.ID,
			Amount:    3000,
			Product:   "Subscription",
			Quantity:  order.Quantity,
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Email:     data.Email,
			CreatedAt: time.Now(),
		}
		err = app.callInvoiceMicroservice(invoice)
		if err != nil {
			app.errorLog.Println(err)
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
	_, err = w.Write(out)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

// callInvoiceMicroservice calls invoice microservice that create invoice
func (app *application) callInvoiceMicroservice(invoice Invoice) error {
	url := "http://localhost:5000/invoice/create-and-send"
	out, err := json.MarshalIndent(invoice, "", "\t")
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(out))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			app.errorLog.Println(err)
		}
	}(resp.Body)

	return nil
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
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// handle get user by email
	user, err := app.DB.GetUserByEmail(userInput.Email)
	if err != nil {
		fmt.Println("error getting user by email")

		err = app.invalidCredentials(w)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	validPassword, err := app.passwordMatchers(user.Password, userInput.Password)
	if err != nil {
		fmt.Println("error password")

		err = app.invalidCredentials(w)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	if !validPassword {
		fmt.Println("error - not valid password")

		err = app.invalidCredentials(w)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// Call the TokenService to create and store the token
	token, err := app.tokenService.CreateToken(ctx, user, 24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	var payload struct {
		Error     bool          `json:"error"`
		Message   string        `json:"message"`
		Token     *models.Token `json:"authentication_token"`
		ID        int           `json:"id"`
		FirstName string        `json:"first_name"`
		LastName  string        `json:"last_name"`
		Email     string        `json:"email"`
		Role      string        `json:"role"`
	}
	payload.Error = false
	payload.Message = fmt.Sprintf("Token for user %s created.", userInput.Email)
	payload.Token = token
	payload.ID = user.ID
	payload.FirstName = user.FirstName
	payload.LastName = user.LastName
	payload.Email = user.Email
	payload.Role = user.Role
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) CheckAuthentication(w http.ResponseWriter, r *http.Request) {
	user, err := app.authenticateToken(r)
	if err != nil {
		fmt.Println("error authentication")

		err = app.invalidCredentials(w)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Error = false
	payload.Message = fmt.Sprintf("Token for user %s created.", user.Email)
	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

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
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	payload := stripePayload{}
	card := mappingPayloadToCard(app, payload)

	pi, err := card.RetrievePaymentIntent(txnData.PaymentIntent)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	pm, err := card.GetPaymentMethod(txnData.PaymentMethod)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
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
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, txnData)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
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
		err = app.writeJSON(w, http.StatusAccepted, resp)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
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
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	resp.Error = false
	err = app.writeJSON(w, http.StatusCreated, resp)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// Verify the signed URL token
	sign := urlsigner.Signer{
		Secret: []byte(app.config.secretkey),
	}

	// Verify that the token is valid
	valid := sign.VerifyToken(payload.Token)
	if !valid {
		var resp struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		resp.Error = true
		resp.Message = "Invalid or expired reset token"
		err = app.writeJSON(w, http.StatusBadRequest, resp)
		if err != nil {
			app.errorLog.Println(err)
		}
		return
	}

	user, err := app.DB.GetUserByEmail(payload.Email)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	err = app.DB.UpdatePasswordForUser(user, string(newHash))
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	resp.Error = false
	resp.Message = "password changed"

	err = app.writeJSON(w, http.StatusCreated, resp)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) AllSales(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		PageSize    int `json:"page_size"`
		CurrentPage int `json:"current_page"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	allSales, lastPage, totalRecords, err := app.DB.GetAllOrders(payload.PageSize, payload.CurrentPage)
	if err != nil {
		fmt.Println("error getting all sales", err)
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
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

	err = app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) AllSubscriptions(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		PageSize    int `json:"page_size"`
		CurrentPage int `json:"current_page"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	allSubscriptions, lastPage, totalRecords, err := app.DB.GetAllSubscriptions(payload.PageSize, payload.CurrentPage)
	if err != nil {
		fmt.Println("error getting all sales", err)

		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
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

	err = app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) GetSale(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	orderID, err := strconv.Atoi(id)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	order, err := app.DB.GetOrderByID(orderID)
	if err != nil {
		fmt.Println("error getting sale", err)

		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, order)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
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
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: chargeToRefund.Currency,
	}

	err = card.Refund(chargeToRefund.PaymentIntent, chargeToRefund.Amount)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	//update status in DB
	err = app.DB.UpdateOrderStatus(chargeToRefund.ID, 2)
	if err != nil {
		app.infoLog.Println("charge was refunded, but error happens while updating order in DB")
		err = app.badRequest(w, r, errors.New("charge was refunded, but error happens while updating order in DB"))
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}
	app.infoLog.Println("status was refunded")

	// response message with error
	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = "refunded successfully"
	err = app.writeJSON(w, http.StatusOK, resp)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	var subscriptionToCancel struct {
		ID            int    `json:"id"`
		PaymentIntent string `json:"pi"`
		Currency      string `json:"currency"`
	}

	err := app.readJSON(w, r, &subscriptionToCancel)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: subscriptionToCancel.Currency,
	}

	err = card.CancelSubscription(subscriptionToCancel.PaymentIntent)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	//update status in DB
	err = app.DB.UpdateOrderStatus(subscriptionToCancel.ID, 3)
	if err != nil {
		err = app.badRequest(w, r, errors.New("subscription was canceled, but error happens while updating order in DB"))
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// response message with error
	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = "successfully cancelled subscription"
	err = app.writeJSON(w, http.StatusOK, resp)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := app.DB.GetAllUsers()
	if err != nil {
		err := app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, allUsers)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) ShowUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)

	user, err := app.DB.GetUserByID(userID)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, user)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) EditUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)

	var user models.User

	err := app.readJSON(w, r, &user)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	if userID > 0 {
		err = app.DB.EditUser(user)
		if err != nil {
			err = app.badRequest(w, r, err)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			return
		}

		if user.Password != "" {
			newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
			if err != nil {
				err = app.badRequest(w, r, err)
				if err != nil {
					app.errorLog.Println(err)
					return
				}
				return
			}

			err = app.DB.UpdatePasswordForUser(user, string(newHash))
			if err != nil {
				err = app.badRequest(w, r, err)
				if err != nil {
					app.errorLog.Println(err)
					return
				}
				return
			}
		}
	} else {
		newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			err = app.badRequest(w, r, err)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			return
		}
		err = app.DB.AddUser(user, string(newHash))
		if err != nil {
			err = app.badRequest(w, r, err)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			return
		}
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	err = app.writeJSON(w, http.StatusOK, resp)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)
	err := app.DB.DeleteUser(userID)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}
	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	resp.Error = false
	err = app.writeJSON(w, http.StatusOK, resp)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

// GetAllUsers returns paginated list of all users (super_admin only)
func (app *application) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page := 1
	pageSize := 10
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	// Default sort values
	if sortBy == "" {
		sortBy = "id"
	}
	if sortOrder == "" {
		sortOrder = "asc"
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr := r.URL.Query().Get("page_size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100 {
			pageSize = s
		}
	}

	offset := (page - 1) * pageSize

	// Get total count (with search filter if provided)
	totalCount, err := app.DB.GetUserCount(search)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// Get paginated users with search and sort
	users, err := app.DB.GetAllUsersPaginated(offset, pageSize, search, sortBy, sortOrder)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// Calculate total pages
	totalPages := (totalCount + pageSize - 1) / pageSize

	var resp struct {
		Error      bool          `json:"error"`
		Users      []models.User `json:"users"`
		TotalCount int           `json:"total_count"`
		Page       int           `json:"page"`
		PageSize   int           `json:"page_size"`
		TotalPages int           `json:"total_pages"`
	}
	resp.Error = false
	resp.Users = users
	resp.TotalCount = totalCount
	resp.Page = page
	resp.PageSize = pageSize
	resp.TotalPages = totalPages

	err = app.writeJSON(w, http.StatusOK, resp)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

// DeleteUserByID deletes a user (with super_admin protection)
func (app *application) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		err = app.badRequest(w, r, errors.New("invalid user ID"))
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// Get user to check if super_admin
	userToDelete, err := app.DB.GetUserByID(userID)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// If trying to delete super_admin, check count
	if userToDelete.Role == "super_admin" {
		superAdminCount, err := app.DB.GetSuperAdminCount()
		if err != nil {
			err = app.badRequest(w, r, err)
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			return
		}

		if superAdminCount < 2 {
			err = app.badRequest(w, r, errors.New("cannot delete the last super admin"))
			if err != nil {
				app.errorLog.Println(err)
				return
			}
			return
		}
	}

	err = app.DB.DeleteUser(userID)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	resp.Error = false
	resp.Message = "User deleted successfully"

	err = app.writeJSON(w, http.StatusOK, resp)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

// CreateUser creates a new user account
func (app *application) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Role      string `json:"role"` // Optional: admin, supporter, or user (default)
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// Validate input
	if payload.Email == "" || payload.Password == "" || payload.FirstName == "" || payload.LastName == "" {
		err = app.badRequest(w, r, errors.New("all fields are required"))
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// Validate role if provided
	validRoles := map[string]bool{
		"super_admin": true,
		"admin":       true,
		"supporter":   true,
		"user":        true,
	}

	if payload.Role != "" && !validRoles[payload.Role] {
		err = app.badRequest(w, r, errors.New("role must be 'super_admin', 'admin', 'supporter', or 'user'"))
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// Default to 'supporter' if role not specified (for support staff creation)
	if payload.Role == "" {
		payload.Role = "supporter"
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// Create user
	user := models.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Role:      payload.Role,
	}

	err = app.DB.AddUser(user, string(hashedPassword))
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// Get the created user to return the ID
	createdUser, err := app.DB.GetUserByEmail(payload.Email)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	var resp struct {
		Error     bool   `json:"error"`
		Message   string `json:"message"`
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Role      string `json:"role"`
	}
	resp.Error = false
	resp.Message = "User created successfully"
	resp.ID = createdUser.ID
	resp.FirstName = createdUser.FirstName
	resp.LastName = createdUser.LastName
	resp.Email = createdUser.Email
	resp.Role = createdUser.Role

	err = app.writeJSON(w, http.StatusCreated, resp)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

// SendMessageViaKafka sends a message to Kafka for processing (e.g., email sending)
func (app *application) SendMessageViaKafka(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		err = app.badRequest(w, r, err)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// Validate input
	if payload.To == "" || payload.Subject == "" || payload.Message == "" {
		err = app.badRequest(w, r, errors.New("to, subject, and message are required"))
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	// Create Kafka producer
	brokers := []string{"kafka:9092"} // Use Kafka broker from Docker network
	topic := "email-notifications"
	producer := messaging.NewProducer(brokers, topic, app.infoLog)
	defer func() {
		if closeErr := producer.Close(); closeErr != nil {
			app.errorLog.Printf("Error closing Kafka producer: %v", closeErr)
		}
	}()

	// Create email message with plain text (no template)
	data := map[string]interface{}{
		"body": payload.Message,
	}

	// Publish message to Kafka using SendEmail method
	ctx := context.Background()
	err = producer.SendEmail(
		ctx,
		"noreply@usualstore.com", // from
		payload.To,               // to
		payload.Subject,          // subject
		"plain",                  // template (we'll need to create this)
		data,                     // data
		messaging.PriorityNormal, // priority
	)
	if err != nil {
		app.errorLog.Printf("Failed to publish message to Kafka: %v", err)
		err = app.badRequest(w, r, errors.New("failed to queue message"))
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		return
	}

	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	resp.Success = true
	resp.Message = "Message queued for delivery"

	err = app.writeJSON(w, http.StatusOK, resp)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}
