package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Kamila3820/hoca-backend/config"
	"github.com/Kamila3820/hoca-backend/modules/account/misc"
	"github.com/Kamila3820/hoca-backend/modules/custom"
	_postModel "github.com/Kamila3820/hoca-backend/modules/post/model"
	_postService "github.com/Kamila3820/hoca-backend/modules/post/service"
	mod "github.com/Kamila3820/hoca-backend/pkg"
	"github.com/Kamila3820/hoca-backend/utils/text"
	"github.com/davecgh/go-spew/spew"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
)

type postControllerImpl struct {
	postService _postService.PostService
}

func NewPostControllerImpl(postService _postService.PostService) PostController {
	return &postControllerImpl{
		postService: postService,
	}
}

func (c *postControllerImpl) FindPostByDistance(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)
	fmt.Println(userID.ID)
	fmt.Println("UserID")

	userLat, err := strconv.ParseFloat(pctx.QueryParam("lat"), 64)
	if err != nil {
		return pctx.JSON(http.StatusBadRequest, "Invalid latitude")
	}

	userLong, err := strconv.ParseFloat(pctx.QueryParam("long"), 64)
	if err != nil {
		return pctx.JSON(http.StatusBadRequest, "Invalid longitude")
	}

	workerPost, err := c.postService.FindPostByDistance(userID.ID, userLat, userLong)
	if err != nil {
		return pctx.String(http.StatusInternalServerError, err.Error())
	}

	return pctx.JSON(http.StatusOK, workerPost)
}

func (c *postControllerImpl) GetOwnPost(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	post, err := c.postService.GetPostByUserID(userID.ID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	if post == nil {
		return pctx.JSON(http.StatusOK, nil)
	}

	return pctx.JSON(http.StatusOK, post)
}

func (c *postControllerImpl) GetPostByPostID(pctx echo.Context) error {
	postID, err := c.getPostID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	workerPost, err := c.postService.FindPostByPostID(postID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, workerPost)
}

func parseUint64Array(input string) ([]uint64, error) {
	var ids []uint64
	if input == "" {
		return ids, nil
	}

	parts := strings.Split(input, ",")
	for _, part := range parts {
		id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (c *postControllerImpl) CreateWorkerPost(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	name := pctx.FormValue("name")
	description := pctx.FormValue("description")

	// Convert placetype to []uint64
	placetypeStr := pctx.FormValue("placetype_ids")
	placetypeIDs, err := parseUint64Array(placetypeStr)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, fmt.Errorf("invalid placetype_ids"))
	}

	categoryStr := pctx.FormValue("category_ids")
	categoryIDs, err := parseUint64Array(categoryStr)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, fmt.Errorf("invalid placetype_ids"))
	}

	location := pctx.FormValue("location")
	latitude := pctx.FormValue("latitude")
	longtitude := pctx.FormValue("longtitude")

	// Convert price to float64
	priceStr := pctx.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, fmt.Errorf("invalid price format"))
	}
	promptPay := pctx.FormValue("prompt_pay")
	phoneNumber := pctx.FormValue("phone_number")
	gender := pctx.FormValue("gender")
	amountFamily := pctx.FormValue("amount_family")
	duration := pctx.FormValue("duration")
	availableStart := pctx.FormValue("available_start")
	availableEnd := pctx.FormValue("available_end")

	// Handle file upload
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

	// Create post request
	postCreatingReq := &_postModel.PostCreatingReq{
		OwnerID:        userID.ID,
		Avatar:         imageUrl,
		Name:           name,
		Description:    description,
		CategoryIDs:    categoryIDs,
		PlaceTypeIDs:   placetypeIDs,
		Location:       location,
		LocationLat:    latitude,
		LocationLong:   longtitude,
		Price:          price,
		PromptPay:      promptPay,
		PhoneNumber:    phoneNumber,
		Gender:         gender,
		AmountFamily:   amountFamily,
		Duration:       duration,
		AvailableStart: availableStart,
		AvailableEnd:   availableEnd,
	}

	workerPost, err := c.postService.CreatingPost(postCreatingReq, userID.ID)
	if err != nil {
		return pctx.String(http.StatusInternalServerError, err.Error())
	}

	return pctx.JSON(http.StatusCreated, workerPost)
}

