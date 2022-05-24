package model

import (
	"bootcamp/config"
	"bootcamp/logger"
	"net/http"
	"strconv"

	// "strings"
	"time"

	// "bootcamp/util"
	"database/sql"
	"fmt"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type ReqUser struct {
	// ID 				int 		`json:"id"`
	FName    string `json:"firstname"`
	Lname    string `json:"lastname"`
	Email    string `json:"email"`
	Dob      string `json:"dob"`
	Password string `json:"password"`
}

type ResUser struct {
	ID         int    `json:"id"`
	FName      string `json:"firstname"`
	Lname      string `json:"lastname"`
	Email      string `json:"email"`
	Dob        string `json:"dob"`
	Created_at string `json:"created_"at`
	Updated_at string `json:"updated_"at`
}

type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"username"`
	jwt.StandardClaims
}

func init() {
	config.Dbinit()
	db = config.GetDb()
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (b *ReqUser) CreatUser() error {
	archived := 0
	password := b.Password
	t, err := time.Parse("2006-Jan-02", b.Dob)
	fmt.Println(t)
	if err != nil {
		logger.ErrorLogger.Println("Time value is not correct")
	}
	hash, _ := HashPassword(password)
	var createdat = time.Now()
	_, err = db.Exec(`INSERT INTO user(firstname,lastname,email,dob,password,created_at,archived)
	VALUES (?,?,?,?,?,?,?);`, b.FName, b.Lname, b.Email, t.Format("2006-01-02"), hash, createdat, archived)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	return err
}

func Query(r *http.Request) string {
	archived := r.URL.Query().Get("archived")
	id := r.URL.Query().Get("id")
	name := r.URL.Query().Get(("name"))
	email := r.URL.Query().Get("email")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	query := "Select id,firstname,lastname,email,dob,created_at from user where archived=0"
	if archived == "true" {
		query = "Select id,firstname,lastname,email,dob,created_at from user where archived=1"
	}

	if id != "" {
		query += " and id =" + id
	}
	if name != "" {
		query += ` and firstname like'%` + name + `%'or lastname like '%` + name + `%'`
	}
	if email != "" {
		query += ` and email like '%` + email + `%'`
	}
	if sort != "" {
		if order != "" {
			query += ` ORDER BY ` + sort + ` ` + order
		} else {
			query += ` ORDER BY ` + sort + ` ASC`
		}
	}
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}
	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)

	query += fmt.Sprintf(` LIMIT %d OFFSET %d`, l, (p-1)*l)

	return query
}

func GetAllUser(r *http.Request) []ResUser {
	query := Query(r)
	fmt.Println(query)
	var data []ResUser
	row, err := db.Query(query)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	defer row.Close()

	for row.Next() {
		var u ResUser
		err := row.Scan(&u.ID, &u.FName, &u.Lname, &u.Email, &u.Dob, &u.Created_at)
		if err != nil {
			logger.GeneralLogger.Println(err)
		}
		data = append(data, u)
	}
	return data
}

func DeleteUser(Id int64) {
	var updated = time.Now()
	var count int
	err := db.QueryRow(`select count(id) from user where id = ?`, Id).Scan(&count)
	if err != nil {
		logger.GeneralLogger.Println("there is no id in here")
	}

	if count == 0 {
		fmt.Println("id is not present")
		return
	}

	_, err = db.Exec(`Update user SET archived=1,updated_at=? where id=?`, updated, Id)
	if err != nil {
		logger.GeneralLogger.Println("Archive Not Set")
		return
	}
}

func GetUserByid(id int) (ReqUser, *sql.DB) {
	var data ReqUser
	row := db.QueryRow(`select firstname,lastname,email,dob,updated_at from user where id = ?`, id)
	err := row.Scan(&data)
	if err != nil {
		fmt.Println(err)
	}

	return data, db
}
