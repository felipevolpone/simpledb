package simpledb

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// DB represents a local database
type DB struct {
	Path string

	db *dbData
}

type row struct {
	Element    interface{} `json:"element"`
	InsertedAt int64       `json:"inserted_at"`
	Hash       string      `json:"hash"`
}

type dbData struct {
	Content gjson.Result
}

// Open opens a database and stabilishes a connection
func Open(path string) (*DB, error) {

	dbFile, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer dbFile.Close()

	b, err := ioutil.ReadAll(dbFile)

	res := gjson.Parse(string(b))
	if !res.IsObject() && string(b) != "" {
		return &DB{}, errors.New("database is not a valid json file")
	}

	return &DB{
		Path: path,
		db:   &dbData{Content: res},
	}, nil
}

// Save saves a struct into the database
func (db *DB) Save(data interface{}) error {
	ref := reflect.ValueOf(data)

	if !ref.IsValid() || ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
		return ErrDataMustBeStructPointer
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	structType := reflect.TypeOf(data)
	namespace := structType.Elem().Name()

	dataHash := hash(string(b))

	newRow := row{
		Element:    data,
		InsertedAt: time.Now().Unix(),
		Hash:       dataHash,
	}

	operation := fmt.Sprintf("%s.-1", namespace)

	value, err := sjson.Set(db.db.Content.Raw, operation, newRow)
	if err != nil {
		return err
	}

	db.db.Content.Raw = value

	return db.write()
}

// FetchList returns a list of items, if the number of available items is
// lower then the limit argument its returned anyway
func (db *DB) FetchList(items interface{}, limit int) error {

	valuePtr := reflect.ValueOf(items)
	elem := valuePtr.Elem()

	sliceType := reflect.Indirect(reflect.ValueOf(items)).Type()
	namespace := sliceType.Elem().Name()

	for _, value := range db.db.Content.Get(fmt.Sprintf("%s.#.element", namespace)).Array() {
		i := reflect.New(sliceType.Elem())
		err := json.Unmarshal([]byte(value.String()), i.Interface())
		if err != nil {
			return err
		}

		fmt.Println(i.Elem())
		elem.Set(reflect.Append(elem, i.Elem()))
	}

	return nil
}

func (db *DB) write() error {
	return ioutil.WriteFile(db.Path, []byte(db.db.Content.Raw), 0644)
}

func hash(s interface{}) string {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(s)
	return base64.StdEncoding.EncodeToString(b.Bytes())
}
