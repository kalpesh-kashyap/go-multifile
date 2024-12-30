package models

type FILEMODEL struct {
	ID         int    `json:"id"`
	FileName   string `json:"fileNam"`
	FilePath   string `json:"filePath"`
	FileType   string `json:"fileType"`
	URL        string `json:"url"`
	CretedDate string `json:"ceatedDate"`
}
