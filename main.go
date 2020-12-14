package main

import (
	"fmt"
	"os"
	"time"
)

type Accounts []Account
type Vrms []Vrm
type Vehicles []Vehicle
type AccountVehicles []Vehicle
type Pcns []Pcn
type Charges []Charge

func main() {
	st := time.Now()

	// set default slice for storing accounts to 100 - theres about 50 so just giving a little bit of breathing capacity
	accounts := make(Accounts, 100)
	vrms := make(Vrms, 200)
	vehicles := make(Vehicles, 200)
	accountVehicles := make(AccountVehicles, 300)
	pcns := make(Pcns, 50)
	charges := make(Charges, 50)
	
	LogInit(os.Stdout)

	Logger.Println("analysing mock data")
	Logger.Println("")
	Logger.Println("The purpose of this log file is to highlight the mock data available for testing.")

	processAccounts(accounts)
	processVrms(vrms)
	processVehicles(vehicles, vrms, accounts, accountVehicles)
	processPcns(pcns)
	processCharges(charges)
	printAccounts(accounts)
	printVrms(vrms)
	printVehicles(vehicles, vrms, accounts, accountVehicles)
	printPcns(pcns)
	printCharges(charges)

	ft := time.Now()
	fmt.Println("Complete", int(ft.Sub(st).Seconds()))
}


