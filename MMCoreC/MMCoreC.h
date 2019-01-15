#ifndef MMCOREC_H_
#define MMCOREC_H_

#include <stdint.h>

#define DllExport __declspec(dllexport)

typedef void *MM_Session;

typedef enum {
    MM_ErrOK = 0,
    MM_ErrGENERIC = 1, // unspecified error
    MM_ErrNoDevice = 2,
    MM_ErrSetPropertyFailed = 3,
    MM_ErrLibraryFunctionNotFound = 4,
    MM_ErrModuleVersionMismatch = 5,
    MM_ErrDeviceVersionMismatch = 6,
    MM_ErrUnknownModule = 7,
    MM_ErrLoadLibraryFailed = 8,
    MM_ErrCreateFailed = 9,
    MM_ErrCreateNotFound = 10,
    MM_ErrDeleteNotFound = 11,
    MM_ErrDeleteFailed = 12,
    MM_ErrUnexpectedDevice = 13,
    MM_ErrDeviceUnloadFailed = 14,
    MM_ErrCameraNotAvailable = 15,
    MM_ErrDuplicateLabel = 16,
    MM_ErrInvalidLabel = 17,
    MM_ErrInvalidStateDevice = 19,
    MM_ErrNoConfiguration = 20,
    MM_ErrInvalidConfigurationIndex = 21,
    MM_ErrDEVICE_GENERIC = 22,
    MM_ErrInvalidPropertyBlock = 23,
    MM_ErrUnhandledException = 24,
    MM_ErrDevicePollingTimeout = 25,
    MM_ErrInvalidShutterDevice = 26,
    MM_ErrInvalidSerialDevice = 27,
    MM_ErrInvalidStageDevice = 28,
    MM_ErrInvalidSpecificDevice = 29,
    MM_ErrInvalidXYStageDevice = 30,
    MM_ErrFileOpenFailed = 31,
    MM_ErrInvalidCFGEntry = 32,
    MM_ErrInvalidContents = 33,
    MM_ErrInvalidCoreProperty = 34,
    MM_ErrInvalidCoreValue = 35,
    MM_ErrNoConfigGroup = 36,
    MM_ErrCameraBufferReadFailed = 37,
    MM_ErrDuplicateConfigGroup = 38,
    MM_ErrInvalidConfigurationFile = 39,
    MM_ErrCircularBufferFailedToInitialize = 40,
    MM_ErrCircularBufferEmpty = 41,
    MM_ErrContFocusNotAvailable = 42,
    MM_ErrAutoFocusNotAvailable = 43,
    MM_ErrBadConfigName = 44,
    MM_ErrCircularBufferIncompatibleImage = 45,
    MM_ErrNotAllowedDuringSequenceAcquisition = 46,
    MM_ErrOutOfMemory = 47,
    MM_ErrInvalidImageSequence = 48,
    MM_ErrNullPointerException = 49,
    MM_ErrCreatePeripheralFailed = 50,
    MM_ErrPropertyNotInCache = 51,
    MM_ErrBadAffineTransform = 52
} MM_Status;

typedef enum {
    MM_UnknownType = 0,
    MM_AnyType,
    MM_CameraDevice,
    MM_ShutterDevice,
    MM_StateDevice,
    MM_StageDevice,
    MM_XYStageDevice,
    MM_SerialDevice,
    MM_GenericDevice,
    MM_AutoFocusDevice,
    MM_CoreDevice,
    MM_ImageProcessorDevice,
    MM_SignalIODevice,
    MM_MagnifierDevice,
    MM_SLMDevice,
    MM_HubDevice,
    MM_GalvoDevice
} MM_DeviceType;

typedef enum { MM_Undef, MM_String, MM_Float, MM_Integer } MM_PropertyType;

