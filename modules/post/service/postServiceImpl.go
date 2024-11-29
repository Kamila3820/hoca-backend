package service

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/Kamila3820/hoca-backend/entities"
	_postModel "github.com/Kamila3820/hoca-backend/modules/post/model"
	_postRepository "github.com/Kamila3820/hoca-backend/modules/post/repository"
)

type postServiceImpl struct {
	postRepository _postRepository.PostRepository
}

func NewPostServiceImpl(postRepository _postRepository.PostRepository) PostService {
	return &postServiceImpl{
		postRepository: postRepository,
	}
}

func (s *postServiceImpl) FindPostByDistance(userID string, userLat, userLong float64) ([]*_postModel.Post, error) {
	posts, err := s.postRepository.FindPost()
	if err != nil {
		return nil, err
	}

	const feePerKm = 3.0 // the rate per kilometer for distance fee

	distancePost := make([]*_postModel.Post, 0)
	for _, post := range posts {
		postLat := parseCoordinate(post.LocationLat)
		postLong := parseCoordinate(post.LocationLong)
		distance := calculateDistance(userLat, userLong, postLat, postLong)

		var distanceFee int
		if distance > 4.0 {
			distanceFee = int(feePerKm * distance) // Convert to integer for whole number
		} else {
			distanceFee = 0
		}

		distanceStr := fmt.Sprintf("%.1f", distance)
		newDistance := parseCoordinate(distanceStr)

		if newDistance <= 15.0 && post.ActiveStatus == true && post.OwnerID != userID && post.IsReserved == false {
			post.Distance = distanceStr
			post.DistanceFee = fmt.Sprintf("%d", distanceFee)
			distancePost = append(distancePost, post.ToPostModel())
		}
	}

	return distancePost, nil
}

func parseCoordinate(coordinate string) float64 {
	value, err := strconv.ParseFloat(coordinate, 64)
	if err != nil {
		return 0.0
	}
	return value
}

func calculateDistance(lat1, long1, lat2, long2 float64) float64 {
	const EarthRadiusKm = 6371.0

	lat1Rad := toRadians(lat1)
	long1Rad := toRadians(long1)
	lat2Rad := toRadians(lat2)
	long2Rad := toRadians(long2)

	deltaLat := lat2Rad - lat1Rad
	deltaLong := long2Rad - long1Rad

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(deltaLong/2)*math.Sin(deltaLong/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusKm * c
}

func toRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

func (s *postServiceImpl) FindPostByPostID(postID uint64) (*_postModel.Post, error) {
	post, err := s.postRepository.FindPostByID(postID)
	if err != nil {
		return nil, errors.New("service: failed to get post by post_id")
	}

	return post.ToPostModel(), nil
}

func (s *postServiceImpl) GetPostByUserID(userID string) (*_postModel.Post, error) {
	post, err := s.postRepository.FindPostByUserID(userID)
	if err != nil {
		return nil, errors.New("service: failed to get post by user_id")
	}
	if post == nil {
		return nil, nil
	}

	return post.ToPostModel(), nil
}

func (s *postServiceImpl) CreatingPost(postCreatingReq *_postModel.PostCreatingReq, userID string) (*_postModel.Post, error) {
	// Check if the user already has the post
	existPost, _ := s.postRepository.FindPostByUserID(userID)

	if existPost != nil {
		return nil, errors.New("service: user alredy has the post exist")
	}

	postEntity := &entities.Post{
		OwnerID:        postCreatingReq.OwnerID,
		Name:           postCreatingReq.Name,
		Description:    postCreatingReq.Description,
		Avatar:         postCreatingReq.Avatar,
		Gender:         postCreatingReq.Gender,
		PhoneNumber:    postCreatingReq.PhoneNumber,
		PromptPay:      postCreatingReq.PromptPay,
		Price:          postCreatingReq.Price,
		Location:       postCreatingReq.Location,
		LocationLat:    postCreatingReq.LocationLat,
		LocationLong:   postCreatingReq.LocationLong,
		AmountFamily:   postCreatingReq.AmountFamily,
		Duration:       postCreatingReq.Duration,
		AvailableStart: postCreatingReq.AvailableStart,
		AvailableEnd:   postCreatingReq.AvailableEnd,
	}

	categories, err := s.postRepository.GetCategoriesByIds(postCreatingReq.CategoryIDs)
	if err != nil {
		return nil, err
	}

	postEntity.Categories = categories

	placeTypes, err := s.postRepository.GetPlaceTypesByIds(postCreatingReq.PlaceTypeIDs)
	if err != nil {
		return nil, err
	}

	postEntity.PlaceTypes = placeTypes

	post, err := s.postRepository.CreatingPost(postEntity)
	if err != nil {
		return nil, err
	}

	return post.ToPostModel(), nil
}

func (s *postServiceImpl) EditingPost(postID uint64, postEditingReq *_postModel.PostEditingReq) (*_postModel.Post, error) {
	_, err := s.postRepository.EditingPost(postID, postEditingReq)
	if err != nil {
		return nil, nil
	}
	fmt.Println("111")

	postEntity, err := s.postRepository.FindPostByID(postID)
	if err != nil {
		return nil, err
	}
	fmt.Println("222")

	return postEntity.ToPostModel(), nil
}

func (s *postServiceImpl) DeletePost(postID uint64, userID string) error {
	postEntity, err := s.postRepository.FindPostByID(postID)
	if err != nil {
		return err
	}

	if postEntity.OwnerID != userID {
		return errors.New("forbidden: you don't have permission to delete this post")
	}

	if err := s.postRepository.DeletePost(postID); err != nil {
		return err
	}

	return nil
}

func (s *postServiceImpl) Activate(postID uint64) error {
	if err := s.postRepository.Activate(postID); err != nil {
		return err
	}

	return nil
}

func (s *postServiceImpl) UnActivate(postID uint64) error {
	if err := s.postRepository.UnActivate(postID); err != nil {
		return err
	}

	return nil
}
