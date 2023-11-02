package mysql

import (
	"database/sql/driver"
	"encoding/json"
)

// Array parse array field
type Array[T any] []T

// Scan .
func (j *Array[T]) Scan(value any) error {
	*j = make([]T, 0)
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	if len(bytes) == 0 {
		return nil
	}
	result := make([]T, 0)
	err := json.Unmarshal(bytes, &result)
	if err == nil {
		*j = result
		return nil
	}
	return nil
}

// Value .
func (j Array[T]) Value() (driver.Value, error) {
	if j == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(j)
}

// Json parse json field
type Json[T any] struct {
	Val *T `json:"json"`
}

// MarshalJSON .
func (j Json[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Value)
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *Json[T]) Scan(value any) error {
	//初始赋值，防止为nil
	*j = Json[T]{}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	//情况1：空字符串
	if len(bytes) == 0 {
		return nil
	}

	//情况2：正常的json
	result := Json[T]{}
	err := json.Unmarshal(bytes, &result.Val)
	if err == nil {
		*j = result
		return nil
	}

	//情况3："[]"
	var arr []any
	err = json.Unmarshal(bytes, &arr)
	if err == nil {
		return nil
	}

	return nil
}

// Value return json value, implement driver.Valuer interface
func (j Json[T]) Value() (driver.Value, error) {
	if j.Val == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(j.Val)
}
