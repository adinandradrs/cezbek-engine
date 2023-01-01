package model

import "io"

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

	S3UploadRequest struct {
		ContentType string
		Source      io.Reader
		Destination string
	}

	SendEmailRequest struct {
		Destination string
		Subject     string
		Content     string
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
