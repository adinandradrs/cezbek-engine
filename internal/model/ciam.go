package model

type (
	CiamOnboardPartnerRequest struct {
		Username    string
		Name        string
		PhoneNumber string
		Email       string
		Picture     string
		Password    string
	}

	CiamAuthenticationRequest struct {
		Username string
		Secret   string
	}
)

type (
	CiamUserResponse struct {
		SubId string
		TransactionResponse
	}

	CiamAuthenticationResponse struct {
		AccessToken  string
		Token        string
		RefreshToken string
		ExpiresIn    int64
	}
)