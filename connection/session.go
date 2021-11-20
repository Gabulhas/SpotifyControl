package connection

type session struct {
	accessToken     string
	clientId        string
	tokenExpiration string
}

func NewSession() session {
	return session{}
}
