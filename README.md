# simpledb

> local JSON database from small and pet projects

simpledb is a local database based on a json file. I wrote this mostly to help myself
in the development of pet projects. This obviouslly should not be used in production.

```golang
db, err := simpledb.Connect("local.json")

type User struct {
    Name          string
    Age           int
}

u := User{
    Name: "someone",
    Age:  28,
}

err = db.Save(&u)
```

## Install
```bash
go get github.com/felipevolpone/simpledb
```

## Features

- Minimalist and easy to use and learn
- Extensible
- Lightweight

## API

#### Save
```golang
db, err := simpledb.Connect("local.json")
u := User{Name: "someone", Age:  28}
err = db.Save(&u)
```

#### FindOne
Get one record based on a very simple query.

```golang
db, err := simpledb.Connect("local.json")
var u user
err := db.FindOne(&uu, "Name", "someone")
```

#### FindWhere
Gets a list of records record based on the simpledb.Where param.

```golang
var users []User
err = db.FindWhere(&b, Where{"Name": "someone", "Age": 99})
```

#### FetchN
Just gets the last N records.

```golang
var users []User
err = db.FetchN(&users, 5)
```

#### Delete
to do
```golang
```

#### Drop
Delete all records of `user` type.

```golang
db, err := simpledb.Connect("local.json")
err = db.Drop(&user{})
```
