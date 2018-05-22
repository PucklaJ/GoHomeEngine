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

if [ ! -f $GOPATH/src/github.com/PucklaMotzer09/gohomeengine/cmd/fsreplacer/fsreplacer.exe ]
then
	echo Compiling string in file replacer ...
	workingDir=$(pwd)
	cd $GOPATH/src/github.com/PucklaMotzer09/gohomeengine/cmd/fsreplacer
	$GOROOT/bin/go build
	cd $workingDir
fi

if [ ! -f $GOPATH/src/github.com/PucklaMotzer09/gohomeengine/cmd/bsreplacer/bsreplacer.exe ]
then
	echo Compiling backslash replacer ...
	workingDir=$(pwd)
	cd $GOPATH/src/github.com/PucklaMotzer09/gohomeengine/cmd/bsreplacer
	$GOROOT/bin/go build
	cd $workingDir
fi

PKG_CONFIG_PATH=${GOPATH}/src/github.com/PucklaMotzer09/gohomeengine/deps/pkg/windows
$GOPATH/src/github.com/PucklaMotzer09/gohomeengine/cmd/fsreplacer/fsreplacer $PKG_CONFIG_PATH/assimp.pc '%GOPATH%' $($GOPATH/src/github.com/PucklaMotzer09/gohomeengine/cmd/bsreplacer/bsreplacer.exe $GOPATH)

wantInstall=0
wantRun=0
wantHelp=0
wantRelease=0
placedPath=""
compileFlags=''
projectDir=$(pwd)
executableName=${projectDir#*/}
for i in {1..20}
do
executableName=${executableName#*\\}
executableName=${executableName#*/}
done

for arg in "$@"
do
	if [ "$arg" = "-install" ]
	then
		wantInstall=1
	elif [ "$arg" = "-run" ]
	then
		wantRun=1
	elif [ "$arg" = "-help" ]
	then
		wantHelp=1
	elif [ "$arg" = "-release" ]
	then
		wantRelease=1
	fi
done

if [ $wantHelp -eq 1 ]
then
	echo build_windows.sh [options]
	echo options:
	echo -install: uses go install
	echo -run: starts program after compilation
	echo -help: shows help page
	echo -release: removes symbols
	exit 0
fi

if [ $wantRelease -eq 1 ]
then
	compileFlags='-ldflags -H=windowsgui'
fi

echo "Compiling $executableName ..."

if	[ $wantInstall -eq 1 ]
then
	exitCode=1
	$GOROOT/bin/go install $compileFlags && exitCode=0
	if [ $exitCode  -eq 0 ]
	then
		placedPath="$GOPATH/bin"
	else
		exit 1
	fi
else
	exitCode=1
	$GOROOT/bin/go build $compileFlags  && exitCode=0
	if [ $exitCode -eq 0 ]
	then
		placedPath=$("pwd")
	else
		exit 1
	fi
fi

if [ $wantRelease -eq 1 ]
then
	strip -s $placedPath/$executableName.exe
fi

$GOPATH/src/github.com/PucklaMotzer09/gohomeengine/cmd/fsreplacer/fsreplacer $PKG_CONFIG_PATH/assimp.pc $($GOPATH/src/github.com/PucklaMotzer09/gohomeengine/cmd/bsreplacer/bsreplacer.exe $GOPATH) '%GOPATH%'

echo Placed executable in $placedPath

if [ $wantRun -eq 1 ]
then
	exitCode=1
	$placedPath/$executableName.exe && exitCode=0
	exit $exitCode
fi

exit 0