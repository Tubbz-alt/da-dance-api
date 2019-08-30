package models

// Score represents a DDA score.
type Score struct {
	Player string `json:"player" db:"player"`
	Game   string `json:"game" db:"game"`
	Points int    `json:"points" db:"points"`
}
