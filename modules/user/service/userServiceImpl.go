package service

import (
	"errors"

	_userModel "github.com/Kamila3820/hoca-backend/modules/user/model"
	_userRepository "github.com/Kamila3820/hoca-backend/modules/user/repository"
)

type userServiceImpl struct {
	userRepository _userRepository.UserRepository
}

func NewUserServiceImpl(userRepository _userRepository.UserRepository) UserService {
	return &userServiceImpl{
		userRepository: userRepository,
	}
}

func (s *userServiceImpl) FindUserByID(userID string) (*_userModel.ProfileUser, error) {
	user, err := s.userRepository.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("service: unable to query user by id")
	}

	return user.ToAccountUserModel(), nil
}

func (s *userServiceImpl) EditingUser(userID string, userEditingReq *_userModel.UserEditingReq) (*_userModel.ProfileUser, error) {
	_, err := s.userRepository.EditingUser(userID, userEditingReq)
	if err != nil {
		return nil, nil
	}

	userEntity, err := s.userRepository.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	return userEntity.ToAccountUserModel(), nil
}

func (s *userServiceImpl) FindLocation(userID string) (*_userModel.UserLocation, error) {
	location, err := s.userRepository.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("service: unable to query user by id")
	}

	return location.ToUserLocationModel(), nil
}

func (s *userServiceImpl) EditingLocation(userID string, userLocationReq *_userModel.UserLocationReq) (*_userModel.UserLocation, error) {
	_, err := s.userRepository.EditingUserLocation(userID, userLocationReq)
	if err != nil {
		return nil, nil
	}

	userEntity, err := s.userRepository.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	return userEntity.ToUserLocationModel(), nil
}
