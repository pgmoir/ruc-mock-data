package main

type Account struct {
	Number int
	AccountType string
	NumberOfUsers int
	Balance float32
	Services []Service
	PaymentMethod string
	Vehicles []Vehicle
}

type Vehicle struct {
	Vrm string
	VrmInt string
	Make string
	Model string
	Colour string
	InAutoPay bool
	IsCc100PcDiscounted bool
	IsUlez100PcDiscounted bool
	IsULEZExempt int
	ULEZVehicleListType string
	IsULEZNonChargeable int
	CcCharge string
	LezCharge string
	UlezCharge string
	Services []Service
}

type Service struct {
	Name string
	ServiceType string
	ServiceStatus string
	DiscountStatus string
	DiscountValue int
	StartDate string
	EndDate string
	IsRenewable bool
}

// type Pcn struct {	
// 	Number string
// 	Vrm string
// 	Status string
// 	PcnFlags PcnFlags
// }
