cmake_minimum_required(VERSION 3.18)

project(MMCoreAPI CXX)
set(CMAKE_CXX_STANDARD 14)

find_package(Boost 1.77 REQUIRED COMPONENTS thread)
include_directories(${Boost_INCLUDE_DIRS})

################################
#  MMDevice
################################
set(MMDEVICE_DIR "mmCoreAndDevices/MMDevice")

set(MMDEVICE_HDRS
    ${MMDEVICE_DIR}/Debayer.h
    ${MMDEVICE_DIR}/DeviceBase.h
    ${MMDEVICE_DIR}/DeviceThreads.h
    ${MMDEVICE_DIR}/DeviceUtils.h
    ${MMDEVICE_DIR}/FixSnprintf.h
    ${MMDEVICE_DIR}/ImageMetadata.h
    ${MMDEVICE_DIR}/ImgBuffer.h
    ${MMDEVICE_DIR}/MMDevice.h
    ${MMDEVICE_DIR}/MMDeviceConstants.h
    ${MMDEVICE_DIR}/ModuleInterface.h
    ${MMDEVICE_DIR}/Property.h
)

set(MMDEVICE_SRCS
    ${MMDEVICE_DIR}/Debayer.cpp
    ${MMDEVICE_DIR}/DeviceUtils.cpp
    ${MMDEVICE_DIR}/ImgBuffer.cpp
    ${MMDEVICE_DIR}/MMDevice.cpp
    ${MMDEVICE_DIR}/ModuleInterface.cpp
    ${MMDEVICE_DIR}/Property.cpp
)

add_library(MMDevice STATIC
    ${MMDEVICE_HDRS}
    ${MMDEVICE_SRCS}
)

target_compile_definitions(MMDevice PRIVATE
    WIN32
    NDEBUG
    _LIB
    MODULE_EXPORTS
)

################################
#  MMCore
################################
set(MMCORE_DIR "mmCoreAndDevices/MMCore")

set(MMCORE_HDRS
    ${MMCORE_DIR}/CircularBuffer.h
    ${MMCORE_DIR}/ConfigGroup.h
    ${MMCORE_DIR}/Configuration.h
    ${MMCORE_DIR}/CoreCallback.h
    ${MMCORE_DIR}/CoreProperty.h
    ${MMCORE_DIR}/CoreUtils.h
    ${MMCORE_DIR}/DeviceManager.h
    ${MMCORE_DIR}/Devices/AutoFocusInstance.h
    ${MMCORE_DIR}/Devices/CameraInstance.h
    ${MMCORE_DIR}/Devices/DeviceInstance.h
    ${MMCORE_DIR}/Devices/DeviceInstanceBase.h
    ${MMCORE_DIR}/Devices/DeviceInstances.h
    ${MMCORE_DIR}/Devices/GalvoInstance.h
    ${MMCORE_DIR}/Devices/GenericInstance.h
    ${MMCORE_DIR}/Devices/HubInstance.h
    ${MMCORE_DIR}/Devices/ImageProcessorInstance.h
    ${MMCORE_DIR}/Devices/MagnifierInstance.h
    ${MMCORE_DIR}/Devices/SerialInstance.h
    ${MMCORE_DIR}/Devices/ShutterInstance.h
    ${MMCORE_DIR}/Devices/SignalIOInstance.h
    ${MMCORE_DIR}/Devices/SLMInstance.h
    ${MMCORE_DIR}/Devices/StageInstance.h
    ${MMCORE_DIR}/Devices/StateInstance.h
    ${MMCORE_DIR}/Devices/XYStageInstance.h
    ${MMCORE_DIR}/Error.h
    ${MMCORE_DIR}/FrameBuffer.h
    ${MMCORE_DIR}/Host.h
    ${MMCORE_DIR}/LibraryInfo/LibraryPaths.h
    ${MMCORE_DIR}/LoadableModules/LoadedDeviceAdapter.h
    ${MMCORE_DIR}/LoadableModules/LoadedModule.h
    ${MMCORE_DIR}/LoadableModules/LoadedModuleImpl.h
    ${MMCORE_DIR}/LoadableModules/LoadedModuleImplWindows.h
    ${MMCORE_DIR}/Logging/GenericEntryFilter.h
    ${MMCORE_DIR}/Logging/GenericLinePacket.h
    ${MMCORE_DIR}/Logging/GenericLogger.h
    ${MMCORE_DIR}/Logging/GenericLoggingCore.h
    ${MMCORE_DIR}/Logging/GenericMetadata.h
    ${MMCORE_DIR}/Logging/GenericPacketArray.h
    ${MMCORE_DIR}/Logging/GenericPacketQueue.h
    ${MMCORE_DIR}/Logging/GenericSink.h
    ${MMCORE_DIR}/Logging/GenericStreamSink.h
    ${MMCORE_DIR}/Logging/Logger.h
    ${MMCORE_DIR}/Logging/Logging.h
    ${MMCORE_DIR}/Logging/Metadata.h
    ${MMCORE_DIR}/Logging/MetadataFormatter.h
    ${MMCORE_DIR}/LogManager.h
    ${MMCORE_DIR}/MMCore.h
    ${MMCORE_DIR}/MMEventCallback.h
    ${MMCORE_DIR}/PluginManager.h
    ${MMCORE_DIR}/Semaphore.h
    ${MMCORE_DIR}/Task.h
    ${MMCORE_DIR}/TaskSet.h
    ${MMCORE_DIR}/TaskSet_CopyMemory.h
    ${MMCORE_DIR}/ThreadPool.h
)

