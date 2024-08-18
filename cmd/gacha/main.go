package main

import (
	"database/sql"
	"fmt"
	"go-mock-test-demo/gacha"
	"go-mock-test-demo/gacha/repository"
	"go-mock-test-demo/random"
	"go-mock-test-demo/tx"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("実行時引数にユーザーのIDを指定してください")
	}

	userIdStr := os.Args[1]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		log.Fatal("ユーザーIDは整数値で指定してください")
	}

	db, err := sql.Open("mysql", "root:password@tcp(localhost:13306)/app")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = db.Close()
	}()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	var (
		userRep     = repository.NewUser(db)
		itemRep     = repository.NewItem(db)
		userItemRep = repository.NewUserItem(db)
		tx          = tx.NewTransaction(db)
		rg          = random.NewRandGenerator()
	)

	gacha := gacha.NewGacha(userRep, itemRep, userItemRep, tx, rg)

	result, err := gacha.Draw(userId)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
