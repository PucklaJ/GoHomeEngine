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
projectDir=$(pwd)
projectDir=${projectDir#${GOPATH}"/src/"}
topDirectory=${projectDir#*/}
for i in {1..20}
do
topDirectory=${topDirectory#*/}
done

for arg in "$@"
do
	if [ "$arg" = "-install" ]
	then
		wantInstall=1
	elif [ "$arg" = "-help" ]
	then
		wantHelp=1
	fi
done

if [ $wantHelp -eq 1 ]
then
	echo build_android.sh [options]
	echo options:
	echo -install: uses gomobile install
	echo -help: shows help page
	exit 0
fi

echo "Compiling $topDirectory ..."

if	[ $wantInstall -eq 1 ]
then
	exitCode=1
	$GOPATH/bin/gomobile install -target=android $projectDir && exitCode=0
	if [ $exitCode  -eq 0 ]
	then
		echo Installed $topDirectory.apk onto the android phone
		echo Placed $topDirectory.apk into $("pwd")
	else
		exit 1
	fi
else
	exitCode=1
	$GOPATH/bin/gomobile build -target=android $projectDir && exitCode=0
	if [ $exitCode -eq 0 ]
	then
		echo Placed $topDirectory.apk into $("pwd")
	else
		exit 1
	fi
fi

exit 0
