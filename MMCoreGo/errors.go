package mmcore

// #cgo CFLAGS: -I../MMCoreC
//
// #include "MMCoreC.h"
import "C"

import (
	"fmt"
)

type Error int

func statusToError(status C.MM_Status) error {
	if int(C.int(status)) == 0 {
		return nil
	}
	return Error(int(C.int(status)))
}

func (e Error) Error() string {
	s := errText[e]
	if s == "s" {
		return fmt.Sprintf("error %d", int(e))
	}
	return s
}

var (
	ErrGENERIC                             Error = 1
	ErrNoDevice                            Error = 2
	ErrSetPropertyFailed                   Error = 3
	ErrLibraryFunctionNotFound             Error = 4
	ErrModuleVersionMismatch               Error = 5
	ErrDeviceVersionMismatch               Error = 6
	ErrUnknownModule                       Error = 7
	ErrLoadLibraryFailed                   Error = 8
	ErrCreateFailed                        Error = 9
	ErrCreateNotFound                      Error = 10
	ErrDeleteNotFound                      Error = 11
	ErrDeleteFailed                        Error = 12
	ErrUnexpectedDevice                    Error = 13
	ErrDeviceUnloadFailed                  Error = 14
	ErrCameraNotAvailable                  Error = 15
	ErrDuplicateLabel                      Error = 16
	ErrInvalidLabel                        Error = 17
	ErrInvalidStateDevice                  Error = 19
	ErrNoConfiguration                     Error = 20
	ErrInvalidConfigurationIndex           Error = 21
	ErrDEVICE_GENERIC                      Error = 22
	ErrInvalidPropertyBlock                Error = 23
	ErrUnhandledException                  Error = 24
	ErrDevicePollingTimeout                Error = 25
	ErrInvalidShutterDevice                Error = 26
	ErrInvalidSerialDevice                 Error = 27
	ErrInvalidStageDevice                  Error = 28
	ErrInvalidSpecificDevice               Error = 29
	ErrInvalidXYStageDevice                Error = 30
	ErrFileOpenFailed                      Error = 31
	ErrInvalidCFGEntry                     Error = 32
	ErrInvalidContents                     Error = 33
	ErrInvalidCoreProperty                 Error = 34
	ErrInvalidCoreValue                    Error = 35
	ErrNoConfigGroup                       Error = 36
	ErrCameraBufferReadFailed              Error = 37
	ErrDuplicateConfigGroup                Error = 38
	ErrInvalidConfigurationFile            Error = 39
	ErrCircularBufferFailedToInitialize    Error = 40
	ErrCircularBufferEmpty                 Error = 41
	ErrContFocusNotAvailable               Error = 42
	ErrAutoFocusNotAvailable               Error = 43
	ErrBadConfigName                       Error = 44
	ErrCircularBufferIncompatibleImage     Error = 45
	ErrNotAllowedDuringSequenceAcquisition Error = 46
	ErrOutOfMemory                         Error = 47
	ErrInvalidImageSequence                Error = 48
	ErrNullPointerException                Error = 49
	ErrCreatePeripheralFailed              Error = 50
	ErrPropertyNotInCache                  Error = 51
	ErrBadAffineTransform                  Error = 52
)

var errText = map[Error]string{
	ErrGENERIC:                             "generic (unspecified) error",
	ErrNoDevice:                            "no device",
	ErrSetPropertyFailed:                   "set property failed",
	ErrLibraryFunctionNotFound:             "library function not found",
	ErrModuleVersionMismatch:               "module version mismatch",
	ErrDeviceVersionMismatch:               "device version mismatch",
	ErrUnknownModule:                       "unknown module",
	ErrLoadLibraryFailed:                   "load library failed",
	ErrCreateFailed:                        "create failed",
	ErrCreateNotFound:                      "create not found",
	ErrDeleteNotFound:                      "delete not found",
	ErrDeleteFailed:                        "delete failed",
	ErrUnexpectedDevice:                    "unexpected device",
	ErrDeviceUnloadFailed:                  "device unload failed",
	ErrCameraNotAvailable:                  "camera not available",
	ErrDuplicateLabel:                      "duplicated label",
	ErrInvalidLabel:                        "invalid label",
	ErrInvalidStateDevice:                  "invalid state device",
	ErrNoConfiguration:                     "no configuration",
	ErrInvalidConfigurationIndex:           "invalid configuration index",
	ErrDEVICE_GENERIC:                      "device generic (unspecified) error",
	ErrInvalidPropertyBlock:                "invalid property block",
	ErrUnhandledException:                  "unhandled exception",
	ErrDevicePollingTimeout:                "device polling timeout",
	ErrInvalidShutterDevice:                "invalid shutter device",
	ErrInvalidSerialDevice:                 "invalid serial device",
	ErrInvalidStageDevice:                  "invalid stage device",
	ErrInvalidSpecificDevice:               "invalid specific device",
	ErrInvalidXYStageDevice:                "invalid XY stage device",
	ErrFileOpenFailed:                      "file open failed",
	ErrInvalidCFGEntry:                     "invalid CFG entry",
	ErrInvalidContents:                     "invalid contents",
	ErrInvalidCoreProperty:                 "invalid core property",
	ErrInvalidCoreValue:                    "invalid core value",
	ErrNoConfigGroup:                       "no config group",
	ErrCameraBufferReadFailed:              "camera buffer read failed",
	ErrDuplicateConfigGroup:                "duplicated config group",
	ErrInvalidConfigurationFile:            "invalid configuration file",
	ErrCircularBufferFailedToInitialize:    "circular buffer failed to initialize",
	ErrCircularBufferEmpty:                 "circular buffer empty",
	ErrContFocusNotAvailable:               "continuous focus not available",
	ErrAutoFocusNotAvailable:               "auto focus not available",
	ErrBadConfigName:                       "bad config name",
	ErrCircularBufferIncompatibleImage:     "circular buffer incompatible image",
	ErrNotAllowedDuringSequenceAcquisition: "not allowed during sequence acquisition",
	ErrOutOfMemory:                         "out of memory",
	ErrInvalidImageSequence:                "invalid image sequence",
	ErrNullPointerException:                "null pointer exception",
	ErrCreatePeripheralFailed:              "create peripheral failed",
	ErrPropertyNotInCache:                  "property not in cache",
	ErrBadAffineTransform:                  "bad affine transform",
}
