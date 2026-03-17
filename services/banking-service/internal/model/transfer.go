package model

type Transfer struct {
	TransferID    uint `gorm:"primaryKey"`
	TransactionID uint `gorm:"not null"`
	ExchangeRate  float64
	Commission    float64

	Transaction Transaction
}
