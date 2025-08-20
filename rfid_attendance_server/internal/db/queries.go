package db

import (
	"fmt"

	"github.com/ronishg27/rfid_attendance/internal/my_queries"
	"gorm.io/gorm"
)

type Queries struct {
	Public *my_queries.Query
}

// tenant query factory
func (q *Queries) Tenant(db *gorm.DB, schema string) (*my_queries.Query, error) {
	tenantDB := db.Session(&gorm.Session{NewDB: true})
	if err := tenantDB.Exec(fmt.Sprintf("SET search_path TO %s", schema)).Error; err != nil {
		return nil, err
	}
	return my_queries.Use(tenantDB), nil

}

func NewQueries(db *gorm.DB) *Queries {
	return &Queries{
		Public: my_queries.Use(db), // bound to public schema
	}
}
