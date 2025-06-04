package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type rssFeed struct {
	ID   int
	Name string

	ScanSlepper int

	LastScan  int
	LastFound int
}

func main() {
	var err error
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	var feedList []rssFeed
	fmt.Println("sql is online")

	result, err := db.Query("SELECT feed.id, feed.name, feed.ttl, feed.lastUpdate, MAX(entry.date) FROM feed LEFT JOIN entry ON feed.id = entry.id_feed GROUP BY feed.id;")

	for result.Next() {
		var newFeed rssFeed
		result.Scan(&newFeed.ID, &newFeed.Name, &newFeed.ScanSlepper, &newFeed.LastScan, &newFeed.LastFound)
		feedList = append(feedList, newFeed)
	}

	for _, v := range feedList {
		UpdateTime := v.LastScan - v.LastFound

		if UpdateTime < 1209600 {
			if v.ScanSlepper != 0 {
				fmt.Println(v.Name + " : " + secToHumanTime(UpdateTime))
				db.Exec("UPDATE feed SET ttl=? WHERE id=?;", 0, v.ID)
			}
		} else if UpdateTime < 2592000 {
			if v.ScanSlepper != 86400 {
				fmt.Println(v.Name + " : " + secToHumanTime(UpdateTime))
				db.Exec("UPDATE feed SET ttl=? WHERE id=?;", 86400, v.ID)
			}
		} else if UpdateTime < 7776000 {
			if v.ScanSlepper != 604800 {
				fmt.Println(v.Name + " : " + secToHumanTime(UpdateTime))
				db.Exec("UPDATE feed SET ttl=? WHERE id=?;", 604800, v.ID)
			}
		} else {
			if v.ScanSlepper != 2592000 {
				fmt.Println(v.Name + " : " + secToHumanTime(UpdateTime))
				db.Exec("UPDATE feed SET ttl=? WHERE id=?;", 2592000, v.ID)
			}
		}
	}

}

func secToHumanTime(sec int) (human string) {
	var found int
	found, sec = returnRemaing(sec, 60)
	human = strconv.Itoa(found)
	if sec <= 0 {
		return
	}
	found, sec = returnRemaing(sec, 3600)
	human = strconv.Itoa(found/60) + ":" + human
	if sec <= 0 {
		return
	}
	found, sec = returnRemaing(sec, 86400)
	human = strconv.Itoa(found/3600) + ":" + human
	if sec <= 0 {
		return
	}
	human = strconv.Itoa(sec/86400) + " D-" + human
	return
}

func returnRemaing(input int, taget int) (found int, remaing int) {
	found = input % taget
	remaing = input - found
	return
}
