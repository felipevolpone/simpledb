# simpledb

> local JSON database for small and pet projects

simpledb is a local database based on a json file. I wrote this mostly to help myself
in the development of pet projects. *This obviouslly should not be used in production.*

Also, don't take this lib too seriously.

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

This project is mostly an interface and some hacks around the great https://github.com/tidwall/gjson project. 

## API

#### Save
```golang
db, err := simpledb.Connect("local.json")
u := User{Name: "someone", Age:  99}
err = db.Save(&u)
```

#### Find
Gets a list of records based on a very simple query.

```golang
db, err := simpledb.Connect("local.json")
var u user
err := db.FindOne(&uu, "Name", "someone")
```

#### FindWhere
Gets a list of records record based on a more advanced query,
using the simpledb.Where param.

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
