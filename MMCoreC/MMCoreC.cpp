#include "MMCoreC.h"

#include <stdlib.h>
#include <string.h>
#include <map>

#include "MMCore.h"
#include "MMEventCallback.h"

extern "C" {

class MM_CPP_EventCallback;
void MM_FreeRegisteredCallback(MM_Session mm);
static std::map<MM_Session, MM_CPP_EventCallback*> mm_registered_callbacks;

void std_to_c_string(std::string std_str, char **c_str) {
    size_t cap_c_str = std_str.size() + 1;
    *c_str = (char *)malloc(cap_c_str);
    strcpy_s(*c_str, cap_c_str, std_str.c_str());
    return;
}

void std_to_c_string_list(std::vector<std::string> str_list,
                          char ***c_str_list) {
    size_t cap_c_str_list = str_list.size() + 1;
    *c_str_list = (char **)malloc(cap_c_str_list * sizeof(char *));
    (*c_str_list)[cap_c_str_list - 1] = NULL;

    for (size_t i = 0; i < str_list.size(); i++) {
        std_to_c_string(str_list[i], &((*c_str_list)[i]));
    }
    return;
}

DllExport void MM_Open(MM_Session *core) {
    *core = reinterpret_cast<MM_Session>(new CMMCore());
    return;
}

DllExport void MM_Close(MM_Session mm) {
    MM_FreeRegisteredCallback(mm);
    delete reinterpret_cast<CMMCore *>(mm);
}

DllExport void MM_GetVersionInfo(MM_Session mm, char **info) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str = core->getVersionInfo();
    std_to_c_string(str, info);
    return;
}

DllExport void MM_GetAPIVersionInfo(MM_Session mm, char **info) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str = core->getAPIVersionInfo();
    std_to_c_string(str, info);
    return;
}

//
//
//

DllExport void MM_StringFree(char *str) {
    if (str != NULL) {
        free(str);
    }
}

DllExport void MM_StringListFree(char **str_list) {
    if (str_list == NULL) {
        return;
    }

    size_t i = 0;
    while (str_list[i]) {
        free(str_list[i]);
        i++;
    }

    free(str_list);
}

//
// Device initialization and setup
//

