#ifndef MMCOREC_H_
#define MMCOREC_H_

#include <stdint.h>

#define DllExport __declspec(dllexport)

typedef void* MM_Session;
typedef int MM_Status;

#ifdef __cplusplus
extern "C" {
#endif

DllExport void MM_Open(MM_Session* mm);
DllExport void MM_Close(MM_Session mm);
DllExport void MM_GetVersionInfo(MM_Session mm, char *info, int len_info);
DllExport void MM_GetAPIVersionInfo(MM_Session mm, char *info, int len_info);

// Device initialization and setup
DllExport MM_Status MM_LoadDevice(MM_Session mm, const char* label, const char* module_name, const char* device_name);
DllExport MM_Status MM_UnloadDevice(MM_Session mm, const char* label);
DllExport MM_Status MM_UnloadAllDevices(MM_Session mm);
DllExport MM_Status MM_InitializeAllDevices(MM_Session mm);
DllExport MM_Status MM_InitializeDevice(MM_Session mm, const char* label);
DllExport MM_Status MM_Reset(MM_Session mm);

// Device listing
DllExport void MM_SetDeviceAdapterSearchPaths(MM_Session mm, const char** paths, size_t len_paths);
DllExport void MM_GetDeviceAdapterSearchPaths(MM_Session mm, char*** paths, size_t *len_paths);
DllExport MM_Status MM_GetDeviceAdapterNames(MM_Session mm, char*** device_adapter_names, size_t *len_paths);

// Generic device control

// Manage current devices
DllExport MM_Status MM_SetCameraDevice(MM_Session mm, const char* label);

// Image acquisition
DllExport MM_Status MM_SetExposure(MM_Session mm, double exp);
DllExport MM_Status MM_GetExposure(MM_Session mm, double *exp);
DllExport MM_Status MM_SnapImage(MM_Session mm);
DllExport MM_Status MM_GetImage(MM_Session mm, uint8_t **ptr_buffer);

DllExport MM_Status MM_GetImageWidth(MM_Session mm, uint16_t *width);
DllExport MM_Status MM_GetImageHeight(MM_Session mm, uint16_t *height);
DllExport MM_Status MM_GetBytesPerPixel(MM_Session mm, uint8_t *bytes);
DllExport MM_Status MM_GetImageBitDepth(MM_Session mm, uint8_t *bit_depth);
DllExport MM_Status MM_GetNumberOfComponents(MM_Session mm, uint8_t *n_components);
DllExport MM_Status MM_GetNumberOfCameraChannels(MM_Session mm, uint8_t *n_channels);
DllExport MM_Status MM_GetImageBufferSize(MM_Session mm, uint32_t *len);

// State device control
DllExport MM_Status MM_SetState(MM_Session mm, const char* label, int32_t state);
DllExport MM_Status MM_GetState(MM_Session mm, const char* label, int32_t *state);
DllExport MM_Status MM_GetNumberOfStates(MM_Session mm, const char* label, int32_t *state);

enum {
    MM_ErrOK                        = 0,
    MM_ErrGENERIC                   = 1, // unspecified error
    MM_ErrNoDevice                  = 2,
    MM_ErrSetPropertyFailed         = 3,
    MM_ErrLibraryFunctionNotFound   = 4,
    MM_ErrModuleVersionMismatch     = 5,
    MM_ErrDeviceVersionMismatch     = 6,
    MM_ErrUnknownModule             = 7,
    MM_ErrLoadLibraryFailed         = 8,
    MM_ErrCreateFailed              = 9,
    MM_ErrCreateNotFound            = 10,
    MM_ErrDeleteNotFound            = 11,
    MM_ErrDeleteFailed              = 12,
    MM_ErrUnexpectedDevice          = 13,
    MM_ErrDeviceUnloadFailed        = 14,
    MM_ErrCameraNotAvailable        = 15,
    MM_ErrDuplicateLabel            = 16,
    MM_ErrInvalidLabel              = 17,
    MM_ErrInvalidStateDevice        = 19,
    MM_ErrNoConfiguration           = 20,
    MM_ErrInvalidConfigurationIndex = 21,
    MM_ErrDEVICE_GENERIC            = 22,
    MM_ErrInvalidPropertyBlock      = 23,
    MM_ErrUnhandledException        = 24,
    MM_ErrDevicePollingTimeout      = 25,
    MM_ErrInvalidShutterDevice      = 26,
    MM_ErrInvalidSerialDevice       = 27,
    MM_ErrInvalidStageDevice        = 28,
    MM_ErrInvalidSpecificDevice     = 29,
    MM_ErrInvalidXYStageDevice      = 30,
    MM_ErrFileOpenFailed            = 31,
    MM_ErrInvalidCFGEntry           = 32,
    MM_ErrInvalidContents           = 33,
    MM_ErrInvalidCoreProperty       = 34,
    MM_ErrInvalidCoreValue          = 35,
    MM_ErrNoConfigGroup             = 36,
    MM_ErrCameraBufferReadFailed    = 37,
    MM_ErrDuplicateConfigGroup      = 38,
    MM_ErrInvalidConfigurationFile  = 39,
    MM_ErrCircularBufferFailedToInitialize = 40,
    MM_ErrCircularBufferEmpty       = 41,
    MM_ErrContFocusNotAvailable     = 42,
    MM_ErrAutoFocusNotAvailable     = 43,
    MM_ErrBadConfigName             = 44,
    MM_ErrCircularBufferIncompatibleImage = 45,
    MM_ErrNotAllowedDuringSequenceAcquisition = 46,
    MM_ErrOutOfMemory              = 47,
    MM_ErrInvalidImageSequence     = 48,
    MM_ErrNullPointerException     = 49,
    MM_ErrCreatePeripheralFailed   = 50,
    MM_ErrPropertyNotInCache       = 51,
    MM_ErrBadAffineTransform       = 52
};

#ifdef __cplusplus
}
#endif

#endif
