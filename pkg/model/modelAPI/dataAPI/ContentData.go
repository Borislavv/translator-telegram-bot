package dataAPI

type ContentData struct {
	Content interface{} `json:"content"`
}

// NewContentData - constructor of ContentData struct.
func NewContentData(data interface{}) *ContentData {
	return &ContentData{
		Content: data,
	}
}

// GetContent - return the data of current context.
func (data *ContentData) GetContent() interface{} {
	return data.Content
}
