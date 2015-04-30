package client_test

import (
	"fmt"
	"log"

	"github.com/cockroachdb/cockroach/client"
	"github.com/cockroachdb/cockroach/server"
)

func ExampleDB() {
	s := server.StartTestServer(nil)
	defer s.Stop()

	db := client.Open("https://root@" + s.ServingAddr())

	if err := db.Put("key", "value"); err != nil {
		log.Fatal(err)
	}
	result := db.Get("key")
	fmt.Println(1, result.Rows[0].String())

	if r := db.Inc("inc", 100); r.Err != nil {
		log.Fatal(r.Err)
	}
	result = db.Get("key", "inc")
	fmt.Println(2, result.Rows[0].String())
	fmt.Println(3, result.Rows[1].String())

	b := client.B.Get("key").Scan("i", "j", 100).Put("foo", "bar")
	if err := db.Run(b); err != nil {
		log.Fatal(err)
	}
	fmt.Println(4, b.Results[0].Rows[0].String())
	fmt.Println(5, b.Results[1].Rows[0].ValueInt())
	fmt.Println(6, b.Results[2].Rows[0].String())

	err := db.Tx(func(tx *client.Tx) error {
		return tx.Commit(client.B.Put("aa", "1").Put("ab", "2"))
	})
	if err != nil {
		log.Fatal(err)
	}

	result = db.Scan("a", "b", 100)
	fmt.Println(7, result.Rows[0].String())
	fmt.Println(8, result.Rows[1].String())

	if err := db.Del("key", "inc"); err != nil {
		log.Fatal(err)
	}
	result = db.Get("key", "aa", "inc")
	fmt.Println(9, result.Rows[0].Exists(), result.Rows[0].String())
	fmt.Println(10, result.Rows[1].Exists(), result.Rows[1].String())
	fmt.Println(11, result.Rows[2].Exists(), result.Rows[2].String())

	// Output:
	// 1 key:value
	// 2 key:value
	// 3 inc:100
	// 4 key:value
	// 5 100
	// 6 foo:bar
	// 7 aa:1
	// 8 ab:2
	// 9 false key:nil
	// 10 true aa:1
	// 11 false inc:nil
}
