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

	donationStatusChannel := make(chan *donationstats.DonationStatus)
	donationStatsChannel := make(chan *donationstats.DonationStats)
	go donationstats.CalculateDonationStats(donationStatsChannel, donationStatusChannel)

	var waitMultiCPU sync.WaitGroup
	numCPU := runtime.NumCPU()
	waitMultiCPU.Add(numCPU)
	for i := 0; i < numCPU; i++ {
		go func() {
			defer waitMultiCPU.Done()
			omisecharge.ChargeChannel(publicKey, secretKey, songPahPaChannel, donationStatusChannel)
		}()
	}
	waitMultiCPU.Wait()

	close(donationStatusChannel)
	donationStats := <-donationStatsChannel

	currency := "THB"
	totalDonation := donationStats.SuccessAmount + donationStats.FailAmount
	averageDonation := float64(totalDonation) / float64(donationStats.Count)
	fmt.Printf("%25s %s %14d.00\n", "total received:", currency, totalDonation)
	fmt.Printf("%25s %s %14d.00\n", "successfully donated:", currency, donationStats.SuccessAmount)
	fmt.Printf("%25s %s %14d.00\n", "faulty donation:", currency, donationStats.FailAmount)
	fmt.Printf("%25s %s %17.2f\n", "average per person:", currency, averageDonation)
	fmt.Printf("%25s", "top donors:\n")
	for _, topDonation := range donationStats.TopDonations {
		fmt.Printf("%25s %s\n", "", topDonation.Name)
	}

	fmt.Println("done.")
}
