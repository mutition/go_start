package adapters

import (
	"context"
	"strconv"
	"sync"
	"time"

	domain "github.com/mutition/go_start/order/domain/order"
	"github.com/sirupsen/logrus"
)

type MemoryOrderRepository struct {
	lock  *sync.RWMutex
	store []*domain.Order
}

var fakeData = []*domain.Order{
	{
		ID: "1",
		CustomerID: "1",
		Status: "pending",
		PaymentLink: "https://payment.com",
		Items: nil,
	},
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	return &MemoryOrderRepository{
		lock:  &sync.RWMutex{},
		store: fakeData,
	}
}



//implement Repository interface
func (m *MemoryOrderRepository) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	res := &domain.Order{
		ID:          strconv.FormatInt(time.Now().UnixNano(), 10),
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
	m.store = append(m.store, res)
	logrus.WithFields(logrus.Fields{
		"orderId": res.ID,
		"customerId": res.CustomerID,
		"status": res.Status,
		"paymentLink": res.PaymentLink,
		"items": res.Items,
	}).Debug("Order created")
	return res, nil
}

func (m *MemoryOrderRepository) Get(ctx context.Context, orderId string, customerId string) (*domain.Order, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, o := range m.store {
		if o.ID == orderId && o.CustomerID == customerId {
			logrus.Info("Order found ", "orderId ", orderId, "customerId ", customerId)
			return o, nil
		}
	}
	return nil, domain.NotFoundError{OrderId: orderId}
}

func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(ctx context.Context, order *domain.Order) (*domain.Order, error)) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for i , o := range m.store {
		if o.ID == order.ID && o.CustomerID == order.CustomerID {
			updatedorder, err := updateFn(ctx, o)
			if err != nil {
				return err
			}
			m.store[i] = updatedorder
			logrus.Info("Order updated", "orderId", order.ID, "customerId", order.CustomerID)
			return nil
		}
	}
	return domain.NotFoundError{OrderId: order.ID}
}