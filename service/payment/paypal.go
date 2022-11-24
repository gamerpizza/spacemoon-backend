package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"moonspace/model"
	"moonspace/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/plutov/paypal"
)

type (
	PaypalAuthResponse struct {
		Scope       string `json:"scope"`
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		AppId       string `json:"app_id"`
		Expiry      int64  `json:"expires_in"`
		Nonce       string `json:"nonce"`
	}

	PaypalClientTokenResp struct {
		ClientToken string `json:"client_token"`
	}

	paypalHATEOAS struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	}

	PaypalCreateOrderResp struct {
		ID     string          `json:"id"`
		Status string          `json:"status"`
		Links  []paypalHATEOAS `json:"links"`
	}

	itemTotal struct {
		CurrencyCode string `json:"currency_code"`
		Value        string `json:"value"`
	}

	breakdown struct {
		ItemTotal itemTotal `json:"item_total"`
	}

	amount struct {
		CurrencyCode string    `json:"currency_code"`
		Value        string    `json:"value"`
		Breakdown    breakdown `json:"breakdown"`
	}

	item struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		UnitAmount  itemTotal `json:"unit_amount"`
		Quantity    string    `json:"quantity"`
	}

	units struct {
		Amount amount `json:"amount"`
		Items  []item `json:"items"`
	}

	paymentBody struct {
		Intent intent  `json:"intent"`
		Units  []units `json:"purchase_units"`
	}

	intent string

	PayPal struct {
		clientId     string
		clientSecret string
		authUrl      string
		baseUrl      string
	}
)

const (
	retryTime      = time.Second * 15
	auth           = "/v1/oauth2/token"
	clientIdentity = "/v1/identity/generate-token"
	createOrder    = "/v2/checkout/orders"
	checkoutOrder  = "/v2/checkout/orders/%s/capture"

	intentCapture = "CAPTURE"
	intentAuth    = "AUTH"
)

func NewPayPalClient(cfg PaypalCfg) Payment {
	return &PayPal{
		clientId:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		authUrl:      cfg.URL + auth,
		baseUrl:      cfg.URL,
	}
}

func (p *PayPal) InitPayment(pr model.PaymentRequest, currency string) (any, error) {
	price := float64(0)
	items := make([]item, 0)

	for _, p := range pr.Order.Cart.Products {
		price += p.Price
		it := item{
			Name:        p.Name,
			Description: p.Description,
			UnitAmount: itemTotal{
				CurrencyCode: currency,
				Value:        utils.FloatToString(0, p.Price),
			},
			Quantity: strconv.FormatInt(1, 10), // Quantities should be implemented
		}
		items = append(items, it)
	}

	paymentInfo := paymentBody{
		Intent: intentCapture,
		Units: []units{
			{
				Amount: amount{
					CurrencyCode: currency,
					Value:        utils.FloatToString(0, price),
					Breakdown: breakdown{
						ItemTotal: itemTotal{
							CurrencyCode: currency,
							Value:        utils.FloatToString(0, price),
						},
					},
				},
				Items: items,
			},
		},
	}

	token, err := p.generateClientToken()
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(paymentInfo)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, p.baseUrl+createOrder, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	ppcor := PaypalCreateOrderResp{}
	if err := p.sendRequestWithAuth(req, token, &ppcor); err != nil {
		return nil, err
	}

	return ppcor, nil
}

func (p *PayPal) CaptureOrder(id string) (any, error) {
	url := fmt.Sprintf(p.baseUrl+checkoutOrder, id)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	token, err := p.generateClientToken()
	if err != nil {
		return nil, err
	}

	resp := paypal.PaymentSource{}
	if err := p.sendRequestWithAuth(req, token, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *PayPal) generateClientToken() (string, error) {
	resp, err := p.tryAuth()
	if err != nil {
		return "", err
	}

	return resp.AccessToken, nil
}

func (p *PayPal) sendRequestWithAuth(req *http.Request, token string, result any) error {
	req.Header.Add(utils.TokenHeader, "Bearer "+token)
	return p.sendAPIRequest(req, result)
}

func (p *PayPal) sendAPIRequest(req *http.Request, result any) error {
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		err = utils.DecodeRequestBody(resp.Body, result)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PayPal) tryAuth() (PaypalAuthResponse, error) {
	resp := PaypalAuthResponse{}
	err := utils.BasicAuthorization(p.clientId, p.clientSecret, p.authUrl, true, &resp)
	if err != nil {
		fmt.Printf("error authorizing with paypal: %v", err)
		return PaypalAuthResponse{}, err
	}

	return resp, nil
}
