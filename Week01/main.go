package main

import "github.com/pkg/errors"
import  _ "github.com/go-sql-driver/mysql"
import "database/sql"
import "fmt"
import "log"
// question 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
// 个人感觉需要处理这个error,并且warp之后抛给上层。go的error不同于其他语言的异常处理,err不是传统意义上的异常，只是一个正常的业务返回值。当通过主键或其他筛选条件无法定位到数据行时,理应返回一个error
// 来告知调用端 “数据无法找到，请不要使用返回的value ”。由于dao层是业务逻辑的最底层，因此需要wrap这个error往上抛。

var db *sql.DB
func main() {
	var initDbErr error
	db, initDbErr = sql.Open("mysql", "root:123@/mydb?charset=utf8")
	if initDbErr!=nil {
		log.Fatal(initDbErr)
	}
	fmt.Printf("%+v",service())
}

func dao(sqlStr string) error {
	rows,err:= db.Query(sqlStr)
	if err!=nil {
		return errors.Wrap(err,"product data not found!!")
	}
	fmt.Println(rows)
	return nil
}

func service() error{
	return dao("select id,name,price from product")
}