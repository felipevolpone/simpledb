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
			Name: fmt.Sprintf("henry %d david throreau", i),
			FavoriteBooks: []string{fmt.Sprintf("book %d", i)},
		}
		err = db.Save(&u)
		time.Sleep(time.Microsecond * 100)
		assert.Nil(t, err)
	}

	var users []user
	err = db.FetchList(&users, 10)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))
	assert.Equal(t, "henry 1 david throreau", users[0].Name)
	assert.Equal(t, "henry 5 david throreau", users[4].Name)
	assert.Equal(t, "book 2", users[1].FavoriteBooks[0])
}
