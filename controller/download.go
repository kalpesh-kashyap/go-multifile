package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"url-file-save/constant"
	"url-file-save/models"
)

func readFile() ([]models.FILEMODEL, error) {
	file, err := os.Open(constant.JSON_FILE_PATH)
	if errors.Is(err, os.ErrNotExist) {
		return []models.FILEMODEL{}, nil
	}
	defer file.Close()

	stat, err := file.Stat()

	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	if stat.Size() == 0 {
		return []models.FILEMODEL{}, nil
	}

	var posts []models.FILEMODEL
	err = json.NewDecoder(file).Decode(&posts)
	if err != nil {
		return nil, err
	}
	return posts, nil

}

func saveFileData(fileData models.FILEMODEL) (models.FILEMODEL, error) {
	file, err := os.OpenFile(constant.JSON_FILE_PATH, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return models.FILEMODEL{}, err
	}
	defer file.Close()
	posts, err := readFile()
	if err != nil {
		return models.FILEMODEL{}, err
	}

	fileData.ID = len(posts) + 1

	posts = append(posts, fileData)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	if err := encoder.Encode(posts); err != nil {
		return models.FILEMODEL{}, err
	}
	return fileData, nil
}

func DownloadAndSaveFile(url string, filePath string) ([]models.FILEMODEL, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", res.Status)
	}

	fileType := res.Header.Get("Content-Type")

	fileName := extractFileName(res, url)
	fullName := filepath.Join(filePath, fileName)

	out, err := os.Create(fullName)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}

	var fileData models.FILEMODEL
	fileData.FileType = fileType
	fileData.FilePath = fullName
	fileData.URL = url
	fileData.FileName = fileName
	fileData.CretedDate = fmt.Sprint(time.Now())

	fileData, err = saveFileData(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	return []models.FILEMODEL{fileData}, nil
}

func extractFileName(res *http.Response, url string) string {
	contentDisposition := res.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err == nil {
			if fileName, ok := params["filename"]; ok {
				return fileName
			}
		}
	}
	return filepath.Base(url)
}
