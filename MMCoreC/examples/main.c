#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "MMCoreC.h"

int main() {
    MM_Session mmc;
    MM_Open(&mmc);

    char *info;
    MM_GetVersionInfo(mmc, &info);
    printf("Version: %s\n", info);
    MM_StringFree(info);

    MM_GetAPIVersionInfo(mmc, &info);
    printf("API Version: %s\n\n", info);
    MM_StringFree(info);

    const char *adapter_path[2] =  {"C:/Program Files/Micro-Manager-2.0", NULL};
    MM_SetDeviceAdapterSearchPaths(mmc, adapter_path);

    char **path;
    MM_GetDeviceAdapterSearchPaths(mmc, &path);
    printf("Adapter path:");
    for (int i = 0; path[i]; i++) {
        printf(" %s", path[i]);
    }
    printf("\n");
    MM_StringListFree(path);

    char **adapter_list;
    MM_Status status;
    status = MM_GetDeviceAdapterNames(mmc, &adapter_list);
    if (status != MM_ErrOK) {
        printf("GetDeviceAdapterNames: error %d\n", status);
        return 1;
    }
    printf("Adapters:");
    for (int i = 0; adapter_list[i]; i++) {
        printf(" %s", adapter_list[i]);
    }
    printf("\n\n");
    MM_StringListFree(adapter_list);
    
    status = MM_LoadDevice(mmc, "Camera", "DemoCamera", "DCam");
    if (status != MM_ErrOK) {
        printf("Load DemoCamera: error %d\n", status);
        return 1;
    }
    printf("DemoCamera loaded.\n");

    status = MM_UnloadAllDevices(mmc);
    if (status != MM_ErrOK) {
        printf("UnloadAllDevices: error %d\n", status);
        return 1;
    }
    printf("All devices unloaded.\n");

    MM_Close(mmc);
    return 0;
}
