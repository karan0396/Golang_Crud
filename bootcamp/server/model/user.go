package model

import (
	"bootcamp/config"
	"bootcamp/logger"
	"errors"
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

type params struct{
	archieved string
	id string
	name string
	email string
	sort string
	order string
	page string
	limit string
}


func init() {
	config.Dbinit()
	db = config.GetDb()
	// fmt.Println(db)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (b *ReqUser) CreatUser() error {
	archived := 0
	password := b.Password
	fmt.Println(b.Dob)
	t, err := time.Parse("2006-01-02", b.Dob)
	fmt.Println(t)
	if err != nil {
		logger.ErrorLogger.Println("Time value is not correct")
	}
	fmt.Println(b.Email)
		var countEmail int
		db.QueryRow(`SELECT COUNT(email) FROM user where email=?`, b.Email).Scan(&countEmail)
		
		fmt.Println(countEmail)
		if countEmail != 0 {
			logger.ErrorLogger.Println("email already exit")
			return errors.New("email already exist")
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
	p:= params{
		archieved : r.URL.Query().Get("archived"),
		id : r.URL.Query().Get("id"),
		name : r.URL.Query().Get(("name")),
		email : r.URL.Query().Get("email"),
		sort : r.URL.Query().Get("sort"),
		order : r.URL.Query().Get("order"),
		page : r.URL.Query().Get("page"),
		limit : r.URL.Query().Get("limit"),

	}


	query := "Select id,firstname,lastname,email,dob,created_at from user where archived=0"
	if p.archieved == "true" {
		query = "Select id,firstname,lastname,email,dob,created_at from user where archived=1"
	}

	if p.id != "" {
		query += " and id =" + p.id
	}
	if p.name != "" {
		query += ` and firstname like'%` + p.name + `%'or lastname like '%` + p.name + `%'`
	}
	if p.email != "" {
		query += ` and email like '%` + p.email + `%'`
	}
	if p.sort != "" {
		if p.order != "" {
			query += ` ORDER BY ` + p.sort + ` ` + p.order
		} else {
			query += ` ORDER BY ` + p.sort + ` ASC`
		}
	}
	if p.page == "" {
		p.page = "1"
	}
	if p.limit == "" {
		p.limit = "10"
	}
	q, _ := strconv.Atoi(p.page)
	l, _ := strconv.Atoi(p.limit)

	query += fmt.Sprintf(` LIMIT %d OFFSET %d`, l, (q-1)*l)

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
