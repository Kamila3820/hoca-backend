package repository

import (
	"github.com/Kamila3820/hoca-backend/entities"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
	"github.com/labstack/echo/v4"
)

type orderRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewOrderRepositoryImpl(db databases.Database, logger echo.Logger) OrderRepository {
	return &orderRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *orderRepositoryImpl) CreatingOrder(orderEntity *entities.Order) (*entities.Order, error) {
	order := new(entities.Order)

	if err := r.db.Connect().Create(orderEntity).Scan(order).Error; err != nil {
		r.logger.Errorf("Failed to create new order entity: %s", err.Error())
		return nil, err
	}

	return order, nil
}

func (r *orderRepositoryImpl) FindPostByID(postID uint64) (*entities.Post, error) {
	post := new(entities.Post)

	if err := r.db.Connect().Preload("PlaceTypes").First(post, postID).Error; err != nil {
		r.logger.Errorf("Failed to find post entity by ID: %s", err.Error())
		return nil, err
	}

	return post, nil
}

func (r *orderRepositoryImpl) FindUserByID(userID string) (*entities.User, error) {
	user := new(entities.User)
	if err := r.db.Connect().Where("id = ?", userID).First(user).Error; err != nil {
		r.logger.Errorf("Failed to find user entity by ID: %s", err.Error())
		return nil, err
	}
	return user, nil
}

func (r *orderRepositoryImpl) FindOrderByID(orderID uint64) (*entities.Order, error) {
	order := new(entities.Order)

	if err := r.db.Connect().Where("id = ?", orderID).Preload("User").Preload("Post").First(&order).Error; err != nil {
		r.logger.Errorf("Failed to find order entity by ID: %s", err.Error())
		return nil, err
	}

	return order, nil
}

func (r *orderRepositoryImpl) UpdateOrder(orderEntity *entities.Order) error {
	if err := r.db.Connect().Save(orderEntity).Error; err != nil {
		r.logger.Errorf("Failed to update order entity: %s", err.Error())
		return err
	}
	return nil
}

func (r *orderRepositoryImpl) CreatingQRpayment(orderPaymentEntity *entities.OrderQrpayment) error {
	if err := r.db.Connect().Create(orderPaymentEntity).Error; err != nil {
		r.logger.Errorf("Failed to create order QR entity: %s", err.Error())
		return err
	}

	return nil
}

func (r *orderRepositoryImpl) FindTransactionByID(transactionID string) (*entities.OrderQrpayment, error) {
	paymentOrder := new(entities.OrderQrpayment)

	if err := r.db.Connect().First(paymentOrder, "transaction_id = ?", transactionID).Error; err != nil {
		r.logger.Errorf("Failed to query order QR entity: %s", err.Error())
		return nil, err
	}

	return paymentOrder, nil
}

func (r *orderRepositoryImpl) UpdateTransactionOrder(paymentOrderEntity *entities.OrderQrpayment) error {
	if err := r.db.Connect().Save(paymentOrderEntity).Error; err != nil {
		r.logger.Errorf("Failed to update payment order entity: %s", err.Error())
		return err
	}

	return nil
}

// Post
func (r *orderRepositoryImpl) UpdatePost(postEntity *entities.Post) error {
	if err := r.db.Connect().Save(postEntity).Error; err != nil {
		r.logger.Errorf("Failed to update post entity: %s", err.Error())
		return err
	}

	return nil
}

// History
func (r *orderRepositoryImpl) CreatingHistory(historyEntity *entities.History) (*entities.History, error) {
	history := new(entities.History)

	if err := r.db.Connect().Create(historyEntity).Scan(history).Error; err != nil {
		r.logger.Errorf("Failed to create order history entity: %s", err.Error())
		return nil, err
	}

	return history, nil
}

// Noti
func (r *orderRepositoryImpl) CreateNotification(notiEntityy *entities.Notification) error {
	if err := r.db.Connect().Create(&notiEntityy).Error; err != nil {
		r.logger.Errorf("Failed to create notification entity: %s", err.Error())
		return err
	}

	return nil
}
