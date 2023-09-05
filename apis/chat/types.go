package chat

type ChatMsgRequest struct {
	Order int64 `json:"order" xml:"order" form:"order"`
	Message string `json:"message" xml:"message" form:"message"`
	Receiver string `json:"receiver" xml:"receiver" form:"receiver"`
	Sender string `json:"sender" xml:"sender" form:"sender"`
	MessageType string `json:"type" xml:"type" form:"type"`
}

type FileTransferRequest struct {
	FileID string `json:"fileID" xml:"fileID" form:"fileID"`
	FileName string `json:"fileName" xml:"fileName" form:"fileName"`
	FileObj string `json:"fileObj" xml:"fileObj" form:"fileObj"`
	FileSize string `json:"fileSize" xml:"fileSize" form:"fileSize"`
	FileType int64 `json:"fileType" xml:"fileType" form:"fileType"`
	MsgID string `json:"msgID" xml:"msgID" form:"msgID"`
	ReceiverType int64 `json:"receiverType" xml:"receiverType" form:"receiverType"`
	XmppMsgData string `json:"xmppMsgData" xml:"xmppMsgData" form:"xmppMsgData"`
}
