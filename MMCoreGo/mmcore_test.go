package mmcore_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	mmcore "github.com/Andeling/MMCoreAPI/MMCoreGo"
)

const microManagerInstallPath = "C:\\Program Files\\Micro-Manager-2.0"

func ExampleSession() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	fmt.Printf("Version Info: %s\n", mmc.VersionInfo())
	fmt.Printf("API Version Info: %s\n", mmc.APIVersionInfo())
	// Output:
	// Version Info: MMCore version 10.2.0
	// API Version Info: Device API version 70, Module API version 10
}

func ExampleSession_SnapImage() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	// Set the search path for device adapters
	// MMCore will use "mmgr_dal_DemoCamera.dll" when we load a device with DemoCamera module.
	mmc.SetDeviceAdapterSearchPaths([]string{microManagerInstallPath})

	// MMCore refers a device by the label name.
	cameraLabel := "Camera"

	// Load device "DCam" with "DemoCamera" module, and assign the label name.
	err := mmc.LoadDevice(cameraLabel, "DemoCamera", "DCam")
	if err != nil {
		log.Fatal(err)
	}

	//
	err = mmc.InitializeAllDevices()
	if err != nil {
		log.Fatal(err)
	}

	// Set the camera as default camera device in the session.
	// SnapImage() and StartContinuousSequenceAcquisition() can only use the default camera.
	// The MMCore C++ API allows StartSequenceAcquisition to use non-default camera,
	// but that is not implemented in the MMCoreC and MMCoreGo API.
	// To access multiple cameras, just create a session for each of the cameras.
	err = mmc.SetCameraDevice(cameraLabel)
	if err != nil {
		log.Fatal(err)
	}

	// Get the parameters which we will need to interprete the raw image data.
	// These will not change without a SetProperty.
	width := mmc.ImageWidth()
	height := mmc.ImageHeight()
	bytesPerPixel := mmc.BytesPerPixel()
	bitDepth := mmc.ImageBitDepth()
	fmt.Printf("width=%d, height=%d, bytesPerPixel=%d, bitDepth=%d\n", width, height, bytesPerPixel, bitDepth)

	// Set exposure time to 100 ms
	err = mmc.SetExposureTime(100)
	if err != nil {
		log.Fatal(err)
	}

	// Get the acutal exposure time
	// When controlling an actual hardware,
	// this number will not be exactly the same as the requested exposure time,
	// and it will be some multiply of the internal clock of the camera.
	// It is good to get and save the actual exposure time when acquiring an image.
	exposure, err := mmc.ExposureTime()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Exposure time: %.6f ms\n", exposure)

	// Start the exposure and wait for the exposure to finish
	err = mmc.SnapImage()
	if err != nil {
		log.Fatal(err)
	}

	// Wait for the read-out and data transfering and get the image data.
	// This may take a while in the case of an actual camera.
	// You will not be able to start another exposure with SnapImage in the meantime.
	//
	buf, err := mmc.GetImage()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("len(buf)=%d\n", len(buf))

	// Output:
	// width=512, height=512, bytesPerPixel=1, bitDepth=8
	// Exposure time: 100.000000 ms
	// len(buf)=262144
}

