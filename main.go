package main

import (
	"log"
	"os"
	"errors"
	"io"
	"net/http"
	"strings"
  	"io/ioutil"
	"strconv"
	"encoding/json"
	"time"
	"bytes"
	"mime/multipart"
	"path/filepath"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/invoice"
)

// Config
var stripeKey = "rk_";
var lexofficeKey = "";
var startdate = "2023-01-01T00:00:00"
var enddate = "2023-01-02T23:59:59"
//End Config

type Response struct {
	ID          string    `json:"id"`
	ResourceURI string    `json:"resourceUri"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
	Version     int       `json:"version"`
}

func main() {

	layout := "2006-01-02T15:04:05"
	starttime, _ := time.Parse(layout, startdate)
	endtime, _ := time.Parse(layout, enddate)
	
	stripe.Key = stripeKey
	params := &stripe.InvoiceListParams{}
	params.Filters.AddFilter("limit", "", "2")
	params.Filters.AddFilter("status", "","paid")
	//To Do - Date to Timestamp
	params.Filters.AddFilter("created[gte]", "",""+strconv.Itoa(int(starttime.Unix()))+"")
	params.Filters.AddFilter("created[lt]", "",""+strconv.Itoa(int(endtime.Unix()))+"")
	i := invoice.List(params)

	for i.Next() {
		in := i.Invoice()

		paidHelper := time.Unix(in.StatusTransitions.PaidAt, 0)
		paidAt := paidHelper.Format("2006-01-02")

		log.Println("Processing invoice: "+ in.Number + " - Date: "+ paidAt)

		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}

		dest := path + "\\invoices\\2023\\"+paidAt+ "-" + in.Number +".pdf"
		if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
			err := DownloadFile(dest, in.InvoicePDF)
			if err != nil {
				panic(err)
			}
			log.Println("Downloaded: " + in.InvoicePDF)
			uploadURL := CreateInvoice(in.Number, in.Subtotal, in.Tax, paidAt)
			UploadPdf(uploadURL, dest)
		}
	}
}


func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}


func CreateInvoice (invoicenumber string, total int64, tax int64, created string) string {
	url := "https://api.lexoffice.io/v1/vouchers"
	method := "POST"
	payload := strings.NewReader(`{
		"type": "salesinvoice",
		"voucherNumber": "`+invoicenumber+`",
		"voucherDate": "`+created+`",
		"totalGrossAmount": `+strconv.FormatFloat(float64(total)/100, 'g', 4, 64)+`,
		"totalTaxAmount": `+strconv.FormatFloat(float64(tax)/100, 'g', 4, 64)+`,
		"taxType": "gross",
		"useCollectiveContact": true,
		"remark": "Rechnung aus Stripe importiert.",
		"voucherItems": [{
			"amount": `+strconv.FormatFloat(float64(total)/100, 'g', 4, 64)+`,
			"taxAmount": `+strconv.FormatFloat(float64(tax)/100, 'g', 4, 64)+`,
			"taxRatePercent": 19,
			"categoryId": "8f8664a1-fd86-11e1-a21f-0800200c9a66"
		}]
	}`)
	client := &http.Client {}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
	   // return err
	}
	req.Header.Add("Authorization", "Bearer "+lexofficeKey+"")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	   // return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	   // return err
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {   // Parse []byte to go struct pointer
		log.Println("Can not unmarshal JSON")
	}
	return result.ResourceURI

}

func UploadPdf (voucherUrl string, dest string) {
	url := ""+voucherUrl+"/files"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open(""+dest+"")
	defer file.Close()
	part1,errFile1 := writer.CreateFormFile("file",filepath.Base(""+dest+""))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		log.Println(errFile1)
		return
	}
	err := writer.Close()
	if err != nil {
		log.Println(err)
		return
	}

	client := &http.Client {}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+lexofficeKey+"")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(body))
}