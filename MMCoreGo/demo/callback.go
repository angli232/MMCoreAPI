package main

import (
	"fmt"
	"log"
	"time"

	mmcore "github.com/Andeling/MMCoreAPI/MMCoreGo"
)

const microManagerInstallPath = "C:\\Program Files\\Micro-Manager-2.0"

func main() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	mmc.SetDeviceAdapterSearchPaths([]string{microManagerInstallPath})

	propertyChangedEvent := make(chan *mmcore.PropertyChangedEvent)
	stagePositionChangedEvents := make(chan *mmcore.StagePositionChangedEvent)
	mmc.NotifyPropertyChanged(propertyChangedEvent)
	mmc.NotifyStagePositionChanged(stagePositionChangedEvents)
	go func() {
		for {
			select {
			case event := <-propertyChangedEvent:
				fmt.Printf("%s [PropertyChangedEvent] %+v\n", time.Now().String(), event)
			case event := <-stagePositionChangedEvents:
				fmt.Printf("%s [StagePositionChangedEvents] %+v\n", time.Now().String(), event)
			}
		}
	}()
	log.Printf("Listening to events...")

	err := mmc.LoadDevice("TIScope", "NikonTI", "TIScope")
	if err != nil {
		log.Fatal(err)
	}

	err = mmc.LoadDevice("TINosePiece", "NikonTI", "TINosePiece")
	if err != nil {
		log.Fatal(err)
	}

	err = mmc.LoadDevice("TIDiaShutter", "NikonTI", "TIDiaShutter")
	if err != nil {
		log.Fatal(err)
	}

	err = mmc.LoadDevice("TIDiaLamp", "NikonTI", "TIDiaLamp")
	if err != nil {
		log.Fatal(err)
	}

	err = mmc.LoadDevice("TIFilterBlock1", "NikonTI", "TIFilterBlock1")
	if err != nil {
		log.Fatal(err)
	}

	err = mmc.LoadDevice("TILightPath", "NikonTI", "TILightPath")
	if err != nil {
		log.Fatal(err)
	}

	err = mmc.LoadDevice("TIZDrive", "NikonTI", "TIZDrive")
	if err != nil {
		log.Fatal(err)
	}

	err = mmc.LoadDevice("TIPFSOffset", "NikonTI", "TIPFSOffset")
	if err != nil {
		log.Fatal(err)
	}

	err = mmc.LoadDevice("TIPFSStatus", "NikonTI", "TIPFSStatus")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Initializing Devices...")
	err = mmc.InitializeAllDevices()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Set default devices...")
	err = mmc.SetFocusDevice("TIZDrive")
	if err != nil {
		log.Fatalf("SetFocusDevice: %v\n", err)
	}

	err = mmc.SetShutterDevice("TIDiaShutter")
	if err != nil {
		log.Fatalf("SetShutterDevice: %v\n", err)
	}

	log.Printf("Done.")

	for {
	}

}
