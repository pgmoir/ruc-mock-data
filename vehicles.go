package main

import (
	"os"
	"io/ioutil"
	"strings"
	"strconv"
	"encoding/xml"
)

// Xml structure
type VehicleXml struct {
	Body    VehicleBody
}

type VehicleBody struct {
	ViewVehicle ViewVehicle `xml:"ViewVehicleResponse"`
}

type ViewVehicle struct {
	VehicleList VehicleList
}

type VehicleList struct {
	CompositeVehicles []CompositeVehicle `xml:"CompositeVehicle"`
}

type CompositeVehicle struct {
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
	ServiceList ServiceList `xml:"ServiceList"`
}

func processVehicles(vehicles Vehicles, vrms Vrms, accounts Accounts, accountVehicles AccountVehicles) {
	vd := GetVehiclesDirectory() 
	verr := GetAllVehicles(vd, &vehicles, &accounts, &accountVehicles)
	if verr != nil {
		panic(verr)
	}
}

func printVehicles(vehicles Vehicles, vrms Vrms, accounts Accounts, accountVehicles AccountVehicles) {
	Logger.Println("==========================================")
	Logger.Println("")
	Logger.Println("Vehicles")
	Logger.Println("")

	for _, vehicle := range vehicles {
		if vehicle.Vrm != "" {
			Logger.Println("Vehicle: ", vehicle.Vrm)
			if vehicle.Vrm != vehicle.VrmInt {
				Logger.Println("ERROR - mismatch of VRMs", vehicle.VrmInt)
			}

			contains(vrms, vehicle)

			if vehicle.Make != "" || vehicle.Model != "" || vehicle.Colour != "" {
				Logger.Println("Make/Model/Colour: ", vehicle.Make, "/", vehicle.Model, "/", vehicle.Colour)
			}
			Logger.Println("CC:", vehicle.CcCharge, "; LEZ:", vehicle.LezCharge, "; ULEZ:", vehicle.UlezCharge)
			for _, svc := range vehicle.Services {
				Logger.Println(  "")
				Logger.Println(  "  Service: ", svc.Name, svc.ServiceType, "(", svc.ServiceStatus, ")")
				Logger.Println(  "    Discount: ", svc.DiscountStatus, svc.DiscountValue)
				Logger.Println(  "    Dates: ", svc.StartDate, "-", svc.EndDate)
				Logger.Println(  "    Renewable: ", svc.IsRenewable)
			}
			Logger.Println("==========================================")
			Logger.Println("")
		}
	}

	Logger.Println("==========================================")
	Logger.Println("")
	Logger.Println("Check Vrms vs Vehicles")
	Logger.Println("This checks that the original VrmLookup response has a matching ViewVehicle response. Any missing are reported.")

	for _, vrm := range vrms {
		if vrm.Vrm != "" {
			inc := containsVrm(vehicles, vrm.Vrm)
			if !inc {
				Logger.Println("ERROR - this Vrm not in Vehicles", vrm.Vrm)
			}
		}
	}	
}

func contains(vrms Vrms, v Vehicle) {
	for _, vrm := range vrms {
		if vrm.Vrm == v.Vrm {
			if vrm.CcCharge == v.CcCharge && vrm.LezCharge == v.LezCharge && vrm.UlezCharge == v.UlezCharge {
				return
			}
			Logger.Println("ERROR - mismnatch", vrm.CcCharge, v.CcCharge, vrm.LezCharge, v.LezCharge, vrm.UlezCharge, v.UlezCharge)
			return
		}
	}
	Logger.Println("ERROR - this Vehicle not in Vrms", v.Vrm)
}

 func containsVrm(vehicles Vehicles, str string) bool {
	for _, v := range vehicles {
	   if v.Vrm == str {
		  return true
	   }
	}
	return false
 }

func GetVehiclesDirectory() string {
	return "C:\\code\\ruc-api-tfl-gov-uk\\Source\\Presentation\\tfl.api.presentation.protectedapi\\Content\\mock\\MockVehicleRepository\\ViewVehicle"
}

