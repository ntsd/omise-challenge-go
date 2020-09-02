package omisecharge

import (
	"math/rand"
	"time"

	"github.com/ntsd/omise-challenge-go/internal/checkerror"
	"github.com/ntsd/omise-challenge-go/internal/donationstats"
	"github.com/ntsd/omise-challenge-go/internal/songpahpa"
	omise "github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

func Charge(publicKey string,
	privateKey string,
	songPahPaChannel <-chan *songpahpa.SongPahPa,
	donationStatusChannel chan<- *donationstats.DonationStatus) {

	for songPahPa := range songPahPaChannel {
		client, err := omise.NewClient(publicKey, privateKey)
		checkerror.CheckError(err)

		isSuccess := callAPI(client, songPahPa)
		donationStatusChannel <- &donationstats.DonationStatus{
			IsSuccess: isSuccess,
			SongPahPa: songPahPa,
		}
	}
}

func callAPI(client *omise.Client, songPahPa *songpahpa.SongPahPa) bool {
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	token, createToken := &omise.Token{}, &operations.CreateToken{
		Name:            songPahPa.Name,
		Number:          songPahPa.CCNumber,
		ExpirationMonth: songPahPa.ExpMonth,
		ExpirationYear:  songPahPa.ExpYear,
		SecurityCode:    songPahPa.CVV,
	}
	if e := client.Do(token, createToken); e != nil {
		checkerror.CheckError(e)
		return false
	}

	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   songPahPa.Amount,
		Currency: "thb",
		Card:     token.ID,
	}
	if e := client.Do(charge, createCharge); e != nil {
		checkerror.CheckError(e)
		return false
	}

	return charge.Paid
}
