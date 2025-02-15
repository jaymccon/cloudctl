# GoogleTest dependency
include(FetchContent)
FetchContent_Declare(
    googletest
    GIT_REPOSITORY https://github.com/google/googletest.git
    GIT_TAG release-1.12.1
)
FetchContent_MakeAvailable(googletest)

# Add test executable
add_executable(${PROJECT_NAME}_tests
    # Add your test source files here
    # example: test_feature1.cpp
)

# Link test executable with gtest and your library
target_link_libraries(${PROJECT_NAME}_tests
    PRIVATE
    GTest::gtest
    GTest::gtest_main
    ${PROJECT_NAME}
)

# Register tests
add_test(NAME ${PROJECT_NAME}_tests COMMAND ${PROJECT_NAME}_tests)