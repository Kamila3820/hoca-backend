package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Kamila3820/hoca-backend/modules/custom"
	_postModel "github.com/Kamila3820/hoca-backend/modules/post/model"
	_postService "github.com/Kamila3820/hoca-backend/modules/post/service"
	"github.com/labstack/echo/v4"
)

type postControllerImpl struct {
	postService _postService.PostService
}

func NewPostontrollerImpl(postService _postService.PostService) PostController {
	return &postControllerImpl{
		postService: postService,
	}
}

func (c *postControllerImpl) FindPostByDistance(pctx echo.Context) error {
	userID := pctx.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		return pctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user ID from context",
		})
	}
	fmt.Println(userIDStr)

	userLat, err := strconv.ParseFloat(pctx.QueryParam("lat"), 64)
	if err != nil {
		return pctx.JSON(http.StatusBadRequest, "Invalid latitude")
	}

	userLong, err := strconv.ParseFloat(pctx.QueryParam("long"), 64)
	if err != nil {
		return pctx.JSON(http.StatusBadRequest, "Invalid longitude")
	}

	workerPost, err := c.postService.FindPostByDistance(userLat, userLong)
	if err != nil {
		return pctx.String(http.StatusInternalServerError, err.Error())
	}

	return pctx.JSON(http.StatusOK, workerPost)
}

func (c *postControllerImpl) CreateWorkerPost(pctx echo.Context) error {
	userID := pctx.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		return pctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user ID from context",
		})
	}
	fmt.Println("user_id = ", userIDStr)
	fmt.Println("1")

	postCreatingReq := new(_postModel.PostCreatingReq)

	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	fmt.Println("2")

	if err := customEchoRequest.Bind(postCreatingReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	postCreatingReq.OwnerID = userIDStr
	fmt.Println("3")

	workerPost, err := c.postService.CreatingPost(postCreatingReq)
	if err != nil {
		return pctx.String(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("4")

	return pctx.JSON(http.StatusCreated, workerPost)
}

func (c *postControllerImpl) EditWorkerPost(pctx echo.Context) error {
	userID := pctx.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		return pctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user ID from context",
		})
	}
	fmt.Println("1")

	postID, err := c.getPostID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	fmt.Println("2")

	postEditingReq := new(_postModel.PostEditingReq)

	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(postEditingReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	fmt.Println("3")

	postEdit, err := c.postService.EditingPost(postID, postEditingReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}
	fmt.Println("4")

	// check the the edited post belong to the user or not
	if postEdit.OwnerID != userIDStr {
		return custom.Error(pctx, http.StatusForbidden, err)
	}
	fmt.Println("5")

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
	userID := pctx.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		return pctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user ID from context",
		})
	}

	postID, err := c.getPostID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	if err := c.postService.DeletePost(postID, userIDStr); err != nil {
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
