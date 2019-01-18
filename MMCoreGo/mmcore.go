// Package mmcore provides Go interface to Micro-Manager Core API for automated microscopy.
package mmcore

// #cgo CFLAGS: -I../MMCoreC
// #cgo LDFLAGS: -L../lib -lMMCoreC
//
// #include <stdlib.h>
//
// #include "MMCoreC.h"
//
// void c_registerCallback(MM_Session mm);
// extern void onPropertyChanged(MM_Session mm, char *label, char *property, char *value);
// extern void onStagePositionChanged(MM_Session mm, char *label, double pos);
import "C"

import (
	"sync"
	"unsafe"
)

var go_session map[C.MM_Session]*Session

type Session struct {
	mmcore C.MM_Session

	// Events
	c_callback_registered bool
	propertyChanged       []chan<- *PropertyChangedEvent
	stagePositionChanged  []chan<- *StagePositionChangedEvent
}

func NewSession() *Session {
	var s Session
	C.MM_Open(&s.mmcore)
	if go_session == nil {
		go_session = make(map[C.MM_Session]*Session)
	}
	go_session[s.mmcore] = &s
	return &s
}

func (s *Session) Close() {
	delete(go_session, s.mmcore)
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
// Event notification
//

type PropertyChangedEvent struct {
	label    string
	property string
	value    string
}

type StagePositionChangedEvent struct {
	label string
	pos   float64
}

func (s *Session) NotifyPropertyChanged(event chan<- *PropertyChangedEvent) {
	// Register the callback on the C side, if not already done.
	if !s.c_callback_registered {
		C.c_registerCallback(s.mmcore)
		s.c_callback_registered = true
	}

	// Save the channel for sending events.
	s.propertyChanged = append(s.propertyChanged, event)
}

func (s *Session) NotifyStagePositionChanged(event chan<- *StagePositionChangedEvent) {
	// Register the callback on the C side, if not already done.
	if !s.c_callback_registered {
		C.c_registerCallback(s.mmcore)
		s.c_callback_registered = true
	}

	// Save the channel for sending events.
	s.stagePositionChanged = append(s.stagePositionChanged, event)
}

//export onPropertyChanged
func onPropertyChanged(mmcore C.MM_Session, label *C.char, property *C.char, value *C.char) {
	event := &PropertyChangedEvent{
		label:    C.GoString(label),
		property: C.GoString(property),
		value:    C.GoString(value),
	}

	// Notify the listeners
	var wg sync.WaitGroup
	for _, ch := range go_session[mmcore].propertyChanged {
		wg.Add(1)
		go func() {
			ch <- event
			wg.Done()
		}()
	}
	wg.Wait()
}

//export onStagePositionChanged
func onStagePositionChanged(mmcore C.MM_Session, label *C.char, pos C.double) {
	event := &StagePositionChangedEvent{
		label: C.GoString(label),
		pos:   float64(pos),
	}

	// Notify the listeners
	var wg sync.WaitGroup
	for _, ch := range go_session[mmcore].stagePositionChanged {
		wg.Add(1)
		go func() {
			ch <- event
			wg.Done()
		}()
	}
	wg.Wait()
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

func (s *Session) HasProperty(label string, property string) (has_property bool, err error) {
	c_label := C.CString(label)
	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_property))

	var c_has_property C.uint8_t
	status := C.MM_HasProperty(s.mmcore, c_label, c_property, &c_has_property)

	has_property = goBool(c_has_property)
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

func (s *Session) GetAllowedPropertyValues(label string, property string) (values []string, err error) {
	c_label := C.CString(label)
	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_property))

	var c_values **C.char
	status := C.MM_GetAllowedPropertyValues(s.mmcore, c_label, c_property, &c_values)
	defer C.MM_StringListFree(c_values)

	values = goStringList(c_values)
	err = statusToError(status)
	return
}

func (s *Session) IsPropertyReadOnly(label string, property string) (read_only bool, err error) {
	c_label := C.CString(label)
	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_property))

	var c_read_only C.uint8_t
	status := C.MM_IsPropertyReadOnly(s.mmcore, c_label, c_property, &c_read_only)

	read_only = goBool(c_read_only)
	err = statusToError(status)
	return
}

func (s *Session) IsPropertyPreInit(label string, property string) (pre_init bool, err error) {
	c_label := C.CString(label)
	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_property))

	var c_pre_init C.uint8_t
	status := C.MM_IsPropertyPreInit(s.mmcore, c_label, c_property, &c_pre_init)

	pre_init = goBool(c_pre_init)
	err = statusToError(status)
	return
}

