cmake_minimum_required(VERSION 3.18)

project(MMCoreC_Demo C)

add_executable(MMCoreC_Demo
    main.c
)

target_link_libraries(MMCoreC_Demo
    MMCoreC
)

add_custom_command(TARGET MMCoreC_Demo POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy_if_different
    $<TARGET_RUNTIME_DLLS:MMCoreC_Demo>
    $<TARGET_FILE_DIR:MMCoreC_Demo>
    COMMAND_EXPAND_LISTS
)
