package dao

import (
	"time"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
)

type OrderDAO struct {
	UUID        uuid.UUID             `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code        string                `gorm:"type:varchar;not null;default:'';uniqueIndex:orders_code_idx"`
	Metadata    string                `gorm:"type:varchar;not null;default:''"`
	UserUUID    uuid.UUID             `gorm:"type:uuid;not null"`
	TotalAmount float64               `gorm:"type:float8;not null;default:0.0"`
	ExpiredAt   time.Time             `gorm:"type:timestamptz;not null"`
	Status      model.OrderStatusType `gorm:"type:order_status_type;not null;default:'created'"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
	UpdatedBy uuid.UUID `gorm:"type:uuid;not null"`

	OrderDetails []OrderDetailDAO `gorm:"foreignKey:OrderUUID;references:UUID"`
}

func (OrderDAO) TableName() string {
	return "orders"
}

type OrderDetailDAO struct {
	UUID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductUUID  uuid.UUID `gorm:"type:uuid;not null"`
	ProductTitle string    `gorm:"type:varchar;not null;default:''"`
	ProductPrice float64   `gorm:"type:float8;not null;default:0.0"`
	Quantity     int       `gorm:"type:int;not null;default:0"`
	SubTotal     float64   `gorm:"type:float8;not null;default:0.0"`
	OrderUUID    uuid.UUID `gorm:"type:uuid;not null;index:orders_detail_order_uuid_idx"`

	Order *OrderDAO `gorm:"foreignKey:OrderUUID;references:UUID"`
}

func (OrderDetailDAO) TableName() string {
	return "orders_detail"
}