func (s *Session) IsPropertySequenceable(label string, property string) (sequenceable bool, err error) {
	c_label := C.CString(label)
	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_property))

	var c_sequenceable C.uint8_t
	status := C.MM_IsPropertySequenceable(s.mmcore, c_label, c_property, &c_sequenceable)

	sequenceable = goBool(c_sequenceable)
	err = statusToError(status)
	return
}

func (s *Session) HasPropertyLimits(label string, property string) (has_limits bool, err error) {
	c_label := C.CString(label)
	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_property))

	var c_has_limits C.uint8_t
	status := C.MM_HasPropertyLimits(s.mmcore, c_label, c_property, &c_has_limits)

	has_limits = goBool(c_has_limits)
	err = statusToError(status)
	return
}

func (s *Session) GetPropertyLowerLimit(label string, property string) (lower_limit float64, err error) {
	c_label := C.CString(label)
	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_property))

	var c_lower_limit C.double
	status := C.MM_GetPropertyLowerLimit(s.mmcore, c_label, c_property, &c_lower_limit)

	lower_limit = float64(c_lower_limit)
	err = statusToError(status)
	return
}

func (s *Session) GetPropertyUpperLimit(label string, property string) (upper_limit float64, err error) {
	c_label := C.CString(label)
	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_property))

	var c_upper_limit C.double
	status := C.MM_GetPropertyUpperLimit(s.mmcore, c_label, c_property, &c_upper_limit)

	upper_limit = float64(c_upper_limit)
	err = statusToError(status)
	return
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

func (s *Session) SetROI(x int, y int, x_size int, y_size int) error {
	status := C.MM_SetROI(s.mmcore, (C.int)(x), (C.int)(y), (C.int)(x_size), (C.int)(y_size))
	return statusToError(status)
}

func (s *Session) GetROI() (x int, y int, x_size int, y_size int, err error) {
	var c_x C.int
	var c_y C.int
	var c_x_size C.int
	var c_y_size C.int
	status := C.MM_GetROI(s.mmcore, &c_x, &c_y, &c_x_size, &c_y_size)

	x = int(c_x)
	y = int(c_y)
	x_size = int(x_size)
	y_size = int(y_size)
	err = statusToError(status)
	return
}

func (s *Session) ClearROI() error {
	status := C.MM_ClearROI(s.mmcore)
	return statusToError(status)
}

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
// Shutter control
//

func (s *Session) SetShutterOpen(label string, is_open bool) error {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetShutterOpen(s.mmcore, c_label, cBool(is_open))
	return statusToError(status)
}

func (s *Session) GetShutterOpen(label string) (is_open bool, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_is_open C.uint8_t
	status := C.MM_GetShutterOpen(s.mmcore, c_label, &c_is_open)

	is_open = goBool(c_is_open)
	err = statusToError(status)
	return
}

//
// Autofocus control
//

func (s *Session) LastFocusScore() (score float64) {
	var c_score C.double
	C.MM_GetLastFocusScore(s.mmcore, &c_score)

	score = float64(c_score)
	return
}

func (s *Session) CurrentFocusScore() (score float64) {
	var c_score C.double
	C.MM_GetCurrentFocusScore(s.mmcore, &c_score)

	score = float64(c_score)
	return
}

func (s *Session) EnableContinuousFocus() error {
	status := C.MM_EnableContinuousFocus(s.mmcore)
	return statusToError(status)
}

func (s *Session) DisableContinuousFocus() error {
	status := C.MM_DisableContinuousFocus(s.mmcore)
	return statusToError(status)
}

func (s *Session) IsContinuousFocusEnabled() (enabled bool, err error) {
	var c_enabled C.uint8_t
	status := C.MM_IsContinuousFocusEnabled(s.mmcore, &c_enabled)

	enabled = goBool(c_enabled)
	err = statusToError(status)
	return
}

func (s *Session) IsContinuousFocusLocked() (locked bool, err error) {
	var c_locked C.uint8_t
	status := C.MM_IsContinuousFocusLocked(s.mmcore, &c_locked)

	locked = goBool(c_locked)
	err = statusToError(status)
	return
}

func (s *Session) IsContinuousFocusDrive(label string) (is_continuous_focus_drive bool, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_is_continuous_focus_drive C.uint8_t
	status := C.MM_IsContinuousFocusDrive(s.mmcore, c_label, &c_is_continuous_focus_drive)

	is_continuous_focus_drive = goBool(c_is_continuous_focus_drive)
	err = statusToError(status)
	return
}

