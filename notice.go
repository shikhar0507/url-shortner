package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

func onNotify(c *pgconn.PgConn, n *pgconn.Notice) {
	fmt.Println("Message:", *n)
}
func main() {

	ctx := context.Background()

	//connectionString := fmt.Sprintf("user=%s password=%s host=%s sslmode=disable dbname=db", "user", "pass", "host")

	conf, err := pgx.ParseConfig("postgres://xanadu:xanadu@localhost:5432/tracker")

	if err != nil {
		fmt.Println(err)
	}

	conf.OnNotice = onNotify

	conn, err := pgx.ConnectConfig(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close(ctx)

	query := `SELECT insertLongUrl($1)`
	conn.Config().OnNotice = onNotify
	var id string
	err = conn.QueryRow(ctx, query, "https://google.com").Scan(&id)

	if err != nil {
		log.Fatal(err)
	}

}
