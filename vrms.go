package main

import (
	"os"
	"io/ioutil"
	"strings"
	"encoding/xml"
)

// Xml structure
type VrmXml struct {
	Body    VrmBody
}

type VrmBody struct {
	VrmLookup VrmLookup `xml:"VrmLookupResponse"`
}

type VrmLookup struct {
	VehicleDetails VehicleDetails
}

type VehicleDetails struct {
	VRM string
	Make string
	Model string
	Colour string
	Chargeability Chargeability
	InAutoPay bool
	IsCc100PcDiscounted bool
	IsUlez100PcDiscounted bool
	IsULEZExempt int
	ULEZVehicleListType string
	IsULEZNonChargeable int
}

type Chargeability struct {
	IsCcChargeable int
	IsLezChargeable int
	IsUlezChargeable int
}

// Local structure
type Vrm struct {
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
	HasViewVehicle bool
}

func processVrms(vrms Vrms) {
	vd := GetVrmsDirectory() 
	verr := GetAllVrms(vd, &vrms)
	if verr != nil {
		panic(verr)
	}
}

func printVrms(vrms Vrms) {
	Logger.Println("==========================================")
	Logger.Println("")
	Logger.Println("VRMS")
	Logger.Println("")

	for _, vrm := range vrms {
		if vrm.Vrm != "" {
			Logger.Println("Vrm: ", vrm.Vrm)
			if vrm.Vrm != vrm.VrmInt {
				Logger.Println("ERROR - mismatch of VRMs", vrm.VrmInt)
			}
			if vrm.Make != "" || vrm.Model != "" || vrm.Colour != "" {
				Logger.Println("Make/Model/Colour: ", vrm.Make, "/", vrm.Model, "/", vrm.Colour)
			}
			Logger.Println("CC:", vrm.CcCharge, "; LEZ:", vrm.LezCharge, "; ULEZ:", vrm.UlezCharge)

		Logger.Println("==========================================")
		Logger.Println("")
		}
	}
}

func GetVrmsDirectory() string {
	return "C:\\code\\ruc-api-tfl-gov-uk\\Source\\Presentation\\tfl.api.presentation.protectedapi\\Content\\mock\\MockVehicleRepository\\VrmLookup"
}

func GetAllVrms(vd string, vrms *Vrms) error {
	files, err := ioutil.ReadDir(vd)
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
				v := strings.Split(fi.Name(), "-")
				vn := v[0]
				i++
				(*vrms)[i].Vrm = vn
				vrm := GetVrmDetails(fi)
				(*vrms)[i].VrmInt = vrm.Body.VrmLookup.VehicleDetails.VRM
				(*vrms)[i].Make = vrm.Body.VrmLookup.VehicleDetails.Make
				(*vrms)[i].Model = vrm.Body.VrmLookup.VehicleDetails.Model
				(*vrms)[i].Colour = vrm.Body.VrmLookup.VehicleDetails.Colour

				if vrm.Body.VrmLookup.VehicleDetails.Chargeability.IsCcChargeable == 0 {
					(*vrms)[i].CcCharge = "-"
				}
				if vrm.Body.VrmLookup.VehicleDetails.Chargeability.IsCcChargeable == 1 {
					(*vrms)[i].CcCharge = "Yes"
				}
				if vrm.Body.VrmLookup.VehicleDetails.Chargeability.IsCcChargeable > 1 {
					(*vrms)[i].CcCharge = "CHECK"
				}

				if vrm.Body.VrmLookup.VehicleDetails.Chargeability.IsLezChargeable == 0 {
					(*vrms)[i].LezCharge = "-"
				}
				if vrm.Body.VrmLookup.VehicleDetails.Chargeability.IsLezChargeable == 1 {
					(*vrms)[i].LezCharge = "Low"
				}
				if vrm.Body.VrmLookup.VehicleDetails.Chargeability.IsLezChargeable == 3 {
					(*vrms)[i].LezCharge = "Medium"
				}
				if vrm.Body.VrmLookup.VehicleDetails.Chargeability.IsLezChargeable == 2 {
					(*vrms)[i].LezCharge = "High"
				}
				if vrm.Body.VrmLookup.VehicleDetails.Chargeability.IsLezChargeable > 3 {
					(*vrms)[i].LezCharge = "CHECK"
				}
				
				if vrm.Body.VrmLookup.VehicleDetails.Chargeability.IsUlezChargeable == 0 {
					(*vrms)[i].UlezCharge = "-"
				}
				if vrm.Body.VrmLookup.VehicleDetails.Chargeability.IsUlezChargeable == 1 {
					(*vrms)[i].UlezCharge = "Low"
				}
				if vrm.Body.VrmLookup.VehicleDetails.Chargeability.IsUlezChargeable == 3 {
					(*vrms)[i].UlezCharge = "Medium"
				}
				if vrm.Body.VrmLookup.VehicleDetails.Chargeability.IsUlezChargeable == 2 {
					(*vrms)[i].UlezCharge = "High"
				}
				if vrm.Body.VrmLookup.VehicleDetails.Chargeability.IsUlezChargeable > 3 {
					(*vrms)[i].UlezCharge = "CHECK"
				}
				
		    }
	}

	return nil
}

func GetVrmDetails(f os.FileInfo) *VrmXml {
	filename := f.Name()
	vd := GetVrmsDirectory()
	fn := vd + "\\" + filename
	data, _ := ioutil.ReadFile(fn)
 
	vrmxml := &VrmXml{}
 
	_ = xml.Unmarshal([]byte(data), &vrmxml)
 
	// fmt.Println(accxml.Body.AcctDetail.AccountType)
	// fmt.Println(accxml.Body.AcctDetail.CurrentNoOfUsers)
	// fmt.Println(accxml.Body.AcctDetail.ServiceList)

	return vrmxml
	
}
