package repository

import (
	"github.com/Kamila3820/hoca-backend/entities"
	_userException "github.com/Kamila3820/hoca-backend/modules/user/exception"
	_userModel "github.com/Kamila3820/hoca-backend/modules/user/model"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
	"github.com/labstack/echo/v4"
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

func (r *userRepositoryImpl) FindUserByID(userID string) (*entities.User, error) {
	user := new(entities.User)

	if err := r.db.Connect().Where("id = ?", userID).First(&user).Error; err != nil {
		r.logger.Errorf("Query user by email failed: %s", err.Error())
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) EditingUser(userID string, userEditingReq *_userModel.UserEditingReq) (string, error) {
	var user entities.User
	if err := r.db.Connect().Where("id = ?", userID).First(&user).Error; err != nil {
		r.logger.Errorf("Fetching user failed: %s", err.Error())
		return "", err
	}

	updates := make(map[string]interface{})

	if userEditingReq.UserName != "" {
		updates["user_name"] = userEditingReq.UserName
	}
	if userEditingReq.PhoneNumber != "" {
		updates["phone_number"] = userEditingReq.PhoneNumber
	}
	if userEditingReq.Avatar != "" {
		updates["avatar"] = userEditingReq.Avatar
	}

	// Update the post with only the selected fields
	if len(updates) > 0 {
		if err := r.db.Connect().Model(&user).Updates(updates).Error; err != nil {
			r.logger.Errorf("Editing user failed: %s", err.Error())
			return "", err
		}
	}

	return userID, nil
}

func (r *userRepositoryImpl) EditingUserLocation(userID string, userLocationReq *_userModel.UserLocationReq) (string, error) {
	var user entities.User
	if err := r.db.Connect().Where("id = ?", userID).First(&user).Error; err != nil {
		r.logger.Errorf("Fetching user failed: %s", err.Error())
		return "", err
	}

	updates := make(map[string]interface{})

	if userLocationReq.Location != "" {
		updates["location"] = userLocationReq.Location
	}
	if userLocationReq.Latitude != "" {
		updates["latitude"] = userLocationReq.Latitude
	}
	if userLocationReq.Longtitude != "" {
		updates["longtitude"] = userLocationReq.Longtitude
	}

	if len(updates) > 0 {
		if err := r.db.Connect().Model(&user).Updates(updates).Error; err != nil {
			r.logger.Errorf("Editing user failed: %s", err.Error())
			return "", err
		}
	}

	return userID, nil
}
