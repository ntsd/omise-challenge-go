package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	"github.com/ntsd/omise-challenge-go/internal/checkerror"
	"github.com/ntsd/omise-challenge-go/internal/donationstats"
	"github.com/ntsd/omise-challenge-go/internal/omisecharge"
	"github.com/ntsd/omise-challenge-go/internal/songpahpa"
	"github.com/omiselabs/challenges/challenge-go/cipher"
)

func main() {
	fmt.Println("performing donations...")

	var wg sync.WaitGroup

	filePath := os.Args[1]
	publicKey := os.Getenv("OMISE_PUBLIC_KEY")
	secretKey := os.Getenv("OMISE_SECRET_KEY")

	data, err := os.Open(filePath)
	checkerror.CheckError(err)

	reader, err := cipher.NewRot128Reader(data)
	checkerror.CheckError(err)

	songPahPaChannel := make(chan *songpahpa.SongPahPa)

	songPahPa := songpahpa.SongPahPaCSVReader(reader, songPahPaChannel)
	go songPahPa.Read()

	donationStatus := make(chan *donationstats.DonationStatus)
	donationStats := &donationstats.DonationStats{
		TotalAmount:   0,
		SuccessAmount: 0,
		FailAmount:    0,
		TopDonations:  []*songpahpa.SongPahPa{},
		Count:         0,
	}
	go donationstats.CalculateDonationStats(donationStats, donationStatus)

	numCPU := runtime.NumCPU()
	wg.Add(numCPU)
	for i := 0; i < numCPU; i++ {
		go func() {
			defer wg.Done()
			omisecharge.ChargeChannel(publicKey, secretKey, songPahPaChannel, donationStatus)
		}()
	}

	wg.Wait()
	close(donationStatus)
	fmt.Println("done.")
}
