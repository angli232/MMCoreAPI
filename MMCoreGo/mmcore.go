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

func (s *Session) VersionInfo() string {
	var buf [256]C.char
	C.MM_GetVersionInfo(s.mmcore, (*C.char)(unsafe.Pointer(&buf[0])), (C.int)(len(buf)))
	return C.GoString((*C.char)(&buf[0]))
}

func (s *Session) APIVersionInfo() string {
	var buf [256]C.char
	C.MM_GetAPIVersionInfo(s.mmcore, (*C.char)(unsafe.Pointer(&buf[0])), (C.int)(len(buf)))
	return C.GoString((*C.char)(&buf[0]))
}

//
// Initialization and setup.
//

func (s *Session) LoadDevice(label, module_name, device_name string) error {
	c_label := C.CString(label)
	c_module_name := C.CString(module_name)
	c_device_name := C.CString(device_name)
	defer C.free(unsafe.Pointer(c_label))
	defer C.free(unsafe.Pointer(c_module_name))
	defer C.free(unsafe.Pointer(c_device_name))

	status := C.MM_LoadDevice(s.mmcore, c_label, c_module_name, c_device_name)
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

func (s *Session) DeviceAdapterSearchPaths() (paths []string) {
	var c_paths **C.char
	var len_c_paths C.size_t
	C.MM_GetDeviceAdapterSearchPaths(s.mmcore, &c_paths, &len_c_paths)

	if len_c_paths == 0 {
		return []string{}
	}

	paths = make([]string, len_c_paths)

	c_paths_slice := (*[1 << 30]*C.char)(unsafe.Pointer(c_paths))[:len_c_paths:len_c_paths]
	for i, c_path := range c_paths_slice {
		paths[i] = C.GoString(c_path)
		C.MM_Free(unsafe.Pointer(c_path))
	}
	C.MM_Free(unsafe.Pointer(c_paths))

	return
}

func (s *Session) SetDeviceAdapterSearchPaths(paths []string) {
	c_paths := make([]*C.char, len(paths))
	for i, path := range paths {
		c_paths[i] = C.CString(path)
	}
	C.MM_SetDeviceAdapterSearchPaths(s.mmcore, &c_paths[0], (C.size_t)(len(c_paths)))
	for _, c_path := range c_paths {
		C.free(unsafe.Pointer(c_path))
	}
}

func (s *Session) GetDeviceAdapterNames() (names []string, err error) {
	var c_names **C.char
	var len_c_names C.size_t
	status := C.MM_GetDeviceAdapterNames(s.mmcore, &c_names, &len_c_names)

	if status != 0 {
		err = Error(int(status))
		return
	}

	if len_c_names == 0 {
		return []string{}, nil
	}

	names = make([]string, int(len_c_names))

	c_names_slice := (*[1 << 30]*C.char)(unsafe.Pointer(c_names))[:len_c_names:len_c_names]
	for i, c_name := range c_names_slice {
		names[i] = C.GoString(c_name)
		C.MM_Free(unsafe.Pointer(c_name))
	}
	C.MM_Free(unsafe.Pointer(c_names))

	return
}

//
// Generic device control
//
func (s *Session) GetDevicePropertyNames(label string) (names []string, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_names **C.char
	var len_c_names C.size_t
	status := C.MM_GetDevicePropertyNames(s.mmcore, c_label, &c_names, &len_c_names)

	if status != 0 {
		err = Error(int(status))
		return
	}

	if len_c_names == 0 {
		return []string{}, nil
	}

	names = make([]string, int(len_c_names))

	c_names_slice := (*[1 << 30]*C.char)(unsafe.Pointer(c_names))[:len_c_names:len_c_names]
	for i, c_name := range c_names_slice {
		names[i] = C.GoString(c_name)
		C.MM_Free(unsafe.Pointer(c_name))
	}
	C.MM_Free(unsafe.Pointer(c_names))

	return
}

func (s *Session) GetProperty(label string, property string) (value string, err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_property))

	var c_value *C.char
	var len_c_value C.size_t

	status := C.MM_GetProperty(s.mmcore, c_label, c_property, &c_value, &len_c_value)
	defer C.MM_Free(unsafe.Pointer(c_value))

	value = C.GoString(c_value)
	err = statusToError(status)
	return
}

