package songpahpa

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"

	"github.com/ntsd/omise-challenge-go/internal/checkerror"
)

type SongPahPaCSV struct {
	CSVReader        *csv.Reader
	SongPahPaChannel chan<- *SongPahPa
}

func SongPahPaCSVReader(reader io.Reader, songPahPaChannel chan<- *SongPahPa) *SongPahPaCSV {
	csvReader := csv.NewReader(reader)

	// Read the header
	_, err := csvReader.Read()
	checkerror.CheckError(err)

	return &SongPahPaCSV{
		CSVReader:        csvReader,
		SongPahPaChannel: songPahPaChannel,
	}
}

func (c *SongPahPaCSV) Read() {
	defer close(c.SongPahPaChannel)

	for {
		record, err := c.CSVReader.Read()
		if err == io.EOF {
			break
		}
		checkerror.CheckError(err)

		c.SongPahPaChannel <- songPahPaParser(record)
	}
}

func songPahPaParser(record []string) *SongPahPa {
	amount, err := strconv.ParseInt(record[1], 10, 64)
	checkerror.CheckError(err)

	month, err := strconv.Atoi(record[4])
	checkerror.CheckError(err)

	year, err := strconv.Atoi(record[5])
	checkerror.CheckError(err)

	return &SongPahPa{
		Name:     record[0],
		Amount:   amount,
		CCNumber: record[2],
		CVV:      record[3],
		ExpMonth: time.Month(month),
		ExpYear:  year,
	}
}
