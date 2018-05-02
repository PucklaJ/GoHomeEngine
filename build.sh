#!/bin/bash

platform=""

for arg in "$@"
do
	if [ "$arg" = "-windows" ]
	then
		platform="windows"
	elif [ "$arg" = "-android" ]
	then
		platform="android"
	elif [ "$arg" = "-linux" ]
	then
		platform="linux"
	elif [ "$arg" = "-help" ]
	then
		echo build.sh [options]
		echo options:
		echo -windows: compiles for windows
		echo -linux: compiles for linux
		echo -android: compiles for android
		echo ""
		echo "-install: uses go install (-desktop only)"
		echo "-run: runs program after compilation (-desktop only)"
		exit 0
	fi
done

if [ "$platform" = "" ]
then
	$(dirname "$0")/build.sh -help
	exit 0
fi

if [ $platform = "linux" ]
then
	bash $(dirname "$0")/build_linux.sh "$@"
elif [ $platform = "windows" ]
then
	bash $(dirname "$0")/build_windows.sh "$@"
elif [ $platform = "android" ]
then
	bash $(dirname "$0")/build_android.sh "$@"
fi