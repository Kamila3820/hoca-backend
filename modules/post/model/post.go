package model

import (
	"time"
)

type (
	Post struct {
		ID           uint64       `json:"id"`
		OwnerID      string       `json:"owner_id"`
		Name         string       `json:"name"`
		Description  string       `json:"description"`
		Avatar       string       `json:"avatar"`
		CategoryID   uint64       `json:"category_id"`
		Location     string       `json:"location"`
		LocationLat  string       `json:"latitude"`
		LocationLong string       `json:"longitude"`
		Price        float64      `json:"price"`
		Distance     string       `json:"distance"`
		PhoneNumber  string       `json:"phone_number"`
		Gender       string       `json:"gender"`
		AmountFamily string       `json:"amount_family"`
		TotalScore   float64      `json:"total_score"`
		ActiveStatus bool         `json:"active_status"`
		PlaceTypes   []PlaceType  `json:"place_types"`
		UserRatings  []UserRating `json:"user_ratings"`
		CreatedAt    time.Time    `json:"created_at"`
		UpdatedAt    time.Time    `json:"updated_at"`
	}

	PlaceType struct {
		ID          uint64 `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	UserRating struct {
		ID            uint64 `json:"id"`
		UserID        string `json:"user_id"`
		Username      string `json:"username"`
		Avatar        string `json:"avatar"`
		WorkScore     int    `json:"work_score"`
		SecurityScore int    `json:"security_score"`
		Comment       string `json:"comment"`
		CreatedAt     string `json:"created_at"`
	}

	PostCreatingReq struct {
		OwnerID      string
		Name         string   `json:"name" validate:"required,max=64"`
		Description  string   `json:"description" validate:"omitempty,max=128"`
		Avatar       string   `json:"avatar" validate:"omitempty"`
		CategoryID   uint64   `json:"category_id" validate:"required"`
		PlaceTypeIDs []uint64 `json:"placetype_ids" validate:"required"`
		Location     string   `json:"location" validate:"required"`
		LocationLat  string   `json:"latitude" validate:"required"`
		LocationLong string   `json:"longtitude" validate:"required"`
		Price        float64  `json:"price" validate:"required"`
		PhoneNumber  string   `json:"phone_number" validate:"required"`
		Gender       string   `json:"gender" validate:"required"`
		AmountFamily string   `json:"amount_family" validate:"required"`
	}

	PostEditingReq struct {
		OwnerID      string
		Name         string   `json:"name" validate:"omitempty,max=64"`
		Description  string   `json:"description" validate:"omitempty,max=128"`
		Avatar       string   `json:"avatar" validate:"omitempty"`
		CategoryID   uint64   `json:"category_id" validate:"omitempty"`
		PlaceTypeIDs []uint64 `json:"placetype_ids" validate:"omitempty"`
		Location     string   `json:"location" validate:"omitempty"`
		LocationLat  string   `json:"latitude" validate:"omitempty"`
		LocationLong string   `json:"longtitude" validate:"omitempty"`
		Price        float64  `json:"price" validate:"omitempty"`
		PhoneNumber  string   `json:"phone_number" validate:"omitempty"`
		Gender       string   `json:"gender" validate:"omitempty"`
		AmountFamily string   `json:"amount_family" validate:"omitempty"`
	}
)
