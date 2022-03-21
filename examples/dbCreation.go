package main

import (
	"fmt"

	"github.com/MikhailBatsin-code/gomemdb/gomemdb"
)

func main() {
	// create new db
	// you need to give to this function filename that will be used for db
	db := gomemdb.NewDb("example.db")

	// basic interactions
	// add Key-Value Pair to db
	// this function returns error if key already in use
	err := db.Add("example", "any type that you want")

	if err != nil {
		fmt.Println("this is bad :(")
	}

	// you can get value from db
	fmt.Println(db.Get("example"))

	// you can change value
	db.Set("example", "just a text")
	fmt.Println(db.Get("example"))

	// you can save it to the file with compression and without
	// uncomment this if you want compression
	// db.NeedCompress()
	// db.ZlibCompressLevel = zlib.BestCompression
	// db.Save()

	// comment this if you uncommented upper lines
	db.Save()

	// you can open file and get db state
	// if it is without compression
	db2, err := gomemdb.Open("example.db", false)

	// you can get map representation of db
	fmt.Println(db2.Representate())

	// you can clear db
	db2.Clear()

	fmt.Println(db2.Get("example"))

	// you can delete key-value pair by it's key
	// use it if you have more than one element
	fmt.Println(db.Delete("example"))
	fmt.Println(db.Get("example")) // does not work

	db.Add("another", "one")
	db.Add("hour", 60)
	db.Add("two hours", 120)

	// you can get all pairs with the same values datatype as map
	fmt.Println(db.GroupByPairDatatype(""))
	fmt.Println(db.GroupByPairDatatype(0))

	// if you really need information about version and license type this
	fmt.Println(gomemdb.Info())
}
