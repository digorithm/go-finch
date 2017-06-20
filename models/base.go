package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"reflect"
)

type InsertResult struct {
	lastInsertId int64
	rowsAffected int64
}

func (ir *InsertResult) LastInsertId() (int64, error) {
	return ir.lastInsertId, nil
}

func (ir *InsertResult) RowsAffected() (int64, error) {
	return ir.rowsAffected, nil
}

type Base struct {
	db    *sqlx.DB
	table string
	hasID bool
}

func (b *Base) rowGetter(tx *sqlx.Tx, rowStruct []interface{}, queryBody string, row_id int64) (interface{}, error){

	//var rowTableResults []interface{}

	fmt.Println("2")
	arr := reflect.ValueOf(rowStruct)
	fmt.Println("3")
	v := reflect.New(reflect.TypeOf(rowStruct))
	rows, err := b.db.Queryx(queryBody, row_id)

	fmt.Println(rows)
	if err != nil {
		//fmt.Println("3")
		fmt.Printf("%v", err)
	}

	for rows.Next() {
		fmt.Println("4")
		err = rows.StructScan(v.Interface())
		fmt.Println("5")
		if err != nil {
			fmt.Println("err is not nil")
			fmt.Printf("%v", err)
		}
		arr.Set(reflect.Append(arr, v.Elem()))
		/*fmt.Println("this is rowTableResults before append: ")
		fmt.Println(rowTableResults)
		rowTableResults = append(rowTableResults, rowStruct)
		fmt.Println("this is rowTableResults after append: ")
		fmt.Println(rowTableResults)*/

	}

	return arr, err

}

func (b *Base) newTransactionIfNeeded(tx *sqlx.Tx) (*sqlx.Tx, bool, error) {
	var err error
	wrapInSingleTransaction := false

	if tx != nil {
		return tx, wrapInSingleTransaction, nil
	}

	tx, err = b.db.Beginx()
	if err == nil {
		wrapInSingleTransaction = true
	}

	if err != nil {
		return nil, wrapInSingleTransaction, err
	}

	return tx, wrapInSingleTransaction, nil
}

func (b *Base) InsertIntoTable(tx *sqlx.Tx, data map[string]interface{}) (sql.Result, error) {
	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	dollarMarks := make([]string, 0)
	values := make([]interface{}, 0)

	loopCounter := 1
	for key, value := range data {
		keys = append(keys, key)
		dollarMarks = append(dollarMarks, fmt.Sprintf("$%v", loopCounter))
		values = append(values, value)

		loopCounter++
	}

	query := fmt.Sprintf(
		"INSERT INTO %v (%v) VALUES (%v)",
		b.table,
		strings.Join(keys, ","),
		strings.Join(dollarMarks, ","))

	result := &InsertResult{}
	result.rowsAffected = 1

	if b.hasID {
		query = query + " RETURNING id"

		var lastInsertId int64
		err = tx.QueryRow(query, values...).Scan(&lastInsertId)
		if err != nil {
			return nil, err
		}

		result.lastInsertId = lastInsertId
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	return result, err
}

func (b *Base) UpdateFromTable(tx *sqlx.Tx, data map[string]interface{}, where string) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	keysWithDollarMarks := make([]string, 0)
	values := make([]interface{}, 0)

	loopCounter := 1
	for key, value := range data {
		keysWithDollarMark := fmt.Sprintf("%v=$%v", key, loopCounter)
		keysWithDollarMarks = append(keysWithDollarMarks, keysWithDollarMark)
		values = append(values, value)

		loopCounter++
	}

	query := fmt.Sprintf(
		"UPDATE %v SET %v WHERE %v",
		b.table,
		strings.Join(keysWithDollarMarks, ","),
		where)

	result, err = tx.Exec(query, values...)

	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	return result, err
}

func (b *Base) UpdateByID(tx *sqlx.Tx, data map[string]interface{}, id int64) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	keysWithDollarMarks := make([]string, 0)
	values := make([]interface{}, 0)

	loopCounter := 1
	for key, value := range data {
		keysWithDollarMark := fmt.Sprintf("%v=$%v", key, loopCounter)
		keysWithDollarMarks = append(keysWithDollarMarks, keysWithDollarMark)
		values = append(values, value)

		loopCounter++
	}

	// Add id as part of values
	values = append(values, id)

	query := fmt.Sprintf(
		"UPDATE %v SET %v WHERE id=$%v",
		b.table,
		strings.Join(keysWithDollarMarks, ","),
		loopCounter)

	result, err = tx.Exec(query, values...)

	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	return result, err
}

func (b *Base) UpdateByKeyValueString(tx *sqlx.Tx, data map[string]interface{}, key, value string) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	keysWithDollarMarks := make([]string, 0)
	values := make([]interface{}, 0)

	loopCounter := 1
	for key, value := range data {
		keysWithDollarMark := fmt.Sprintf("%v=$%v", key, loopCounter)
		keysWithDollarMarks = append(keysWithDollarMarks, keysWithDollarMark)
		values = append(values, value)

		loopCounter++
	}

	// Add value as part of values
	values = append(values, value)

	query := fmt.Sprintf(
		"UPDATE %v SET %v WHERE %v=$%v",
		b.table,
		strings.Join(keysWithDollarMarks, ","),
		key,
		loopCounter)

	result, err = tx.Exec(query, values...)

	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	return result, err
}

func (b *Base) DeleteFromTable(tx *sqlx.Tx, where string) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("DELETE FROM %v", b.table)

	if where != "" {
		query = query + " WHERE " + where
	}

	result, err = tx.Exec(query)

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	if err != nil {
		return nil, err
	}

	return result, err
}

func (b *Base) DeleteById(tx *sqlx.Tx, id int64) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("DELETE FROM %v WHERE id=$1", b.table)

	result, err = tx.Exec(query, id)

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	if err != nil {
		return nil, err
	}

	return result, err
}