DllExport MM_Status MM_LoadDevice(MM_Session mm, const char *label,
                                  const char *module_name,
                                  const char *device_name) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->loadDevice(label, module_name, device_name);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_UnloadDevice(MM_Session mm, const char *label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->unloadDevice(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_UnloadAllDevices(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->unloadAllDevices();
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_InitializeAllDevices(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->initializeAllDevices();
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_InitializeDevice(MM_Session mm, const char *label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->initializeDevice(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_Reset(MM_Session mm) {
    MM_FreeRegisteredCallback(mm);
    
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->reset();
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

//
// Event callback
//

class MM_CPP_EventCallback: public MMEventCallback {
private:
    struct MM_EventCallback *callback;
    MM_Session mm;
public:
    MM_CPP_EventCallback(struct MM_EventCallback *callback, MM_Session mm) {
        this->callback = callback;
        this->mm = mm;
    }
    virtual ~MM_CPP_EventCallback() {}

    void onPropertiesChanged(){
        if (this->callback->onPropertiesChanged != NULL) {
            this->callback->onPropertiesChanged(this->mm);
        }
    }

    void onPropertyChanged(const char *label, const char *propName, const char *propValue) {
        if (this->callback->onPropertyChanged != NULL) {
            this->callback->onPropertyChanged(this->mm, label, propName, propValue);
        }
    }

    void onConfigGroupChanged(const char *groupName, const char *newConfigName) {
        if (this->callback->onConfigGroupChanged != NULL) {
            this->callback->onConfigGroupChanged(this->mm, groupName, newConfigName);
        }
    }

    void onSystemConfigurationLoaded() {
        if (this->callback->onSystemConfigurationLoaded != NULL) {
            this->callback->onSystemConfigurationLoaded(this->mm);
        }
    }

    void onPixelSizeChanged(double newPixelSizeUm) {
        if (this->callback->onPixelSizeChanged != NULL) {
            this->callback->onPixelSizeChanged(this->mm, newPixelSizeUm);
        }
    }

    void onStagePositionChanged(char *label, double pos) {
        if (this->callback->onStagePositionChanged != NULL) {
            this->callback->onStagePositionChanged(this->mm, label, pos);
        }
    }

    void onXYStagePositionChanged(char *label, double xpos, double ypos) {
        if (this->callback->onXYStagePositionChanged != NULL) {
            this->callback->onXYStagePositionChanged(this->mm, label, xpos, ypos);
        }
    }

    void onExposureChanged(char *label, double newExposure) {
        if (this->callback->onExposureChanged != NULL) {
            this->callback->onExposureChanged(this->mm, label, newExposure);
        }
    }

    void onSLMExposureChanged(char *label, double newExposure) {
        if (this->callback->onSLMExposureChanged != NULL) {
            this->callback->onSLMExposureChanged(this->mm, label, newExposure);
        }
    }
};

DllExport void MM_RegisterCallback(MM_Session mm, struct MM_EventCallback *callback) {
    // If a callback has been registered, free it
    MM_FreeRegisteredCallback(mm);
    
    // Create and register the call back from C struct
    MM_CPP_EventCallback *cb = new MM_CPP_EventCallback(callback, mm);
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    core->registerCallback(cb);

    // Save the pointer
    mm_registered_callbacks[mm] = cb;
    return;
}

void MM_FreeRegisteredCallback(MM_Session mm) {
    std::map<MM_Session, MM_CPP_EventCallback*>::iterator it = mm_registered_callbacks.begin();
    while(it != mm_registered_callbacks.end()) {
        if (it->first == mm) {
            delete it->second;
            mm_registered_callbacks.erase(it);
        }
        it++;
    }
}

//
// Device listing
//

DllExport void MM_SetDeviceAdapterSearchPaths(MM_Session mm,
                                              const char **paths) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);

    std::vector<std::string> str_list;
    size_t i = 0;
    while (paths[i]) {
        str_list.push_back(std::string(paths[i]));
        i++;
    }

    core->setDeviceAdapterSearchPaths(str_list);
    return;
}

DllExport void MM_GetDeviceAdapterSearchPaths(MM_Session mm, char ***paths) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::vector<std::string> str_list = core->getDeviceAdapterSearchPaths();
    std_to_c_string_list(str_list, paths);
    return;
}

DllExport MM_Status MM_GetDeviceAdapterNames(MM_Session mm, char ***names) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::vector<std::string> str_list;
    try {
        str_list = core->getDeviceAdapterNames();
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }

    std_to_c_string_list(str_list, names);
    return MM_ErrOK;
}

DllExport MM_Status MM_GetAvailableDevices(MM_Session mm, const char *library,
                                           char ***names) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::vector<std::string> str_list;
    try {
        str_list = core->getAvailableDevices(library);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }

    std_to_c_string_list(str_list, names);
    return MM_ErrOK;
}

DllExport MM_Status MM_GetAvailableDeviceDescriptions(MM_Session mm,
                                                      const char *library,
                                                      char ***descriptions) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::vector<std::string> str_list;
    try {
        str_list = core->getAvailableDeviceDescriptions(library);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }

    std_to_c_string_list(str_list, descriptions);
    return MM_ErrOK;
}

DllExport MM_Status MM_GetAvailableDeviceTypes(MM_Session mm,
                                               const char *library,
                                               MM_DeviceType **types,
                                               size_t *len_types) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::vector<long> type_list;
    try {
        type_list = core->getAvailableDeviceTypes(library);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }

    *len_types = type_list.size();
    for (int i = 0; i < type_list.size(); i++) {
        (*types)[i] = (MM_DeviceType)(type_list[i]);
    }
    return MM_ErrOK;
}
//
// Generic device control
//
DllExport MM_Status MM_GetDevicePropertyNames(MM_Session mm, const char *label,
                                              char ***names) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::vector<std::string> str_list;
    try {
        str_list = core->getDevicePropertyNames(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }

    std_to_c_string_list(str_list, names);
    return MM_ErrOK;
}

