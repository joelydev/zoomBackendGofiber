package schemas

// Chatmsg schema structure
type Videofilelog struct {
	FileName string  `json:"fileName"`
	UserId string  `json:"userId"`
	UserName string `json:"username" bson:"username"`
	ID      string `json:"id,omitempty" bson:"_id,omitempty"`
	EndTime int64  `json:"created"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`
}