func (s *Session) FullFocus() (err error) {
	status := C.MM_FullFocus(s.mmcore)
	return statusToError(status)
}

func (s *Session) IncrementalFocus() (err error) {
	status := C.MM_IncrementalFocus(s.mmcore)
	return statusToError(status)
}

func (s *Session) SetAutoFocusOffset(offset float64) error {
	status := C.MM_SetAutoFocusOffset(s.mmcore, (C.double)(offset))
	return statusToError(status)
}

func (s *Session) GetAutoFocusOffset() (offset float64, err error) {
	var c_offset C.double
	status := C.MM_GetAutoFocusOffset(s.mmcore, &c_offset)

	offset = float64(c_offset)
	err = statusToError(status)
	return
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

func (s *Session) SetStateLabel(label string, state_label string) error {
	c_label := C.CString(label)
	c_state_label := C.CString(state_label)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_state_label))

	status := C.MM_SetStateLabel(s.mmcore, c_label, c_state_label)
	return statusToError(status)
}

func (s *Session) GetStateLabel(label string) (state_label string, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_state_label *C.char
	status := C.MM_GetStateLabel(s.mmcore, c_label, &c_state_label)
	defer C.MM_StringFree(c_state_label)

	state_label = C.GoString(c_state_label)
	err = statusToError(status)
	return
}

func (s *Session) DefineStateLabel(label string, state int, state_label string) error {
	c_label := C.CString(label)
	c_state_label := C.CString(state_label)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_state_label))

	status := C.MM_DefineStateLabel(s.mmcore, c_label, (C.int32_t)(state), c_state_label)
	return statusToError(status)
}

func (s *Session) GetStateLabels(label string) (state_labels []string, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_state_labels **C.char
	status := C.MM_GetStateLabels(s.mmcore, c_label, &c_state_labels)
	defer C.MM_StringListFree(c_state_labels)

	state_labels = goStringList(c_state_labels)
	err = statusToError(status)
	return
}

func (s *Session) GetStateFromLabel(label string, state_label string) (state int, err error) {
	c_label := C.CString(label)
	c_state_label := C.CString(state_label)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_state_label))

	var c_state C.int32_t
	status := C.MM_GetStateFromLabel(s.mmcore, c_label, c_state_label, &c_state)

	state = int(c_state)
	err = statusToError(status)
	return
}

//
// Focus (Z) stage control
//

func (s *Session) SetPosition(label string, position float64) error {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetPosition(s.mmcore, c_label, (C.double)(position))
	return statusToError(status)
}

func (s *Session) SetRelativePosition(label string, delta float64) error {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetRelativePosition(s.mmcore, c_label, (C.double)(delta))
	return statusToError(status)
}

func (s *Session) GetPosition(label string) (position float64, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_GetPosition(s.mmcore, c_label, (*C.double)(&position))
	err = statusToError(status)
	return
}

func (s *Session) SetOrigin(label string) (err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetOrigin(s.mmcore, c_label)
	err = statusToError(status)
	return
}

func (s *Session) SetAdapterOrigin(label string, new_z_um float64) (err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetAdapterOrigin(s.mmcore, c_label, (C.double)(new_z_um))
	err = statusToError(status)
	return
}

func (s *Session) SetFocusDirection(label string, sign int) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	C.MM_SetFocusDirection(s.mmcore, c_label, (C.int8_t)(sign))
	return
}

func (s *Session) GetFocusDirection(label string) (sign int, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_sign C.int8_t
	status := C.MM_GetFocusDirection(s.mmcore, c_label, &c_sign)

	sign = int(c_sign)
	err = statusToError(status)
	return
}

//
// XY stage control
//

func (s *Session) SetXYPosition(label string, x float64, y float64) (err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetXYPosition(s.mmcore, c_label, (C.double)(x), (C.double)(y))
	err = statusToError(status)
	return
}

func (s *Session) SetRelativeXYPosition(label string, dx float64, dy float64) (err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetRelativeXYPosition(s.mmcore, c_label, (C.double)(dx), (C.double)(dy))
	err = statusToError(status)
	return
}

func (s *Session) GetXYPosition(label string) (x float64, y float64, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_x C.double
	var c_y C.double
	status := C.MM_GetXYPosition(s.mmcore, c_label, &c_x, &c_y)

	x = float64(c_x)
	y = float64(c_y)
	err = statusToError(status)
	return
}

