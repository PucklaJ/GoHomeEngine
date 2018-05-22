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

PKG_CONFIG_PATH=${GOPATH}/src/github.com/PucklaMotzer09/gohomeengine/deps/pkg/windows

wantInstall=0
wantRun=0
wantHelp=0
wantRelease=0
placedPath=""
compileFlags=''
projectDir=$(pwd)
projectDir=${projectDir#${GOPATH}"/src/"}
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
		wantRun=0
		echo -run does not work on windows
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
	compileFlags='-ldflags "-s"'
fi

echo "Compiling $executableName ..."

if	[ $wantInstall -eq 1 ]
then
	exitCode=1
	$GOROOT/bin/go install $compileFlags $projectDir && exitCode=0
	if [ $exitCode  -eq 0 ]
	then
		placedPath="$GOPATH/bin"
	else
		exit 1
	fi
else
	exitCode=1
	$GOROOT/bin/go build $compileFlags $projectDir  && exitCode=0
	if [ $exitCode -eq 0 ]
	then
		placedPath=$("pwd")
	else
		exit 1
	fi
fi

echo Placed executable in $placedPath

if [ $wantRun -eq 1 ]
then
	exitCode=1
	$placedPath/$executableName && exitCode=0
	exit $exitCode
fi

exit 0