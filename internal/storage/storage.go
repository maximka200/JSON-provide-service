package storage

type Storage interface {
	NewJSON(js string) (id int, err error)
	GetJSON(id int) (js string, err error)
	DeleteJSON(id int) (err error)
}
