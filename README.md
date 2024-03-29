# MMCoreAPI

MMCoreAPI provides a portable C DLL interface to MMCore C++ API of Micro-Manager 2.0.

To use the C interface, include header `MMCoreC.h`, and link to `MMCoreC.dll`. A prebuilt DLL is provided at `lib/MMCoreC.dll`, which is built with `vc142` toolset.

Examples can be find in `MMCoreC/examples` folder.

This project was started because Micro-Manager 1.4 uses VC 2010 compiler, and a C++ DLL generated by VC 2010 is not binary-compatible with MSVC 2015-2022. Micro-Manager 2.0 has migrated to VC 2019, so it is fine to build and use a C++ DLL of `MMCore` to use it with MSVC 2015-2022. This C DLL not only works with later versions of MSVC, but also works with MinGW.

A Go interface (based on Cgo) is provided. However, it is no longer actively used by the author.

## Build
### Dependencies
* MSVC 2019 or 2022
  * Install vc142 if using MSVC 2022
* [vcpkg](https://vcpkg.io) installed to `C:\vcpkg` (to get `boost-1.77.0`)
* (Optional) [Ninja](https://ninja-build.org) (for faster build)

### Checkout the code
```
git clone --recursive https://github.com/Andeling/MMCoreAPI.git
```
No additional repositories are needed. 

### Build manually
* Open Developer Command Prompt for VS 2019.
  * (If using VS 2022, to select vc142 toolset, open Command Prompt and enter `"C:\Program Files\Microsoft Visual Studio\2022\Community\VC\Auxiliary\Build\vcvarsall.bat" x64 -vcvars_ver=14.2`)

* To install boost-1.77.0 with vcpkg and build:
  ```
  cd MMCoreAPI
  mkdir build
  cd build
  cmake .. -GNinja -DVCPKG_TARGET_TRIPLET:STRING=x64-windows-static-md -DCMAKE_TOOLCHAIN_FILE:STRING=C:/vcpkg/scripts/buildsystems/vcpkg.cmake
  cmake --build . --config Release
  ```

### Build in VS Code
It should be straightforward if VS 2019 is installed.

However, VS 2022 needs additional steps to select vc142 toolset.

#### Select vc142 toolset installed with VS 2022 in VS Code
It be should fine to use vc143. However, because Micro-Manager 2.0 officially uses vc142, it may be needed to use vc142 under certain conditions.

You can install vc142 with VS 2022.

To set up the kit in VS Code. Launch `CMake: Edit User-Local CMake Kits` in `Command Palette`. Make a copy of the VS 2022 `amd64` kit and add `vc142` in toolset parameter. This will pass `-T v142,host=x64` to `cmake.exe`

When configurating and building the project with CMake, you need to use `Visual Studio 17 2022` generator, because `Ninja` generator does not support this way of toolset selection. And VS Code does not support passing `-vcvars_ver=14.2` to `vcvarsall.bat`.

Example:
```json
{
    "name": "Visual Studio Community 2022 Release - amd64 (v142)",
    "visualStudio": "afb0c136",
    "visualStudioArchitecture": "x64",
    "preferredGenerator": {
      "name": "Visual Studio 17 2022",
      "platform": "x64",
      "toolset": "v142,host=x64"
    }
  }
```
