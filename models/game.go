package models

// Game represents a game of DDA.
type Game struct {
	ID        string `json:"id" db:"id"`
	Song      string `json:"song" db:"song"`
	HomeID    string `json:"home_id" db:"home_id"`
	HomeScore int    `json:"home_score" db:"home_score"`
	HomeReady bool   `json:"home_ready" db:"home_ready"`
	AwayID    string `json:"away_id" db:"away_id"`
	AwayScore int    `json:"away_score" db:"away_score"`
	AwayReady bool   `json:"away_ready" db:"away_ready"`
	Started   int64  `json:"started" db:"started"`
	Finished  int64  `json:"finished" db:"finished"`
}
