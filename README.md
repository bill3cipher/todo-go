# todo-go
todo api with golang

### use

* gin
* gorm
* dotenv
* jwt
* sqlite && mariaDB

### How to use
> Create local.env and add like this
```
PORT=8081
SIGN=jwtSign
DB_CON=root:my-secret-pw@tcp(127.0.0.1:3306)/myapp?charset=utf8mb4&parseTime=True&loc=Local
```

### How to run
```
go rum main.go
```
