package model

type (
	CiamSignUpPartnerRequest struct {
		Username    string
		Name        string
		PhoneNumber string
		Email       string
		Picture     string
		Password    string
	}

	CiamSignInRequest struct {
		Username string
		Password string
	}
)

type (
	CiamUserResponse struct {
		SubId string
		TransactionResponse
	}
)
