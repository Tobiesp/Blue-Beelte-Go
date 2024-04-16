package database

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"
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
	Name    string
	Columns []ColumnSchema
}

func GetStructSchema(s interface{}) (TableSchema, error) {
	var ts TableSchema
	sch, err := schema.Parse(s, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return ts, err
	}

	ts.Name = sch.Table
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

func BuildSelectQuery(s interface{}) (string, error) {
	tableSchema, err := GetStructSchema(s)
	if err != nil {
		return "", err
	}
	formatName := strings.TrimSpace(strings.ToLower(tableSchema.Name))
	if !strings.HasSuffix(formatName, "s") {
		formatName = formatName + "s"
	}
	return "SELECT (.+) FROM \"" + formatName + "\" WHERE (.+)", nil
}

func BuildInsertQuery(s interface{}) (string, error) {
	tableSchema, err := GetStructSchema(s)
	if err != nil {
		return "", err
	}
	formatName := strings.TrimSpace(strings.ToLower(tableSchema.Name))
	if !strings.HasSuffix(formatName, "s") {
		formatName = formatName + "s"
	}
	var strBuilder strings.Builder
	strBuilder.WriteString("INSERT INTO \"")
	strBuilder.WriteString(formatName)
	strBuilder.WriteString("\" (")
	var valuesStr string = ""
	for idx, column := range tableSchema.Columns {
		strBuilder.WriteString("\"")
		strBuilder.WriteString(column.DBColumnName)
		strBuilder.WriteString("\"")
		valuesStr = valuesStr + "$" + strconv.Itoa((idx + 1))
		if idx < (len(tableSchema.Columns) - 1) {
			strBuilder.WriteString(",")
			valuesStr = valuesStr + ","
		}
	}
	strBuilder.WriteString(") VALUES (")
	strBuilder.WriteString(valuesStr)
	strBuilder.WriteString(")")
	return strBuilder.String(), nil
}

func GetSliceOfItem(data any) ([][]driver.Value, error) {
	rt := reflect.TypeOf(data)
	log.Default().Println("Data Type: " + rt.Kind().String())
	if rt.Kind() != reflect.Array && rt.Kind() != reflect.Slice {
		return nil, errors.New("must supply a list of data")
	}

	rv := reflect.ValueOf(data)
	if rv.IsNil() {
		return nil, errors.New("data is nil")
	}
	if rv.Len() == 0 {
		return nil, errors.New("data must  have at least 1 element in the array")
	}
	val := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		val[i] = rv.Index(i).Interface()
	}
	var rows [][]driver.Value
	for _, row := range val {
		rValue := reflect.ValueOf(row)
		rType := rValue.Type()
		if rType.Kind() == reflect.Struct {
			var values []driver.Value
			for i := 0; i < rType.NumField(); i++ {
				fieldValue := rValue.Field(i).Interface()
				srValue := reflect.ValueOf(fieldValue)
				srType := srValue.Type()

				if srType.Kind() == reflect.Struct {
					if srValue.FieldByName("ID").IsValid() {
						values = append(values, srValue.FieldByName("ID").Interface())
					} else {
						values = append(values, fieldValue)
					}
				} else {
					values = append(values, fieldValue)
				}
			}
			log.Println("value: ", values)
			rows = append(rows, [][]driver.Value{values}...)
		}
	}
	return rows, nil
}

func BuildMockDBSelectRows(data any) (*sqlmock.Rows, error) {
	rt := reflect.TypeOf(data)
	log.Default().Println("Data Type: " + rt.Kind().String())
	if rt.Kind() != reflect.Array && rt.Kind() != reflect.Slice {
		return nil, errors.New("must supply a list of data")
	}

	rv := reflect.ValueOf(data)
	if rv.IsNil() {
		return nil, errors.New("data is nil")
	}
	if rv.Len() == 0 {
		return nil, errors.New("data must  have at least 1 element in the array")
	}
	val := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		val[i] = rv.Index(i).Interface()
	}

	rows, err := BuildMockRowsTableHeader(data, false)
	if err != nil {
		return nil, err
	}

	for _, row := range val {
		rValue := reflect.ValueOf(row)
		rType := rValue.Type()
		if rType.Kind() == reflect.Struct {
			var values []driver.Value
			for i := 0; i < rType.NumField(); i++ {
				fieldValue := rValue.Field(i).Interface()
				srValue := reflect.ValueOf(fieldValue)
				srType := srValue.Type()

				if srType.Kind() == reflect.Struct {
					if srValue.FieldByName("ID").IsValid() {
						values = append(values, srValue.FieldByName("ID").Interface())
					} else {
						values = append(values, fieldValue)
					}
				} else {
					values = append(values, fieldValue)
				}
			}
			log.Println("value: ", values)
			rows = rows.AddRow(values...)
		}
	}
	return rows, nil
}

func BuildMockRowsTableHeader(data any, insertQuery bool) (*sqlmock.Rows, error) {
	rt := reflect.TypeOf(data)
	log.Default().Println("Data Type: " + rt.Kind().String())
	if rt.Kind() != reflect.Array && rt.Kind() != reflect.Slice {
		return nil, errors.New("must supply a list of data")
	}

	rv := reflect.ValueOf(data)
	if rv.IsNil() {
		return nil, errors.New("data is nil")
	}
	if rv.Len() == 0 {
		return nil, errors.New("data must  have at least 1 element in the array")
	}
	val := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		val[i] = rv.Index(i).Interface()
	}
	tableSchema, err := GetStructSchema((val)[0])
	if err != nil {
		return nil, err
	}
	var headers []string
	for _, column := range tableSchema.Columns {
		if insertQuery {
			if strings.ToLower(column.DBColumnName) == "id" {
				headers = append(headers, column.DBColumnName)
				break
			}
		} else {
			headers = append(headers, column.DBColumnName)
		}
	}

	log.Println("Headers: ", headers)

	rows := sqlmock.NewRows(headers[:])

	return rows, nil
}

func BuildMockDBInsertRows(data any) (*sqlmock.Rows, error) {
	rt := reflect.TypeOf(data)
	log.Default().Println("Data Type: " + rt.Kind().String())
	if rt.Kind() != reflect.Array && rt.Kind() != reflect.Slice {
		return nil, errors.New("must supply a list of data")
	}

	rv := reflect.ValueOf(data)
	if rv.IsNil() {
		return nil, errors.New("data is nil")
	}
	if rv.Len() == 0 {
		return nil, errors.New("data must  have at least 1 element in the array")
	}
	val := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		val[i] = rv.Index(i).Interface()
	}

	rows, err := BuildMockRowsTableHeader(data, true)
	if err != nil {
		return nil, err
	}

	for _, row := range val {
		rValue := reflect.ValueOf(row)
		rType := rValue.Type()
		if rType.Kind() == reflect.Struct {
			var values []driver.Value
			for i := 0; i < rType.NumField(); i++ {
				if rValue.FieldByName("ID").IsValid() {
					fieldValue := rValue.Field(i).Interface()
					values = append(values, fieldValue)
					break
				}
			}
			log.Println("value: ", values)
			rows = rows.AddRow(values...)
		}
	}
	return rows, nil
}
