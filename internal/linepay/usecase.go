package linepay

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

/*
	{
	  "returnCode": "0000",
	  "returnMessage": "OK",
	  "info": {
	    "orderId": "EXAMPLE_ORDER_20230422_1000001",
	    "transactionId": 2023042201206549310,
	    "payInfo": [
	      {
	        "method": "BALANCE",
	        "amount": 100
	      }
	    ]
	  }
	}
*/
type ConfirmPaymentResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Info          struct {
		OrderID       string `json:"orderId"`
		TransactionID int    `json:"transactionId"`
		PayInfo       []struct {
			Method string `json:"method"`
			Amount int    `json:"amount"`
		} `json:"payInfo"`
	} `json:"info"`
}

/*
	{
	  "returnCode": "0000",
	  "returnMessage": "Success.",
	  "info": {
	    "paymentUrl": {
	      "web": "https://sandbox-web-pay.line.me/web/payment/wait?transactionReserveId=REpEWEttQ0F2RmFnaFFzVndIdjl6Z0lqbGpPemZjOHpNWTFZTmdibUlRNlEzOG50N2VSRmdGU2IxcnVjMHZ1NQ",
	      "app": "line://pay/payment/REpEWEttQ0F2RmFnaFFzVndIdjl6Z0lqbGpPemZjOHpNWTFZTmdibUlRNlEzOG50N2VSRmdGU2IxcnVjMHZ1NQ"
	    },
	    "transactionId": 2023042201206549310,
	    "paymentAccessToken": "056579816895"
	  }
	}
*/
type RequestPaymentResponse struct {
	ReturnCode    string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	Info          struct {
		PaymentURL struct {
			Web string `json:"web"`
			App string `json:"app"`
		} `json:"paymentUrl"`
		TransactionID      int    `json:"transactionId"`
		PaymentAccessToken string `json:"paymentAccessToken"`
	} `json:"info"`
}

func (lp *LinePay) RequestPayment(
	ctx context.Context,
	amount int,
	orderID string,
	productID string,
	productName string,
	productImageURL string,
	productQuantity int,
	productOriginalPrice int,
	productPrice int,
	confirmURL string,
	cancelURL string,
) (RequestPaymentResponse, error) {
	type Product struct {
		ID            string `json:"id,omitempty"`
		Name          string `json:"name"`
		ImageURL      string `json:"imageUrl,omitempty"`
		Quantity      int    `json:"quantity"`
		OriginalPrice int    `json:"originalPrice,omitempty"`
		Price         int    `json:"price"`
	}

	type Package struct {
		ID       string    `json:"id"`
		Amount   int       `json:"amount"`
		Name     string    `json:"name,omitempty"`
		Products []Product `json:"products"`
	}

	type RedirectURLs struct {
		ConfirmURLType string `json:"confirmUrlType"`
		ConfirmURL     string `json:"confirmUrl,omitempty"`
		CancelURL      string `json:"cancelUrl,omitempty"`
	}

	type Request struct {
		Amount       int          `json:"amount"`
		Currency     string       `json:"currency"`
		OrderID      string       `json:"orderId"`
		Packages     []Package    `json:"packages"`
		RedirectURLs RedirectURLs `json:"redirectUrls"`
	}

	paymentRequest := Request{
		Amount:   amount,
		Currency: "TWD",
		OrderID:  orderID,
		Packages: []Package{{
			ID:     "1",
			Amount: amount,
			Products: []Product{{
				ID:       productID,
				Name:     productName,
				ImageURL: productImageURL,
				Quantity: productQuantity,
				Price:    productPrice,
			}},
		}},
		RedirectURLs: RedirectURLs{
			// ConfirmURLType: "NONE",
			ConfirmURLType: "CLIENT",
			ConfirmURL:     confirmURL,
			CancelURL:      cancelURL,
		},
	}

	req := &APIRequest{
		Method:  "POST",
		APIPath: "/v3/payments/request",
		Data:    paymentRequest,
	}

	a, err := lp.requestOnlineAPI(ctx, req, func(b []byte) (any, error) {
		var res RequestPaymentResponse
		_ = json.Unmarshal(b, &res)
		return res, nil
	})
	if err != nil {
		return RequestPaymentResponse{}, errors.Errorf("request online api, err: %+v", err)
	}

	return a.(RequestPaymentResponse), nil
}

func (lp *LinePay) ConfirmPayment(ctx context.Context, requestTransactionID string, amount int) (ConfirmPaymentResponse, error) {
	type Request struct {
		Amount   int    `json:"amount"`
		Currency string `json:"currency"`
	}

	paymentRequest := Request{
		Amount:   amount,
		Currency: "TWD",
	}

	req := &APIRequest{
		Method:  "POST",
		APIPath: "/v3/payments/" + requestTransactionID + "/confirm",
		Data:    paymentRequest,
	}

	a, err := lp.requestOnlineAPI(ctx, req, func(b []byte) (any, error) {
		var res ConfirmPaymentResponse
		_ = json.Unmarshal(b, &res)
		return res, nil
	})
	if err != nil {
		return ConfirmPaymentResponse{}, errors.Errorf("request online api, err: %+v", err)
	}

	return a.(ConfirmPaymentResponse), nil
}
