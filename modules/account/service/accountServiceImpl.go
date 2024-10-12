package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Kamila3820/hoca-backend/entities"
	"github.com/Kamila3820/hoca-backend/modules/account/misc"
	_accountModel "github.com/Kamila3820/hoca-backend/modules/account/model"
	_accountRepository "github.com/Kamila3820/hoca-backend/modules/account/repository"
	"github.com/Kamila3820/hoca-backend/utils/crypto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type accountServiceImpl struct {
	accountRepository _accountRepository.AccountRepository
}

func NewAccountServiceImpl(accountRepository _accountRepository.AccountRepository) AccountService {
	return &accountServiceImpl{
		accountRepository: accountRepository,
	}
}

func (s *accountServiceImpl) Register(registerReq *_accountModel.RegisterRequest) (*_accountModel.RegisterResponse, error) {
	if err := s.accountRepository.CheckDuplicateEmail(*registerReq.Email); err != nil {
		return nil, errors.New("service: cannot use the duplicate email")
	}

	if !strings.Contains(*registerReq.Email, "@") {
		return nil, errors.New("service: invalid email address")
	}
	if len(*registerReq.Username) < 4 {
		return nil, errors.New("service: required 4 character or more")
	}
	if len(*registerReq.Password) < 4 {
		return nil, errors.New("service: password required 8 character or more")
	}
	if *registerReq.Password != *registerReq.ConfirmPassword {
		return nil, errors.New("service: the password confirmation does not match")
	}

	hashedPassword, err := crypto.HashPassword(*registerReq.Password)
	if err != nil {
		return nil, errors.New("service: unable to hash password")
	}

	fmt.Println(*registerReq.IDcard)
	key := "mysecretencryptionkey123456789012"
	encryptedIDCard, err := crypto.EncryptString(*registerReq.IDcard, key)
	if err != nil {
		return nil, errors.New("service: unable to encrypt id card")
	}

	user := &entities.User{
		ID:           uuid.New().String(),
		UserName:     *registerReq.Username,
		PhoneNumber:  *registerReq.PhoneNumber,
		Email:        *registerReq.Email,
		Password:     hashedPassword,
		IDCard:       encryptedIDCard,
		VerifyStatus: true,
	}

	_, err = s.accountRepository.Creating(user)
	if err != nil {
		return nil, errors.New("service: unable to hash password")
	}

	return &_accountModel.RegisterResponse{
		UserID: &user.ID,
	}, nil
}

func (s *accountServiceImpl) Login(loginReq *_accountModel.LoginRequest) (*_accountModel.LoginResponse, error) {
	user, err := s.accountRepository.FindUserByEmail(*loginReq.Email)
	if err != nil {
		return nil, errors.New("service: unable to find user by email")
	}

	if !crypto.ComparePassword(user.Password, *loginReq.Password) {
		return nil, errors.New("service: incorrect password")
	}

	claim := &misc.UserClaim{
		ID:    user.ID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // Set expiration
		},
	}

	token, err := crypto.SignJwt(claim)
	if err != nil {
		return nil, errors.New("service: unable to sign token")
	}

	loginRes := &_accountModel.LoginResponse{
		Token: &token,
	}

	return loginRes, nil
}
