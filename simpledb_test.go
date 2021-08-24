package simpledb

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type book struct {
	Title       string
	ReleaseYear int
	Genre       string
	Highlights  []string
}

func Test_Open(t *testing.T) {
	db, err := Connect("local_test.json")

	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, db.db.Content)
	assert.Equal(t, "bar", db.db.Content.Get("foo").String())
}

func Test_Save_EmptyData(t *testing.T) {
	os.Remove("testing.json")
	db, err := Connect("testing.json")
	assert.Nil(t, err)

	for _, value := range []interface{}{
		"",
		nil,
		[]string{},
	} {
		err = db.Save(value)
		assert.NotNil(t, err)
		assert.Equal(t, ErrDataMustBeStructPointer, err)
	}
}

func Test_Save(t *testing.T) {
	os.Remove("testing.json")
	db, err := Connect("testing.json")
	assert.Nil(t, err)

	b := book{
		Title:       "walden",
		ReleaseYear: 1854,
	}

	err = db.Save(&b)
	assert.Nil(t, err)

	assert.Equal(t, "walden", db.db.Content.Get("book.0.element.Title").String(), db.db.Content.String())
	assert.Equal(t, int64(1854), db.db.Content.Get("book.0.element.ReleaseYear").Int(), db.db.Content.String())

	err = db.Save(&book{Title: "lotr", Highlights: []string{"something"}})
	assert.Nil(t, err)

	numberOfKeys := len(db.db.Content.Get("book.#.element").Array())
	assert.Equal(t, 2, numberOfKeys)
}

func Test_FetchN(t *testing.T) {
	os.Remove("testing.json")
	db, err := Connect("testing.json")
	assert.Nil(t, err)

	type user struct {
		Name          string
		Age           int
		FavoriteBooks []string
	}

	for _, i := range []int{1, 2, 3, 4, 5} {
		u := user{
			Name:          fmt.Sprintf("henry %d david throreau", i),
			FavoriteBooks: []string{fmt.Sprintf("book %d", i)},
		}
		err = db.Save(&u)
		time.Sleep(time.Microsecond * 100)
		assert.Nil(t, err)
	}

	var users []user
	err = db.FetchN(&users, 3)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(users))
	assert.Equal(t, "henry 1 david throreau", users[0].Name)
	assert.Equal(t, "henry 3 david throreau", users[2].Name)
	assert.Equal(t, "book 2", users[1].FavoriteBooks[0])
}

func Test_Drop(t *testing.T) {
	os.Remove("testing.json")
	db, err := Connect("testing.json")
	assert.Nil(t, err)

	type user struct {
		Name string
	}

	u := &user{
		Name: "something",
	}
	err = db.Save(u)
	assert.Nil(t, err)

	type book struct {
		Title string
	}
	b := &book{
		Title: "lotr",
	}
	err = db.Save(b)
	assert.Nil(t, err)

	err = db.Drop(&user{})
	assert.Nil(t, err)
	assert.Equal(t, "", db.db.Content.Get("user").Raw)
	assert.NotEmpty(t, db.db.Content.Get("book"))
}

func Test_FindOne(t *testing.T) {
	os.Remove("testing.json")
	db, err := Connect("testing.json")
	assert.Nil(t, err)

	notPointer := []book{}
	err = db.Find(notPointer, "Name", "Something")
	assert.Equal(t, err, ErrDataMustBeSlicePointer)

	for _, i := range []int{1, 2, 3, 4, 5} {
		u := book{
			Title:       fmt.Sprintf("harry potter %d", i),
			ReleaseYear: i,
		}
		err = db.Save(&u)
		time.Sleep(time.Microsecond * 100)
		assert.Nil(t, err)
	}

	var b []book
	err = db.Find(&b, "Title", "harry potter 3")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(b))
	assert.Equal(t, "harry potter 3", b[0].Title)

	var anotherBook []book
	err = db.Find(&anotherBook, "ReleaseYear", 5)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(anotherBook))
	assert.Equal(t, 5, anotherBook[0].ReleaseYear)

	err = db.Find(&anotherBook, "Title", "harry potter 10")
	assert.Equal(t, err, ErrNotFound)
}

func Test_FindWhere(t *testing.T) {
	os.Remove("testing.json")
	db, err := Connect("testing.json")
	assert.Nil(t, err)

	notPointer := []book{}
	err = db.FindWhere(notPointer, Where{"Name": "Something"})
	assert.Equal(t, err, ErrDataMustBeSlicePointer)

	for _, i := range []int{1, 2, 3, 4, 5} {
		u := book{
			Title:       fmt.Sprintf("harry potter %d", i),
			ReleaseYear: i * 10,
			Genre:       "fiction",
		}
		err = db.Save(&u)
		time.Sleep(time.Microsecond * 100)
		assert.Nil(t, err)
	}

	var b []book
	err = db.FindWhere(&b, Where{"Title": "harry potter 3", "ReleaseYear": 30})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(b))
	assert.Equal(t, "harry potter 3", b[0].Title)
	assert.Equal(t, 30, b[0].ReleaseYear)

	err = db.FindWhere(&b, Where{"Title": "harry potter 10"})
	assert.Equal(t, err, ErrNotFound)

	var allBooks []book
	err = db.FindWhere(&allBooks, Where{"Genre": "fiction"})
	assert.Equal(t, 5, len(allBooks))
}
