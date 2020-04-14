package main

import (
	"os"
	"io/ioutil"
	"strings"
	"strconv"
	"encoding/xml"
)

// Xml structure
type AccountXml struct {
	Body    AccountBody `xml:"Body"`
}

type AccountBody struct {
	AcctDetail AcctDetail `xml:"ViewAccountDetailResponse"`
}

type AcctDetail struct {
	AccountType string
	CurrentNoOfUsers int
	Balance float32
	ServiceList ServiceList `xml:"ServiceList"`
	PaymentMethod string
}

type ServiceList struct {
	Services []ServiceXml `xml:"Service"`
}

type ServiceXml struct {
	Name string
	ServiceType string
	ServiceStatus string
	DiscountStatus string
	DiscountValue int
	StartDate string
	EndDate string
	IsRenewable bool
}

func processAccounts(accounts Accounts) {
	ad := GetAccountsDirectory() 
	err := GetAllAccounts(ad, &accounts)
	if err != nil {
		panic(err)
	}
}

func printAccounts(accounts Accounts) {
	Logger.Println("==========================================")
	Logger.Println("")
	Logger.Println("ACCOUNTS")
	Logger.Println("")

	for _, ac := range accounts {
		if ac.Number > 0 {
		Logger.Println("Account No: ", ac.Number)
		Logger.Println("Account Type: ", ac.AccountType)
		if ac.NumberOfUsers > 0 {
			Logger.Println("Number of Users: ", ac.NumberOfUsers)
		}
		if ac.Balance > 0 {
			Logger.Println("Balance: ", ac.Balance)
		}
		if ac.PaymentMethod != "" {
			Logger.Println("PaymentMethod: ", ac.PaymentMethod)
		}
		Logger.Println(  "")

		for _, svc := range ac.Services {
			Logger.Println(  "  Service: ", svc.Name, svc.ServiceType, "(", svc.ServiceStatus, ")")
			Logger.Println(  "    Discount: ", svc.DiscountStatus, svc.DiscountValue)
			Logger.Println(  "    Dates: ", svc.StartDate, "-", svc.EndDate)
			Logger.Println(  "    Renewable: ", svc.IsRenewable)
			Logger.Println(  "")
		}

		for _, vehicle := range ac.Vehicles {
			if vehicle.Vrm != "" {
				Logger.Println("    Vehicle: ", vehicle.Vrm)
				if vehicle.Make != "" || vehicle.Model != "" || vehicle.Colour != "" {
					Logger.Println("      Make/Model/Colour: ", vehicle.Make, "/", vehicle.Model, "/", vehicle.Colour)
				}
				Logger.Println("      CC:", vehicle.CcCharge, "; LEZ:", vehicle.LezCharge, "; ULEZ:", vehicle.UlezCharge)
				Logger.Println(  "")
				for _, svc := range vehicle.Services {
					Logger.Println(  "      Service: ", svc.Name, svc.ServiceType, "(", svc.ServiceStatus, ")")
					Logger.Println(  "        Discount: ", svc.DiscountStatus, svc.DiscountValue)
					Logger.Println(  "        Dates: ", svc.StartDate, "-", svc.EndDate)
					Logger.Println(  "        Renewable: ", svc.IsRenewable)
					Logger.Println(  "")
				}
			}
		}

		Logger.Println("==========================================")
		Logger.Println("")
		}
	}
}

func GetAccountsDirectory() string {
	return "C:\\code\\ruc-api-tfl-gov-uk\\Source\\Presentation\\tfl.api.presentation.protectedapi\\Content\\mock\\MockAccountRepository\\ViewAccountDetail"
}

func GetAllAccounts(ad string, accounts *Accounts) error {
	files, err := ioutil.ReadDir(ad)
	if err != nil {
		return err
	}

	i := -1
	for _, fi := range files {
		switch mode := fi.Mode(); {
			case mode.IsDir():
				// do directory stuff
				Logger.Println("directory", fi.Name())
		    case mode.IsRegular():
				// do file stuff
				a := strings.Split(fi.Name(), "-")
				an, _ := strconv.Atoi(a[0])
				i++
				(*accounts)[i].Number = an
				acc := GetAccountDetails(fi)
				(*accounts)[i].AccountType = acc.Body.AcctDetail.AccountType
				(*accounts)[i].NumberOfUsers = acc.Body.AcctDetail.CurrentNoOfUsers
				(*accounts)[i].Balance = acc.Body.AcctDetail.Balance
				(*accounts)[i].PaymentMethod = acc.Body.AcctDetail.PaymentMethod
				for _, svc := range acc.Body.AcctDetail.ServiceList.Services {
					//Logger.Println(svc)
					newSvc := Service{ Name: svc.Name,
						ServiceType: svc.ServiceType,
						ServiceStatus: svc.ServiceStatus,
						DiscountStatus: svc.DiscountStatus,
						DiscountValue: svc.DiscountValue,
						StartDate: svc.StartDate,
						EndDate: svc.EndDate,
						IsRenewable: svc.IsRenewable }
					(*accounts)[i].Services = append((*accounts)[i].Services, newSvc)
				}
				//Logger.Println((*accounts)[i])
		    }
	}

	return nil
}

func GetAccountDetails(f os.FileInfo) *AccountXml {
	filename := f.Name()
	ad := GetAccountsDirectory()
	fn := ad + "\\" + filename
	data, _ := ioutil.ReadFile(fn)
 
	accxml := &AccountXml{}
 
	_ = xml.Unmarshal([]byte(data), &accxml)

	return accxml	
}


