package routes

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Cloudstorage(serv *fiber.App) {
	r := serv.Group("/fortnite/api/cloudstorage")

	r.Get("/system", func(c *fiber.Ctx) error {
		files, err := os.ReadDir("./hotfixes/")
		if err != nil {
			return err
		}

		result := []CloudStorageSystemEntry{}

		for _, file := range files {
			fi, err := file.Info()
			if err != nil {
				return err
			}

			fileData, err := os.ReadFile(fmt.Sprintf("./hotfixes/%v", file.Name()))
			if err != nil {
				return err
			}

			hash := sha1.New()
			hash.Write([]byte(fileData))

			hash256 := sha256.New()
			hash256.Write([]byte(fileData))

			cs := CloudStorageSystemEntry{
				UniqueFilename: file.Name(),
				Filename:       file.Name(),
				Hash:           hashToHexStr(hash),
				Hash256:        hashToHexStr(hash256),
				Length:         int(fi.Size()),
				ContentType:    "application/octet-stream",
				Uploaded:       fi.ModTime().Format("2006-01-02T15:04:05.999Z"),
				StorageType:    "S3",
				DoNotCache:     false,
			}

			result = append(result, cs)
		}

		return c.JSON(result)
	})

	r.Get("/system/config", func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	r.Get("/system/:filename", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/octet-stream")

		filename := c.Params("filename")

		data, err := os.ReadFile(fmt.Sprintf("./hotfixes/%v", filename))
		if err != nil {
			c.Status(404)
			return c.JSON(ResponseError{
				ErrorCode:          "inferno.cloudstorage.file_not_found",
				ErrorMessage:       fmt.Sprintf("Sorry, we couldn't find a system file for %v", filename),
				MessageVars:        []string{filename},
				NumericErrorCode:   12004,
				OriginatingService: "inferno",
				Intent:             "dev",
			})
		}

		return c.Send(data)
	})
}

func hashToHexStr(hash hash.Hash) string {
	return string(hex.EncodeToString(hash.Sum(nil)))
}

type CloudStorageSystemEntry struct {
	UniqueFilename string `json:"uniqueFilename"`
	Filename       string `json:"filename"`
	Hash           string `json:"hash"`
	Hash256        string `json:"hash256"`
	Length         int    `json:"length"`
	ContentType    string `json:"contentType"`
	Uploaded       string `json:"uploaded"`
	StorageType    string `json:"storageType"`
	DoNotCache     bool   `json:"doNotCache"`
}

type ResponseError struct {
	ErrorCode          string   `json:"errorCode"`
	ErrorMessage       string   `json:"errorMessage"`
	MessageVars        []string `json:"messageVars"`
	NumericErrorCode   int      `json:"numericErrorCode"`
	OriginatingService string   `json:"originatingService"`
	Intent             string   `json:"intent"`
}
