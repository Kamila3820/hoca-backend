package repository

import (
	"github.com/Kamila3820/hoca-backend/entities"
	_userException "github.com/Kamila3820/hoca-backend/modules/user/exception"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
	"github.com/labstack/echo/v4"
)

type accountRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewAccountRepositoryImpl(db databases.Database, logger echo.Logger) AccountRepository {
	return &accountRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *accountRepositoryImpl) CheckDuplicateEmail(email string) error {
	checkDuplicateUser := new(entities.User)

	if result := r.db.Connect().First(&checkDuplicateUser, "email = ?", email); result.Error == nil {
		r.logger.Errorf("Duplicated email: %s", result.Error)
		return result.Error
	}

	return nil
}

func (r *accountRepositoryImpl) Creating(userEntity *entities.User) (*entities.User, error) {
	user := new(entities.User)

	if err := r.db.Connect().Create(userEntity).Scan(user).Error; err != nil {
		r.logger.Errorf("Creating user failed: %s", err.Error())
		return nil, &_userException.UserCreating{UserID: userEntity.ID}
	}

	return user, nil
}

func (r *accountRepositoryImpl) FindUserByEmail(email string) (*entities.User, error) {
	user := new(entities.User)

	if err := r.db.Connect().Where("email = ?", email).First(&user).Error; err != nil {
		r.logger.Errorf("Query user by email failed: %s", err.Error())
		return nil, err
	}

	return user, nil
}
