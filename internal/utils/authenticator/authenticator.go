package authenticator

type Authenticator interface {
	ValidateToken(token string) (int64, error)
	GenerateToken(userId int64) (string, error)
}