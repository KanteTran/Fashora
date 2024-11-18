package external

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime/multipart"
)

func SimulateExternalAPI(image1, image2, image3 []byte) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	addFileToWriter := func(fieldName string, data []byte) error {
		part, err := writer.CreateFormFile(fieldName, fieldName+".jpg")
		if err != nil {
			return err
		}
		_, err = part.Write(data)
		return err
	}

	if err := addFileToWriter("person_image", image1); err != nil {
		return "", err
	}
	if err := addFileToWriter("cloth_image", image2); err != nil {
		return "", err
	}
	if err := addFileToWriter("mask", image3); err != nil {
		return "", err
	}

	writer.Close()

	processedImage := base64.StdEncoding.EncodeToString(body.Bytes())
	return processedImage, nil
}

func ReadImageFile(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file content into a byte slice
	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}
