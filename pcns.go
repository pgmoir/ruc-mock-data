package main

import (
	"os"
	"io/ioutil"
	"strings"
	"encoding/xml"
)

// Xml structure
type ViewPcnDetailImagesXml struct {
	Body    ViewPcnDetailImagesBodyXml
}

type ViewPcnDetailImagesBodyXml struct {
	ViewPcnDetailXml ViewPcnDetail `xml:"ViewPCNDetailResponse"`
}

type ViewPcnDetail struct {
	PCNNumber string
	ContraventionDescription string
	ContraventionType string
	PCNStatus string
	AmountDue float32
	PCNFlags PcnFlags
	Vehicle PcnVehicle
	ContraventionDetailsList ContraventionDetailsList
}

type PcnFlags struct {
	OnHold bool
	Open bool
	FullyPaid bool
	ProgressionStage string
	Cancelled bool
	CanBePaid bool 
	CanNotBePaidReason bool
	CanMakeInformalRep bool
	CanMakeRep bool
	CanMakeLateRep bool
	CanMakeAppeal bool
	CanMakeWitnessStatement bool
	CanMakeStatDec bool
	HasOpenRep bool
	HasPreviousClosedRep bool
	HasOpenAppeal bool
	AdditionalEvidenceRequested bool
}

type PcnVehicle struct {
	VRM string
	Make string
	Model string
	Colour string
	Type string
}

type ContraventionDetailsList struct {
	ContraventionDetailItems []ContraventionDetails `xml:"ContraventionDetails"`
}

type ContraventionDetails struct {
	ContraventionId string
	Location string
	DateTime string
}

func processPcns(pcns Pcns) {
	d := GetViewPcnDirectory() 
	err := GetAllViewPcns(d, &pcns)
	if err != nil {
		panic(err)
	}
}

func printPcns(pcns Pcns) {
	Logger.Println("==========================================")
	Logger.Println("")
	Logger.Println("Pcns")
	Logger.Println("")

	for _, pcn := range pcns {
		 if pcn.Number != "" {
			Logger.Println("PCN: ", pcn.Number, "/", pcn.Vrm)
			Logger.Println("Status: ", pcn.Status)
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

func GetViewPcnDirectory() string {
	return "C:\\code\\ruc-api-tfl-gov-uk\\Source\\Presentation\\tfl.api.presentation.protectedapi\\Content\\mock\\MockVehicleRepository\\ViewVehicle"
}

func GetAllViewPcns(d string, pcns *Pcns) error {
	files, err := ioutil.ReadDir(d)
	if err != nil {
		return err
	}

	i := -1
	for _, f := range files {
		switch mode := f.Mode(); {
			case mode.IsRegular(): {
				pcn, vrm, ftype := GetPcnParts(f.Name())
				if ftype == "images" {
					i++
					(*pcns)[i].Number = pcn
					(*pcns)[i].Vrm = vrm
					viewPcnDetail := GetViewPcnDetail(f)
					(*pcns)[i].Status = viewPcnDetail.PCNStatus
					(*pcns)[i].PcnFlags = viewPcnDetail.PCNFlags
				}
			}
		}
	}

	return nil
}

func GetPcnParts(fn string) (string, string, string) {
	fnp := strings.Split(fn, "-")
	pcn := fnp[0]
	vrm := fnp[1]
	ftype := fnp[2]
	return pcn, vrm, ftype
}


func GetViewPcnDetail(f os.FileInfo) ViewPcnDetail {
	fn := GetViewPcnDirectory() + "\\" + f.Name()
	data, _ := ioutil.ReadFile(fn)
	viewPcnDetailImagesXml := &ViewPcnDetailImagesXml{}
	_ = xml.Unmarshal([]byte(data), &viewPcnDetailImagesXml)
	return viewPcnDetailImagesXml.Body.ViewPcnDetailXml
}