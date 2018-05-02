#!/bin/bash

platform="desktop"

for arg in "$@"
do
	if [ "$arg" = "-desktop" ]
	then
		platform="desktop"
	elif [ "$arg" = "-android" ]
	then
		platform="android"
	elif [ "$arg" = "-help" ]
	then
		echo build.sh [options]
		echo options:
		echo -desktop: compiles for linux
		echo -android: compiles for android
		echo ""
		echo "-install: uses go install (-desktop only)"
		echo "-run: runs program after compilation (-desktop only)"
		exit 0
	fi
done

if [ $platform = "desktop" ]
then
	sh $(dirname "$0")/build_linux.sh "$@"
elif [ $platform = "android" ]
then
	sh $(dirname "$0")/build_android.sh "$@"
fi