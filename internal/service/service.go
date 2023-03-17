package service

type CashURL interface {
	WriteURLInCash(string2 string) (string, error)
	ReadURLFromCash(string string) (string, error)
}

type Service struct {
	CashURL
}

func NewService(cash CashURL) *Service {
	return &Service{
		CashURL: NewCashURLService(cash),
	}
}
