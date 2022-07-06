This is a quick example to list device adapters and load DemoCamera.

This example will be built when building the whole project with CMake and MSVC.

To build this example manually with GCC compiler from MinGW:
```
gcc -o MMCoreC_Demo -I.. -L../../lib -lMMCoreC main.c
```

`MMCoreC_Demo.exe` will look for `MMCoreC.dll` at runtime. One way to run it is to copy the dll to the same folder of the exe.
