package seed

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

var futuresContracts = []model.FuturesContract{
	// Taken from 2023/24 csv file
	// Settlement April 2026 (J26)
	{Ticker: "ZCJ26", Name: "Corn Futures", ContractSize: 5000, ContractUnit: "bushel", SettlementDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "ZSJ26", Name: "Soybean Futures", ContractSize: 5000, ContractUnit: "bushel", SettlementDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "ZLJ26", Name: "Soybean Oil Futures", ContractSize: 60000, ContractUnit: "pound", SettlementDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "ZMJ26", Name: "Soybean Meal Futures", ContractSize: 180000, ContractUnit: "pound", SettlementDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "ZWJ26", Name: "Chicago Wheat Futures", ContractSize: 5000, ContractUnit: "bushel", SettlementDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "LEJ26", Name: "Live Cattle Futures", ContractSize: 40000, ContractUnit: "pound", SettlementDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "KEJ26", Name: "Wheat Futures", ContractSize: 5000, ContractUnit: "bushel", SettlementDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "GFJ26", Name: "Feeder Cattle Futures", ContractSize: 50000, ContractUnit: "pound", SettlementDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "HEJ26", Name: "Lean Hog Futures", ContractSize: 40000, ContractUnit: "pound", SettlementDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "GPJ26", Name: "Pork Cutout Futures", ContractSize: 40000, ContractUnit: "pound", SettlementDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "GNJ26", Name: "Nonfat Dry Milk Futures", ContractSize: 44000, ContractUnit: "pound", SettlementDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "DCJ26", Name: "Class III Milk Futures", ContractSize: 200000, ContractUnit: "pound", SettlementDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "GDKJ26", Name: "Class IV Milk Futures", ContractSize: 200000, ContractUnit: "pound", SettlementDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)},

	// Settlement June 2026 (M26)
	{Ticker: "CLM26", Name: "Crude Oil Futures", ContractSize: 1000, ContractUnit: "barrel", SettlementDate: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "NGM26", Name: "Natural Gas Futures", ContractSize: 10000, ContractUnit: "MMBtu", SettlementDate: time.Date(2026, 6, 27, 0, 0, 0, 0, time.UTC)},
	{Ticker: "RBM26", Name: "Gasoline Futures", ContractSize: 42000, ContractUnit: "gallon", SettlementDate: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "QMM26", Name: "E-Mini Crude Oil Futures", ContractSize: 500, ContractUnit: "barrel", SettlementDate: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "HOM26", Name: "NY Harbor ULSD Futures", ContractSize: 42000, ContractUnit: "gallon", SettlementDate: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "MCLM26", Name: "Micro Crude Oil Futures", ContractSize: 100, ContractUnit: "barrel", SettlementDate: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "HHM26", Name: "Henry Hub Natural Gas Futures", ContractSize: 10000, ContractUnit: "MMBtu", SettlementDate: time.Date(2026, 6, 27, 0, 0, 0, 0, time.UTC)},
	{Ticker: "BKJM26", Name: "Buckeye Jet Fuel Futures", ContractSize: 42000, ContractUnit: "gallon", SettlementDate: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "GCM26", Name: "Gold Futures", ContractSize: 100, ContractUnit: "troy ounce", SettlementDate: time.Date(2026, 6, 27, 0, 0, 0, 0, time.UTC)},
	{Ticker: "SIM26", Name: "Silver Futures", ContractSize: 5000, ContractUnit: "troy ounce", SettlementDate: time.Date(2026, 6, 27, 0, 0, 0, 0, time.UTC)},
	{Ticker: "PLM26", Name: "Platinum Futures", ContractSize: 50, ContractUnit: "troy ounce", SettlementDate: time.Date(2026, 6, 27, 0, 0, 0, 0, time.UTC)},
	{Ticker: "HGM26", Name: "Copper Futures", ContractSize: 25000, ContractUnit: "pound", SettlementDate: time.Date(2026, 6, 27, 0, 0, 0, 0, time.UTC)},
	{Ticker: "ALIM26", Name: "Aluminum Futures", ContractSize: 50000, ContractUnit: "pound", SettlementDate: time.Date(2026, 6, 27, 0, 0, 0, 0, time.UTC)},
	{Ticker: "QCM26", Name: "E-Mini Copper Futures", ContractSize: 12500, ContractUnit: "pound", SettlementDate: time.Date(2026, 6, 27, 0, 0, 0, 0, time.UTC)},
	{Ticker: "GDBM26", Name: "Cash-Settled Butter Futures", ContractSize: 20000, ContractUnit: "pound", SettlementDate: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "GDCM26", Name: "Cash-Settled Cheese Futures", ContractSize: 20000, ContractUnit: "pound", SettlementDate: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "BCM26", Name: "Block Cheese Futures", ContractSize: 2000, ContractUnit: "pound", SettlementDate: time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)},

	// Settlement September 2026 (U26)
	{Ticker: "ZOU26", Name: "Oats Futures", ContractSize: 5000, ContractUnit: "bushel", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "ZRU26", Name: "Rough Rice Futures", ContractSize: 180000, ContractUnit: "pound", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "URMU26", Name: "Urea Futures", ContractSize: 200000, ContractUnit: "pound", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "LBU26", Name: "Lumber Futures", ContractSize: 27500, ContractUnit: "board feet", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "CPMU26", Name: "Copper Mini Futures", ContractSize: 12500, ContractUnit: "pound", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "SIMU26", Name: "Silver Mini Futures", ContractSize: 1000, ContractUnit: "troy ounce", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "PLMU26", Name: "Platinum Mini Futures", ContractSize: 10, ContractUnit: "troy ounce", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "OGU26", Name: "Gold Options Futures", ContractSize: 100, ContractUnit: "troy ounce", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "SOU26", Name: "Silver Options Futures", ContractSize: 5000, ContractUnit: "troy ounce", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "PAOU26", Name: "Palladium Options Futures", ContractSize: 100, ContractUnit: "troy ounce", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "CTU26", Name: "Cotton Futures", ContractSize: 50000, ContractUnit: "pound", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "KCU26", Name: "Coffee Futures", ContractSize: 37500, ContractUnit: "pound", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "SBU26", Name: "Sugar Futures", ContractSize: 112000, ContractUnit: "pound", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "CCU26", Name: "Cocoa Futures", ContractSize: 10, ContractUnit: "metric ton", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "OJU26", Name: "Orange Juice Futures", ContractSize: 15000, ContractUnit: "pound", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},
	{Ticker: "LOU26", Name: "Lumber Options Futures", ContractSize: 1100, ContractUnit: "board feet", SettlementDate: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)},

	// Settlement December 2026 (Z26)
	{Ticker: "LHOZ26", Name: "Lean Hog Options Futures", ContractSize: 40000, ContractUnit: "pound", SettlementDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)},
	{Ticker: "LCOZ26", Name: "Live Cattle Options Futures", ContractSize: 40000, ContractUnit: "pound", SettlementDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)},
	{Ticker: "FCOZ26", Name: "Feeder Cattle Options Futures", ContractSize: 50000, ContractUnit: "pound", SettlementDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)},
	{Ticker: "BOOZ26", Name: "Butter Options Futures", ContractSize: 20000, ContractUnit: "pound", SettlementDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)},
	{Ticker: "CHOZ26", Name: "Cheese Options Futures", ContractSize: 20000, ContractUnit: "pound", SettlementDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)},
	{Ticker: "PBOZ26", Name: "Pork Belly Options Futures", ContractSize: 40000, ContractUnit: "pound", SettlementDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)},
}

func seedFuturesContracts(db *gorm.DB) error {
	for _, fc := range futuresContracts {
		var existing model.FuturesContract
		err := db.Where("ticker = ?", fc.Ticker).First(&existing).Error
		if err == nil {
			continue // Skip if contract with that ticker already exists
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := db.Create(&fc).Error; err != nil {
			return err
		}
	}
	return nil
}

func Run(db *gorm.DB) error {
	if err := seedFuturesContracts(db); err != nil {
		return err
	}
	return nil
}
