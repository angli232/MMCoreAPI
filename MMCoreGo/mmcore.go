package mmcore

// #cgo CFLAGS: -I../MMCoreC
// #cgo LDFLAGS: -L../lib -lMMCoreC
//
// #include <stdlib.h>
//
// #include "MMCoreC.h"
import "C"

import (
	"unsafe"
)

type Session struct {
	mmcore C.MM_Session
}

func NewSession() *Session {
	var s Session
	C.MM_Open(&s.mmcore)
	return &s
}

func (s *Session) Close() {
	C.MM_Close(s.mmcore)
}

// VersionInfo returns the version of Micro-Manager Core.
func (s *Session) VersionInfo() string {
	var c_str *C.char
	C.MM_GetVersionInfo(s.mmcore, &c_str)
	defer C.MM_StringFree(c_str)

	str := C.GoString(c_str)
	return str
}

// APIVersionInfo returns module and device interface versions.
func (s *Session) APIVersionInfo() string {
	var c_str *C.char
	C.MM_GetAPIVersionInfo(s.mmcore, &c_str)
	defer C.MM_StringFree(c_str)

	str := C.GoString(c_str)
	return str
}

//
// Initialization and setup.
//

// LoadDevice loads a device with specified device adapter module
// and assigns a label name in the session.
func (s *Session) LoadDevice(label, module_name, dev_name string) error {
	c_label := C.CString(label)
	c_module_name := C.CString(module_name)
	c_dev_name := C.CString(dev_name)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_module_name))
	defer C.free(unsafe.Pointer(c_dev_name))

	status := C.MM_LoadDevice(s.mmcore, c_label, c_module_name, c_dev_name)
	return statusToError(status)
}

func (s *Session) UnloadDevice(label string) error {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_UnloadDevice(s.mmcore, c_label)
	return statusToError(status)
}

func (s *Session) UnloadAllDevices() error {
	status := C.MM_UnloadAllDevices(s.mmcore)
	return statusToError(status)
}

func (s *Session) InitializeAllDevices() error {
	status := C.MM_InitializeAllDevices(s.mmcore)
	return statusToError(status)
}

func (s *Session) InitializeDevice(label string) error {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_InitializeDevice(s.mmcore, c_label)
	return statusToError(status)
}

func (s *Session) Reset() error {
	status := C.MM_Reset(s.mmcore)
	return statusToError(status)
}

//
// Device listing.
//

// DeviceAdapterSearchPaths returns the search path for device adapter modules.
//
// Device adapters are usually in the installation path of Micro-Manager,
// for example, "C:\\Program Files\\Micro-Manager-1.4"
func (s *Session) DeviceAdapterSearchPaths() (paths []string) {
	var c_paths **C.char
	C.MM_GetDeviceAdapterSearchPaths(s.mmcore, &c_paths)
	defer C.MM_StringListFree(c_paths)

	paths = goStringList(c_paths)
	return
}

// SetDeviceAdapterSearchPaths sets the search path for device adapter modules.
//
// Device adapters are usually in the installation path of Micro-Manager,
// for example, "C:\\Program Files\\Micro-Manager-1.4"
func (s *Session) SetDeviceAdapterSearchPaths(paths []string) {
	c_paths := make([]*C.char, len(paths)+1)
	for i, path := range paths {
		c_paths[i] = C.CString(path)
	}
	c_paths[len(paths)] = (*C.char)(C.NULL)

	C.MM_SetDeviceAdapterSearchPaths(s.mmcore, &c_paths[0])
	for i := 0; i < len(paths); i++ {
		C.free(unsafe.Pointer(c_paths[i]))
	}
}

// GetDeviceAdapterNames returns the names of discoverable device adapter modules.
//
// The list is  constructed based on filename matching in the current search paths,
// and it does not check whether the files are valid and compatible device adapters.
func (s *Session) GetDeviceAdapterNames() (names []string, err error) {
	var c_names **C.char
	status := C.MM_GetDeviceAdapterNames(s.mmcore, &c_names)
	defer C.MM_StringListFree(c_names)

	names = goStringList(c_names)
	err = statusToError(status)
	return
}

func (s *Session) GetAvailableDevices(module_name string) (dev_names []string, err error) {
	c_module_name := C.CString(module_name)
	defer C.free(unsafe.Pointer(c_module_name))

	var c_dev_names **C.char
	status := C.MM_GetAvailableDevices(s.mmcore, c_module_name, &c_dev_names)
	defer C.MM_StringListFree(c_dev_names)

	dev_names = goStringList(c_dev_names)
	err = statusToError(status)
	return
}

