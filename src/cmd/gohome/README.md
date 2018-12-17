# gohome

This is the build command for building a GoHome application.

## Description

+ gohome build|install|run OS={place_os_here} ARCH={place_arch_here} FRAME={GLFW|SDL2|GTK} RENDER={OpenGL|OpenGLES2|OpenGLES3|OpenGLES31} START={package_path_of_start_scene_struct} RELEASE|DEBUG
    - build builds the application for the choosen OS and ARCH
    - install uses go install or adb install on android
    - run runs the build application or runs the app on android
    - a main.go file is generated when not already there
        * In there the framework, renderer and start scene gets set
    - android
        * There are additional parameters: API={api_level} KEYSTORE={path_to_keystore_file} KEYALIAS={name_of_keystore_alias} KEYPWD={keystore_password} STOREPWD={store_password}
        * on android the build files are generated when not already there
    - all the parameters are optional and have default values.
    - some parameters (like KEYSTORE) can't have a default value and so are questioned later when they are needed
    - the default build configuration is debug
+ gohome generate OS=... FRAME=... RENDER=... START=...
    - generates the files needed for building
+ gohome set OS=... FRAME= ... etc.
    - sets all the paramters
+ gohome reset
    - resets all the parameters
+ If the command is used the first time a .gohomebuild file is generated in the working directory in which are all the configurations like OS, ARCH, FRAME, etc.