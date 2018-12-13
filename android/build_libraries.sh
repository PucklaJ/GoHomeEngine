export ANDROID_API_LEVEL=24
export ANDROID_SYSROOT=$ANDROID_NDK_HOME/platforms/android-$ANDROID_API_LEVEL/arch-arm
export LIB_NAME=gohome

export CC=arm-linux-androideabi-gcc
export CGO_CFLAGS="-w -D__ANDROID_API__=${ANDROID_API_LEVEL} -I${ANDROID_NDK_HOME}/sysroot/usr/include -I${ANDROID_NDK_HOME}/sysroot/usr/include/arm-linux-androideabi --sysroot=${ANDROID_SYSROOT}"
export CGO_LDFLAGS="-L${ANDROID_NDK_HOME}/sysroot/usr/lib -L${ANDROID_NDK_HOME}/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/lib/gcc/arm-linux-androideabi/4.9.x/ --sysroot=${ANDROID_SYSROOT}"
export CGO_ENABLED=1
export GOOS=android
export GOARCH=arm
echo API Level: $ANDROID_API_LEVEL
echo SYSROOT: $ANDROID_SYSROOT
go build -v -tags static -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-soname,lib${lIB_NAME}.so" -o=android/libs/armeabi-v7a/lib${LIB_NAME}.so