package entity

// ErrorURL - ошибка дублирования url при запросе.
type ErrorURL struct {
	URL string
	Err error
}

func (eu *ErrorURL) Error() string {
	return eu.URL
}

func (eu *ErrorURL) Unwrap() error {
	return eu.Err
}

func NewErrorURL(err error, url string) error {
	return &ErrorURL{
		Err: err,
		URL: url,
	}
}
