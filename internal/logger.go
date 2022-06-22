package internal

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/xuri/excelize/v2"
)

type Log struct {
	Bucket string
	Item   *s3.Object
	Error  string
}

type logger struct {
	file           *excelize.File
	currentRow     int64
	totalItems     int64
	totalItemSizes int64
}

func NewLogger(file *excelize.File) *logger {
	return &logger{file: file, currentRow: 1, totalItems: 0, totalItemSizes: 0}
}

func (logger *logger) Log(log Log) {
	if log.hasItem() {
		logger.totalItems++
		logger.totalItemSizes += *log.Item.Size
	}

	if log.hasError() {
		_ = logger.file.SetCellValue("Sheet1", fmt.Sprintf("A%d", logger.currentRow), log.Bucket)
		_ = logger.file.SetCellValue("Sheet1", fmt.Sprintf("C%d", logger.currentRow), log.Error)

		if log.hasItem() {
			_ = logger.file.SetCellValue("Sheet1", fmt.Sprintf("B%d", logger.currentRow), *log.Item.Key)
		}

		logger.currentRow++
	}
}

func (logger *logger) SaveFile(fileName string) error {
	logger.file.SetActiveSheet(logger.file.NewSheet("Sheet1"))

	_ = logger.file.SetCellValue("Sheet1", "A1", "Bucket")
	_ = logger.file.SetCellValue("Sheet1", "B1", "Key")
	_ = logger.file.SetCellValue("Sheet1", "C1", "Error")
	_ = logger.file.SetCellValue("Sheet1", "O10", fmt.Sprintf("Total objects: %d", logger.totalItems))
	_ = logger.file.SetCellValue("Sheet1", "O11", fmt.Sprintf("Total item sizes: %d bytes", logger.totalItemSizes))

	return logger.file.SaveAs(fileName)
}

func (l *Log) hasError() bool {
	return l.Error != ""
}

func (l *Log) hasItem() bool {
	return l.Item != nil
}
