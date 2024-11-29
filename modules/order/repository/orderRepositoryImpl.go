package repository

import (
	"errors"
	"time"

	"github.com/Kamila3820/hoca-backend/entities"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

func (r *orderRepositoryImpl) FindLastWeekOrderByPostID(postID uint64) ([]*entities.Order, error) {
	startDate := time.Now().AddDate(0, 0, -7) // 7 days before now
	endDate := time.Now()                     // current time

	var orders []*entities.Order
	if err := r.db.Connect().
		Where("worker_post_id = ?", postID).
		Where("payment_type = ?", "cash").
		Where("paid = ?", false).
		Where("order_status = ?", "complete").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Find(&orders).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Infof("No active order found for post %s", postID)
			return nil, nil
		}

		r.logger.Errorf("Failed to find last week's cash orders: %s", err.Error())
		return nil, err
	}

	return orders, nil
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

func (r *orderRepositoryImpl) FindActiveOrder(userID string) (*entities.Order, error) {
	order := new(entities.Order)

	if err := r.db.Connect().
		Where("user_id = ?", userID).
		Where("order_status IN (?)", []string{"confirmation", "preparing", "working"}).
		Where("order_status NOT IN (?)", []string{"cancelled", "complete"}).
		Preload("User").
		Preload("Post").
		First(&order).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Infof("No active order found for user %s", userID)
			return nil, nil
		}
		r.logger.Errorf("Failed to find active order for user %s: %s", userID, err.Error())
		return nil, err
	}

	return order, nil
}

func (r *orderRepositoryImpl) FindWorkerOrder(userID string) (*entities.Order, error) {
	order := new(entities.Order)

	if err := r.db.Connect().
		Joins("JOIN posts ON posts.id = orders.worker_post_id").
		Where("posts.owner_id = ?", userID).
		Where("order_status IN (?)", []string{"confirmation", "preparing", "working"}).
		Where("order_status NOT IN (?)", []string{"cancelled", "complete"}).
		Preload("User").
		Preload("Post").
		First(&order).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Infof("No active order found for user %s", userID)
			return nil, nil
		}
		r.logger.Errorf("Failed to find active order for user %s: %s", userID, err.Error())
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

func (r *orderRepositoryImpl) CreatingWorkerFeePayment(orderPaymentEntity *entities.WorkerFeepayment) error {
	if err := r.db.Connect().Create(orderPaymentEntity).Error; err != nil {
		r.logger.Errorf("Failed to create worker QR entity: %s", err.Error())
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

func (r *orderRepositoryImpl) FindWorkerTransactionByID(transactionID string) (*entities.WorkerFeepayment, error) {
	paymentOrder := new(entities.WorkerFeepayment)

	if err := r.db.Connect().First(paymentOrder, "transaction_id = ?", transactionID).Error; err != nil {
		r.logger.Errorf("Failed to query worker QR entity: %s", err.Error())
		return nil, err
	}

	return paymentOrder, nil
}

func (r *orderRepositoryImpl) UpdateTransactionOrder(paymentOrderEntity *entities.OrderQrpayment) error {
	if err := r.db.Connect().Where("order_id = ?", paymentOrderEntity.OrderID).Updates(paymentOrderEntity).Error; err != nil {
		r.logger.Errorf("Failed to update payment order entity: %s", err.Error())
		return err
	}
	return nil
}

func (r *orderRepositoryImpl) UpdateTransactionFee(paymentOrderEntity *entities.WorkerFeepayment) error {
	if err := r.db.Connect().Where("transaction_id = ?", paymentOrderEntity.TransactionID).Updates(paymentOrderEntity).Error; err != nil {
		r.logger.Errorf("Failed to update payment fee entity: %s", err.Error())
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