func ExampleSession_StartContinuousSequenceAcquisition() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	// Set the search path for device adapters
	// MMCore will use "mmgr_dal_DemoCamera.dll" when we load a device with DemoCamera module.
	mmc.SetDeviceAdapterSearchPaths([]string{microManagerInstallPath})

	// MMCore refers a device by the label name.
	cameraLabel := "Camera"

	// Load device "DCam" with "DemoCamera" module, and assign the label name.
	err := mmc.LoadDevice(cameraLabel, "DemoCamera", "DCam")
	if err != nil {
		log.Fatal(err)
	}

	//
	err = mmc.InitializeAllDevices()
	if err != nil {
		log.Fatal(err)
	}

	// Set the camera as default camera device in the session.
	// SnapImage() and StartContinuousSequenceAcquisition() can only use the default camera.
	// The MMCore C++ API allows StartSequenceAcquisition to use non-default camera,
	// but that is not implemented in the MMCoreC and MMCoreGo API.
	// To access multiple cameras, just create a session for each of the cameras.
	err = mmc.SetCameraDevice(cameraLabel)
	if err != nil {
		log.Fatal(err)
	}

	// Get the parameters which we will need to interprete the raw image data.
	// These will not change without a SetProperty.
	width := mmc.ImageWidth()
	height := mmc.ImageHeight()
	bytesPerPixel := mmc.BytesPerPixel()
	bitDepth := mmc.ImageBitDepth()
	fmt.Printf("width=%d, height=%d, bytesPerPixel=%d, bitDepth=%d\n", width, height, bytesPerPixel, bitDepth)

	// Set exposure time to 100 ms
	err = mmc.SetExposureTime(100)
	if err != nil {
		log.Fatal(err)
	}

	// Get the acutal exposure time
	// When controlling an actual hardware,
	// this number will not be exactly the same as the requested exposure time,
	// and it will be some multiply of the internal clock of the camera.
	// It is good to get and save the actual exposure time when acquiring an image.
	exposure, err := mmc.ExposureTime()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Exposure time: %.6f ms\n", exposure)

	err = mmc.StartContinuousSequenceAcquisition(0)
	if err != nil {
		log.Fatal(err)
	}

	n_images := 0
	for {
		if mmc.GetRemainingImageCount() > 0 {
			_, err := mmc.PopNextImage()
			if err != nil {
				mmc.StopSequenceAcquisition()
				log.Fatal(err)
			}
			n_images++
			if n_images == 10 {
				break
			}
		}
	}

	err = mmc.StopSequenceAcquisition()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Finished acquiring %d images with ContinuousSequenceAcquisition.\n", n_images)
	// Output:
	// width=512, height=512, bytesPerPixel=1, bitDepth=8
	// Exposure time: 100.000000 ms
	// Finished acquiring 10 images with ContinuousSequenceAcquisition.
}

func ExampleSession_GetAvailableDevices() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	// Set the search path for device adapters
	mmc.SetDeviceAdapterSearchPaths([]string{microManagerInstallPath})

	// This function is not very useful, but you can get the search path.
	fmt.Printf("DeviceAdapterSearchPaths: %v\n", mmc.DeviceAdapterSearchPaths())

	// List the discovered device adapter modules.
	names, err := mmc.GetDeviceAdapterNames()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d Device Adapters have been discovered.", len(names))

	// We can confirm that DemoCamera module has been discovered
	var discovered bool
	for _, name := range names {
		if name == "DemoCamera" {
			discovered = true
			break
		}
	}
	if !discovered {
		log.Fatal("DemoCamera is not discovered")
	}
	fmt.Println(" Including DemoCamera.")

	// Get available devices from the module.
	dev_names, err := mmc.GetAvailableDevices("DemoCamera")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("AvailableDevices from DemoCamera:\n  %v\n", dev_names)

	// Get available device descriptions from the module.
	descriptions, err := mmc.GetAvailableDeviceDescriptions("DemoCamera")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("AvailableDeviceDescriptions from DemoCamera:\n")
	for _, description := range descriptions {
		fmt.Printf("  %s\n", description)
	}

	// Output:
	// DeviceAdapterSearchPaths: [C:\Program Files\Micro-Manager-2.0]
	// 212 Device Adapters have been discovered. Including DemoCamera.
	// AvailableDevices from DemoCamera:
	//   [DCam DWheel DStateDevice DObjective DStage DXYStage DLightPath DAutoFocus DShutter D-DA D-DA2 DOptovar DGalvo TransposeProcessor ImageFlipX ImageFlipY MedianFilter DHub]
	// AvailableDeviceDescriptions from DemoCamera:
	//   Demo camera
	//   Demo filter wheel
	//   Demo State Device
	//   Demo objective turret
	//   Demo stage
	//   Demo XY stage
	//   Demo light path
	//   Demo auto focus
	//   Demo shutter
	//   Demo DA
	//   Demo DA-2
	//   Demo Optovar
	//   Demo Galvo
	//   TransposeProcessor
	//   ImageFlipX
	//   ImageFlipY
	//   MedianFilter
	//   DHub
}

