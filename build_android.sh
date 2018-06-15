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

if [ -z ${GOHOMEENGINE_ROOT+x} ]
then
	export GOHOMEENGINE_ROOT=$GOPATH/src/github.com/PucklaMotzer09/gohomeengine
fi

wantInstall=0
wantHelp=0
wantRelease=0
wantRun=0
wantLog=0
wantBig=0
compileFlags=''
projectDir=$(pwd)
projectDir=${projectDir#${GOPATH}"/src/"}
executableName=${projectDir#*/}
CMDDIR=${GOHOMEENGINE_ROOT}/src/cmd
for i in {1..20}
do
executableName=${executableName#*/}
done

if [ ! -f $CMDDIR/fcbreplacer/fcbreplacer -a ! -f $CMDDIR/fcbreplacer/fcbreplacer.exe ]
then
	echo Compiling first character big replacer ...
	workingDir=$(pwd)
	cd $CMDDIR/fcbreplacer
	$GOROOT/bin/go build
	cd $workingDir
fi

appName=$($CMDDIR/fcbreplacer/fcbreplacer $executableName)
packageName=${executableName}p

if [ ! -f $CMDDIR/fcnreplacer/fcnreplacer -a ! -f $CMDDIR/fcnreplacer/fcnreplacer.exe ]
then
	echo Compiling first char number replacer ...
	workingDir=$(pwd)
	cd $CMDDIR/fcnreplacer
	$GOROOT/bin/go build
	cd $workingDir
fi

packageName=$($CMDDIR/fcnreplacer/fcnreplacer $packageName)

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
	elif [ "$arg" = "-big" ]
	then
		wantBig=1
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
	echo -big: sets the name of the app with a big first character
	exit 0
fi

if [ $wantRelease -eq 1 ]
then
	compileFlags='-ldflags "-w"'
fi

# Setup AndroidMainfest.xml
copiedManifest=0
if [ ! -f $(pwd)/AndroidMainfest.xml ]
then
	if [ ! -f $CMDDIR/fsreplacer/fsreplacer -a ! -f $CMDDIR/fsreplacer/fsreplacer.exe ]
	then
		echo Compiling string in file replacer ...
		workingDir=$(pwd)
		cd $CMDDIR/fsreplacer
		$GOROOT/bin/go build
		cd $workingDir
	fi
# Copy default AndroidManifest.xml
	cp $GOHOMEENGINE_ROOT/AndroidManifest.xml $(pwd)/AndroidManifest.xml
	copiedManifest=1
	$CMDDIR/fsreplacer/fsreplacer $(pwd)/AndroidManifest.xml '%APPNAME%' $appName
	$CMDDIR/fsreplacer/fsreplacer $(pwd)/AndroidManifest.xml '%PACKAGENAME%' $packageName
fi

echo "Compiling $executableName ..."
compilationFailed=1
if	[ $wantInstall -eq 1 ]
then
	$GOPATH/bin/gomobile install $compileFlags -target=android $projectDir && compilationFailed=0
	if [ $compilationFailed  -eq 0 ]
	then
		echo Installed $executableName.apk onto the android phone
		echo Placed $executableName.apk into $("pwd")
	fi
else
	$GOPATH/bin/gomobile build $compileFlags -target=android $projectDir && compilationFailed=0
	if [ $compilationFailed -eq 0 ]
	then
		echo Placed $executableName.apk into $("pwd")
	fi
fi

if [ $copiedManifest -eq 1 ]
then
	rm $(pwd)/AndroidManifest.xml
fi

if [ $compilationFailed -eq 1 ]
then
	exit 1
fi

if [ $wantRun -eq 1 ]
then
	couldRun=0
	adb shell am start -n org.golang.todo.$packageName/org.golang.app.GoNativeActivity && couldRun=1
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
