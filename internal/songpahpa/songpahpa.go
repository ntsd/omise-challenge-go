package songpahpa

import "time"

type SongPahPa struct {
	Name     string
	Amount   int64
	CCNumber string
	CVV      string
	ExpMonth time.Month
	ExpYear  int
}
