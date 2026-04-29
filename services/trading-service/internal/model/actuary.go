package model

type Actuary struct {
	ActuaryID   uint    `gorm:"primaryKey"`
	FirstName   string
	LastName    string
	EmployeeID  uint
	ProfitRSD   float64
	IsSupervisor bool
}