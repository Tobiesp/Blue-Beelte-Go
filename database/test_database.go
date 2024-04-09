package database

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		t.Fatal(err)
	}
	return sqldb, gormdb, mock
}

type ColumnSchema struct {
	DBColumnName  string
	DBColumnIndex int
	FieldName     string
	FieldIndex    int
}

type TableSchema struct {
	Columns []ColumnSchema
}

func GetStructSchema(s interface{}) (TableSchema, error) {
	var ts TableSchema
	sch, err := schema.Parse(s, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return ts, err
	}

	for idx, field := range sch.Fields {
		var cs ColumnSchema
		cs.DBColumnIndex = idx
		cs.DBColumnName = field.DBName
		cs.FieldIndex = idx
		cs.FieldName = field.Name
		ts.Columns = append(ts.Columns, cs)
	}
	return ts, nil
}

func BuildMockDBRows(data []interface{}) (*sqlmock.Rows, error) {
	if len(data) == 0 {
		return nil, errors.New("must supply a list of data")
	}
	tableSchema, err := GetStructSchema(data[0])
	if err != nil {
		return nil, err
	}
	var headers []string
	for _, column := range tableSchema.Columns {
		headers = append(headers, column.DBColumnName)
	}
	rows := sqlmock.NewRows(headers)
	for _, row := range data {
		rValue := reflect.ValueOf(row)
		rType := rValue.Type()
		if rType.Kind() == reflect.Struct {
			var values []driver.Value
			for i := 0; i < rType.NumField(); i++ {
				fieldValue := rValue.Field(i)
				values = append(values, fieldValue)
				//TODO: Add in logic to get ID for struct in struct
			}
			rows.AddRow(values)
		}
	}
	return rows, nil
}
