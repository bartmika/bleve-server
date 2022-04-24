package controller

type Controller interface {
	Index(filename string, identifier string, data interface{}) error
	Query(filename string, search string) ([]string, error)
	Close()
}