DllExport MM_Status MM_HasProperty(MM_Session mm, const char *label,
                                   const char *prop_name,
                                   uint8_t *has_property) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *has_property = (bool)(core->hasProperty(label, prop_name));
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetProperty(MM_Session mm, const char *label,
                                   const char *prop_name, char **value) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str;
    try {
        str = core->getProperty(label, prop_name);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    std_to_c_string(str, value);
    return MM_ErrOK;
}

DllExport MM_Status MM_SetPropertyString(MM_Session mm, const char *label,
                                         const char *prop_name,
                                         const char *value) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setProperty(label, prop_name, value);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetPropertyBool(MM_Session mm, const char *label,
                                       const char *prop_name,
                                       const uint8_t value) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setProperty(label, prop_name, (const bool)(value !=0));
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetPropertyInt(MM_Session mm, const char *label,
                                      const char *prop_name,
                                      const int32_t value) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setProperty(label, prop_name, (const long)value);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetPropertyFloat(MM_Session mm, const char *label,
                                        const char *prop_name,
                                        const float value) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setProperty(label, prop_name, value);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetPropertyDouble(MM_Session mm, const char *label,
                                         const char *prop_name,
                                         const double value) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setProperty(label, prop_name, value);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetAllowedPropertyValues(MM_Session mm,
                                                const char *label,
                                                const char *prop_name,
                                                char ***values) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::vector<std::string> list;
    try {
        list = core->getAllowedPropertyValues(label, prop_name);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    std_to_c_string_list(list, values);
    return MM_ErrOK;
}

