#include <string.h>
#include "_cgo_export.h"

void c_onPropertyChanged(MM_Session mm, const char* label, const char* property, const char* value) {
    onPropertyChanged(mm, (char *)label, (char *)property, (char *)value);
}

void c_onStagePositionChanged(MM_Session mm, char* label, double pos) {
    onStagePositionChanged(mm, label, pos);
}

void c_registerCallback(MM_Session mm) {
    struct MM_EventCallback *callback = (struct MM_EventCallback *)malloc(sizeof(struct MM_EventCallback));
    memset(callback, 0, sizeof(*callback));

    callback->onPropertyChanged = &c_onPropertyChanged;
    callback->onStagePositionChanged = &c_onStagePositionChanged;
    MM_RegisterCallback(mm, callback);
}