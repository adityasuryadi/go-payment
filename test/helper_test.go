package test

// func TestGenerateBillNo(t *testing.T) {
// 	var payment entity.Payment
// 	var billNo string
// 	configuration := config.New("../.env")
// 	db := config.NewTestPostgresDB(configuration)
// 	today := time.Now().Format("2006-01-02")
// 	curdate := time.Now().Format("20060102")
// 	result := db.Where(clause.Expr{SQL: "DATE(created_at) = ?", Vars: []interface{}{today}}).
// 		Order("bill_no_counter desc").First(&payment)
// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		billNo = fmt.Sprintf("INV-%s%d", curdate, 1)
// 	} else {
// 		billNo = fmt.Sprintf("INV-%s%d", curdate, 2)
// 	}
// 	fmt.Println(billNo)
// }
