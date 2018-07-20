package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

func openDB(dbuser, dbpass, dbhost, dbname, socket string, port int) (*sql.DB, error) {
	userpass := fmt.Sprintf("%s:%s", dbuser, dbpass)
	var conn string
	if socket != "" {
		conn = fmt.Sprintf("unix(%s)", socket)
	} else {
		conn = fmt.Sprintf("tcp(%s:%d)", dbhost, port)
	}

	return sql.Open("mysql", fmt.Sprintf("%s@%s/%s", userpass, conn, dbname))
}

func convVal(data []byte) (interface{}, error) {
	var v interface{}
	var err error
	v, err = strconv.ParseInt(string(data), 10, 64)
	if err == nil {
		return v, nil
	}

	v, err = strconv.ParseFloat(string(data), 64)
	if err == nil {
		return v, nil
	}

	v = string(data)

	return v, nil
}

func convKey(data interface{}, pkType string) (interface{}, error) {
	var v interface{}
	var err error

	switch val := data.(type) {
	case int64:
		switch pkType {
		case "int":
			v = data
		case "float":
			v = float64(data.(int64))
		default:
			v = fmt.Sprint(val)
		}
	case float64:
		switch pkType {
		case "int":
			v = int64(data.(float64))
		case "float":
			v = data
		default:
			v = fmt.Sprint(val)
		}
	case string:
		switch pkType {
		case "int":
			v, err = strconv.ParseInt(data.(string), 10, 64)
		case "float":
			v, err = strconv.ParseFloat(data.(string), 64)
		default:
			v = data
		}
	}

	return v, err
}

func createData(pkType string) reflect.Value {
	var data reflect.Value
	switch pkType {
	case "int":
		data = reflect.ValueOf(make(map[int64]interface{}))
	case "float":
		data = reflect.ValueOf(make(map[float64]interface{}))
	default:
		data = reflect.ValueOf(make(map[string]interface{}))
	}
	return data
}

func Marshal(data interface{}, format string) ([]byte, error) {
	var b []byte
	var err error
	switch format {
	case "json":
		b, err = json.Marshal(data)
	case "yaml":
		b, err = yaml.Marshal(data)
	default:
		err = fmt.Errorf("%s is not supproted", format)
	}

	return b, err
}

func main() {
	var app = kingpin.New("db2data", "Database dump to json / yaml")

	var dbuser = app.Flag("dbuser", "Database user").Default("root").String()
	var dbpass = app.Flag("dbpass", "Database password").String()
	var dbhost = app.Flag("dbhost", "Database host").Default("localhost").String()
	var dbport = app.Flag("dbport", "Database port").Default("3306").Int()
	var dbsock = app.Flag("dbsock", "Database socket").String()
	var dbname = app.Flag("dbname", "Database name").Required().String()
	var query = app.Flag("query", "SQL").Required().String()
	var pkey = app.Flag("pkey", "Primary key").String()
	var pkeyType = app.Flag("pkey-type", "Primary key type [int, float, string]").Default("string").Enum("int", "float", "string")
	var outFormat = app.Flag("output", "Output file format [json, yaml]").Default("json").Enum("json", "yaml")

	app.Version("0.1.0")

	kingpin.MustParse(app.Parse(os.Args[1:]))

	db, err := openDB(*dbuser, *dbpass, *dbhost, *dbname, *dbsock, *dbport)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(*query)
	if err != nil {
		log.Fatal(err)
	}
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	colNames := make(map[string]struct{})
	for _, col := range cols {
		colNames[col] = struct{}{}
	}

	values := make([][]byte, len(cols))

	row := make([]interface{}, len(cols))
	for i, _ := range values {
		row[i] = &values[i]
	}

	data := createData(*pkeyType)

	for rows.Next() {
		if err := rows.Scan(row...); err != nil {
			log.Fatal(err)
		}

		r := make(map[string]interface{})
		for i, val := range values {
			v, err := convVal(val)
			if err != nil {
				log.Fatal(err)
			}

			r[cols[i]] = v
		}

		key, err := convKey(r[*pkey], *pkeyType)
		if err != nil {
			log.Fatal(err)
		}

		data.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(r))
	}

	b, err := Marshal(data.Interface(), *outFormat)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewReader(b)
	_, err = io.Copy(os.Stdout, buf)
	if err != nil {
		log.Fatal(err)
	}
}
