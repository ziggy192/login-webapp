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

type ProfileViewPage struct {
	Profile
	*ErrorPage
}

type ProfileEditPage struct {
	Profile
	*ErrorPage
}
