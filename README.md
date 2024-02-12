# Matpay</br>
A simple payment api solution written in go to help in the processing of payments.
It queries the pesapal api to enable successful transactions.</br>
A developer that will use this will just have to provide environemt variables in their deployment environment to use this api.</br>
# API Endpoints</br>
YOUR_DOMAIN_HERE/ipn - this endpoint is used to register your url with pesapal servers. It is a POST request and you provide your url in the body as json.</br></br>
YOUR_DOMAIN_HERE/getipn - this endpoint will fetch all your registered urls on pesapal servers. This is a simple GET request.</br></br>
YOUR_DOMAIN_HERE/pay - this enpoint  submits all payment data such as user email, amount etc. to the pesapal servers.This is a POST request isssued to this endpooint containing the user data and returns a json response with a url which can then be rendered in an iframe on your website to complete the payment.</br></br>

# Example
Here is an example to make a payment request on the /pay endpoint.</br>
curl --location 'YOUR_DOMAIN_HERE/pay' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "AA1122-3333ZZ",
    "currency": "KES",
    "amount": 10.00,
    "description": "Payment description goes here",
    "callback_url": "YOUR_DOMAIN_HERE/callback",
    "redirect_mode": "",
    "notification_id": "YOUR_NOTIFICATION_ID",
    "branch": "Store Name - HQ",
    "billing_address": {
        "email_address": "test@example.com",
        "phone_number": "070XXXX",
        "country_code": "KE",
        "first_name": "John",
        "middle_name": "",
        "last_name": "Doe",
        "line_1": "Eaglebots",
        "line_2": "",
        "city": "",
        "state": "",
        "postal_code": "",
        "zip_code": ""
    }
}    
'