#ifdef __cplusplus
extern "C" {
#endif

DllExport void MM_Open(MM_Session *mm);
DllExport void MM_Close(MM_Session mm);
DllExport void MM_GetVersionInfo(MM_Session mm, char **info);
DllExport void MM_GetAPIVersionInfo(MM_Session mm, char **info);

DllExport void MM_StringFree(char *str);
DllExport void MM_StringListFree(char **str_list);

// Device initialization and setup
DllExport MM_Status MM_LoadDevice(MM_Session mm, const char *label,
                                  const char *module_name,
                                  const char *device_name);
DllExport MM_Status MM_UnloadDevice(MM_Session mm, const char *label);
DllExport MM_Status MM_UnloadAllDevices(MM_Session mm);
DllExport MM_Status MM_InitializeAllDevices(MM_Session mm);
DllExport MM_Status MM_InitializeDevice(MM_Session mm, const char *label);
DllExport MM_Status MM_Reset(MM_Session mm);

// Device listing
DllExport void MM_SetDeviceAdapterSearchPaths(MM_Session mm,
                                              const char **paths);
DllExport void MM_GetDeviceAdapterSearchPaths(MM_Session mm, char ***paths);
DllExport MM_Status MM_GetDeviceAdapterNames(MM_Session mm, char ***names);
DllExport MM_Status MM_GetAvailableDevices(MM_Session mm, const char *library,
                                           char ***names);
DllExport MM_Status MM_GetAvailableDeviceDescriptions(MM_Session mm,
                                                      const char *library,
                                                      char ***descriptions);
DllExport MM_Status MM_GetAvailableDeviceTypes(MM_Session mm,
                                               const char *library,
                                               MM_DeviceType **types,
                                               size_t *len_types);

// Generic device control
DllExport MM_Status MM_GetDevicePropertyNames(MM_Session mm, const char *label,
                                              char ***names);
DllExport MM_Status MM_HasProperty(MM_Session mm, const char *label,
                                   const char *prop_name,
                                   uint8_t *has_property);
DllExport MM_Status MM_GetProperty(MM_Session mm, const char *label,
                                   const char *prop_name, char **value);
DllExport MM_Status MM_SetPropertyString(MM_Session mm, const char *label,
                                         const char *prop_name,
                                         const char *value);
DllExport MM_Status MM_SetPropertyBool(MM_Session mm, const char *label,
                                       const char *prop_name,
                                       const uint8_t value);
DllExport MM_Status MM_SetPropertyInt(MM_Session mm, const char *label,
                                      const char *prop_name,
                                      const int32_t value);
DllExport MM_Status MM_SetPropertyFloat(MM_Session mm, const char *label,
                                        const char *prop_name,
                                        const float value);
DllExport MM_Status MM_SetPropertyDouble(MM_Session mm, const char *label,
                                         const char *prop_name,
                                         const double value);

DllExport MM_Status MM_GetAllowedPropertyValues(MM_Session mm,
                                                const char *label,
                                                const char *prop_name,
                                                char ***values);
DllExport MM_Status MM_IsPropertyReadOnly(MM_Session mm, const char *label,
                                          const char *prop_name,
                                          uint8_t *read_only);
DllExport MM_Status MM_IsPropertyPreInit(MM_Session mm, const char *label,
                                         const char *prop_name,
                                         uint8_t *pre_init);
DllExport MM_Status MM_IsPropertySequenceable(MM_Session mm, const char *label,
                                              const char *prop_name,
                                              uint8_t *sequenceable);
DllExport MM_Status MM_HasPropertyLimits(MM_Session mm, const char *label,
                                         const char *prop_name,
                                         uint8_t *has_limits);
DllExport MM_Status MM_GetPropertyLowerLimit(MM_Session mm, const char *label,
                                             const char *prop_name,
                                             double *lower_limit);
DllExport MM_Status MM_GetPropertyUpperLimit(MM_Session mm, const char *label,
                                             const char *prop_name,
                                             double *upper_limit);
DllExport MM_Status MM_GetPropertyType(MM_Session mm, const char *label,
                                       const char *prop_name,
                                       MM_PropertyType *type);

DllExport MM_Status MM_DeviceBusy(MM_Session mm, const char *label,
                                  uint8_t *busy);
DllExport MM_Status MM_DeviceTypeBusy(MM_Session mm, MM_DeviceType type,
                                      uint8_t *busy);

// Manage current devices
DllExport MM_Status MM_SetCameraDevice(MM_Session mm, const char *label);
DllExport MM_Status MM_SetShutterDevice(MM_Session mm, const char *label);
DllExport MM_Status MM_SetFocusDevice(MM_Session mm, const char *label);
DllExport MM_Status MM_SetXYStageDevice(MM_Session mm, const char *label);
DllExport MM_Status MM_SetAutoFocusDevice(MM_Session mm, const char *label);

DllExport void MM_GetCameraDevice(MM_Session mm, char **label);
DllExport void MM_GetShutterDevice(MM_Session mm, char **label);
DllExport void MM_GetFocusDevice(MM_Session mm, char **label);
DllExport void MM_GetXYStageDevice(MM_Session mm, char **label);
DllExport void MM_GetAutoFocusDevice(MM_Session mm, char **label);

// Image acquisition settings
DllExport MM_Status MM_SetROI(MM_Session mm, int x, int y, int x_size,
                              int y_size);
DllExport MM_Status MM_GetROI(MM_Session mm, int *x, int *y, int *x_size,
                              int *y_size);
DllExport MM_Status MM_ClearROI(MM_Session mm);

DllExport MM_Status MM_SetExposure(MM_Session mm, double exp);
DllExport MM_Status MM_GetExposure(MM_Session mm, double *exp);

DllExport void MM_GetImageWidth(MM_Session mm, uint16_t *width);
DllExport void MM_GetImageHeight(MM_Session mm, uint16_t *height);
DllExport void MM_GetBytesPerPixel(MM_Session mm, uint8_t *bytes);
DllExport void MM_GetImageBitDepth(MM_Session mm, uint8_t *bit_depth);
DllExport void MM_GetNumberOfComponents(MM_Session mm,
                                        uint8_t *n_components);
DllExport void MM_GetNumberOfCameraChannels(MM_Session mm,
                                            uint8_t *n_channels);
DllExport void MM_GetImageBufferSize(MM_Session mm, uint32_t *len);

// Image acquisition
DllExport MM_Status MM_SnapImage(MM_Session mm);
DllExport MM_Status MM_GetImage(MM_Session mm, uint8_t **ptr_buffer);
DllExport MM_Status MM_GetImageOfChannel(MM_Session mm, uint16_t channel, uint8_t **ptr_buffer);

// Image sequence acquisition
DllExport MM_Status MM_StartSequenceAcquisition(MM_Session mm,
                                                int16_t num_images,
                                                double interval_ms,
                                                uint8_t stop_on_overflow);
DllExport MM_Status MM_StartContinuousSequenceAcquisition(MM_Session mm,
                                                          double interval_ms);
DllExport MM_Status MM_StopSequenceAcquisition(MM_Session mm);
DllExport void MM_IsSequenceRunning(MM_Session mm, uint8_t *status);

// Image circular buffer
DllExport MM_Status MM_GetLastImage(MM_Session mm, uint8_t **ptr_buffer);
DllExport MM_Status MM_PopNextImage(MM_Session mm, uint8_t **ptr_buffer);

DllExport void MM_GetRemainingImageCount(MM_Session mm, int16_t *count);
DllExport void MM_GetBufferTotalCapacity(MM_Session mm, int16_t *capacity);
DllExport void MM_GetBufferFreeCapacity(MM_Session mm, int16_t *capacity);
DllExport void MM_IsBufferOverflowed(MM_Session mm, uint8_t *overflowed);

DllExport MM_Status MM_SetCircularBufferMemoryFootprint(MM_Session mm,
                                                        uint32_t size_MB);
DllExport void MM_GetCircularBufferMemoryFootprint(MM_Session mm,
                                                   uint32_t *size_MB);
DllExport MM_Status MM_InitializeCircularBuffer(MM_Session mm);
DllExport MM_Status MM_ClearCircularBuffer(MM_Session mm);

// Shutter control
DllExport MM_Status MM_SetShutterOpen(MM_Session mm, const char *label,
                                  uint8_t is_open);
DllExport MM_Status MM_GetShutteOpenr(MM_Session mm, const char *label,
                                  uint8_t *is_open);

// Autofocus control
DllExport void MM_GetLastFocusScore(MM_Session mm, double *score);
DllExport void MM_GetCurrentFocusScore(MM_Session mm, double *score);
DllExport MM_Status MM_EnableContinuousFocus(MM_Session mm);
DllExport MM_Status MM_DisableContinuousFocus(MM_Session mm);
DllExport MM_Status MM_IsContinuousFocusEnabled(MM_Session mm, uint8_t *status);
DllExport MM_Status MM_IsContinuousFocusLocked(MM_Session mm, uint8_t *status);
DllExport MM_Status MM_IsContinuousFocusDrive(
    MM_Session mm, const char *label, uint8_t *is_continuous_focus_drive);
DllExport MM_Status MM_FullFocus(MM_Session mm);
DllExport MM_Status MM_IncrementalFocus(MM_Session mm);
DllExport MM_Status MM_SetAutoFocusOffset(MM_Session mm, double offset);
DllExport MM_Status MM_GetAutoFocusOffset(MM_Session mm, double *offset);

// State device control
DllExport MM_Status MM_SetState(MM_Session mm, const char *label,
                                int32_t state);
DllExport MM_Status MM_GetState(MM_Session mm, const char *label,
                                int32_t *state);
DllExport MM_Status MM_GetNumberOfStates(MM_Session mm, const char *label,
                                         int32_t *state);
DllExport MM_Status MM_SetStateLabel(MM_Session mm, const char *label,
                                     const char *state_label);
DllExport MM_Status MM_GetStateLabel(MM_Session mm, const char *label,
                                     char **state_label);
DllExport MM_Status MM_DefineStateLabel(MM_Session mm, const char *label,
                                        int32_t state, const char *state_label);
DllExport MM_Status MM_GetStateLabels(MM_Session mm, const char *label,
                                      char ***state_labels);
DllExport MM_Status MM_GetStateFromLabel(MM_Session mm, const char *label,
                                         const char *state_label,
                                         int32_t *state);

// Focus (Z) stage control
DllExport MM_Status MM_SetPosition(MM_Session mm, const char *label,
                                   double position);
DllExport MM_Status MM_GetPosition(MM_Session mm, const char *label,
                                   double *position);
DllExport MM_Status MM_SetRelativePosition(MM_Session mm, const char *label,
                                           double delta);
DllExport MM_Status MM_SetOrigin(MM_Session mm, const char *label);
DllExport MM_Status MM_SetAdapterOrigin(MM_Session mm, const char *label,
                                        double new_z_um);
DllExport MM_Status MM_SetFocusDirection(MM_Session mm, const char *label,
                                         int8_t sign);
DllExport MM_Status MM_GetFocusDirection(MM_Session mm, const char *label,
                                         int8_t *sign);

// XY stage control
DllExport MM_Status MM_SetXYPosition(MM_Session mm, const char *label, double x,
                                     double y);
DllExport MM_Status MM_SetRelativeXYPosition(MM_Session mm, const char *label,
                                             double dx, double dy);

DllExport MM_Status MM_GetXYPosition(MM_Session mm, const char *label,
                                     double *x, double *y);
DllExport MM_Status MM_GetXPosition(MM_Session mm, const char *label,
                                    double *x);
DllExport MM_Status MM_GetYPosition(MM_Session mm, const char *label,
                                    double *y);

DllExport MM_Status MM_Stop(MM_Session mm, const char *label);
DllExport MM_Status MM_Home(MM_Session mm, const char *label);
DllExport MM_Status MM_SetOriginXY(MM_Session mm, const char *label);
DllExport MM_Status MM_SetOriginX(MM_Session mm, const char *label);
DllExport MM_Status MM_SetOriginY(MM_Session mm, const char *label);
DllExport MM_Status MM_SetAdpaterOriginXY(MM_Session mm, const char *label,
                                          double new_x_um, double new_y_um);

// Hub and peripheral devices
DllExport MM_Status MM_SetParentLabel(MM_Session mm, const char *label,
                                      const char *parent_label);
DllExport MM_Status MM_GetParentLabel(MM_Session mm, const char *label,
                                      char **parent_label);

DllExport MM_Status MM_GetInstalledDevices(MM_Session mm, const char *hub_label,
                                           char ***names);
DllExport MM_Status MM_GetInstalledDeviceDescription(MM_Session mm,
                                                     const char *hub_label,
                                                     const char *name,
                                                     char **descriptions);
DllExport MM_Status MM_GetLoadedPeripheralDevices(MM_Session mm,
                                                  const char *hub_label,
                                                  char ***labels);

// Miscellaneous
DllExport void MM_GetUserId(MM_Session mm, char **userid);
DllExport void MM_GetHostName(MM_Session mm, char **hostname);
DllExport void MM_GetMACAddresses(MM_Session mm, char ***addresses);

#ifdef __cplusplus
}
#endif

#endif
