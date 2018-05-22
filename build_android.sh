#!/bin/bash

if [ -z ${GOPATH+x} ]
then
	echo GOPATH is not set
	exit 1
fi

if [ -z ${GOROOT+x} ]
then
	echo GOROOT is not set
	exit 1
fi

wantInstall=0
wantHelp=0
wantRelease=0
wantRun=0
wantLog=0
compileFlags=''
projectDir=$(pwd)
projectDir=${projectDir#${GOPATH}"/src/"}
executableName=${projectDir#*/}
for i in {1..20}
do
executableName=${executableName#*/}
done

for arg in "$@"
do
	if [ "$arg" = "-install" ]
	then
		wantInstall=1
	elif [ "$arg" = "-help" ]
	then
		wantHelp=1
	elif [ "$arg" = "-release" ]
	then
		wantRelease=1
	elif [ "$arg" = "-run" ]
	then
		wantRun=1
	elif [ "$arg" = "-log" ]
	then
		wantLog=1
	fi
done

if [ $wantHelp -eq 1 ]
then
	echo build_android.sh [options]
	echo options:
	echo -install: uses gomobile install
	echo -help: shows help page
	echo -release: removes debugging symbols
	echo -run: runs app on device
	echo -log: runs adb logcat after run
	exit 0
fi

if [ $wantRelease -eq 1 ]
then
	compileFlags='-ldflags "-w"'
fi

echo "Compiling $executableName ..."

if	[ $wantInstall -eq 1 ]
then
	exitCode=1
	$GOPATH/bin/gomobile install $compileFlags -target=android $projectDir && exitCode=0
	if [ $exitCode  -eq 0 ]
	then
		echo Installed $executableName.apk onto the android phone
		echo Placed $executableName.apk into $("pwd")
	else
		exit 1
	fi
else
	exitCode=1
	$GOPATH/bin/gomobile build $compileFlags -target=android $projectDir && exitCode=0
	if [ $exitCode -eq 0 ]
	then
		echo Placed $executableName.apk into $("pwd")
	else
		exit 1
	fi
fi

if [ $wantRun -eq 1 ]
then
	couldRun=0
	adb shell am start -n org.golang.todo.$executableName/org.golang.app.GoNativeActivity && couldRun=1
	if [ $wantLog -eq 1 -a $couldRun -eq 1 ]
	then
		adb logcat | grep GoLog
	fi

	if [ $couldRun -eq 0 ]
	then
		exit 1
	fi
fi



exit 0