func GetAllVehicles(vd string, vehicles *Vehicles, accounts *Accounts, accountVehicles *AccountVehicles) error {
	files, err := ioutil.ReadDir(vd)
	if err != nil {
		return err
	}

	i := -1
	//av := -1
	for _, fi := range files {
		switch mode := fi.Mode(); {
		    case mode.IsRegular():
				// do file stuff
				v := strings.Split(fi.Name(), "-")
				vnumb, _ := strconv.Atoi(v[0])
				if vnumb == 0 && len(v) > 3 {
					vn := v[3]
					i++
					(*vehicles)[i].Vrm = vn
					vehicle := GetVehicleDetails(fi)

					(*vehicles)[i].VrmInt = vehicle.Body.ViewVehicle.VehicleList.CompositeVehicles[0].VRM
					(*vehicles)[i].Make = vehicle.Body.ViewVehicle.VehicleList.CompositeVehicles[0].Make
					(*vehicles)[i].Model = vehicle.Body.ViewVehicle.VehicleList.CompositeVehicles[0].Model
					(*vehicles)[i].Colour = vehicle.Body.ViewVehicle.VehicleList.CompositeVehicles[0].Colour
					(*vehicles)[i].CcCharge = GetCcCharge(vehicle.Body.ViewVehicle.VehicleList.CompositeVehicles[0].Chargeability.IsCcChargeable)
					(*vehicles)[i].LezCharge = GetLezUlezCharge(vehicle.Body.ViewVehicle.VehicleList.CompositeVehicles[0].Chargeability.IsLezChargeable)
					(*vehicles)[i].UlezCharge = GetLezUlezCharge(vehicle.Body.ViewVehicle.VehicleList.CompositeVehicles[0].Chargeability.IsUlezChargeable)

					if len(vehicle.Body.ViewVehicle.VehicleList.CompositeVehicles) > 0 {
						for _, svc := range vehicle.Body.ViewVehicle.VehicleList.CompositeVehicles[0].ServiceList.Services {
							newSvc := Service{ Name: svc.Name,
								ServiceType: svc.ServiceType,
								ServiceStatus: svc.ServiceStatus,
								DiscountStatus: svc.DiscountStatus,
								DiscountValue: svc.DiscountValue,
								StartDate: svc.StartDate,
								EndDate: svc.EndDate,
								IsRenewable: svc.IsRenewable }
							(*vehicles)[i].Services = append((*vehicles)[i].Services, newSvc)
						}
					}
				}
				if vnumb > 0 && v[1] == "10" && v[2] == "1" && v[3] == "Response.xml" {
					found := Find(*accounts, vnumb)
					if found > -1 {
						vehicle := GetVehicleDetails(fi)
						for _, cv := range vehicle.Body.ViewVehicle.VehicleList.CompositeVehicles {
							newVehicle := Vehicle { 
								Vrm: cv.VRM,
								Make: cv.Make,
								Model: cv.Model,
								Colour: cv.Colour,
								CcCharge: GetCcCharge(cv.Chargeability.IsCcChargeable),
								LezCharge: GetLezUlezCharge(cv.Chargeability.IsLezChargeable),
								UlezCharge: GetLezUlezCharge(cv.Chargeability.IsUlezChargeable) }
			
							if len(cv.ServiceList.Services) > 0 {
								for _, svc := range cv.ServiceList.Services {
									newSvc := Service{ Name: svc.Name,
										ServiceType: svc.ServiceType,
										ServiceStatus: svc.ServiceStatus,
										DiscountStatus: svc.DiscountStatus,
										DiscountValue: svc.DiscountValue,
										StartDate: svc.StartDate,
										EndDate: svc.EndDate,
										IsRenewable: svc.IsRenewable }
									newVehicle.Services = append(newVehicle.Services, newSvc)
								}
							}
							(*accounts)[found].Vehicles = append((*accounts)[found].Vehicles, newVehicle)
						}
					}
				}
			}
	}

	return nil
}

func Find(a Accounts, x int) int {
    for i, n := range a {
        if x == n.Number {
            return i
        }
    }
    return -1
}

func GetVehicleDetails(f os.FileInfo) *VehicleXml {
	filename := f.Name()
	vd := GetVehiclesDirectory()
	fn := vd + "\\" + filename
	data, _ := ioutil.ReadFile(fn)
 
	vehiclexml := &VehicleXml{}
	_ = xml.Unmarshal([]byte(data), &vehiclexml)
	return vehiclexml
}

func GetCcCharge(charge int) string {
	if charge == 0 {
		return "-"
	}
	if charge == 1 {
		return "Yes"
	}
	return "ERROR"
}

func GetLezUlezCharge(charge int) string {
	if charge == 0 {
		return "-"
	}
	if charge == 1 {
		return "Low"
	}
	if charge == 3 {
		return "Medium"
	}
	if charge == 2 {
		return "High"
	}
	return "ERROR"
}