DllExport MM_Status MM_IsPropertyReadOnly(MM_Session mm, const char *label,
                                          const char *prop_name,
                                          uint8_t *read_only) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *read_only = (bool)core->isPropertyReadOnly(label, prop_name);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_IsPropertyPreInit(MM_Session mm, const char *label,
                                         const char *prop_name,
                                         uint8_t *pre_init) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *pre_init = (bool)core->isPropertyPreInit(label, prop_name);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_IsPropertySequenceable(MM_Session mm, const char *label,
                                              const char *prop_name,
                                              uint8_t *sequenceable) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *sequenceable = (bool)core->isPropertySequenceable(label, prop_name);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_HasPropertyLimits(MM_Session mm, const char *label,
                                         const char *prop_name,
                                         uint8_t *has_limit) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *has_limit = (bool)core->hasPropertyLimits(label, prop_name);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetPropertyLowerLimit(MM_Session mm, const char *label,
                                             const char *prop_name,
                                             double *lower_limit) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *lower_limit = core->getPropertyLowerLimit(label, prop_name);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetPropertyUpperLimit(MM_Session mm, const char *label,
                                             const char *prop_name,
                                             double *upper_limit) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *upper_limit = core->getPropertyUpperLimit(label, prop_name);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetPropertyType(MM_Session mm, const char *label,
                                       const char *prop_name,
                                       MM_PropertyType *type) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *type = (MM_PropertyType)core->getPropertyType(label, prop_name);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_DeviceBusy(MM_Session mm, const char *label,
                                  uint8_t *busy) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *busy = (bool)core->deviceBusy(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_DeviceTypeBusy(MM_Session mm, MM_DeviceType type,
                                      uint8_t *busy) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *busy = (bool)core->deviceTypeBusy((MM::DeviceType)type);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

//
// Manage current devices
//
DllExport MM_Status MM_SetCameraDevice(MM_Session mm, const char *label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setCameraDevice(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetShutterDevice(MM_Session mm, const char *label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setShutterDevice(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetFocusDevice(MM_Session mm, const char *label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setFocusDevice(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetXYStageDevice(MM_Session mm, const char *label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setXYStageDevice(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetAutoFocusDevice(MM_Session mm, const char *label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setXYStageDevice(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport void MM_GetCameraDevice(MM_Session mm, char **label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str = core->getCameraDevice();
    std_to_c_string(str, label);
    return;
}

DllExport void MM_GetShutterDevice(MM_Session mm, char **label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str = core->getShutterDevice();
    std_to_c_string(str, label);
    return;
}

DllExport void MM_GetFocusDevice(MM_Session mm, char **label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str = core->getFocusDevice();
    std_to_c_string(str, label);
    return;
}

DllExport void MM_GetXYStageDevice(MM_Session mm, char **label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str = core->getXYStageDevice();
    std_to_c_string(str, label);
    return;
}

DllExport void MM_GetAutoFocusDevice(MM_Session mm, char **label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str = core->getAutoFocusDevice();
    std_to_c_string(str, label);
    return;
}

//
// Image acquisition
//
DllExport MM_Status MM_SetROI(MM_Session mm, int x, int y, int x_size,
                              int y_size) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setROI(x, y, x_size, y_size);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetROI(MM_Session mm, int *x, int *y, int *x_size,
                              int *y_size) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->getROI(*x, *y, *x_size, *y_size);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_ClearROI(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->clearROI();
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetExposure(MM_Session mm, double exp) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setExposure(exp);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetExposure(MM_Session mm, double *exp) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *exp = core->getExposure();
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SnapImage(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->snapImage();
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetImage(MM_Session mm, uint8_t **ptr_buffer) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *ptr_buffer = (uint8_t *)(core->getImage());
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetImageOfChannel(MM_Session mm, uint16_t channel, uint8_t **ptr_buffer) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *ptr_buffer = (uint8_t *)(core->getImage(channel));
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport void MM_GetImageWidth(MM_Session mm, uint16_t *width) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *width = (uint16_t)(core->getImageWidth());
    return;
}

DllExport void MM_GetImageHeight(MM_Session mm, uint16_t *height) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *height = (uint16_t)(core->getImageHeight());
    return;
}

DllExport void MM_GetBytesPerPixel(MM_Session mm, uint8_t *bytes) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *bytes = (uint8_t)(core->getBytesPerPixel());
    return;
}

DllExport void MM_GetImageBitDepth(MM_Session mm, uint8_t *bit_depth) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *bit_depth = (uint8_t)(core->getImageBitDepth());
    return;
}

DllExport void MM_GetNumberOfComponents(MM_Session mm,
                                        uint8_t *n_components) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *n_components = (uint8_t)(core->getNumberOfComponents());
    return;
}

DllExport void MM_GetNumberOfCameraChannels(MM_Session mm,
                                            uint8_t *n_channels) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *n_channels = (uint8_t)(core->getNumberOfCameraChannels());
    return;
}

DllExport void MM_GetImageBufferSize(MM_Session mm, uint32_t *len) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *len = (uint32_t)(core->getImageBufferSize());
    return;
}

DllExport MM_Status MM_SetCircularBufferMemoryFootprint(MM_Session mm,
                                                        uint32_t size_MB) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setCircularBufferMemoryFootprint((unsigned)size_MB);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport void MM_GetCircularBufferMemoryFootprint(MM_Session mm,
                                                   uint32_t *size_MB) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *size_MB = (uint32_t)core->getCircularBufferMemoryFootprint();
    return;
}

DllExport MM_Status MM_InitializeCircularBuffer(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->initializeCircularBuffer();
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_ClearCircularBuffer(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->clearCircularBuffer();
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_StartSequenceAcquisition(MM_Session mm,
                                                int16_t num_images,
                                                double interval_ms,
                                                uint8_t stop_on_overflow) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->startSequenceAcquisition(num_images, interval_ms,
                                       (bool)(stop_on_overflow != 0));
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_StartContinuousSequenceAcquisition(MM_Session mm,
                                                          double interval_ms) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->startContinuousSequenceAcquisition(interval_ms);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_StopSequenceAcquisition(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->stopSequenceAcquisition();
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport void MM_IsSequenceRunning(MM_Session mm, uint8_t *status) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *status = (bool)(core->isSequenceRunning());
    return;
}

DllExport MM_Status MM_GetLastImage(MM_Session mm, uint8_t **ptr_buffer) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *ptr_buffer = (uint8_t *)(core->getLastImage());
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_PopNextImage(MM_Session mm, uint8_t **ptr_buffer) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *ptr_buffer = (uint8_t *)(core->popNextImage());
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport void MM_GetRemainingImageCount(MM_Session mm, int16_t *count) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *count = (int16_t)core->getRemainingImageCount();
    return;
}

DllExport void MM_GetBufferTotalCapacity(MM_Session mm, int16_t *capacity) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *capacity = (int16_t)core->getBufferTotalCapacity();
    return;
}

DllExport void MM_GetBufferFreeCapacity(MM_Session mm, int16_t *capacity) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *capacity = (int16_t)core->getBufferFreeCapacity();
    return;
}

DllExport void MM_IsBufferOverflowed(MM_Session mm, uint8_t *overflowed) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *overflowed = (bool)(core->isBufferOverflowed());
    return;
}

//
// Shutter control
//

DllExport MM_Status MM_SetShutterOpen(MM_Session mm, const char *label,
                                  uint8_t is_open) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setShutterOpen(label, (bool)(is_open != 0));
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetShutterOpen(MM_Session mm, const char *label,
                                  uint8_t *is_open) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *is_open = (bool)(core->getShutterOpen(label));
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }

    return MM_ErrOK;
}

//
// Autofocus control
//
DllExport void MM_GetLastFocusScore(MM_Session mm, double *score) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *score = core->getLastFocusScore();
    return;
}

