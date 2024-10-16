package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func (runTimeData *RunTimeData) startTicker() {
	checkTicker = CustomTicker{
		ticker:  time.NewTicker(1 * time.Second),
		running: true,
	}
	go runTimeData.Check()
	fmt.Printf("SUCCESS @ %s : START MONITORING\n", time.Now().Format("2006-01-02 15:04:05.00000"))
}

func stopTicker() {
	if checkTicker.running {
		checkClose <- true
		checkTicker.running = false
		checkTicker.ticker.Stop()
		fmt.Printf("SUCCESS @ %s : STOP MONITORING\n", time.Now().Format("2006-01-02 15:04:05.00000"))
	}
}

func (runTimeData *RunTimeData) Check() {
	var lastNotification time.Time
	for {
		select {
		case <-checkTicker.ticker.C:
			resp, err := http.Get(runTimeData.Url)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			closed := strings.Contains(string(body), "no longer accepting responses")
			timePassed := int(time.Since(lastNotification).Seconds())
			if !closed && timePassed >= runTimeData.Cooldown {
				fmt.Printf("SUCCESS @ %s : COOLDOWN PASSED - ALERTING USERS\n", time.Now().Format("2006-01-02 15:04:05.00000"))
				for memberID, seen := range runTimeData.GetSubscribed() {
					if !seen {
						user, _ := s.User(memberID)
						dmChannel, _ := s.UserChannelCreate(memberID)
						err := SendComplexMessage(user, dmChannel, OpenEmbed(), runTimeData.FormLink(), runTimeData.SeenButton())
						if err != nil {
							fmt.Printf("FAILURE @ %s : USER %s HASN'T ENABLED DIRECT MESSAGES\n", time.Now().Format("2006-01-02 15:04:05.00000"), memberID)
						}
						fmt.Printf("SUCCESS @ %s : NOTIFIED USER %s\n", time.Now().Format("2006-01-02 15:04:05.00000"), memberID)
					}
				}
				lastNotification = time.Now()
			}
		case <-checkClose:
			for memberID := range runTimeData.GetSubscribed() {
				runTimeData.UpdateSeen(memberID, false)
			}
			return
		}
	}
}
