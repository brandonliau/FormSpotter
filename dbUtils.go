package main

func (runTimeData *RunTimeData) InitDb() {
	query := "CREATE TABLE IF NOT EXISTS subscribed (memberID TEXT, seen BOOLEAN)"
	runTimeData.Db.Exec(query)
}

func (runTimeData *RunTimeData) InsertIntoDb(memberID string) {
	query := "INSERT INTO subscribed (memberID, seen) VALUES (?, ?)"
	runTimeData.Db.Exec(query, memberID, false)
}

func (runTimeData *RunTimeData) RemoveFromDb(memberID string) {
	query := "DELETE FROM subscribed WHERE memberID = ?"
	runTimeData.Db.Exec(query, memberID)
}

func (runTimeData *RunTimeData) UpdateSeen(memberID string, val bool) {
	query := "UPDATE subscribed SET seen = ? WHERE memberID = ?"
	runTimeData.Db.Exec(query, val, memberID)
}

func (runTimeData *RunTimeData) GetSeen(memberID string) bool {
	var seen bool
	query := "SELECT FROM subscribed WHERE memberID = ?"
	row := runTimeData.Db.QueryRow(query, memberID)
	row.Scan(&seen)
	return seen
}

func (runTimeData *RunTimeData) GetSubscribed() map[string]bool {
	var id string
	var seen bool
	registered := make(map[string]bool)
	query := "SELECT memberID, seen FROM subscribed"
	rows, _ := runTimeData.Db.Query(query)
	for rows.Next() {
		rows.Scan(&id, &seen)
		registered[id] = seen
	}
	return registered
}
