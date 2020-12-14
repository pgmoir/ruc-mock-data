package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Xml structure
type FindCharge struct {
	Body FindChargeBody
}

type FindChargeBody struct {
	FindChargeResponse FindChargeResponse `xml:"FindChargeResponse"`
}

type FindChargeResponse struct {
	Result          FindChargeResult `xml:"Result"`
	TransactionList TransactionList
}

type FindChargeResult struct {
	Errors ChargeErrorList `xml:"Errors"`
}

type ChargeErrorList struct {
	ChargeErrors []ChargeError `xml:"Error"`
}

type ChargeError struct {
	Code        int
	Description string
}

type TransactionList struct {
	FinancialTransactions []FinancialTransaction `xml:"FinancialTransaction"`
}

type FinancialTransaction struct {
	ItemId      string
	Date        string
	ReceiptId   string
	TotalAmount float32
	Description string
}

func processCharges(charges Charges) {
	d := GetFindChargeDirectory()
	err := GetAllFindCharges(d, &charges)
	if err != nil {
		panic(err)
	}
}

func printCharges(charges Charges) {
	Logger.Println("==========================================")
	Logger.Println("")
	Logger.Println("Charges")
	Logger.Println("")

	for _, charge := range charges {
		if charge.Account != "" {
			if len(charge.ChargeErrors) > 0 {
				for _, e := range charge.ChargeErrors {
					Logger.Println("Error: ", e.Code, "/", e.Description)
				}
			} else {
				Logger.Println("Account: ", charge.Account, " / Receipt: ", charge.Receipt, " / VRM: ", charge.Vrm, " / Period: ", charge.Period)
				for _, t := range charge.FinancialTransactions {
					ta := fmt.Sprintf("%.2f", t.TotalAmount)
					taf, _ := strconv.ParseFloat(ta, 2)
					Logger.Println("  Trans: ", t.ItemId, "/", "£", taf, "/", t.Description)
				}
			}
			Logger.Println("==========================================")
			Logger.Println("")
		}
	}
}

// func contains(vrms Vrms, v Vehicle) {
// 	for _, vrm := range vrms {
// 		if vrm.Vrm == v.Vrm {
// 			if vrm.CcCharge == v.CcCharge && vrm.LezCharge == v.LezCharge && vrm.UlezCharge == v.UlezCharge {
// 				return
// 			}
// 			Logger.Println("ERROR - mismnatch", vrm.CcCharge, v.CcCharge, vrm.LezCharge, v.LezCharge, vrm.UlezCharge, v.UlezCharge)
// 			return
// 		}
// 	}
// 	Logger.Println("ERROR - this Vehicle not in Vrms", v.Vrm)
// }

//  func containsVrm(vehicles Vehicles, str string) bool {
// 	for _, v := range vehicles {
// 	   if v.Vrm == str {
// 		  return true
// 	   }
// 	}
// 	return false
//  }

func GetFindChargeDirectory() string {
	return "C:\\code\\ruc-api-tfl-gov-uk\\Source\\Presentation\\tfl.api.presentation.protectedapi\\Content\\mock\\MockChargeRepository\\FindCharge"
}

func GetAllFindCharges(d string, charges *Charges) error {
	files, err := ioutil.ReadDir(d)
	if err != nil {
		return err
	}

	i := -1
	for _, f := range files {
		switch mode := f.Mode(); {
		case mode.IsRegular():
			{
				account, receipt, vrm, period := GetChargeParts(f.Name())
				i++
				(*charges)[i].Account = account
				(*charges)[i].Receipt = receipt
				(*charges)[i].Vrm = vrm
				(*charges)[i].Period = period
				findCharge := GetFindCharge(f)
				if len(findCharge.Result.Errors.ChargeErrors) > 0 {
					(*charges)[i].ChargeErrors = findCharge.Result.Errors.ChargeErrors
				}
				(*charges)[i].FinancialTransactions = findCharge.TransactionList.FinancialTransactions
			}
		}
	}

	return nil
}

func GetChargeParts(fn string) (string, string, string, string) {
	fnp := strings.Split(fn, "-")
	account := fnp[0]
	receipt := fnp[1]
	vrm := fnp[2]
	period := fnp[3]
	return account, receipt, vrm, period
}

func GetFindCharge(f os.FileInfo) FindChargeResponse {
	fn := GetFindChargeDirectory() + "\\" + f.Name()
	data, _ := ioutil.ReadFile(fn)
	findChargeDetail := &FindCharge{}
	_ = xml.Unmarshal([]byte(data), &findChargeDetail)
	return findChargeDetail.Body.FindChargeResponse
}
