package users

import (
	"chatonetoone/internals/models"
	"errors"

	"gorm.io/gorm"
)

type SqlUserRepository struct{
	db *gorm.DB
}

func NewSqlUserRepository(db *gorm.DB) (*SqlUserRepository, error) {
	if db == nil{
		return nil, errors.New("argument should must contain db")
	}
	return &SqlUserRepository{db: db}, nil
}

func(ur *SqlUserRepository) FindUserByEmail(email string) (*models.User, error){
	var user models.User
	result := ur.db.Where(models.User{Email: email}).First(&user)
	if(result.Error != nil){
		if(errors.Is(result.Error, gorm.ErrRecordNotFound)){
			return nil, errors.New("Record not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func(ur *SqlUserRepository) CreateUser(user *models.User) error{
	return ur.db.Create(user).Error
}