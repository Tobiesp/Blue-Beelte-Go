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

func BuildMockDBRows(data any) (*sqlmock.Rows, error) {
	rt := reflect.TypeOf(data)
	if rt.Kind() != reflect.Array {
		return nil, errors.New("must supply a list of data")
	}
	val, ok := data.([]interface{})
	if !ok {
		return nil, errors.New("data must be an array of interfaces")
	}
	if len(val) == 0 {
		return nil, errors.New("array can not be empty")
	}
	tableSchema, err := GetStructSchema((val)[0])
	if err != nil {
		return nil, err
	}
	var headers []string
	for _, column := range tableSchema.Columns {
		headers = append(headers, column.DBColumnName)
	}
	rows := sqlmock.NewRows(headers)
	for _, row := range val {
		rValue := reflect.ValueOf(row)
		rType := rValue.Type()
		if rType.Kind() == reflect.Struct {
			var values []driver.Value
			for i := 0; i < rType.NumField(); i++ {
				fieldValue := rValue.Field(i)
				srValue := reflect.ValueOf(fieldValue)
				srType := srValue.Type()
				if srType.Kind() == reflect.Struct {
					values = append(values, srValue.FieldByName("ID"))
				} else {
					values = append(values, fieldValue)
				}
			}
			rows.AddRow(values)
		}
	}
	return rows, nil
}