func (s *Session) GetAvailableDeviceDescriptions(module_name string) (descriptions []string, err error) {
	c_module_name := C.CString(module_name)
	defer C.free(unsafe.Pointer(c_module_name))

	var c_descriptions **C.char
	status := C.MM_GetAvailableDeviceDescriptions(s.mmcore, c_module_name, &c_descriptions)
	defer C.MM_StringListFree(c_descriptions)

	descriptions = goStringList(c_descriptions)
	err = statusToError(status)
	return
}

//
// Generic device control
//

// GetDevicePropertyNames returns all property names supported by the device.
func (s *Session) GetDevicePropertyNames(label string) (names []string, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_names **C.char
	status := C.MM_GetDevicePropertyNames(s.mmcore, c_label, &c_names)
	defer C.MM_StringListFree(c_names)

	names = goStringList(c_names)
	err = statusToError(status)
	return
}

// GetProperty returns the property value of the device as a string.
func (s *Session) GetProperty(label string, property string) (value string, err error) {
	c_label := C.CString(label)
	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_property))

	var c_value *C.char
	status := C.MM_GetProperty(s.mmcore, c_label, c_property, &c_value)
	defer C.MM_StringFree(c_value)

	value = C.GoString(c_value)
	err = statusToError(status)
	return
}

func (s *Session) SetProperty(label string, property string, state interface{}) (err error) {
	c_label := C.CString(label)
	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_property))

	var status C.MM_Status
	switch state.(type) {
	case bool:
		var c_state C.uint8_t
		if state.(bool) {
			c_state = 1
		} else {
			c_state = 0
		}
		status = C.MM_SetPropertyBool(s.mmcore, c_label, c_property, c_state)
	case int:
		status = C.MM_SetPropertyInt(s.mmcore, c_label, c_property, (C.int32_t)(state.(int)))
	case float32:
		status = C.MM_SetPropertyFloat(s.mmcore, c_label, c_property, (C.float)(state.(float32)))
	case float64:
		status = C.MM_SetPropertyDouble(s.mmcore, c_label, c_property, (C.double)(state.(float64)))
	case string:
		c_state := C.CString(state.(string))
		status = C.MM_SetPropertyString(s.mmcore, c_label, c_property, c_state)
		C.free(unsafe.Pointer(c_state))
	}
	return statusToError(status)
}

//
// Manage current devices.
//

func (s *Session) SetCameraDevice(label string) error {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetCameraDevice(s.mmcore, c_label)
	return statusToError(status)
}

func (s *Session) SetShutterDevice(label string) error {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetShutterDevice(s.mmcore, c_label)
	return statusToError(status)
}

func (s *Session) SetFocusDevice(label string) error {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetFocusDevice(s.mmcore, c_label)
	return statusToError(status)
}

func (s *Session) SetXYStageDevice(label string) error {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetXYStageDevice(s.mmcore, c_label)
	return statusToError(status)
}

func (s *Session) SetAutoFocusDevice(label string) error {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetAutoFocusDevice(s.mmcore, c_label)
	return statusToError(status)
}

func (s *Session) CameraDevice() (label string) {
	var c_str *C.char
	C.MM_GetCameraDevice(s.mmcore, &c_str)
	defer C.MM_StringFree(c_str)

	label = C.GoString(c_str)
	return label
}

func (s *Session) ShutterDevice() (label string) {
	var c_str *C.char
	C.MM_GetShutterDevice(s.mmcore, &c_str)
	defer C.MM_StringFree(c_str)

	label = C.GoString(c_str)
	return label
}

func (s *Session) FocusDevice() (label string) {
	var c_str *C.char
	C.MM_GetFocusDevice(s.mmcore, &c_str)
	defer C.MM_StringFree(c_str)

	label = C.GoString(c_str)
	return label
}

func (s *Session) XYStageDevice() (label string) {
	var c_str *C.char
	C.MM_GetXYStageDevice(s.mmcore, &c_str)
	defer C.MM_StringFree(c_str)

	label = C.GoString(c_str)
	return label
}

func (s *Session) AutoFocusDevice() (label string) {
	var c_str *C.char
	C.MM_GetAutoFocusDevice(s.mmcore, &c_str)
	defer C.MM_StringFree(c_str)

	label = C.GoString(c_str)
	return label
}

//
// Image acquisition settings
//

// SetExposureTime sets the exposure time of the current camera in milliseconds.
func (s *Session) SetExposureTime(exposure_ms float64) error {
	status := C.MM_SetExposure(s.mmcore, (C.double)(exposure_ms))
	return statusToError(status)
}

