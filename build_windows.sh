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

CMDPATH=${GOPATH}/src/github.com/PucklaMotzer09/gohomeengine/src/cmd

if [ ! -f $CMDPATH/fsreplacer/fsreplacer.exe ]
then
	echo Compiling string in file replacer ...
	workingDir=$(pwd)
	cd $CMDPATH/fsreplacer
	$GOROOT/bin/go build
	cd $workingDir
fi

if [ ! -f $CMDPATH/bsreplacer/bsreplacer.exe ]
then
	echo Compiling backslash replacer ...
	workingDir=$(pwd)
	cd $CMDPATH/bsreplacer
	$GOROOT/bin/go build
	cd $workingDir
fi

GOHOME_PKG_CONFIG_PATH=${GOPATH}/src/github.com/PucklaMotzer09/gohomeengine/deps/pkg/windows
PKG_CONFIG_PATH=${PKG_CONFIG_PATH}:${GOHOME_PKG_CONFIG_PATH}
$CMDPATH/fsreplacer/fsreplacer.exe $GOHOME_PKG_CONFIG_PATH/assimp.pc '%GOPATH%' $($CMDPATH/bsreplacer/bsreplacer.exe $GOPATH)

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

$CMDPATH/fsreplacer/fsreplacer.exe $GOHOME_PKG_CONFIG_PATH/assimp.pc $($CMDPATH/bsreplacer/bsreplacer.exe $GOPATH) '%GOPATH%'

echo Placed executable in $placedPath

if [ $wantRun -eq 1 ]
then
	exitCode=1
	$placedPath/$executableName.exe && exitCode=0
	exit $exitCode
fi

exit 0
