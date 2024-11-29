package controller

import (
	"fmt"
	"net/http"

	"github.com/Kamila3820/hoca-backend/config"
	_accountModel "github.com/Kamila3820/hoca-backend/modules/account/model"
	_accountService "github.com/Kamila3820/hoca-backend/modules/account/service"
	"github.com/Kamila3820/hoca-backend/modules/custom"
	mod "github.com/Kamila3820/hoca-backend/pkg"
	"github.com/Kamila3820/hoca-backend/utils/text"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
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

	fileHeader, err := pctx.FormFile("file")
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	filename := text.GenerateRandomString(12) + "_" + fileHeader.Filename
	imageUrl := "https://" + config.ConfigGetting().Minio.BucketEndpoint + "/" + config.ConfigGetting().Minio.BucketName + "/" + filename

	info, err := mod.Minio.PutObject(
		pctx.Request().Context(),
		"cs-hoca",
		filename,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: fileHeader.Header.Get("Content-Type")},
	)
	if err != nil {
		fmt.Println("Error in Minio.PutObject:", err) // Log the specific error
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	fmt.Println("File uploaded successfully:", info)

	// * Dump uploaded info
	spew.Dump(info)

	registerReq := &_accountModel.RegisterRequest{
		Username:        &username,
		Password:        &password,
		PhoneNumber:     &phoneNumber,
		Email:           &email,
		ConfirmPassword: &confirmPassword,
		IDcard:          &imageUrl,
	}

	register, err := c.accountService.Register(registerReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

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
