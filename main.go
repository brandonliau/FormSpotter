package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
)

var s *discordgo.Session
var checkTicker CustomTicker
var checkClose = make(chan bool)
var scheduler = cron.New()

type CustomTicker struct {
	ticker  *time.Ticker
	running bool
}

type RunTimeData struct {
	Token     string `yaml:"token"`
	Guild     string `yaml:"guild"`
	Url       string `yaml:"url"`
	Cooldown  int    `yaml:"cooldown"`
	StartTime string `yaml:"startTime"`
	StopTime  string `yaml:"stopTime"`
	Db        *sql.DB
}

func (runTimeData *RunTimeData) LoadConfig(filename string) {
	rawYaml, _ := os.ReadFile(filename)
	_ = yaml.Unmarshal(rawYaml, &runTimeData)
	fmt.Printf("SUCCESS @ %s : LOAD CONFIG\n", time.Now().Format("2006-01-02 15:04:05.00000"))
}

func (runTimeData *RunTimeData) InitRunTimeData() {
	// open connection to db
	db, _ := sql.Open("sqlite", "./database.db")
	runTimeData.Db = db
	fmt.Printf("SUCCESS @ %s : CONNECTED TO DATABASE\n", time.Now().Format("2006-01-02 15:04:05.00000"))
	// intialize database
	runTimeData.InitDb()
	fmt.Printf("SUCCESS @ %s : INITIALIZE DATABASE\n", time.Now().Format("2006-01-02 15:04:05.00000"))
}

func main() {
	var runTimeData RunTimeData
	runTimeData.LoadConfig("./config.yml")
	runTimeData.InitRunTimeData()
	runTimeData.InitHandlers()

	s, _ = discordgo.New("Bot " + runTimeData.Token)
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	s.AddHandler(runTimeData.OnReady)
	s.AddHandler(runTimeData.InteractionHandler)

	s.Open()
	fmt.Printf("SUCCESS @ %s : ESTABLISH WEBSOCKET CONNECTION\n", time.Now().Format("2006-01-02 15:04:05.00000"))
	s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", Commands)
	fmt.Printf("SUCCESS @ %s : REGISTER ALL COMMANDS\n", time.Now().Format("2006-01-02 15:04:05.00000"))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	stopTicker()

	s.Close()
	fmt.Printf("SUCCESS @ %s : CLOSE WEBSOCKET CONNECTION\n", time.Now().Format("2006-01-02 15:04:05.00000"))
	s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", nil)
	fmt.Printf("SUCCESS @ %s : REMOVE ALL COMMANDS\n", time.Now().Format("2006-01-02 15:04:05.00000"))
}
