package auth

type authTools interface {
	GetReviewer(userName string) (reviewerName string, err error)
	IsDba(userName string) (isDba bool, err error)
}

type MockAuth struct {
}

var MockAuthService MockAuth

func (MockAuth) GetReviewer(userName string) (reviewerName string, err error) {
	return "nobody", nil
}

func (MockAuth) IsDba(userName string) (isDba bool, err error) {
	return true, nil
}
