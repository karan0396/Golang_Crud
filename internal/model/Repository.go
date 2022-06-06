package model

import (
	"api/config"
	"api/pkg/logger"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var Db *sql.DB

type Params struct {
	Archieved string
	Id        string
	Name      string
	Email     string
	Sort      string
	Order     string
	Page      string
	Limit     string
}

func init() {
	logger.IntializeLogger()

}

type repository struct {
	db  *sql.DB
	con *config.Config
}

type Repository interface {
	Create(b *User) error
	Getuserbyid(id int) (User, *sql.DB)
	Getpasswordbyemail(data Credential) (string, error)
	Get(p Params) ([]User, error)
	Delete(Id int) error
	Update(data User) error
	DeleteUserbynTimes() error
}
//password Hashing
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}
// initializing database and bounding intreface to struct
func NewRepository(con *config.Config) Repository {
	data, err := sql.Open(con.Database.Driver, con.Database.Dsn)
	if err != nil {
		logger.Logger.DPanic("database is not connected", zap.Error(err))
		return nil
	}

	if err := data.Ping(); err != nil {
		logger.Logger.DPanic("Connecion is down", zap.Error(err))
	}

	Db = data

	fmt.Println("Database is Connected")
	return &repository{Db, con}
}
//entering data to databse
func (r repository) Create(b *User) error {
	archived := 0
	password := b.Password
	t, err := time.Parse("2006-01-02", b.Dob)
	if err != nil {
		logger.Logger.Warn("Dob is not parsing", zap.Error(err))
		return err
	}
	var countEmail int
	r.db.QueryRow(`SELECT COUNT(email) FROM user1 where email=?`, b.Email).Scan(&countEmail)

	fmt.Println(countEmail)
	if countEmail != 0 {
		logger.Logger.DPanic("Email is alreadty exit", zap.Int("total email", countEmail))
		return errors.New("email is alreadty exit")
	}
	hash, _ := HashPassword(password)
	var createdat = time.Now()

	_, err = r.db.Exec(`INSERT INTO user1(firstname,lastname,email,dob,password,created_at,archived)
		VALUES (?,?,?,?,?,?,?);`, b.FName, b.Lname, b.Email, t.Format("2006-01-02"), hash, createdat, archived)
	if err != nil {
		logger.Logger.DPanic("error in pushing data", zap.Error(err))
		return err
	}
	return nil
}
//get data by id
func (r repository) Getuserbyid(id int) (User, *sql.DB) {
	var data User
	row := r.db.QueryRow(`select id,firstname,lastname,email,dob from user1 where id = ?`, id)
	err := row.Scan(&data)
	if err != nil {
		logger.Logger.DPanic("data is not added to struct", zap.Error(err))
	}

	return data, r.db
}
//get password by email and checking it is soft delete or not
func (r repository) Getpasswordbyemail(data Credential) (string, error) {
	var password string
	err := r.db.QueryRow(`select password from user1 where email=? and archived=0`, data.Email).Scan(&password)
	if err != nil {
		logger.Logger.DPanic("email is not found")
		return "", err
	}
	return password, nil

}
// read data from database
func (r repository) Get(p Params) ([]User, error) {
	var data []User
	query := "Select id,firstname,lastname,email,dob from user1 where archived=0"
	if p.Archieved == "true" {
		query = "Select id,firstname,lastname,email,dob from user1 where archived=1"
	}
	if p.Id != "" {
		query += " and id =" + p.Id
	}
	if p.Name != "" {
		query += ` and firstname like'%` + p.Name + `%'or lastname like '%` + p.Name + `%'`
	}
	if p.Email != "" {
		query += ` and email like '%` + p.Email + `%'`
	}
	if p.Sort != "" {
		if p.Order != "" {
			query += ` ORDER BY ` + p.Sort + ` ` + p.Order
		} else {
			query += ` ORDER BY ` + p.Sort + ` ASC`
		}
	}

	if p.Page == "" {
		p.Page = r.con.Pagination.Page
	}

	if p.Limit == "" {
		p.Limit = r.con.Pagination.Limit
	}

	page, _ := strconv.Atoi(p.Page)
	limit, _ := strconv.Atoi(p.Limit)

	query += fmt.Sprintf(` LIMIT %d OFFSET %d`, limit, (page-1)*limit)

	row, err := r.db.Query(query)
	if err != nil {
		logger.Logger.DPanic("query data", zap.Error(err))
		return []User{},err
	}
	defer row.Close()

	for row.Next() {
		var u User
		err := row.Scan(&u.ID, &u.FName, &u.Lname, &u.Email, &u.Dob)
		if err != nil {
			logger.Logger.Debug("scanning data from query", zap.Error(err))
			return []User{},err
		}
		data = append(data, u)
	}
	return data, nil
}
// soft delete from database
func (r repository) Delete(Id int) error {
	var updated = time.Now()
	var count int
	err := r.db.QueryRow(`select count(id) from user1 where id = ?`, Id).Scan(&count)
	if err != nil {
		logger.Logger.DPanic("querying id cause error", zap.Error(err))
		return err
	}

	if count == 0 {
		fmt.Println("id is not present")
		return err
	}

	_, err = r.db.Exec(`Update user1 SET archived=1,updated_at=? where id=?`, updated, Id)
	if err != nil {
		logger.Logger.DPanic("executing delete", zap.Error(err))
		return err
	}

	return nil
}
// update in database
func (r repository) Update(data User) error {
	var updated = time.Now()
	var email string
	err := r.db.QueryRow("SELECT email FROM user1 where id=?", data.ID).Scan(&email)
	if err != nil {
		logger.Logger.DPanic("checking id", zap.Error(err))
		return err
	}

	valid := false
	if len(data.FName)+len(data.Lname) < 30 && len(data.Email) < 20 {
		valid = true
	}
	if !valid {
		logger.Logger.DPanic("length of name exceeds minimum length", zap.Error(err))
		return err
	}

	//check if new email is already in use
	if email != data.Email {
		var countEmail int
		r.db.QueryRow("SELECT COUNT(email) FROM user1 where email=?", data.Email).Scan(&countEmail)
		if countEmail != 0 {
			logger.Logger.Warn("email not present", zap.Error(err))
			return err
		}
	}

	_, err = r.db.Exec(`Update user1 Set firstname=?,lastname=?,email=?,dob=?,updated_at=? where id=?`, data.FName, data.Lname, data.Email, data.Dob, updated, data.ID)
	if err != nil {
		logger.Logger.DPanic("update not working", zap.Error(err))
	}
	return nil
}
//hard delete user by n time
func (r repository) DeleteUserbynTimes() error {
	//logging user name and id
	var data []SoftDelete
	rows,_:=r.db.Query(`
	Select id,firstname From user1
	where archived=1 and
	updated_at < DATE_ADD(CURDATE(),INTERVAL ? day)`,r.con.Delete.HardDelete)

	defer rows.Close()

	for rows.Next() {
		var u SoftDelete
		err := rows.Scan(&u.ID, &u.FName)
		if err != nil {
			logger.Logger.Debug("scanning data from query", zap.Error(err))
			
		}
		logger.Logger.Info("data:",zap.Int("Id",u.ID),zap.String("Name",u.FName))
		fmt.Println(u.ID)
		data = append(data, u)
		
	}
	if data!=nil{
		logger.Logger.Info("Hard Delete",zap.Any("Data",data))
	}
	//hard delete data 
	_, err := r.db.Exec(`
	Delete From user1
	where archived=1 and
	updated_at < DATE_ADD(CURDATE(),INTERVAL ? day)`,r.con.Delete.HardDelete)//putting value in negative(-1 represent 2 days before)
	if err != nil {
		logger.Logger.DPanic("soft delete is not working", zap.Error(err))
	}
	
	return nil
}
