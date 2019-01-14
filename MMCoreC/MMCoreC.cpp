#include "MMCoreC.h"

#include <stdlib.h>
#include <string.h>
#include "MMCore.h"

extern "C" {

DllExport void MM_Open(MM_Session* core) {
    *core = reinterpret_cast<MM_Session>(new CMMCore());
    return;
}

DllExport void MM_Close(MM_Session mm) {
    delete reinterpret_cast<CMMCore *>(mm);
}

DllExport void MM_GetVersionInfo(MM_Session mm, char *info, int len_info) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str = core->getVersionInfo();
    strcpy_s(info, len_info, str.c_str());
    return;
}

DllExport void MM_GetAPIVersionInfo(MM_Session mm, char *info, int len_info) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::string str = core->getAPIVersionInfo();
    strcpy_s(info, len_info, str.c_str());
    return;
}

DllExport void MM_Free(void *ptr) {
	free(ptr);
}

//
// Device initialization and setup
//

DllExport MM_Status MM_LoadDevice(MM_Session mm, const char* label, const char* module_name, const char* device_name) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->loadDevice(label, module_name, device_name);
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_UnloadDevice(MM_Session mm, const char* label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->unloadDevice(label);
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_UnloadAllDevices(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->unloadAllDevices();
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_InitializeAllDevices(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->initializeAllDevices();
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_InitializeDevice(MM_Session mm, const char* label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->initializeDevice(label);
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_Reset(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->reset();
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

//
// Device listing
//

DllExport void MM_SetDeviceAdapterSearchPaths(MM_Session mm, const char** paths, size_t len_paths) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);

    std::vector<std::string> vpath;
    vpath.reserve(len_paths);

    for (size_t i = 0; i < len_paths; i++) {
        vpath.push_back(std::string(paths[i]));
    }

    core->setDeviceAdapterSearchPaths(vpath);
}

DllExport void MM_GetDeviceAdapterSearchPaths(MM_Session mm, char*** paths, size_t *len_paths) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::vector<std::string> vpath = core->getDeviceAdapterSearchPaths();

    *len_paths = vpath.size();
    *paths = (char**)malloc(vpath.size() * sizeof(char*));
    for (size_t i = 0; i < vpath.size(); i++) {
        *paths[i] = (char*)malloc(vpath[i].size() + 1);
        strcpy_s(*paths[i], vpath[i].size() + 1, vpath[i].c_str());
    }
}

DllExport MM_Status MM_GetDeviceAdapterNames(MM_Session mm, char*** names, size_t *len_names) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    std::vector<std::string> vnames;
    try {
        vnames = core->getDeviceAdapterNames();
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }

    *len_names = vnames.size();
    *names = (char**)malloc(vnames.size() * sizeof(char*));
    for (size_t i = 0; i < vnames.size(); i++) {
        (*names)[i] = (char*)malloc(vnames[i].size() + 1);
        strcpy_s((*names)[i], vnames[i].size()+1, vnames[i].c_str());
    }
    return MM_ErrOK;
}


//
// Generic device control
//
DllExport MM_Status MM_GetDevicePropertyNames(MM_Session mm, const char* label, char*** names, size_t *len_names) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
	std::vector<std::string> vname;
    try {
         vname = core->getDevicePropertyNames(label);
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }

    *len_names = vname.size();
    *names = (char**)malloc(vname.size() * sizeof(char*));
    for (size_t i = 0; i < vname.size(); i++) {
        *names[i] = (char*)malloc(vname[i].size() + 1);
        strcpy_s(*names[i], vname[i].size() + 1, vname[i].c_str());
    }
}

DllExport MM_Status MM_HasProperty(MM_Session mm, const char* label, const char* prop_name, int* has_property) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *has_property = (int)core->hasProperty(label, prop_name);
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetProperty(MM_Session mm, const char* label, const char* prop_name, char** value, size_t* len_value) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
	std::string str;
    try {
        str = core->getProperty(label, prop_name);
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    *len_value = str.size();
    *value = (char*)malloc(str.size() + 1);
    strcpy_s(*value, str.size() + 1, str.c_str());
    return MM_ErrOK;
}

DllExport MM_Status MM_SetPropertyString(MM_Session mm, const char* label, const char* prop_name, const char* value) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setProperty(label, prop_name, value);
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetPropertyBool(MM_Session mm, const char* label, const char* prop_name, const int value) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setProperty(label, prop_name, (const bool)value);
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetPropertyInt(MM_Session mm, const char* label, const char* prop_name, const int32_t value) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setProperty(label, prop_name, (const long)value);
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetPropertyFloat32(MM_Session mm, const char* label, const char* prop_name, const float value) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setProperty(label, prop_name, value);
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SetPropertyFloat64(MM_Session mm, const char* label, const char* prop_name, const double value) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setProperty(label, prop_name, value);
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

//
// Manage current devices
//
DllExport MM_Status MM_SetCameraDevice(MM_Session mm, const char* label) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setCameraDevice(label);
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

//
// Image acquisition
//

DllExport MM_Status MM_SetExposure(MM_Session mm, double exp) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setExposure(exp);
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetExposure(MM_Session mm, double *exp) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *exp = core->getExposure();
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_SnapImage(MM_Session mm) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->snapImage();
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetImage(MM_Session mm, uint8_t **ptr_buffer) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *ptr_buffer = (uint8_t *)(core->getImage());
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetImageWidth(MM_Session mm, uint16_t *width) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *width = (uint16_t)(core->getImageWidth());
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetImageHeight(MM_Session mm, uint16_t *height) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *height = (uint16_t)(core->getImageHeight());
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetBytesPerPixel(MM_Session mm, uint8_t *bytes) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *bytes = (uint8_t)(core->getBytesPerPixel());
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetImageBitDepth(MM_Session mm, uint8_t *bit_depth) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *bit_depth = (uint8_t)(core->getImageBitDepth());
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetNumberOfComponents(MM_Session mm, uint8_t *n_components) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *n_components = (uint8_t)(core->getNumberOfComponents());
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetNumberOfCameraChannels(MM_Session mm, uint8_t *n_channels) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *n_channels = (uint8_t)(core->getNumberOfCameraChannels());
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetImageBufferSize(MM_Session mm, uint32_t *len) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *len = (uint32_t)(core->getImageBufferSize());
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

//
// State device control
//

DllExport MM_Status MM_SetState(MM_Session mm, const char* label, int32_t state) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        core->setState(label, (long)(state));
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetState(MM_Session mm, const char* label, int32_t *state) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *state = (int32_t)(core->getState(label));
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

DllExport MM_Status MM_GetNumberOfStates(MM_Session mm, const char* label, int32_t *state) {
    CMMCore *core = reinterpret_cast<CMMCore *>(mm);
    try {
        *state = (int32_t)(core->getNumberOfStates(label));
    } catch (CMMError& e) {
        return MM_Status(e.getCode());
    }
    return MM_ErrOK;
}

}
