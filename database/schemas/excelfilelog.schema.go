package schemas

// Chatmsg schema structure
type Excelfilelog struct {
	FileName string `json:"filename"`
	FilePath string `json:"filepath"`
	Creator string  	`json:"creator"`
	CreatorId string  	`json:"creatorid"`
	Created int64  	`json:"created"`
	Updated int64  	`json:"updated"`
	ID      string 	`json:"id,omitempty" bson:"_id,omitempty"`
}
