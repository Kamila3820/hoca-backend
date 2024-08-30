package service

import (
	"github.com/Kamila3820/hoca-backend/entities"
	_userModel "github.com/Kamila3820/hoca-backend/modules/user/model"
	_userRepository "github.com/Kamila3820/hoca-backend/modules/user/repository"
)

type googleOAuth2Service struct {
	userRepository _userRepository.UserRepository
}

func NewGoogleOAuth2Service(userRepository _userRepository.UserRepository) OAuth2Service {
	return &googleOAuth2Service{
		userRepository: userRepository,
	}
}

func (s *googleOAuth2Service) UserAccountCreating(userCreatingReq *_userModel.UserCreatingReq) error {
	if !s.IsUserExists(userCreatingReq.ID) {
		userEntity := &entities.User{
			ID:       userCreatingReq.ID,
			UserName: userCreatingReq.UserName,
			Email:    userCreatingReq.Email,
			Avatar:   userCreatingReq.Avatar,
		}

		_, err := s.userRepository.Creating(userEntity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *googleOAuth2Service) IsUserExists(userID string) bool {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return false
	}

	return user != nil
}
