package repository

type UserRepo interface {
	CreateUser() error
	GetUserByID() error
	//DeleteUser()
}
