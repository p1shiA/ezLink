package zarinpal

type PaymentRequest struct {
	MerchantID  string            `json:"merchant_id"`
	Amount      int64             `json:"amount"`
	Currency    string            `json:"currency"`
	Description string            `json:"description"`
	CallbackURL string            `json:"callback_url"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type PaymentResponseData struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Authority string `json:"authority"`
	FeeType   string `json:"fee_type"`
	Fee       int64  `json:"fee"`
}

type PaymentResponse struct {
	Data   *PaymentResponseData `json:"data"`
	Errors interface{}          `json:"errors"`
}

type VerifyRequest struct {
	MerchantID string `json:"merchant_id"`
	Amount     int64  `json:"amount"`
	Authority  string `json:"authority"`
}

type VerifyResponseData struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	CardPan  string      `json:"card_pan"`
	CardHash string      `json:"card_hash"`
	RefID    interface{} `json:"ref_id"`
	FeeType  string      `json:"fee_type"`
	Fee      int64       `json:"fee"`
}

type VerifyResponse struct {
	Data   *VerifyResponseData `json:"data"`
	Errors interface{}         `json:"errors"`
}
