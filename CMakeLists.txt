cmake_minimum_required(VERSION 3.22)
project(cloudctl VERSION 0.1)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

if(CMAKE_CXX_COMPILER_ID MATCHES "GNU|Clang")
    add_compile_options(-fmodules-ts)
endif()

add_executable(${PROJECT_NAME} src/main.cpp)

# Add include directories
# include_directories(include)

# # Add source files
# add_library(${PROJECT_NAME}
#     src/feature1.cpp
#     src/feature2.cpp
# )

# Enable testing
# enable_testing()
# add_subdirectory(tests)