func ExampleSession_GetProperty() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	// Set the search path for device adapters
	// MMCore will use "mmgr_dal_DemoCamera.dll" when we load a device with DemoCamera module.
	mmc.SetDeviceAdapterSearchPaths([]string{microManagerInstallPath})

	// MMCore refers a device by the label name.
	cameraLabel := "Camera"

	// Load device "DCam" with "DemoCamera" module, and assign the label name.
	err := mmc.LoadDevice(cameraLabel, "DemoCamera", "DCam")
	if err != nil {
		log.Fatal(err)
	}

	//
	err = mmc.InitializeAllDevices()
	if err != nil {
		log.Fatal(err)
	}

	property_names, err := mmc.GetDevicePropertyNames(cameraLabel)
	if err != nil {
		log.Fatal(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)

	fmt.Fprintf(w, "property_name\tvalue\t\tlimit\tallowed_values\n")
	for _, property_name := range property_names {
		value, err := mmc.GetProperty(cameraLabel, property_name)
		if err != nil {
			log.Fatal(err)
		}

		var bool_value_strs []string
		readonly, err := mmc.IsPropertyReadOnly(cameraLabel, property_name)
		if err != nil {
			log.Fatal(err)
		}
		if readonly {
			bool_value_strs = append(bool_value_strs, "readonly")
		}

		preinit, err := mmc.IsPropertyPreInit(cameraLabel, property_name)
		if err != nil {
			log.Fatal(err)
		}
		if preinit {
			bool_value_strs = append(bool_value_strs, "preinit")
		}

		sequenceable, err := mmc.IsPropertySequenceable(cameraLabel, property_name)
		if err != nil {
			log.Fatal(err)
		}
		if sequenceable {
			bool_value_strs = append(bool_value_strs, "sequenceable")
		}

		has_limit, err := mmc.HasPropertyLimits(cameraLabel, property_name)
		if err != nil {
			log.Fatal(err)
		}
		var limit_str string
		if has_limit {
			lower_limit, err := mmc.GetPropertyLowerLimit(cameraLabel, property_name)
			if err != nil {
				log.Fatal(err)
			}
			upper_limit, err := mmc.GetPropertyUpperLimit(cameraLabel, property_name)
			if err != nil {
				log.Fatal(err)
			}
			limit_str = fmt.Sprintf("[%g, %g]", lower_limit, upper_limit)
		}

		allowed_values, err := mmc.GetAllowedPropertyValues(cameraLabel, property_name)
		if err != nil {
			log.Fatal(err)
		}
		var allowed_values_str string
		if len(allowed_values) > 0 {
			allowed_values_str = fmt.Sprintf("%q", allowed_values)
		}

		fmt.Fprintf(w, "%s\t%q\t%s\t%s\t%s\n", property_name, value, strings.Join(bool_value_strs, " "), limit_str, allowed_values_str)
	}

	w.Flush()
	// Output:
	//                     property_name|                        value|         |           limit|allowed_values
	//                     AllowMultiROI|                          "0"|         |                |["0" "1"]
	//              AsyncPropertyDelayMS|                       "2000"|         |                |
	//             AsyncPropertyFollower|                           ""| readonly|                |
	//               AsyncPropertyLeader|                           ""|         |                |
	//                           Binning|                          "1"|         |                |["1" "2" "4" "8"]
	//                          BitDepth|                          "8"|         |                |["10" "12" "14" "16" "32" "8"]
	//                    CCDTemperature|                     "0.0000"|         |      [-100, 10]|
	//                 CCDTemperature RO|                     "0.0000"| readonly|                |
	//                          CameraID|                       "V1.0"| readonly|                |
	//                        CameraName|       "DemoCamera-MultiMode"| readonly|                |
	//                       Description| "Demo Camera Device Adapter"| readonly|                |
	//                DisplayImageNumber|                          "0"|         |                |["0" "1"]
	//                        DropPixels|                          "0"|         |                |["0" "1"]
	//                          Exposure|                    "10.0000"|         |      [0, 10000]|
	//                         FastImage|                          "0"|         |                |["0" "1"]
	//  FractionOfPixelsToDropOrSaturate|                     "0.0020"|         |        [0, 0.1]|
	//                              Gain|                          "0"|         |         [-5, 8]|
	//                             HubID|                           ""| readonly|                |
	//                 MaximumExposureMs|                 "10000.0000"|  preinit|                |
	//                              Mode|           "Artificial Waves"|         |                |["Artificial Waves" "Color Test Pattern" "Noise"]
	//                 MultiROIFillValue|                          "0"|         |      [0, 65536]|
	//                              Name|                       "DCam"| readonly|                |
	//                            Offset|                          "0"|         |                |
	//                  OnCameraCCDXSize|                        "512"|         |                |
	//                  OnCameraCCDYSize|                        "512"|         |                |
	//          Photon Conversion Factor|                     "1.0000"|         |      [0.01, 10]|
	//                       Photon Flux|                    "50.0000"|         |       [2, 5000]|
	//                         PixelType|                       "8bit"|         |                |["16bit" "32bit" "32bitRGB" "64bitRGB" "8bit"]
	//             ReadNoise (electrons)|                     "2.5000"|         |      [0.25, 50]|
	//                       ReadoutTime|                     "0.0000"|         |                |
	//                      RotateImages|                          "0"|         |                |["0" "1"]
	//                    SaturatePixels|                          "0"|         |                |["0" "1"]
	//                          ScanMode|                          "1"|         |                |["1" "2" "3"]
	//                     SimulateCrash|                           ""|         |                |["" "Dereference Null Pointer" "Divide by Zero"]
	//                       StripeWidth|                     "1.0000"|         |         [0, 10]|
	//                     TestProperty1|                     "0.0000"|         |     [-0.1, 0.1]|
	//                     TestProperty2|                     "0.0000"|         |     [-200, 200]|
	//                     TestProperty3|                     "0.0000"|         |      [0, 0.003]|
	//                     TestProperty4|                     "0.0000"|         | [-40000, 40000]|
	//                     TestProperty5|                     "0.0000"|         |                |
	//                     TestProperty6|                     "0.0000"|         |      [0, 6e+06]|
	//               TransposeCorrection|                          "0"|         |                |["0" "1"]
	//                  TransposeMirrorX|                          "0"|         |                |["0" "1"]
	//                  TransposeMirrorY|                          "0"|         |                |["0" "1"]
	//                       TransposeXY|                          "0"|         |                |["0" "1"]
	//                     TriggerDevice|                           ""|         |                |
	//              UseExposureSequences|                         "No"|         |                |["No" "Yes"]
}

