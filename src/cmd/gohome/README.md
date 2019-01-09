# gohome

This is the build command for building a GoHome application.
To install the tool go into cmd/gohome and execute go install

## Description

+ **gohome build|install|run|generate|clean|env|set|reset|export|help OS={place_os_here} ARCH={place_arch_here} FRAME={GLFW|SDL2|GTK|JS} RENDER={OpenGL|OpenGLES2|OpenGLES3|OpenGLES31|WebGL} START={name_of_start_scene_struct} RELEASE|DEBUG**
    - android
        * There are additional parameters: API={api_level} KEYSTORE={path_to_keystore_file} KEYALIAS={name_of_keystore_alias} KEYPWD={keystore_password} STOREPWD={store_password}
        * on android the build files are generated when not already there
    - all the parameters are optional and have default values.
    - some parameters (like KEYSTORE) can't have a default value and so are questioned later when they are needed
    - the default build configuration is debug

+ If the command is used the first time a .gohome.config file is generated in the working directory in which are all the configurations like OS, ARCH, FRAME, etc.
+ build, install and run automatically call generate if needed
    - if you want to regenerate the files needed for building
    - just execute gohome generate

### Command description

+ **gohome build**
    - builds the application for the choosen OS and ARCH
+ **gohome install**
    - uses go install or adb install on android
+ **gohome run**
    - runs the built application or runs the app on android
    - if OS=browser a server starts using python and a browser
    - starts with localhost:8000
+ **gohome generate OS=... FRAME=... RENDER=... START=...**
    - generates the files needed for building
    - a main.go file is generated when not already there
        * In there the framework, renderer and start scene gets set
+ **gohome clean**
    - executes go clean -r --cache and deletes all build files
+ **gohome env**
    - prints all set values for OS,ARCH etc.
    - if --all or -a is provided go env is executed additionally
+ **gohome set OS=... FRAME= ... etc.**
    - sets all the paramters
+ **gohome reset**
    - resets all the parameters and deletes the .gohome.config file
+ **gohome export**
    - builds the game packages all files into an export folder so that it can
    - be published
