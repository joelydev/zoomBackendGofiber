package schemas

// Chatmsg schema structure
type Uploadfilelog struct {
	FileName     string `json:"fileName"`
	FileObj      string `json:"fileObj"`
	FileSize     string `json:"fileSize"`
	FileType     string `json:"fileType"`
	UserId       string `json:"userId"`
	ReceiverType string `json:"receiverType"`
	TransferType string `json:"type" bson:"type"`
	UserName     string `json:"username" bson:"username"`
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	Created      int64  `json:"created"`
	Updated      int64  `json:"updated"`
}
