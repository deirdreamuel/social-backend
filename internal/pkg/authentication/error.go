package authentication

type AuthenticationError struct {
	Code   int
	Reason string
}
