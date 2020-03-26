package gocord

// Emoji - Holds an emoji's data.
type Emoji struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	User          User   `json:"user"`
	RequireColons bool   `json:"require_colons"`
	Managed       bool   `json:"managed"`
	Animated      bool   `json:"animated"`
	Available     bool   `json:"available"`
}
