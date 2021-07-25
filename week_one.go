package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"log"
)
import _ "github.com/go-sql-driver/mysql"

/*
	第二周作业:
	1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？9
	答：个人认为, sql.ErrNoRows 不应该再往上抛了，原因有以下几点：
		1、sql.ErrNoRows 基础包将该错误定义为：ErrNoRows is returned by Scan when QueryRow doesn't return a row. In such a case,
			QueryRow returns a placeholder *Row value that defers this error until a Scan.
			出现该错误的场景大致是 sql.QueryRow 没有查询到数据，我认为没有查询到数据不应该再向调用方返回错误;
		2、sql.ErrNoRows 应该作为 偏数据层错误，应该在 dao 层得到解决
		3、查询数据为空不属于异常; 并且也不需要专门独立记录日志
		4、如果将 sql.ErrNoRows 继续向上抛出, 上层能够通过 error 判断数据是否为空，为空做其他业务操作，
		如果并不需要记录日志，为什么不直接返回空对象(非 nil, 空对象)作为查询数据为空的结果？


 */
func main() {
	// Service.
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer func(db1 *sql.DB) {
		log.Println("正在执行释放资源操作...")
		err := db1.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	// 向上抛出异常模式
	list, err := queryListReError(db)
	if err != nil {
		//fmt.Printf("%+v", err)
		fmt.Println("错误，找不到数据")
		return
	}
	fmt.Println(list)

	// 返回空数据模式
	total := queryListNoError(db)
	if total <= 0 {
		fmt.Println("错误，找不到数据")
		return
	}
}

/// Helper
func connect() (*sql.DB, error) {
	db1, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/test")
	//check err
	if err != nil {
		return nil, errors.Wrap(err, "[DB One] 数据库连接失败!")
	}
	if err = db1.Ping(); err != nil {
		return nil, errors.Wrap(err, "[DB One] 数据库连接失败!")
	}
	db1.SetConnMaxLifetime(0)
	return db1, nil
}


/// Dao
func queryListReError(db *sql.DB) (int64, error) {
	row := db.QueryRow(`SELECT sign_date FROM table_name WHERE create_time >= '2021-12-12 00:00:00'`)
	var total int64
	err := row.Scan(&total)
	if err != nil {
		return -1, errors.Wrap(err, "404 not found!")
	}
	return total, err
}

func queryListNoError(db *sql.DB) int64 {
	row := db.QueryRow(`SELECT sign_date FROM table_name WHERE create_time >= '2021-12-12 00:00:00'`)
	var total int64
	err := row.Scan(&total)
	if err != nil {
		return -1
	}
	return total
}
