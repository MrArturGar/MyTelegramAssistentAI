package services

type IService interface {
	Send(string) (string, err error)
}