func (s *Session) SetProperty(label string, property string, state interface{}) (err error) {
	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_property))

	var status C.int
	switch state.(type) {
	case bool:
		var c_state C.int
		if state.(bool) {
			c_state = 1
		} else {
			c_state = 0
		}
		status = C.MM_SetPropertyBool(s.mmcore, c_label, c_property, c_state)
	case int:
		status = C.MM_SetPropertyInt(s.mmcore, c_label, c_property, (C.int32_t)(state.(int)))
	case float32:
		status = C.MM_SetPropertyFloat32(s.mmcore, c_label, c_property, (C.float)(state.(float32)))
	case float64:
		status = C.MM_SetPropertyFloat64(s.mmcore, c_label, c_property, (C.double)(state.(float64)))
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

//
// Manage current devices.
//

func (s *Session) SetExposureTime(exposure_ms float64) error {
	status := C.MM_SetExposure(s.mmcore, (C.double)(exposure_ms))
	return statusToError(status)
}

func (s *Session) ExposureTime() (exposure_ms float64, err error) {
	status := C.MM_GetExposure(s.mmcore, (*C.double)(&exposure_ms))
	err = statusToError(status)
	return
}

func (s *Session) SnapImage() error {
	status := C.MM_SnapImage(s.mmcore)
	return statusToError(status)
}

func (s *Session) ImageBufferSize() (len int, err error) {
	var c_len C.uint32_t
	status := C.MM_GetImageBufferSize(s.mmcore, &c_len)
	len = int(c_len)
	err = statusToError(status)
	return
}

func (s *Session) GetImage() (buf []byte, err error) {
	len, err := s.ImageBufferSize()
	if err != nil {
		return
	}

	var c_pbuf *C.uint8_t
	status := C.MM_GetImage(s.mmcore, &c_pbuf)
	if status != 0 {
		err = Error(int(status))
	}

	buf = C.GoBytes(unsafe.Pointer(c_pbuf), C.int(len))
	return
}

func (s *Session) ImageWidth() (width int, err error) {
	var c_width C.uint16_t
	status := C.MM_GetImageWidth(s.mmcore, &c_width)
	width = int(c_width)
	err = statusToError(status)
	return
}

func (s *Session) ImageHeight() (height int, err error) {
	var c_height C.uint16_t
	status := C.MM_GetImageHeight(s.mmcore, &c_height)
	height = int(c_height)
	err = statusToError(status)
	return
}

func (s *Session) BytesPerPixel() (bytes_per_pixel int, err error) {
	var c_bytes_per_pixel C.uint8_t
	status := C.MM_GetBytesPerPixel(s.mmcore, &c_bytes_per_pixel)
	bytes_per_pixel = int(c_bytes_per_pixel)
	err = statusToError(status)
	return
}

func (s *Session) ImageBitDepth() (bit_depth int, err error) {
	var c_bit_depth C.uint8_t
	status := C.MM_GetImageBitDepth(s.mmcore, &c_bit_depth)
	bit_depth = int(c_bit_depth)
	err = statusToError(status)
	return
}

func (s *Session) NumberOfComponents() (n_components int, err error) {
	var c_n_components C.uint8_t
	status := C.MM_GetNumberOfComponents(s.mmcore, &c_n_components)
	n_components = int(c_n_components)
	err = statusToError(status)
	return
}

func (s *Session) NumberOfCameraChannels() (n_channels int, err error) {
	var c_n_channels C.uint8_t
	status := C.MM_GetNumberOfCameraChannels(s.mmcore, &c_n_channels)
	n_channels = int(c_n_channels)
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

//
// Device discovery.
//

// func (s *Session) SupportsDeviceDetection(label string) (bool, error) {
// 	c_label := C.CString(label)
// 	defer C.free(unsafe.Pointer(c_label))

// 	var support uint8_t
// 	status := C.MM_SupportsDeviceDetection(s.mmcore, c_label, &support)
// 	return statusToError(status)
// }

//
// Hub and peripheral devices.
//
