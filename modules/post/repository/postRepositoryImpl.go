package repository

import (
	"fmt"

	"github.com/Kamila3820/hoca-backend/entities"
	_postModel "github.com/Kamila3820/hoca-backend/modules/post/model"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
	"github.com/labstack/echo/v4"
)

type postRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewPostRepositoryImpl(db databases.Database, logger echo.Logger) PostRepository {
	return &postRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *postRepositoryImpl) FindPost() ([]*entities.Post, error) {
	postList := make([]*entities.Post, 0)

	// Preload PlaceTypes or other associations
	if err := r.db.Connect().Preload("PlaceTypes").Preload("UserRatings.User").Find(&postList).Error; err != nil {
		r.logger.Errorf("Failed to find posts: %s", err.Error())
		return nil, err
	}

	return postList, nil
}

func (r *postRepositoryImpl) FindPostByID(postID uint64) (*entities.Post, error) {
	post := new(entities.Post)

	if err := r.db.Connect().Preload("PlaceTypes").Preload("UserRatings.User").First(post, postID).Error; err != nil {
		r.logger.Errorf("Failed to find post by ID: %s", err.Error())
		return nil, err
	}

	return post, nil
}

func (r *postRepositoryImpl) FindPostByUserID(userID string) (*entities.Post, error) {
	post := new(entities.Post)

	if err := r.db.Connect().Where("owner_id = ?", userID).Preload("PlaceTypes").Preload("UserRatings.User").First(post).Error; err != nil {
		r.logger.Errorf("Failed to find post by UserID: %s", err.Error())
		return nil, err
	}

	return post, nil
}

func (r *postRepositoryImpl) CreatingPost(postEntity *entities.Post) (*entities.Post, error) {
	post := new(entities.Post)

	if err := r.db.Connect().Create(postEntity).Scan(post).Error; err != nil {
		r.logger.Errorf("Creating worker post failed: %s", err.Error())
		return nil, err
	}

	for _, placeType := range postEntity.PlaceTypes {
		// Check if the combination already exists
		exists := r.db.Connect().Where("post_id = ? AND place_type_id = ?", postEntity.ID, placeType.ID).First(&entities.PostPlaceType{}).RowsAffected > 0
		if !exists {
			postPlaceType := &entities.PostPlaceType{
				PostID:      postEntity.ID,
				PlaceTypeID: placeType.ID,
			}
			if err := r.db.Connect().Create(postPlaceType).Error; err != nil {
				r.logger.Errorf("Creating place type for post failed: %s", err.Error())
				return nil, err
			}
		}
	}

	if err := r.db.Connect().Preload("PlaceTypes").First(post, postEntity.ID).Error; err != nil {
		r.logger.Errorf("Fetching post with place types failed: %s", err.Error())
		return nil, err
	}

	fmt.Printf("Post with PlaceTypes: %+v\n", post)

	return post, nil
}

func (r *postRepositoryImpl) GetPlaceTypesByIds(placeTypeIDs []uint64) ([]*entities.PlaceType, error) {
	var placeTypes []*entities.PlaceType
	if err := r.db.Connect().Where("id IN ?", placeTypeIDs).Find(&placeTypes).Error; err != nil {
		r.logger.Errorf("Fetching place types failed: %s", err.Error())
		return nil, err
	}
	return placeTypes, nil
}

func (r *postRepositoryImpl) EditingPost(postID uint64, postEditingReq *_postModel.PostEditingReq) (uint64, error) {
	// Retrieve the existing post
	var post entities.Post
	if err := r.db.Connect().First(&post, postID).Error; err != nil {
		r.logger.Errorf("Fetching post failed: %s", err.Error())
		return 0, err
	}

	// Dynamically build the map for updates
	updates := make(map[string]interface{})

	if postEditingReq.Name != "" {
		updates["name"] = postEditingReq.Name
	}
	if postEditingReq.Description != "" {
		updates["description"] = postEditingReq.Description
	}
	if postEditingReq.Avatar != "" {
		updates["avatar"] = postEditingReq.Avatar
	}
	if postEditingReq.CategoryID != 0 {
		updates["category_id"] = postEditingReq.CategoryID
	}
	if postEditingReq.Location != "" {
		updates["location"] = postEditingReq.Location
	}
	if postEditingReq.LocationLat != "" {
		updates["location_lat"] = postEditingReq.LocationLat
	}
	if postEditingReq.LocationLong != "" {
		updates["location_long"] = postEditingReq.LocationLong
	}
	if postEditingReq.Price != 0 {
		updates["price"] = postEditingReq.Price
	}
	if postEditingReq.PhoneNumber != "" {
		updates["phone_number"] = postEditingReq.PhoneNumber
	}
	if postEditingReq.Gender != "" {
		updates["gender"] = postEditingReq.Gender
	}
	if postEditingReq.AmountFamily != "" {
		updates["amount_family"] = postEditingReq.AmountFamily
	}

	// Update the post with only the selected fields
	if len(updates) > 0 {
		if err := r.db.Connect().Model(&post).Updates(updates).Error; err != nil {
			r.logger.Errorf("Editing worker post failed: %s", err.Error())
			return 0, err
		}
	}

	// Handle PlaceTypeIDs if provided
	if len(postEditingReq.PlaceTypeIDs) > 0 {
		// Delete existing place types
		if err := r.db.Connect().Model(&post).Association("PlaceTypes").Clear(); err != nil {
			r.logger.Errorf("Clearing existing place types failed: %s", err.Error())
			return 0, err
		}

		// Retrieve new place types
		var placeTypes []entities.PlaceType
		if err := r.db.Connect().Where("id IN ?", postEditingReq.PlaceTypeIDs).Find(&placeTypes).Error; err != nil {
			r.logger.Errorf("Finding new place types failed: %s", err.Error())
			return 0, err
		}

		// Update place types association
		if err := r.db.Connect().Model(&post).Association("PlaceTypes").Replace(placeTypes); err != nil {
			r.logger.Errorf("Updating place types failed: %s", err.Error())
			return 0, err
		}
	}

	return postID, nil
}

func (r *postRepositoryImpl) DeletePost(postID uint64) error {
	// Delete all related user ratings
	if err := r.db.Connect().Where("worker_post_id = ?", postID).Delete(&entities.UserRating{}).Error; err != nil {
		r.logger.Errorf("Failed to delete related user ratings: %s", err.Error())
		return err
	}

	if err := r.db.Connect().Delete(&entities.Post{}, postID).Error; err != nil {
		r.logger.Errorf("Delete worker post failed: %s", err.Error())
		return err
	}

	return nil
}

func (r *postRepositoryImpl) Activate(postID uint64) error {
	if err := r.db.Connect().Table("posts").Where("id = ?", postID).Update("active_status", true).Error; err != nil {
		r.logger.Errorf("Activate worker post failed: %s", err.Error())
		return err
	}

	return nil
}

func (r *postRepositoryImpl) UnActivate(postID uint64) error {
	if err := r.db.Connect().Table("posts").Where("id = ?", postID).Update("active_status", false).Error; err != nil {
		r.logger.Errorf("UnActivate worker post failed: %s", err.Error())
		return err
	}

	return nil
}
