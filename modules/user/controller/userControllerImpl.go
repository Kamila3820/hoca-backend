package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Kamila3820/hoca-backend/config"
	"github.com/Kamila3820/hoca-backend/modules/account/misc"
	"github.com/Kamila3820/hoca-backend/modules/custom"
	_userModel "github.com/Kamila3820/hoca-backend/modules/user/model"
	_userService "github.com/Kamila3820/hoca-backend/modules/user/service"
	mod "github.com/Kamila3820/hoca-backend/pkg"
	"github.com/Kamila3820/hoca-backend/utils/text"
	"github.com/davecgh/go-spew/spew"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
)

type userControllerImpl struct {
	userService _userService.UserService
}

func NewUserControllerImpl(userService _userService.UserService) UserController {
	return &userControllerImpl{
		userService: userService,
	}
}

func (c *userControllerImpl) GetUserByID(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	user, err := c.userService.FindUserByID(userID.ID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, user)
}

func (c *userControllerImpl) EditUserProfile(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	userEditingReq := new(_userModel.UserEditingReq)
	contentType := pctx.Request().Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") || strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		userEditingReq.UserName = pctx.FormValue("user_name")
		userEditingReq.PhoneNumber = pctx.FormValue("phone_number")

		fileHeader, err := pctx.FormFile("file")
		if err == http.ErrMissingFile {
			// No file provided, do nothing for Avatar (keep existing)
			fmt.Println("No new file provided, keeping existing avatar")
		} else if err != nil {
			// Some other error occurred
			return custom.Error(pctx, http.StatusBadRequest, fmt.Errorf("failed to retrieve file: %w", err))
		} else {
			// If the file is provided, update the Avatar
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

			// * Dump uploaded info
			spew.Dump(info)
			userEditingReq.Avatar = imageUrl
		}
	} else if strings.HasPrefix(contentType, "application/json") {
		// Handle raw JSON
		customEchoRequest := custom.NewCustomEchoRequest(pctx)
		if err := customEchoRequest.Bind(userEditingReq); err != nil {
			return custom.Error(pctx, http.StatusBadRequest, err)
		}
	}

	userEdit, err := c.userService.EditingUser(userID.ID, userEditingReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, userEdit)
}

func (c *userControllerImpl) GetUserLocation(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	location, err := c.userService.FindLocation(userID.ID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	if location.Location == "" {
		pctx.JSON(http.StatusOK, nil)
	}

	return pctx.JSON(http.StatusOK, location)
}

func (c *userControllerImpl) UpdateUserLocation(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	userLocationReq := new(_userModel.UserLocationReq)
	contentType := pctx.Request().Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") || strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		userLocationReq.Location = pctx.FormValue("location")
		userLocationReq.Latitude = pctx.FormValue("latitude")
		userLocationReq.Longtitude = pctx.FormValue("longtitude")

	} else if strings.HasPrefix(contentType, "application/json") {
		// Handle raw JSON
		customEchoRequest := custom.NewCustomEchoRequest(pctx)
		if err := customEchoRequest.Bind(userLocationReq); err != nil {
			return custom.Error(pctx, http.StatusBadRequest, err)
		}
	}

	userEdit, err := c.userService.EditingLocation(userID.ID, userLocationReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, userEdit)
}
