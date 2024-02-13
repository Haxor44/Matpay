package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type BillingAddress struct {
	EmailAddress string `json:"email_address"`
	PhoneNumber  string `json:"phone_number"`
	CountryCode  string `json:"country_code"`
	FirstName    string `json:"first_name"`
	MiddleName   string `json:"middle_name"`
	LastName     string `json:"last_name"`
	Line1        string `json:"line_1"`
	Line2        string `json:"line_2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	ZipCode      string `json:"zip_code"`
}

type paymentInfo struct {
	ID             string         `json:"id"`
	Currency       string         `json:"currency"`
	Amount         float64        `json:"amount"`
	Description    string         `json:"description"`
	CallbackURL    string         `json:"callback_url"`
	NotificationID string         `json:"notification_id"`
	Branch         string         `json:"branch"`
	BillingAddress BillingAddress `json:"billing_address"`
}

type Token struct {
	Error      interface{} `json:"error"`
	ExpiryDate string      `json:"expiryDate"`
	Message    string      `json:"message"`
	Status     string      `json:"status"`
	Token      string      `json:"token"`
}

type ipnUrl struct {
	Url string `json:"url"`
}

func main() {
	// Register with pesapal servers

	mux := http.NewServeMux()
	mux.HandleFunc("/ipn", registerIpn)
	mux.HandleFunc("/getipn", getRegisteredIpn)
	mux.HandleFunc("/pay", submitOrder)
	mux.HandleFunc("/callback", callbackUrl)
	mux.HandleFunc("/test", test)
	err := http.ListenAndServe(":8085", mux)
	fmt.Printf("Starting server on port 8085!!!")
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server %s\n", err)
		os.Exit(1)
	}
}

func callbackUrl(w http.ResponseWriter, r *http.Request) {
	token := getAcessToken()
	fmt.Printf("This is the callback: ")
	orderID := r.URL.Query().Get("OrderTrackingId")
	url := "https://pay.pesapal.com/v3/api/Transactions/GetTransactionStatus?orderTrackingId=" + orderID
	//send a get request to pesapal endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("%s", err)
	}

	var bearer = "Bearer " + token
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	// sending the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("Response is: %s\n", response)
	w.Write(response)
}

func submitOrder(w http.ResponseWriter, r *http.Request) {
	//get accessToken
	token := getAcessToken()
	//accept json data
	decoder := json.NewDecoder(r.Body)
	var payInfo paymentInfo
	err := decoder.Decode(&payInfo)

	if err != nil {
		fmt.Printf("%s", err)
	}

	payment := paymentInfo{
		ID:             payInfo.ID,
		Currency:       payInfo.Currency,
		Amount:         payInfo.Amount,
		Description:    payInfo.Description,
		CallbackURL:    payInfo.CallbackURL,
		NotificationID: payInfo.NotificationID,
		Branch:         payInfo.Branch,
		BillingAddress: BillingAddress{
			EmailAddress: payInfo.BillingAddress.EmailAddress,
			PhoneNumber:  payInfo.BillingAddress.PhoneNumber,
			CountryCode:  payInfo.BillingAddress.CountryCode,
			FirstName:    payInfo.BillingAddress.FirstName,
			MiddleName:   payInfo.BillingAddress.MiddleName,
			LastName:     payInfo.BillingAddress.LastName,
			Line1:        payInfo.BillingAddress.Line1,
			Line2:        payInfo.BillingAddress.Line2,
			City:         payInfo.BillingAddress.City,
			State:        payInfo.BillingAddress.State,
			PostalCode:   payInfo.BillingAddress.PostalCode,
			ZipCode:      payInfo.BillingAddress.ZipCode,
		},
	}

	paramBody, _ := json.Marshal(payment)

	res := bytes.NewBuffer(paramBody)
	client := &http.Client{}
	resp, err := http.NewRequest("POST", "https://pay.pesapal.com/v3/api/Transactions/SubmitOrderRequest", res)

	if err != nil {
		fmt.Printf("%s", err)
	}
	var bearer = "Bearer " + token
	resp.Header.Add("Accept", "application/json")
	resp.Header.Add("Content-Type", "application/json")
	resp.Header.Add("Authorization", bearer)

	req, err := client.Do(resp)
	if err != nil {
		fmt.Printf("%s", err)
	}

	response, err := ioutil.ReadAll(req.Body)

	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("%s", response)
	w.Write(response)
	defer req.Body.Close()

}

func test(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payment paymentInfo
	err := decoder.Decode(&payment)
	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("Amount paid is:%f", payment.Amount)
}
func getRegisteredIpn(w http.ResponseWriter, r *http.Request) {

	//get accessToken
	token := getAcessToken()

	req, err := http.NewRequest("GET", "https://pay.pesapal.com/v3/api/URLSetup/GetIpnList", nil)

	if err != nil {
		fmt.Printf("%s", err)
	}
	var bearer = "Bearer " + token
	req.Header.Set("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("%s", err)
	}

	response, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("%s", response)

	defer resp.Body.Close()
}
func registerIpn(w http.ResponseWriter, r *http.Request) {
	// get url from the user
	decoder := json.NewDecoder(r.Body)
	var ipn_url ipnUrl
	err := decoder.Decode(&ipn_url)
	if err != nil {
		fmt.Printf("%s", err)
	}
	//get accessToken
	token := getAcessToken()

	params, _ := json.Marshal(map[string]string{
		"url":                   ipn_url.Url,
		"ipn_notification_type": "POST",
	})
	var bearer1 = "Bearer " + token
	var bearer = bearer1 //fmt.Sprintf("%q", bearer1)
	fmt.Printf("Bearer is: %s\n", bearer)
	resBody := bytes.NewBuffer(params)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://pay.pesapal.com/v3/api/URLSetup/RegisterIPN", resBody)

	if err != nil {
		fmt.Printf("error occured: %s", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearer)
	for name, values := range req.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("%s", err)
	}

	response, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("ipn registered:%s\n", response)

	defer req.Body.Close()
}

func getAcessToken() string {
	godotenv.Load(".env")
	client_key := os.Getenv("CLIENT_KEY")
	client_secret := os.Getenv("CLIENT_SECRET")
	paramBody, _ := json.Marshal(map[string]string{
		"consumer_key":    client_key,
		"consumer_secret": client_secret,
	})

	resBody := bytes.NewBuffer(paramBody)

	resp, err := http.Post("https://pay.pesapal.com/v3/api/Auth/RequestToken", "application/json", resBody)

	if err != nil {
		fmt.Printf("Error sending data %s", err)
	}
	//fmt.Printf("Access token is:%s", token.Token)

	token, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response is: %s", token)
	if err != nil {
		fmt.Printf("Error reading Body %s", err)
	}
	var accessToken Token
	json.Unmarshal(token, &accessToken)
	fmt.Printf("Token is:%s\n", accessToken.Token)
	//tkn, _ := json.Marshal(token)
	//accessToken := bytes.Trim(tkn, "\"")
	//fmt.Printf("Access token is:%s\n", accessToken)

	defer resp.Body.Close()
	return string(accessToken.Token)
}