// ExposureTime returns the exposure time of the current camera in milliseconds.
func (s *Session) ExposureTime() (exposure_ms float64, err error) {
	status := C.MM_GetExposure(s.mmcore, (*C.double)(&exposure_ms))
	err = statusToError(status)
	return
}

// ImageBufferSize returns the size of the image buffer.
//
// The size is consistent with values returned by ImageWidth(), ImageHeight() and ImageBytesPerPixel().
// The camera never changes the size of image buffer on its own.
// The buffer size changes only when appropriate properties are set (such as binning, pixel type, etc.)
func (s *Session) ImageBufferSize() (len int) {
	var c_len C.uint32_t
	C.MM_GetImageBufferSize(s.mmcore, &c_len)

	len = int(c_len)
	return
}

// ImageWidth returns the width of the image.
func (s *Session) ImageWidth() (width int) {
	var c_width C.uint16_t
	C.MM_GetImageWidth(s.mmcore, &c_width)

	width = int(c_width)
	return
}

// ImageWidth returns the height of the image.
func (s *Session) ImageHeight() (height int) {
	var c_height C.uint16_t
	C.MM_GetImageHeight(s.mmcore, &c_height)

	height = int(c_height)
	return
}

// ImageWidth returns number of bytes in a pixel in image buffer data.
func (s *Session) BytesPerPixel() (bytes_per_pixel int) {
	var c_bytes_per_pixel C.uint8_t
	C.MM_GetBytesPerPixel(s.mmcore, &c_bytes_per_pixel)

	bytes_per_pixel = int(c_bytes_per_pixel)
	return
}

// ImageBitDepth returns the the bit depth of the pixel to indicate the dynamic range.
//
// It does not directly affect the image buffer size, and just gives a guideline on how to interpret pixel values.
func (s *Session) ImageBitDepth() (bit_depth int) {
	var c_bit_depth C.uint8_t
	C.MM_GetImageBitDepth(s.mmcore, &c_bit_depth)

	bit_depth = int(c_bit_depth)
	return
}

// NumberOfComponents returns the number of comopnents in the image. 1 for monochrome cameras, 4 for RGB cameras.
func (s *Session) NumberOfComponents() (n_components int) {
	var c_n_components C.uint8_t
	C.MM_GetNumberOfComponents(s.mmcore, &c_n_components)

	n_components = int(c_n_components)
	return
}

// NumberOfCameraChannels returns the number of simultaneous channels that camera is capable of.
//
// This is not used by color cameras, which use NumberOfComponents() .
func (s *Session) NumberOfCameraChannels() (n_channels int) {
	var c_n_channels C.uint8_t
	C.MM_GetNumberOfCameraChannels(s.mmcore, &c_n_channels)

	n_channels = int(c_n_channels)
	return
}

//
// Image acquisition
//

// SnapImage starts the exposure of a single image and returns when the exposure is finished.
//
// It does not wait for the read-out and data transfering.
func (s *Session) SnapImage() error {
	status := C.MM_SnapImage(s.mmcore)
	return statusToError(status)
}

// GetImage returns the image buffer data.
//
// GetImage is called after SnapImage returns.
// It waits for the camera read-out and data transfering.
//
// In the case of multi-channel camera, image data of the first channel is returned.
func (s *Session) GetImage() (buf []byte, err error) {
	len := s.ImageBufferSize()

	var c_pbuf *C.uint8_t
	status := C.MM_GetImage(s.mmcore, &c_pbuf)

	buf = C.GoBytes(unsafe.Pointer(c_pbuf), C.int(len))
	err = statusToError(status)

	return
}

func (s *Session) GetImageOfChannel(channel int) (buf []byte, err error) {
	len := s.ImageBufferSize()

	var c_pbuf *C.uint8_t
	status := C.MM_GetImageOfChannel(s.mmcore, (C.uint16_t)(channel), &c_pbuf)

	buf = C.GoBytes(unsafe.Pointer(c_pbuf), C.int(len))
	err = statusToError(status)

	return
}

//
// Image sequence acquisition
//

func (s *Session) StartSequenceAcquisition(num_images int16, interval_ms float64, stop_on_overflow bool) error {
	var c_stop_on_overflow C.uint8_t
	if stop_on_overflow {
		c_stop_on_overflow = 1
	} else {
		c_stop_on_overflow = 0
	}

	status := C.MM_StartSequenceAcquisition(s.mmcore, (C.int16_t)(num_images), (C.double)(interval_ms), c_stop_on_overflow)
	return statusToError(status)
}

func (s *Session) StartContinuousSequenceAcquisition(interval_ms float64) error {
	status := C.MM_StartContinuousSequenceAcquisition(s.mmcore, (C.double)(interval_ms))
	return statusToError(status)
}

