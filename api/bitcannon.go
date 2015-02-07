package main

import (
	"bufio"
	"github.com/antonholmquist/jason"
	"io/ioutil"
	"log"
	"os"
)

var trackers []string
var archives []*jason.Object
var torrentDB *TorrentDB
var err error

const resultLimit int = 100

func main() {
	// Get mongo url from config.json, otherwise default to 127.0.0.1
	mongo := "127.0.0.1"
	bitcannonPort := "1337"
	f, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("[!!!] Config not loaded")
	} else {
		json, err := jason.NewObjectFromBytes(f)
		if err == nil {
			// Get mongo connection details
			val, err := json.GetString("mongo")
			if err == nil {
				mongo = val
			}
			// Get desired port
			val, err = json.GetString("bitcannonPort")
			if err == nil {
				bitcannonPort = val
			}
			// Get archive sources
			arc, err := json.GetObjectArray("archives")
			if err == nil {
				archives = arc
			}
			// Get trackers
			trac, err := json.GetStringArray("trackers")
			if err == nil {
				trackers = trac
			}
		}
	}
	// Try to connect to the database
	log.Println("[OK!] Connecting to Mongo at " + mongo)
	torrentDB, err = NewTorrentDB(mongo)
	if err != nil {
		log.Println("[ERR] I'm sorry! I Couldn't connect to Mongo.")
		log.Println("      Please make sure it is installed and running.")
		return
	}
	defer torrentDB.Close()

	if len(os.Args) > 1 {
		importFile(os.Args[1])
		enterExit()
	} else {
		runServer(bitcannonPort)
	}
}

func runServer(bitcannonPort string) {
	log.Println("[OK!] BitCannon is live at http://127.0.0.1:" + bitcannonPort + "/")
	api := NewAPI()
	api.AddRoutes()
	runScheduler()
	api.Run(":" + bitcannonPort)
}

func enterExit() {
	log.Println("\n\nPress enter to quit...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