set(MMCORE_SRCS
    ${MMCORE_DIR}/CircularBuffer.cpp
    ${MMCORE_DIR}/Configuration.cpp
    ${MMCORE_DIR}/CoreCallback.cpp
    ${MMCORE_DIR}/CoreProperty.cpp
    ${MMCORE_DIR}/DeviceManager.cpp
    ${MMCORE_DIR}/Devices/AutoFocusInstance.cpp
    ${MMCORE_DIR}/Devices/CameraInstance.cpp
    ${MMCORE_DIR}/Devices/DeviceInstance.cpp
    ${MMCORE_DIR}/Devices/GalvoInstance.cpp
    ${MMCORE_DIR}/Devices/HubInstance.cpp
    ${MMCORE_DIR}/Devices/ImageProcessorInstance.cpp
    ${MMCORE_DIR}/Devices/MagnifierInstance.cpp
    ${MMCORE_DIR}/Devices/SerialInstance.cpp
    ${MMCORE_DIR}/Devices/ShutterInstance.cpp
    ${MMCORE_DIR}/Devices/SignalIOInstance.cpp
    ${MMCORE_DIR}/Devices/SLMInstance.cpp
    ${MMCORE_DIR}/Devices/StageInstance.cpp
    ${MMCORE_DIR}/Devices/StateInstance.cpp
    ${MMCORE_DIR}/Devices/XYStageInstance.cpp
    ${MMCORE_DIR}/Error.cpp
    ${MMCORE_DIR}/FrameBuffer.cpp
    ${MMCORE_DIR}/Host.cpp
    ${MMCORE_DIR}/LibraryInfo/LibraryPathsWindows.cpp
    ${MMCORE_DIR}/LoadableModules/LoadedDeviceAdapter.cpp
    ${MMCORE_DIR}/LoadableModules/LoadedModule.cpp
    ${MMCORE_DIR}/LoadableModules/LoadedModuleImpl.cpp
    ${MMCORE_DIR}/LoadableModules/LoadedModuleImplWindows.cpp
    ${MMCORE_DIR}/Logging/Metadata.cpp
    ${MMCORE_DIR}/LogManager.cpp
    ${MMCORE_DIR}/MMCore.cpp
    ${MMCORE_DIR}/PluginManager.cpp
    ${MMCORE_DIR}/Semaphore.cpp
    ${MMCORE_DIR}/Task.cpp
    ${MMCORE_DIR}/TaskSet.cpp
    ${MMCORE_DIR}/TaskSet_CopyMemory.cpp
    ${MMCORE_DIR}/ThreadPool.cpp
)


add_library(MMCore STATIC
    ${MMCORE_HDRS}
    ${MMCORE_SRCS}
)

target_compile_definitions(MMCore PRIVATE
    WIN32
    NDEBUG
    _LIB
    _WINDOWS
)

target_link_libraries(MMCore
    MMDevice
    Iphlpapi
    ${Boost_LIBRARIES}
)

################################
#  MMCoreC
################################

add_library(MMCoreC SHARED
    MMCoreC/MMCoreC.h
    MMCoreC/MMCoreC.cpp
    MMCoreC/dllmain.cpp
)

target_include_directories(MMCoreC PUBLIC
    MMCoreC
)

set_target_properties(MMCoreC PROPERTIES
    PUBLIC_HEADER   MMCoreC/MMCoreC.h
)

target_include_directories(MMCoreC PRIVATE
    ${MMCORE_DIR}
)

target_link_libraries(MMCoreC
    MMCore
)

################################
#  MMCoreC_Demo
################################
add_subdirectory(MMCoreC/examples)
