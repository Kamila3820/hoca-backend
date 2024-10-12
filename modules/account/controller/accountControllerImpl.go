package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	_accountModel "github.com/Kamila3820/hoca-backend/modules/account/model"
	_accountService "github.com/Kamila3820/hoca-backend/modules/account/service"
	"github.com/Kamila3820/hoca-backend/modules/custom"
	"github.com/labstack/echo/v4"
)

type accountControllerImpl struct {
	accountService _accountService.AccountService
}

func NewAccountControllerImpl(accountService _accountService.AccountService) AccountController {
	return &accountControllerImpl{
		accountService: accountService,
	}
}

func (c *accountControllerImpl) Register(pctx echo.Context) error {
	username := pctx.FormValue("username")
	password := pctx.FormValue("password")
	email := pctx.FormValue("email")
	phoneNumber := pctx.FormValue("phone_number")
	confirmPassword := pctx.FormValue("confirmPassword")

	var fileType, fileName string

	file, err := pctx.FormFile("file")
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	} else {
		src, err := file.Open()
		if err != nil {
			return custom.Error(pctx, http.StatusBadRequest, err)
		} else {
			fileByte, _ := ioutil.ReadAll(src)
			fileType = http.DetectContentType(fileByte)

			if fileType == "application/pdf" {
				fileName = "uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + ".pdf"
			} else {
				fileName = "uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + ".jpg"
			}

			err = ioutil.WriteFile(fileName, fileByte, 0777)
			if err != nil {
				return custom.Error(pctx, http.StatusBadRequest, err)
			}
		}
		defer src.Close()
	}

	registerReq := &_accountModel.RegisterRequest{
		Username:        &username,
		Password:        &password,
		PhoneNumber:     &phoneNumber,
		Email:           &email,
		ConfirmPassword: &confirmPassword,
		IDcard:          &fileName,
	}

	register, err := c.accountService.Register(registerReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}
	fmt.Println("3")

	return pctx.JSON(http.StatusOK, register)
}

func (c *accountControllerImpl) Login(pctx echo.Context) error {
	loginReq := new(_accountModel.LoginRequest)

	customEchoRequest := custom.NewCustomEchoRequest(pctx)

	if err := customEchoRequest.Bind(loginReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	token, err := c.accountService.Login(loginReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	pctx.Response().Header().Set("x-auth-token", "Bearer "+*token.Token)
	return pctx.JSON(http.StatusOK, token)
}
