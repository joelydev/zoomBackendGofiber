package schemas

// Chatmsg schema structure
type Chatmsg struct {
	Message string  `json:"message"`
	MessageType string  `json:"messagetype"`
	FilePath string  `json:"filepath"`
	From string  `json:"from"`
	To string  `json:"to"`
	UserId string `json:"userId" bson:"userId"`
	UserName string `json:"username" bson:"username"`
	Created int64  `json:"created"`
	ID      string `json:"id,omitempty" bson:"_id,omitempty"`
	Updated int64  `json:"updated"`
}
