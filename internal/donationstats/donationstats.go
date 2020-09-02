package donationstats

import (
	"github.com/ntsd/omise-challenge-go/internal/songpahpa"
)

type DonationStatus struct {
	SongPahPa *songpahpa.SongPahPa
	IsSuccess bool
}

type TopDonation struct {
	Name   string
	Amount int64
}

type DonationStats struct {
	SuccessAmount int64
	FailAmount    int64
	TopDonations  []TopDonation
	Count         uint
}

func CalculateDonationStats(donationStatsChannel chan<- *DonationStats, donationStatusChannel <-chan *DonationStatus) {
	defer close(donationStatsChannel)

	donationStats := &DonationStats{
		SuccessAmount: 0,
		FailAmount:    0,
		TopDonations: []TopDonation{
			TopDonation{Name: "", Amount: 0},
			TopDonation{Name: "", Amount: 0},
			TopDonation{Name: "", Amount: 0},
		},
		Count: 0,
	}

	for donationStatus := range donationStatusChannel {
		donationStats.Count++

		if donationStatus.IsSuccess {
			donationStats.SuccessAmount += donationStatus.SongPahPa.Amount
		} else {
			donationStats.FailAmount += donationStatus.SongPahPa.Amount
		}

		sortTopDonations(donationStats, donationStatus)
	}
	donationStatsChannel <- donationStats
}

func sortTopDonations(donationStats *DonationStats, donationStatus *DonationStatus) {
	// To sort the top donation I assume that every SongPahPa have a unique name.
	topDonationsLength := len(donationStats.TopDonations)
	i := topDonationsLength
	for i > 0 {
		if donationStats.TopDonations[i-1].Amount > donationStatus.SongPahPa.Amount {
			break
		}
		i--
	}
	if i < topDonationsLength {
		donationStats.TopDonations = append(donationStats.TopDonations, TopDonation{})
		copy(donationStats.TopDonations[i+1:], donationStats.TopDonations[i:])
		donationStats.TopDonations[i] = TopDonation{Name: donationStatus.SongPahPa.Name, Amount: donationStatus.SongPahPa.Amount}
		donationStats.TopDonations = donationStats.TopDonations[:3]
	}
}
