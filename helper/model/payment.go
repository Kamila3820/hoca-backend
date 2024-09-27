package model

type ScbGetTokenRequest struct {
	ApplicationKey    string `json:"applicationKey"`
	ApplicationSecret string `json:"applicationSecret"`
}

type ScbStatusResponse struct {
	Code        int    `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

type ScbGetTokenResponse struct {
	Status ScbStatusResponse       `json:"status,omitempty"`
	Data   ScbGetTokenDataResponse `json:"data,omitempty"`
}

type ScbGetTokenDataResponse struct {
	AccessToken string `json:"accessToken,omitempty"`
	ExpiresIn   int    `json:"expiresIn,omitempty"`
	TokenType   string `json:"tokenType,omitempty"`
	ExpiresAt   int    `json:"expiresAt,omitempty"`
}

type ScbCreateQrPaymentRequest struct {
	QrType string `json:"qrType"`
	Amount string `json:"amount"`
	PpType string `json:"ppType"`
	PpId   string `json:"ppId"`
	Ref1   string `json:"ref1"`
	Ref2   string `json:"ref2"`
	Ref3   string `json:"ref3"`
}

type ScbCreateQrResponse struct {
	Status ScbStatusResponse       `json:"status,omitempty"`
	Data   ScbCreateQrDataResponse `json:"data,omitempty"`
}

type ScbCreateQrDataResponse struct {
	QrRawData string `json:"qrRawData,omitempty"`
	QrImage   string `json:"qrImage,omitempty"`
}

// Inquiry QR payment

type PaymentInquiryRequest struct {
	TransactionId string `json:"transactionId,omitempty" validate:"required"`
}

// SCB Inquiry QR payment response

type ScbInquiryQrResponse struct {
	Status ScbStatusResponse          `json:"status,omitempty"`
	Data   []ScbInquiryQrDataResponse `json:"data,omitempty"`
}

type ScbInquiryQrDataResponse struct {
	EventCode            string  `json:"eventCode,omitempty"`
	TransactionType      string  `json:"transactionType,omitempty"`
	ReverseFlag          string  `json:"reverseFlag,omitempty"`
	PayeeProxyId         string  `json:"payeeProxyId,omitempty"`
	PayeeProxyType       string  `json:"payeeProxyType,omitempty"`
	PayeeAccountNumber   string  `json:"payeeAccountNumber,omitempty"`
	PayeeName            *string `json:"payeeName,omitempty"`
	PayerProxyId         string  `json:"payerProxyId,omitempty"`
	PayerProxyType       string  `json:"payerProxyType,omitempty"`
	PayerAccountNumber   string  `json:"payerAccountNumber,omitempty"`
	PayerName            string  `json:"payerName,omitempty"`
	SendingBankCode      string  `json:"sendingBankCode,omitempty"`
	ReceivingBankCode    string  `json:"receivingBankCode,omitempty"`
	Amount               string  `json:"amount,omitempty"`
	TransactionId        string  `json:"transactionId,omitempty"`
	FastEasySlipNumber   string  `json:"fastEasySlipNumber,omitempty"`
	TransactionDate      string  `json:"transactionDate,omitempty"`
	BillPaymentRef1      string  `json:"billPaymentRef1,omitempty"`
	BillPaymentRef2      string  `json:"billPaymentRef2,omitempty"`
	BillPaymentRef3      string  `json:"billPaymentRef3,omitempty"`
	CurrencyCode         string  `json:"currencyCode,omitempty"`
	EquivalentAmount     string  `json:"equivalentAmount,omitempty"`
	ExchangeRate         string  `json:"exchangeRate,omitempty"`
	ChannelCode          string  `json:"channelCode,omitempty"`
	PartnerTransactionId string  `json:"partnerTransactionId,omitempty"`
	TepaCode             string  `json:"tepaCode,omitempty"`
}

// Inquiry payment response

type PaymentInquiryResponse struct {
	PaymentSuccess bool   `json:"paymentSuccess"`
	Message        string `json:"message,omitempty"`
}
