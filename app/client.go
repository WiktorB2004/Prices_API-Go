package app

// Client model
type Client struct {
	ID      string         `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string         `json:"username" bson:"username"`
	AuthKey string         `json:"authKey" bson:"authKey"`
	ApiKeys map[string]int `json:"apiKeys" bson:"apiKeys"`
}
