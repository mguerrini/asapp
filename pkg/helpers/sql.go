package helpers

import (
	"database/sql"
	"time"
)

//type SqlRow map[string]interface[]

func GetMapFromReader2(rows *sql.Rows) (map[string]interface{}, error) {
	cols, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	vals := make([]interface{}, len(cols))

	for i, _ := range cols {
		vals[i] = new(sql.RawBytes)
	}

	rows.Scan(vals...)

	rowMap := make(map[string]interface{}, len(cols))
	for i, val := range cols {
		rowMap[val] = vals[i]
	}

	return rowMap, nil
}

func GetMapFromReader(rows *sql.Rows) (map[string]interface{}, error) {
	cols, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i, _ := range columns {
		columnPointers[i] = &columns[i]
	}

	if err := rows.Scan(columnPointers...); err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		m[colName] = *val
	}

	return m, nil
}

func GetInt(key string, fields map[string]interface{}) int {
	return GetIntOr(key, fields, 0)
}

func GetIntOr(key string, fields map[string]interface{}, nilValue int) int {
	val, ok := fields[key]
	if  !ok {
		panic("Invalid field name " + key)
	}

	if val == nil {
		return nilValue
	}

	iVal := val.(int64)
	return int(iVal)
}

func GetString(key string, fields map[string]interface{}) string {
	return GetStringOr(key, fields, "")
}

func GetStringOr(key string, fields map[string]interface{}, nilValue string) string {
	val, ok := fields[key]
	if  !ok {
		panic("Invalid field name " + key)
	}

	if val == nil {
		return nilValue
	}

	return val.(string)
}


func GetTimeFromString(key string, fields map[string]interface{}) time.Time {
	val, ok := fields[key]
	if  !ok {
		panic("Invalid field name " + key)
	}

	t, err := time.Parse("2006-01-02T15:04:05Z", val.(string))
	if err != nil {
		panic (err)
	}
	return t
}

func GetIntOrNullInt(value *int) interface{} {
	if value == nil {
		return sql.NullInt32{}
	}

	return value
}

func GetStringOrNullString(value *string) interface{} {
	if value == nil {
		return sql.NullString{}
	}

	return value
}