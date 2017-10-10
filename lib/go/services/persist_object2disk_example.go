package services

import (
	"time"
	"log"
)

type obj struct {
	Name string
	Number int
	When time.Time
}

func main() {
	o := &obj{
		Name:   "Mat",
		Number: 47,
		When:   time.Now(),
	}
	if err := Save("./file.tmp", o); err != nil {
		log.Fatalln(err)
	}
	// load it back
	var o2 obj
	if err := Load("./file.tmp", o2); err != nil {
		log.Fatalln(err)
	}
	// o and o2 are now the same
	// and check out file.tmp - you'll see the JSON file
}