DllExport void MM_GetCurrentFocusScore(MM_Session mm, double *score) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    *score = core->getCurrentFocusScore();
    return;
}

DllExport MM_Status MM_EnableContinuousFocus(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->enableContinuousFocus(true);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_DisableContinuousFocus(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->enableContinuousFocus(false);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_IsContinuousFocusEnabled(MM_Session mm,
                                                uint8_t *enabled) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *enabled = (bool)(core->isContinuousFocusEnabled());
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_IsContinuousFocusLocked(MM_Session mm, uint8_t *locked) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *locked = (bool)(core->isContinuousFocusLocked());
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_IsContinuousFocusDrive(
    MM_Session mm, const char *label, uint8_t *is_continuous_focus_drive) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *is_continuous_focus_drive =
            (bool)(core->isContinuousFocusDrive(label));
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_FullFocus(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->fullFocus();
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_IncrementalFocus(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->incrementalFocus();
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetAutoFocusOffset(MM_Session mm, double offset) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setAutoFocusOffset(offset);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetAutoFocusOffset(MM_Session mm, double *offset) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *offset = core->getAutoFocusOffset();
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

//
// State device control
//

DllExport MM_Status MM_SetState(MM_Session mm, const char *label,
                                int32_t state) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setState(label, (long)(state));
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetState(MM_Session mm, const char *label,
                                int32_t *state) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *state = (int32_t)(core->getState(label));
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetNumberOfStates(MM_Session mm, const char *label,
                                         int32_t *state) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *state = (int32_t)(core->getNumberOfStates(label));
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetStateLabel(MM_Session mm, const char *label,
                                     const char *state_label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setStateLabel(label, state_label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetStateLabel(MM_Session mm, const char *label,
                                     char **state_label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str;
    try {
        str = core->getStateLabel(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    std_to_c_string(str, state_label);
    return MM_ErrOK;
}

DllExport MM_Status MM_DefineStateLabel(MM_Session mm, const char *label,
                                        int32_t state,
                                        const char *state_label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->defineStateLabel(label, state, state_label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetStateLabels(MM_Session mm, const char *label,
                                      char ***state_labels) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::vector<std::string> list;
    try {
        list = core->getStateLabels(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    std_to_c_string_list(list, state_labels);
    return MM_ErrOK;
}

DllExport MM_Status MM_GetStateFromLabel(MM_Session mm, const char *label,
                                         const char *state_label,
                                         int32_t *state) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *state = core->getStateFromLabel(label, state_label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

//
// Focus (Z) stage control
//

DllExport MM_Status MM_SetPosition(MM_Session mm, const char *label,
                                   double position) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setPosition(label, position);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetPosition(MM_Session mm, const char *label,
                                   double *position) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *position = core->getPosition(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetRelativePosition(MM_Session mm, const char *label,
                                           double delta) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setRelativePosition(label, delta);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetOrigin(MM_Session mm, const char *label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setOrigin(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetAdapterOrigin(MM_Session mm, const char *label,
                                        double new_z_um) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setAdapterOrigin(label, new_z_um);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport void MM_SetFocusDirection(MM_Session mm, const char *label,
                                         int8_t sign) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    core->setFocusDirection(label, (int)sign);
    return;
}

DllExport MM_Status MM_GetFocusDirection(MM_Session mm, const char *label,
                                         int8_t *sign) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *sign = (int8_t)core->getFocusDirection(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

//
// XY stage control
//
DllExport MM_Status MM_SetXYPosition(MM_Session mm, const char *label, double x,
                                     double y) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setXYPosition(label, x, y);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetRelativeXYPosition(MM_Session mm, const char *label,
                                             double dx, double dy) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setXYPosition(label, dx, dy);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetXYPosition(MM_Session mm, const char *label,
                                     double *x, double *y) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->getXYPosition(label, *x, *y);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetXPosition(MM_Session mm, const char *label,
                                    double *x) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *x = core->getXPosition(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetYPosition(MM_Session mm, const char *label,
                                    double *y) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *y = core->getYPosition(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_Stop(MM_Session mm, const char *label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->stop(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_Home(MM_Session mm, const char *label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->home(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetOriginXY(MM_Session mm, const char *label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setOriginXY(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetOriginX(MM_Session mm, const char *label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setOriginX(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetOriginY(MM_Session mm, const char *label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setOriginY(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetAdpaterOriginXY(MM_Session mm, const char *label,
                                          double new_x_um, double new_y_um) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setAdapterOriginXY(label, new_x_um, new_y_um);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

//
// Hub and peripheral devices
//
DllExport MM_Status MM_SetParentLabel(MM_Session mm, const char *label,
                                      const char *parent_label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setParentLabel(label, parent_label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetParentLabel(MM_Session mm, const char *label,
                                      char **parent_label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str;
    try {
        str = core->getParentLabel(label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    std_to_c_string(str, parent_label);
    return MM_ErrOK;
}

DllExport MM_Status MM_GetInstalledDevices(MM_Session mm, const char *hub_label,
                                           char ***names) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::vector<std::string> list;
    try {
        list = core->getInstalledDevices(hub_label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    std_to_c_string_list(list, names);
    return MM_ErrOK;
}

DllExport MM_Status MM_GetInstalledDeviceDescription(MM_Session mm,
                                                     const char *hub_label,
                                                     const char *name,
                                                     char **descriptions) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str;
    try {
        str = core->getInstalledDeviceDescription(hub_label, name);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    std_to_c_string(str, descriptions);
    return MM_ErrOK;
}

DllExport MM_Status MM_GetLoadedPeripheralDevices(MM_Session mm,
                                                  const char *hub_label,
                                                  char ***labels) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::vector<std::string> list;
    try {
        list = core->getLoadedPeripheralDevices(hub_label);
    } catch (CMMError &e) {
        return MM_Status(e.getCode());
    }
    std_to_c_string_list(list, labels);
    return MM_ErrOK;
}

//
// Miscellaneous
//

DllExport void MM_GetUserId(MM_Session mm, char **userid) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str = core->getUserId();
    std_to_c_string(str, userid);
    return;
}

DllExport void MM_GetHostName(MM_Session mm, char **hostname) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str = core->getHostName();
    std_to_c_string(str, hostname);
    return;
}

DllExport void MM_GetMACAddresses(MM_Session mm, char ***addresses) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::vector<std::string> list = core->getMACAddresses();
    std_to_c_string_list(list, addresses);
    return;
}
}
