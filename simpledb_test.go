package simpledb

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Open(t *testing.T) {
	db, err := Open("local_test.json")

	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, db.db.Content)
	assert.Equal(t, "bar", db.db.Content.Get("foo").String())
}

func Test_Save_EmptyData(t *testing.T) {
	os.Remove("empty.json")
	db, err := Open("empty.json")
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
	os.Remove("empty.json")
	db, err := Open("empty.json")
	assert.Nil(t, err)

	type user struct {
		Name          string
		Age           int
		FavoriteBooks []string
	}

	u := user{
		Name: "henry david throreau",
		Age:  44,
	}

	err = db.Save(&u)
	assert.Nil(t, err)

	assert.Equal(t, "henry david throreau", db.db.Content.Get("user.0.element.Name").String(), db.db.Content.String())
	assert.Equal(t, int64(44), db.db.Content.Get("user.0.element.Age").Int(), db.db.Content.String())

	err = db.Save(&user{Name: "j r r tolkien", FavoriteBooks: []string{"the hobbit"}})
	assert.Nil(t, err)

	numberOfKeys := len(db.db.Content.Get("user.#.element").Array())
	assert.Equal(t, 2, numberOfKeys)
}

func Test_FetchList(t *testing.T) {
	os.Remove("empty.json")
	db, err := Open("empty.json")
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
	err = db.FetchList(&users, 3)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(users))
	assert.Equal(t, "henry 1 david throreau", users[0].Name)
	assert.Equal(t, "henry 3 david throreau", users[2].Name)
	assert.Equal(t, "book 2", users[1].FavoriteBooks[0])
}

func Test_Drop(t *testing.T) {
	os.Remove("empty.json")
	db, err := Open("empty.json")
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
	os.Remove("empty.json")
	db, err := Open("empty.json")
	assert.Nil(t, err)

	type user struct {
		Name string
		Age  int
	}

	notPointer := user{}
	err = db.FindOne(notPointer, "Name", "Something")
	assert.Equal(t, err, ErrDataMustBeStructPointer)

	for _, i := range []int{1, 2, 3, 4, 5} {
		u := user{
			Name: fmt.Sprintf("harry potter %d", i),
			Age:  i,
		}
		err = db.Save(&u)
		time.Sleep(time.Microsecond * 100)
		assert.Nil(t, err)
	}

	var uu user
	err = db.FindOne(&uu, "Name", "harry potter 3")
	assert.Nil(t, err)
	assert.Equal(t, "harry potter 3", uu.Name)

	var auu user
	err = db.FindOne(&auu, "Age", 5)
	assert.Nil(t, err)
	assert.Equal(t, 5, auu.Age)

	err = db.FindOne(&uu, "Name", "harry potter 10")
	assert.Equal(t, err, ErrNotFound)
}