func (s *Session) StopSequenceAcquisition() error {
	status := C.MM_StopSequenceAcquisition(s.mmcore)
	return statusToError(status)
}

func (s *Session) IsSequenceRunning() bool {
	var c_status C.uint8_t
	C.MM_IsSequenceRunning(s.mmcore, &c_status)
	if uint8(c_status) == 0 {
		return false
	}
	return true
}

//
// Image circular buffer
//

// GetLastImage gets the last image from the circular buffer. It returns nil if the buffer is empty.
func (s *Session) GetLastImage() (buf []byte, err error) {
	var c_pbuf *C.uint8_t
	status := C.MM_GetLastImage(s.mmcore, &c_pbuf)

	if unsafe.Pointer(c_pbuf) == C.NULL {
		buf = nil
	} else {
		len := s.ImageBufferSize()
		buf = C.GoBytes(unsafe.Pointer(c_pbuf), C.int(len))
	}
	err = statusToError(status)
	return
}

// PopNextImage gets the removes the next image from the circular buffer. It returns nil if the buffer is empty.
func (s *Session) PopNextImage() (buf []byte, err error) {
	var c_pbuf *C.uint8_t
	status := C.MM_PopNextImage(s.mmcore, &c_pbuf)

	if unsafe.Pointer(c_pbuf) == C.NULL {
		buf = nil
	} else {
		len := s.ImageBufferSize()
		buf = C.GoBytes(unsafe.Pointer(c_pbuf), C.int(len))
	}
	err = statusToError(status)
	return
}

func (s *Session) GetRemainingImageCount() (count int) {
	var c_count C.int16_t
	C.MM_GetRemainingImageCount(s.mmcore, &c_count)
	return int(c_count)
}

func (s *Session) GetBufferTotalCapacity() (capacity int) {
	var c_capacity C.int16_t
	C.MM_GetBufferTotalCapacity(s.mmcore, &c_capacity)
	return int(c_capacity)
}

func (s *Session) GetBufferFreeCapacity() (capacity int) {
	var c_capacity C.int16_t
	C.MM_GetBufferFreeCapacity(s.mmcore, &c_capacity)
	return int(c_capacity)
}

func (s *Session) IsBufferOverflowed() (overflowed bool) {
	var c_overflowed C.uint8_t
	C.MM_IsBufferOverflowed(s.mmcore, &c_overflowed)
	return goBool(c_overflowed)
}

func (s *Session) SetCircularBufferMemoryFootprint(size_MB uint32) error {
	status := C.MM_SetCircularBufferMemoryFootprint(s.mmcore, (C.uint32_t)(size_MB))
	return statusToError(status)
}

func (s *Session) GetCircularBufferMemoryFootprint() (size_MB uint32) {
	var c_size_MB C.uint32_t
	C.MM_GetCircularBufferMemoryFootprint(s.mmcore, &c_size_MB)
	return uint32(c_size_MB)
}

func (s *Session) InitializeCircularBuffer() error {
	status := C.MM_InitializeCircularBuffer(s.mmcore)
	return statusToError(status)
}

func (s *Session) ClearCircularBuffer() error {
	status := C.MM_ClearCircularBuffer(s.mmcore)
	return statusToError(status)
}

//
// State device control.
//

func (s *Session) SetState(label string, state int) error {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetState(s.mmcore, c_label, (C.int32_t)(state))
	return statusToError(status)
}

func (s *Session) GetState(label string) (state int, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_state C.int32_t
	status := C.MM_GetState(s.mmcore, c_label, &c_state)

	state = int(c_state)
	err = statusToError(status)
	return
}

func (s *Session) NumberOfStates(label string) (n_states int, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_n_states C.int32_t
	status := C.MM_GetNumberOfStates(s.mmcore, c_label, &c_n_states)

	n_states = int(c_n_states)
	err = statusToError(status)
	return
}

//
// Hub and peripheral devices.
//

//
// Helper function
//

// goStringList converts a NULL terminated array of strings to []string in Go.
func goStringList(c_str_list **C.char) []string {
	strs := make([]string, 0)

	c_str_slice := (*[1 << 30]*C.char)(unsafe.Pointer(c_str_list))
	for _, c_str := range c_str_slice {
		if unsafe.Pointer(c_str) == C.NULL {
			break
		}
		str := C.GoString(c_str)
		strs = append(strs, str)
	}
	return strs
}

func goBool(c_bool C.uint8_t) bool {
	if c_bool != 0 {
		return true
	}
	return false
}