//
func ExampleSession_GetState() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	// Set the search path for device adapters
	// MMCore will use "mmgr_dal_DemoCamera.dll" when we load a device with DemoCamera module.
	mmc.SetDeviceAdapterSearchPaths([]string{microManagerInstallPath})

	// MMCore refers a device by the label name.
	wheelLabel := "DWheel"

	// Load device "DWheel" with "DemoCamera" module, and assign the label name.
	err := mmc.LoadDevice(wheelLabel, "DemoCamera", "DWheel")
	if err != nil {
		log.Fatal(err)
	}

	//
	err = mmc.InitializeAllDevices()
	if err != nil {
		log.Fatal(err)
	}

	// GetState returns the current state of the state device.
	state, err := mmc.GetState(wheelLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current State: %d\n", state)

	// NumberOfStates returns total number of states.
	n_states, err := mmc.NumberOfStates(wheelLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of States: %d\n", n_states)

	// SetState sets the state to the requested state.
	// As you noticed from the previous GetState,
	// states are numbered from 0 to n_states-1.
	err = mmc.SetState(wheelLabel, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Set state to 1")

	// Let's check that the state has been set.
	state, err = mmc.GetState(wheelLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current State: %d\n\n", state)

	// Each state can have a human-readable text label.
	// When controlling an automated microscope,
	// device adapter may have read these labels from
	// the configuration of the microscope,
	// and the labels may be the names of the objectives
	// or filter cube.
	state_labels, err := mmc.GetStateLabels(wheelLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("State Labels: %q\n", state_labels)

	// Instead of getting the current state in the form of a number,
	// we can also get the label of the current state.
	state_label, err := mmc.GetStateLabel(wheelLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current state label: %q\n", state_label)

	// DefineStateLabel gives a name to the state.
	// In this filter wheel example, we can call number 1 state "GFP".
	err = mmc.DefineStateLabel(wheelLabel, 1, "GFP")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Defined label of state 1 as \"GFP\"")

	// Let's get all the labels again to see the changes.
	state_labels, err = mmc.GetStateLabels(wheelLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("State Labels: %q\n", state_labels)

	state_label, err = mmc.GetStateLabel(wheelLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Label of current state: %q\n", state_label)

	// We can also get the state number by the label.
	state_of_label, err := mmc.GetStateFromLabel(wheelLabel, "GFP")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("State labeled as GFP is: %d\n", state_of_label)

	// Instead of setting state by the number, we can also set the state by the label.
	err = mmc.SetStateLabel(wheelLabel, "State-5")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("StateStateLabel to State-5")

	state, err = mmc.GetState(wheelLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current State: %d\n\n", state)

	err = mmc.SetStateLabel(wheelLabel, "GFP")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("StateStateLabel to GFP")

	state, err = mmc.GetState(wheelLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current State: %d\n\n", state)

	// State devices can also be controlled through
	// GetProperty and SetProperty.
	property_names, err := mmc.GetDevicePropertyNames(wheelLabel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Properties")
	for _, property_name := range property_names {
		value, err := mmc.GetProperty(wheelLabel, property_name)
		if err != nil {
			log.Fatal(err)
		}
		readonly, err := mmc.IsPropertyReadOnly(wheelLabel, property_name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %s: %s", property_name, value)
		if readonly {
			fmt.Printf(" (readonly)\n")
		} else {
			fmt.Printf("\n")
		}
	}
	fmt.Println()

	err = mmc.SetProperty(wheelLabel, "State", 4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SetProperty State=4")

	state, err = mmc.GetState(wheelLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current State: %d\n", state)

	err = mmc.SetProperty(wheelLabel, "Label", "State-7")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SetProperty Label=\"State-7\"")

	state, err = mmc.GetState(wheelLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current State: %d\n", state)

	// Output:
	// Current State: 0
	// Number of States: 10
	// Set state to 1
	// Current State: 1
	//
	// State Labels: ["State-0" "State-1" "State-2" "State-3" "State-4" "State-5" "State-6" "State-7" "State-8" "State-9"]
	// Current state label: "State-1"
	// Defined label of state 1 as "GFP"
	// State Labels: ["State-0" "GFP" "State-2" "State-3" "State-4" "State-5" "State-6" "State-7" "State-8" "State-9"]
	// Label of current state: "GFP"
	// State labeled as GFP is: 1
	// StateStateLabel to State-5
	// Current State: 5
	//
	// StateStateLabel to GFP
	// Current State: 1
	//
	// Properties
	//   ClosedPosition: 0
	//   Description: Demo filter wheel driver (readonly)
	//   HubID:  (readonly)
	//   Label: GFP
	//   Name: DWheel (readonly)
	//   State: 1
	//
	// SetProperty State=4
	// Current State: 4
	// SetProperty Label="State-7"
	// Current State: 7
}

func ExampleSession_GetPosition() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	// Set the search path for device adapters
	// MMCore will use "mmgr_dal_DemoCamera.dll" when we load a device with DemoCamera module.
	mmc.SetDeviceAdapterSearchPaths([]string{microManagerInstallPath})

	// MMCore refers a device by the label name.
	focusDriveLabel := "DStage"

	// Load device "DStage" with "DemoCamera" module, and assign the label name.
	err := mmc.LoadDevice(focusDriveLabel, "DemoCamera", "DStage")
	if err != nil {
		log.Fatal(err)
	}

	//
	err = mmc.InitializeAllDevices()
	if err != nil {
		log.Fatal(err)
	}

	property_names, err := mmc.GetDevicePropertyNames(focusDriveLabel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Properties of DStage")
	for _, property_name := range property_names {
		value, err := mmc.GetProperty(focusDriveLabel, property_name)
		if err != nil {
			log.Fatal(err)
		}
		readonly, err := mmc.IsPropertyReadOnly(focusDriveLabel, property_name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %s: %s", property_name, value)
		if readonly {
			fmt.Printf(" (readonly)\n")
		} else {
			fmt.Printf("\n")
		}
	}
	fmt.Println()

	pos, err := mmc.GetPosition(focusDriveLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Position: %g\n", pos)

	err = mmc.SetPosition(focusDriveLabel, 10.1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Set Position to 10.1\n")

	pos, err = mmc.GetPosition(focusDriveLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Position: %g\n", pos)

	err = mmc.SetRelativePosition(focusDriveLabel, -1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Set Relative Position: -1\n")

	pos, err = mmc.GetPosition(focusDriveLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Position: %g\n", pos)

	fmt.Println()
	focus_direction, err := mmc.GetFocusDirection(focusDriveLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("FocusDirection: %d\n", focus_direction)

	mmc.SetFocusDirection(focusDriveLabel, 1)
	fmt.Println("Set Focus Direction to 1")

	focus_direction, err = mmc.GetFocusDirection(focusDriveLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("FocusDirection: %d\n", focus_direction)

	// Output:
	// Properties of DStage
	//   Description: Demo stage driver (readonly)
	//   HubID:  (readonly)
	//   Name: DStage (readonly)
	//   Position: 0.0000
	//   UseSequences: No
	//
	// Position: 0
	// Set Position to 10.1
	// Position: 10.1
	// Set Relative Position: -1
	// Position: 9.1
	//
	// FocusDirection: 0
	// Set Focus Direction to 1
	// FocusDirection: 1
}

func ExampleSession_GetInstalledDevices() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	// Set the search path for device adapters
	// MMCore will use "mmgr_dal_DemoCamera.dll" when we load a device with DemoCamera module.
	mmc.SetDeviceAdapterSearchPaths([]string{microManagerInstallPath})

	// MMCore refers a device by the label name.
	hubLabel := "DHub"

	// Load device "DHub" with "DemoCamera" module, and assign the label name.
	err := mmc.LoadDevice(hubLabel, "DemoCamera", "DHub")
	if err != nil {
		log.Fatal(err)
	}

	err = mmc.InitializeAllDevices()
	if err != nil {
		log.Fatal(err)
	}

	// Discover devices under the hub
	names, err := mmc.GetInstalledDevices(hubLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("InstalledDevices under DHub: %q\n", names)

	// You can get the descriptions, although there is nothing interesting in this case.
	fmt.Printf("Descriptions:\n")
	for _, name := range names {
		descriptions, err := mmc.GetInstalledDeviceDescription(hubLabel, name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %s: %q\n", name, descriptions)
	}

	// Nothing under the hub has been loaded
	loaded, err := mmc.GetLoadedPeripheralDevices(hubLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("LoadedPeripheralDevices under DHub: %q\n", loaded)

	// Let's load a device
	err = mmc.LoadDevice("wheel", "DemoCamera", "DWheel")
	if err != nil {
		log.Fatal(err)
	}

	// and initialize it
	err = mmc.InitializeDevice("wheel")
	if err != nil {
		log.Fatal(err)
	}

	// Now we can see the loaded device under the hub
	fmt.Println()
	fmt.Println("After loading DWheel with label \"wheel\":")
	loaded, err = mmc.GetLoadedPeripheralDevices(hubLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("LoadedPeripheralDevices under DHub: %q\n", loaded)

	// We can get the hub of a device
	parent_label, err := mmc.GetParentLabel("wheel")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ParentLabel of wheel: %q\n", parent_label)

	// Not sure what this SetParentLabel really does.
	err = mmc.SetParentLabel("wheel", "random_hub")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	fmt.Println("Set ParentLabel of wheel as random_hub.")

	// You can get the label back
	parent_label, err = mmc.GetParentLabel("wheel")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ParentLabel of wheel: %q\n", parent_label)

	// Of course it does not affect InstalledDevices
	loaded, err = mmc.GetInstalledDevices(hubLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("InstalledDevices under DHub: %q\n", loaded)

	// But the device will disappear in the loaded devices list under the hub.
	loaded, err = mmc.GetLoadedPeripheralDevices(hubLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("LoadedPeripheralDevices under DHub: %q\n", loaded)

	// There is no such hub called "random_hub", so this will fail
	names, err = mmc.GetInstalledDevices("random_hub")
	if err != nil {
		fmt.Printf("GetInstalledDevices under random_hub: Err %d (%s)\n", err, err)
	} else {
		fmt.Printf("InstalledDevices under random_hub: %q\n", names)
	}

	// It does not show up under the loaded devices of the new label either.
	loaded, err = mmc.GetLoadedPeripheralDevices("random_hub")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("LoadedPeripheralDevices under random_hub: %q\n", loaded)

	// Of source, the device is really loaded, and shows up if we get all the loaded devices.
	loaded, err = mmc.GetLoadedDevices()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("LoadedDevices: %q\n", loaded)

	// Output:
	// InstalledDevices under DHub: ["DCam" "DWheel" "DStateDevice" "DObjective" "DStage" "DXYStage" "DLightPath" "DAutoFocus" "DShutter" "D-DA" "D-DA2" "DOptovar" "DGalvo" "TransposeProcessor" "ImageFlipX" "ImageFlipY" "MedianFilter"]
	// Descriptions:
	//   DCam: "N/A"
	//   DWheel: "N/A"
	//   DStateDevice: "N/A"
	//   DObjective: "N/A"
	//   DStage: "N/A"
	//   DXYStage: "N/A"
	//   DLightPath: "N/A"
	//   DAutoFocus: "N/A"
	//   DShutter: "N/A"
	//   D-DA: "N/A"
	//   D-DA2: "N/A"
	//   DOptovar: "N/A"
	//   DGalvo: "N/A"
	//   TransposeProcessor: "N/A"
	//   ImageFlipX: "N/A"
	//   ImageFlipY: "N/A"
	//   MedianFilter: "N/A"
	// LoadedPeripheralDevices under DHub: []
	//
	// After loading DWheel with label "wheel":
	// LoadedPeripheralDevices under DHub: ["wheel"]
	// ParentLabel of wheel: "DHub"
	//
	// Set ParentLabel of wheel as random_hub.
	// ParentLabel of wheel: "random_hub"
	// InstalledDevices under DHub: ["DCam" "DWheel" "DStateDevice" "DObjective" "DStage" "DXYStage" "DLightPath" "DAutoFocus" "DShutter" "D-DA" "D-DA2" "DOptovar" "DGalvo" "TransposeProcessor" "ImageFlipX" "ImageFlipY" "MedianFilter"]
	// LoadedPeripheralDevices under DHub: []
	// GetInstalledDevices under random_hub: Err 1 (generic (unspecified) error)
	// LoadedPeripheralDevices under random_hub: []
	// LoadedDevices: ["DHub" "wheel" "Core"]
}

func ExampleSession_GetXYPosition() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	// Set the search path for device adapters
	// MMCore will use "mmgr_dal_DemoCamera.dll" when we load a device with DemoCamera module.
	mmc.SetDeviceAdapterSearchPaths([]string{microManagerInstallPath})

	// MMCore refers a device by the label name.
	xyStageLabel := "DXYStage"

	// Load device "DXYStage" with "DemoCamera" module, and assign the label name.
	err := mmc.LoadDevice(xyStageLabel, "DemoCamera", "DXYStage")
	if err != nil {
		log.Fatal(err)
	}

	//
	err = mmc.InitializeAllDevices()
	if err != nil {
		log.Fatal(err)
	}

	property_names, err := mmc.GetDevicePropertyNames(xyStageLabel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Properties of DXYStage")
	for _, property_name := range property_names {
		value, err := mmc.GetProperty(xyStageLabel, property_name)
		if err != nil {
			log.Fatal(err)
		}
		readonly, err := mmc.IsPropertyReadOnly(xyStageLabel, property_name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %s: %s", property_name, value)
		if readonly {
			fmt.Printf(" (readonly)\n")
		} else {
			fmt.Printf("\n")
		}
	}
	fmt.Println()

	pos_x, pos_y, err := mmc.GetXYPosition(xyStageLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Position: %g, %g\n", pos_x, pos_y)

	err = mmc.SetXYPosition(xyStageLabel, 10.2, 20)
	if err != nil {
		log.Fatal(err)
	}

	pos_x, pos_y, err = mmc.GetXYPosition(xyStageLabel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Position: %g, %g\n", pos_x, pos_y)

	// Output:
	// Properties of DXYStage
	//   Description: Demo XY stage driver (readonly)
	//   HubID:  (readonly)
	//   Name: DXYStage (readonly)
	//   TransposeMirrorX: 0
	//   TransposeMirrorY: 0
	//
	// Position: -0, -0
	// Position: 10.2, 19.995
}