func (c *postControllerImpl) EditWorkerPost(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	postID, err := c.getPostID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	postEditingReq := new(_postModel.PostEditingReq)
	contentType := pctx.Request().Header.Get("Content-Type")

	// Check if the content type is form-data or raw JSON
	if strings.HasPrefix(contentType, "multipart/form-data") || strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		// Handle form-data
		postEditingReq.Name = pctx.FormValue("name")
		postEditingReq.Description = pctx.FormValue("description")
		postEditingReq.Location = pctx.FormValue("location")
		postEditingReq.LocationLat = pctx.FormValue("latitude")
		postEditingReq.LocationLong = pctx.FormValue("longtitude")
		postEditingReq.Price, _ = strconv.ParseFloat(pctx.FormValue("price"), 64)
		postEditingReq.PromptPay = pctx.FormValue("prompt_pay")
		postEditingReq.PhoneNumber = pctx.FormValue("phone_number")
		postEditingReq.Gender = pctx.FormValue("gender")
		postEditingReq.AmountFamily = pctx.FormValue("amount_family")
		postEditingReq.Duration = pctx.FormValue("duration")
		postEditingReq.AvailableStart = pctx.FormValue("available_start")
		postEditingReq.AvailableEnd = pctx.FormValue("available_end")

		if categoryIDsStr := pctx.FormValue("category_ids"); categoryIDsStr != "" {
			categoryIDsStrArr := strings.Split(categoryIDsStr, ",")
			var categoryIDs []uint64
			for _, idStr := range categoryIDsStrArr {
				id, _ := strconv.ParseUint(idStr, 10, 64)
				categoryIDs = append(categoryIDs, id)
			}
			postEditingReq.CategoryIDs = categoryIDs
		}

		if placeTypeIDsStr := pctx.FormValue("placetype_ids"); placeTypeIDsStr != "" {
			placeTypeIDsStrArr := strings.Split(placeTypeIDsStr, ",")
			var placeTypeIDs []uint64
			for _, idStr := range placeTypeIDsStrArr {
				id, _ := strconv.ParseUint(idStr, 10, 64)
				placeTypeIDs = append(placeTypeIDs, id)
			}
			postEditingReq.PlaceTypeIDs = placeTypeIDs
		}

		// Handle file upload only if the file is provided
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
			postEditingReq.Avatar = imageUrl
		}

	} else if strings.HasPrefix(contentType, "application/json") {
		// Handle raw JSON
		customEchoRequest := custom.NewCustomEchoRequest(pctx)
		if err := customEchoRequest.Bind(postEditingReq); err != nil {
			return custom.Error(pctx, http.StatusBadRequest, err)
		}
	} else {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	// Now call the service to edit the post
	postEdit, err := c.postService.EditingPost(postID, postEditingReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	// Check if the edited post belongs to the user
	if postEdit.OwnerID != userID.ID {
		return custom.Error(pctx, http.StatusForbidden, err)
	}

	return pctx.JSON(http.StatusOK, postEdit)
}

func (c *postControllerImpl) getPostID(pctx echo.Context) (uint64, error) {
	postIDStr := pctx.Param("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		return 0, nil
	}

	return postID, nil
}

func (c *postControllerImpl) DeleteWorkerPost(pctx echo.Context) error {
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	postID, err := c.getPostID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	if err := c.postService.DeletePost(postID, userID.ID); err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, map[string]string{
		"message": "Post deleted successfully",
	})
}

func (c *postControllerImpl) ActivatePost(pctx echo.Context) error {
	postID, err := c.getPostID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	if err := c.postService.Activate(postID); err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, map[string]string{
		"message": "Activate post successfully",
	})
}

func (c *postControllerImpl) UnActivatePost(pctx echo.Context) error {
	postID, err := c.getPostID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	if err := c.postService.UnActivate(postID); err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, map[string]string{
		"message": "UnActivate post successfully",
	})
}
