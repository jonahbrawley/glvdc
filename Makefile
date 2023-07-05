buildcmd:
	rsrc.exe -manifest main.manifest -o rsrc.syso
	go build

buildgui:
	rsrc.exe -manifest main.manifest -o rsrc.syso
	go build -ldflags="-H windowsgui"