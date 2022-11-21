package model

type ErrorPage struct {
	ErrorMessage string
}

type LoginPage struct {
	*ErrorPage
	LoginURI       string
	GoogleClientID string
}

type SignupPage struct {
	*ErrorPage
	LoginURI       string
	GoogleClientID string
}
