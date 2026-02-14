#----------------------------------------------------------------
# Generated CMake target import file for configuration "RelWithDebInfo".
#----------------------------------------------------------------

# Commands may need to know the format version.
set(CMAKE_IMPORT_FILE_VERSION 1)

# Import target "Vulkan::Loader" for configuration "RelWithDebInfo"
set_property(TARGET Vulkan::Loader APPEND PROPERTY IMPORTED_CONFIGURATIONS RELWITHDEBINFO)
set_target_properties(Vulkan::Loader PROPERTIES
  IMPORTED_LOCATION_RELWITHDEBINFO "${_IMPORT_PREFIX}/lib/libvulkan.so.1.4.341"
  IMPORTED_SONAME_RELWITHDEBINFO "libvulkan.so.1"
  )

list(APPEND _cmake_import_check_targets Vulkan::Loader )
list(APPEND _cmake_import_check_files_for_Vulkan::Loader "${_IMPORT_PREFIX}/lib/libvulkan.so.1.4.341" )

# Commands beyond this point should not need to know the version.
set(CMAKE_IMPORT_FILE_VERSION)
