package gocord

// Channel holds a discord channel's data.
type Channel struct {
	CreatedTimestamp int
	Deleted          bool
	ID               string
	Type             string
}
