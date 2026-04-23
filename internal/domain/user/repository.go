package user

type UserRepository interface {
	GetByID(id uint) (*User, error)
	SaveUser(user *User) error
	GetUserVarieties(userID uint) ([]*UserVariety, error)
	SaveUserVariety(userVariety *UserVariety) error
}
