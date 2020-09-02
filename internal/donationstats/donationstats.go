package donationstats

import (
	"fmt"

	"github.com/ntsd/omise-challenge-go/internal/songpahpa"
)

type DonationStatus struct {
	SongPahPa *songpahpa.SongPahPa
	IsSuccess bool
}

type DonationStats struct {
	TotalAmount   int64
	SuccessAmount int64
	FailAmount    int64
	TopDonations  []*songpahpa.SongPahPa
	Count         uint
}

func CalculateDonationStats(donationStats *DonationStats, donationStatusChannel <-chan *DonationStatus) {
	for donationStatus := range donationStatusChannel {
		donationStats.Count++

		donationStats.TotalAmount += donationStatus.SongPahPa.Amount
		if donationStatus.IsSuccess {
			donationStats.SuccessAmount += donationStatus.SongPahPa.Amount
		} else {
			donationStats.FailAmount += donationStatus.SongPahPa.Amount
		}

		fmt.Println(donationStats)
	}
}
