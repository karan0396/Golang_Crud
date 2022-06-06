# Folder Structure :-
***
cmd
├── config
│   └── config.go
├── config.json
├── go.mod
├── go.sum
├── internal
│   ├── controller
│   │   └── userController.go
│   ├── middleware
│   │   └── midddleware.go
│   ├── model
│   │   ├── Repository.go
│   │   └── user.go
│   └── service
│       └── service.go
├── log.json
├── pkg
│   ├── logger
│   │   └── logger.go
│   ├── parse
│   │   └── parse.go
│   ├── token
│   │   └── createToken.go
│   └── validation
│       └── validation.go
├── Readme.md
├── route
│   └── routes.go
└── server
    └── main.go

# To change DRIVER and DSN in Config.Json
***

>POSTGRES
Driver:"postgres"
Dsn : "postgres://employee:123456789@localhost:5432/employee?sslmode=disable"

>SQL
Driver:"mysql",
Dsn":"root:123456789@tcp(localhost:3306)/employee?charset=utf8"
