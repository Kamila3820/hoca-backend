package repository

import (
	"github.com/Kamila3820/hoca-backend/entities"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
	"github.com/labstack/echo"

	_userException "github.com/Kamila3820/hoca-backend/modules/user/exception"
)

type userRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewUserRepositoryImpl(db databases.Database, logger echo.Logger) UserRepository {
	return &userRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *userRepositoryImpl) Creating(userEntity *entities.User) (*entities.User, error) {
	user := new(entities.User)

	if err := r.db.Connect().Create(userEntity).Scan(user).Error; err != nil {
		r.logger.Errorf("Creating user failed: %s", err.Error())
		return nil, &_userException.UserCreating{UserID: userEntity.ID}
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByID(userID string) (*entities.User, error) {
	user := new(entities.User)

	if err := r.db.Connect().Where("id = ?", userID).First(user).Error; err != nil {
		r.logger.Errorf("Find user by ID failed: %s", err.Error())
		return nil, &_userException.UserNotFound{UserID: userID}
	}

	return user, nil
}
