package helper

import (
	"ANTRIQUE/payment/entity"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GenerateBillNo(tx *gorm.DB) (billNo string, billNoCounter int) {
	var payment entity.Payment
	today := time.Now().Format("2006-01-02")
	curdate := time.Now().Format("20060102")
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(clause.Expr{SQL: "DATE(created_at) = ?", Vars: []interface{}{today}}).
		Order("bill_no_counter desc").First(&payment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		billNo = fmt.Sprintf("INV-%s%d", curdate, 1)
		billNoCounter = 1
	} else {
		billNoCounter += payment.BillNoCounter + 1
		billNo = fmt.Sprintf("INV-%s%d", curdate, billNoCounter)
	}
	return billNo, billNoCounter
}
