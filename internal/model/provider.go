package model

import "github.com/shopspring/decimal"

type (
	LinksajaFundTransferRequest struct {
		Bearer string `json:"bearer,omitempty"`
		Amount string `json:"amount"`
		Msisdn string `json:"msisdn"`
		Notes  string `json:"notes"`
	}

	GopaidTopUpRequest struct {
		Receipient string          `json:"receipient"`
		AddBalance decimal.Decimal `json:"add_balance"`
	}

	MiddletransWalletTransferRequest struct {
		Wallet  string          `json:"wallet"`
		Amount  decimal.Decimal `json:"amount"`
		Account string          `json:"account"`
	}

	XenitWalletTopupRequest struct {
		Wallet      string          `json:"wallet"`
		Beneficiary string          `json:"beneficiary"`
		Amount      decimal.Decimal `json:"amount"`
		RefCode     string          `json:"ref_code"`
	}

	JosvoAccountTransferRequest struct {
		PhoneNo       string          `json:"phone_no"`
		Amount        decimal.Decimal `json:"amount"`
		ClientRefCode string          `json:"client_ref_code"`
	}
)

type (
	LinksajaAuthorizationResponse struct {
		Token          string `json:"token"`
		ThirdPartyName string `json:"thirdPartyName"`
	}

	LinksajaFundTransferResponse struct {
		Amount decimal.Decimal `json:"amount"`
		Msisdn string          `json:"msisdn"`
		Notes  string          `json:"notes"`
	}

	GopaidTopupResponse struct {
		RefCode   string `json:"ref_code"`
		Timestamp int64  `json:"timestamp"`
	}

	MiddletransWalletTransferResponse struct {
		IsSuccess      bool   `json:"isSuccess"`
		StatusCode     string `json:"statusCode"`
		TransactionRef string `json:"transactionRef"`
		Message        string `json:"message"`
	}

	XenitWalletTopupResponse struct {
		TopupRef     string `json:"topup_ref"`
		TopupTime    int64  `json:"topup_time"`
		TopupStatus  string `json:"topup_status"`
		TopupMessage string `json:"topup_message"`
	}

	JosvoAccountTransferResponse struct {
		Code  string `json:"code"`
		Notes string `json:"notes"`
	}
)
