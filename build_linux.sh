#!/bin/sh

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
wantRun=0
wantHelp=0
placedPath=""

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
	fi
done

if [ $wantHelp -eq 1 ]
then
	echo build_linux.sh [options]
	echo options:
	echo -install: uses go install
	echo -run: starts program after compilation
	echo -help: shows help page
	exit 0
fi

echo "Compiling ..."

if	[ $wantInstall -eq 1 ]
then
	exitCode=1
	$GOROOT/bin/go install github.com/PucklaMotzer09/gohomeengine && exitCode=0
	if [ $exitCode  -eq 0 ]
	then
		placedPath="$GOPATH/bin"
	else
		exit 1
	fi
else
	exitCode=1
	$GOROOT/bin/go build github.com/PucklaMotzer09/gohomeengine && exitCode=0
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
	echo Running ...
	exitCode=1
	$placedPath/gohomeengine && exitCode=0
	exit $exitCode
fi

exit 0
