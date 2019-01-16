package mmcore_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	mmcore "github.com/Andeling/MMCoreAPI/MMCoreGo"
)

func ExampleSession() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	fmt.Printf("Version Info: %s\n", mmc.VersionInfo())
	fmt.Printf("API Version Info: %s\n", mmc.APIVersionInfo())
	// Output:
	// Version Info: MMCore version 8.6.0
	// API Version Info: Device API version 68, Module API version 10
}

func ExampleSession_SnapImage() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	// Set the search path for device adapters
	// MMCore will use "mmgr_dal_DemoCamera.dll" when we load a device with DemoCamera module.
	mmc.SetDeviceAdapterSearchPaths([]string{"C:\\Program Files\\Micro-Manager-1.4"})

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
	mmc.SetDeviceAdapterSearchPaths([]string{"C:\\Program Files\\Micro-Manager-1.4"})

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
			_, err := mmc.GetLastImage()
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
	path := []string{"C:\\Program Files\\Micro-Manager-1.4"}
	mmc.SetDeviceAdapterSearchPaths(path)

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
	// DeviceAdapterSearchPaths: [C:\Program Files\Micro-Manager-1.4]
	// 177 Device Adapters have been discovered. Including DemoCamera.
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
	mmc.SetDeviceAdapterSearchPaths([]string{"C:\\Program Files\\Micro-Manager-1.4"})

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
	//          Photon Conversion Factor|                     "1.0000"|         |        [0.4, 4]|
	//                         PixelType|                       "8bit"|         |                |["16bit" "32bit" "32bitRGB" "64bitRGB" "8bit"]
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
