package test

import (
	"payment/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestFindPaymentByBillNo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("error opening GORM database: %v", err)
	}
	rows := sqlmock.NewRows([]string{"id", "name", "phone", "email", "booking_date", "redirect_url", "qty", "bill_no", "bill_total", "status_id", "bill_no_counter", "trx_id", "payment_channel_uid", "payment_channel", "siganture", "created_at", "updated_at"}).
		AddRow(uuid.MustParse("33a61a3d-88e8-484d-8061-3db0bff92e3a"), "name 1", "081234567", "adit@mail.com", "2023-01-01", "http://antrique.com", 1, "INV-2023080811", 2000, 0, 11, "235476234572", 402, "Permata Bank", "286587236578235", time.Now(), time.Now())
	query := `SELECT * FROM "payments" WHERE bill_no = $1 ORDER BY "payments"."id" LIMIT 1`
	mock.ExpectQuery(query).WillReturnRows(rows)
	repo := repository.NewPaymentRepository(gormDB)
	payment, err := repo.FindPaymentByBillNo("INV-2023080811")
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, payment)
}
