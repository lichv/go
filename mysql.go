package lichv

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type MysqlDriver struct {
	DriverName     string
	DataSourceName string
	Show           bool
	DB             *sql.DB
	SQLTX          *sql.Tx
}

func InitMysqlDriver(host string, port int, user, password, dbname string) *MysqlDriver {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&autocommit=true", user, password, host, port, dbname, "utf8")
	db := &MysqlDriver{
		DriverName:     "mysql",
		DataSourceName: dataSourceName,
	}

	if err := db.Open(); err != nil {
		log.Panicln("Init mysql pool failed.", err.Error())
	}
	return db
}

func (db *MysqlDriver) Open() (err error) {
	db.DB, err = sql.Open(db.DriverName, db.DataSourceName)
	if err != nil {
		return err
	}
	if err = db.DB.Ping(); err != nil {
		return err
	}
	db.DB.SetMaxOpenConns(20)
	db.DB.SetMaxIdleConns(10)
	db.DB.SetConnMaxLifetime(time.Second * 10)
	return nil
}

func (db *MysqlDriver) Close() error {
	return db.DB.Close()
}

func (db *MysqlDriver) ShowSql() error {
	db.Show = true
	return nil
}

func (db *MysqlDriver) HideSql() error {
	db.Show = false
	return nil
}

func (db *MysqlDriver) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if db.Show {
		fmt.Println(query)
		fmt.Println(args...)
	}
	return db.DB.Query(query, args...)
}

func (db *MysqlDriver) Exec(query string, args ...interface{}) (sql.Result, error) {
	if db.Show {
		fmt.Println(query)
		fmt.Println(args...)
	}
	return db.DB.Exec(query, args...)
}

func (db *MysqlDriver) QueryMap(tableName string, query map[string]interface{}) (*sql.Rows, error) {
	s := "select * from " + tableName
	where, _ := WhereFromQuery(query)
	if db.Show {
		fmt.Println(s + where)
	}
	rows, err := db.Query(s + where)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (db *MysqlDriver) FindById(tableName string, id int64) (*sql.Rows, error) {
	s := "select * from " + tableName + " where id = ? limit 1 "
	if db.Show {
		fmt.Println(s)
		fmt.Println(id)
	}
	rows, err := db.Query(s, id)
	if err != nil {

		return nil, err
	}
	return rows, nil
}

func (db *MysqlDriver) FindOne(tableName string, query map[string]interface{}, orderBy string) (*sql.Rows, error) {
	s := "select * from " + tableName
	if !CheckOrderBy(orderBy) {
		orderBy = ""
	}
	where, _ := WhereFromQuery(query)
	if db.Show {
		fmt.Println(s + where)
	}
	rows, err := db.DB.Query(s + where)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (db *MysqlDriver) GetList(tableName string, query map[string]interface{}, orderBy string) (*sql.Rows, error) {
	s := "select * from " + tableName
	if !CheckOrderBy(orderBy) {
		orderBy = ""
	}
	where, _ := WhereFromQuery(query)
	if db.Show {
		fmt.Println(s + where)
	}
	rows, err := db.DB.Query(s + where)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (db *MysqlDriver) GetPage(tableName string, query map[string]interface{}, orderBy string, page, size int64) (*sql.Rows, *Page, error) {
	var total, last, prev, next int64
	total, _ = db.Count(tableName, query)
	last = total/size + 1
	prev = 1
	if page > 2 {
		prev = page - 1
	} else {
		page = 1
	}
	next = last
	if page < last-1 {
		next = page + 1
	}
	offset := (page - 1) * size
	s := "select * from " + tableName
	if !CheckOrderBy(orderBy) {
		orderBy = ""
	}
	where, _ := WhereFromQuery(query)
	sql2 := s + where
	if orderBy != "" {
		sql2 += " order by " + orderBy
	}
	sql2 += " limit ? offset ?"
	if db.Show {
		fmt.Println(sql2)
		fmt.Println(size,offset)
	}
	rows, err := db.DB.Query(sql2, size, offset)
	if err != nil {
		return nil, nil, err
	}
	return rows, &Page{First: 1, Prev: prev, Page: page, Next: next, Last: last, Size: size, Total: total}, nil
}

func (db *MysqlDriver) Count(tableName string, query map[string]interface{}) (int64, error) {
	var count int64 = 0
	s := "select count(1) as number from " + tableName
	where, _ := WhereFromQuery(query)
	if db.Show {
		fmt.Println(s + where)
	}
	rows, err := db.DB.Query(s + where)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	_ = rows.Scan(&count)
	return count, nil
}

func (db *MysqlDriver) Exists(tableName string, query map[string]interface{}) bool {
	c, err := db.Count(tableName, query)
	if err != nil {
		return false
	}
	return c > 0
}

func (db *MysqlDriver) Insert(tableName string, post map[string]interface{}) (int64, error) {
	s, _ := GetInsertSql(tableName, post,"mysql")
	if db.Show {
		fmt.Println(s)
	}
	exec, err := db.DB.Exec(s)
	if err != nil {
		return 0, err
	}
	return exec.LastInsertId()
}

func (db *MysqlDriver) Update(tableName string, post map[string]interface{}, query map[string]interface{}) (int64, error) {
	s, _ := GetUpdateSQL(tableName, post, query,"mysql")
	if db.Show {
		fmt.Println(s)
	}
	exec, err := db.DB.Exec(s)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (db *MysqlDriver) Save(tableName string, post map[string]interface{}) (int64, error) {
	id, ok := post["id"]
	if ok {
		delete(post, "id")
		return db.Update(tableName, post, map[string]interface{}{"id": id})
	} else {
		return db.Insert(tableName, post)
	}
}

func (db *MysqlDriver) Delete(tableName string, query map[string]interface{}) (int64, error) {
	where, _ := WhereFromQuery(query)
	if where != "" {
		s := "delete from " + tableName + where
		if db.Show {
			fmt.Println(s)
		}
		exec, err := db.DB.Exec(s)
		if err != nil {
			return 0, err
		}
		return exec.RowsAffected()
	} else {
		return 0, nil
	}
}

func (db *MysqlDriver) DeleteById(tableName string, id int64) (int64, error) {
	if id != 0 {
		s := "delete from " + tableName + " where id = ?"
		if db.Show {
			fmt.Println(s)
			fmt.Println(id)
		}
		exec, err := db.DB.Exec(s, id)
		if err != nil {
			return 0, err
		}
		return exec.RowsAffected()
	} else {
		return 0, nil
	}
}

func (db *MysqlDriver) Begin() error {
	err := db.DB.Ping()
	if err != nil {
		return nil
	}
	db.SQLTX, err = db.DB.Begin()
	if err != nil {
		return err
	}
	return nil
}

func (db *MysqlDriver) RollBack() error {
	return db.SQLTX.Rollback()
}

func (db *MysqlDriver) Commit() error {
	return db.SQLTX.Commit()
}

func (db *MysqlDriver) QueryTX(query string, args ...interface{}) (*sql.Rows, error) {
	return db.SQLTX.Query(query, args...)
}

func (db *MysqlDriver) ExecTX(query string, args ...interface{}) (sql.Result, error) {
	return db.SQLTX.Exec(query, args...)
}