func (s *Session) GetXPosition(label string) (x float64, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_x C.double
	status := C.MM_GetXPosition(s.mmcore, c_label, &c_x)

	x = float64(c_x)
	err = statusToError(status)
	return
}

func (s *Session) GetYPosition(label string) (y float64, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_y C.double
	status := C.MM_GetYPosition(s.mmcore, c_label, &c_y)

	y = float64(c_y)
	err = statusToError(status)
	return
}

func (s *Session) Stop(label string) (err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_Stop(s.mmcore, c_label)
	return statusToError(status)
}

func (s *Session) Home(label string) (err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_Home(s.mmcore, c_label)
	return statusToError(status)
}

func (s *Session) SetOriginXY(label string) (err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetOriginXY(s.mmcore, c_label)
	return statusToError(status)
}

func (s *Session) SetOriginX(label string) (err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetOriginX(s.mmcore, c_label)
	return statusToError(status)
}

func (s *Session) SetOriginY(label string) (err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetOriginY(s.mmcore, c_label)
	return statusToError(status)
}

func (s *Session) SetAdpaterOriginXY(label string, new_x_um float64, new_y_um float64) (err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	status := C.MM_SetAdpaterOriginXY(s.mmcore, c_label, (C.double)(new_x_um), (C.double)(new_y_um))
	err = statusToError(status)
	return
}

//
// Hub and peripheral devices
//

func (s *Session) SetParentLabel(label string, parent_label string) (err error) {
	c_label := C.CString(label)
	c_parent_label := C.CString(parent_label)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_parent_label))

	status := C.MM_SetParentLabel(s.mmcore, c_label, c_parent_label)
	return statusToError(status)
}

func (s *Session) GetParentLabel(label string) (parent_label string, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_parent_label *C.char
	status := C.MM_GetParentLabel(s.mmcore, c_label, &c_parent_label)
	defer C.MM_StringFree(c_parent_label)

	parent_label = C.GoString(c_parent_label)
	err = statusToError(status)
	return
}

func (s *Session) GetInstalledDevices(hub_label string) (names []string, err error) {
	c_hub_label := C.CString(hub_label)
	defer C.free(unsafe.Pointer(c_hub_label))

	var c_names **C.char
	status := C.MM_GetInstalledDevices(s.mmcore, c_hub_label, &c_names)
	defer C.MM_StringListFree(c_names)

	names = goStringList(c_names)
	err = statusToError(status)
	return
}

func (s *Session) GetInstalledDeviceDescription(hub_label string, name string) (descriptions string, err error) {
	c_hub_label := C.CString(hub_label)
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_hub_label))
	defer C.free(unsafe.Pointer(c_name))

	var c_descriptions *C.char
	status := C.MM_GetInstalledDeviceDescription(s.mmcore, c_hub_label, c_name, &c_descriptions)
	defer C.MM_StringFree(c_descriptions)

	descriptions = C.GoString(c_descriptions)
	err = statusToError(status)
	return
}

func (s *Session) GetLoadedPeripheralDevices(hub_label string) (labels []string, err error) {
	c_hub_label := C.CString(hub_label)
	defer C.free(unsafe.Pointer(c_hub_label))

	var c_labels **C.char
	status := C.MM_GetLoadedPeripheralDevices(s.mmcore, c_hub_label, &c_labels)
	defer C.MM_StringListFree(c_labels)

	labels = goStringList(c_labels)
	err = statusToError(status)
	return
}

//
// Miscellaneous
//

func (s *Session) UserId() (userid string) {
	var c_userid *C.char
	C.MM_GetUserId(s.mmcore, &c_userid)
	defer C.MM_StringFree(c_userid)

	userid = C.GoString(c_userid)
	return
}

func (s *Session) HostName() (hostname string) {
	var c_hostname *C.char
	C.MM_GetHostName(s.mmcore, &c_hostname)
	defer C.MM_StringFree(c_hostname)

	hostname = C.GoString(c_hostname)
	return
}

func (s *Session) MACAddresses() (addresses []string) {
	var c_addresses **C.char
	C.MM_GetMACAddresses(s.mmcore, &c_addresses)
	defer C.MM_StringListFree(c_addresses)

	addresses = goStringList(c_addresses)
	return
}

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

func cBool(go_bool bool) C.uint8_t {
	if go_bool {
		return 1
	}
	return 0
}
