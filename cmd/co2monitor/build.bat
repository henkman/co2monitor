rsrc -ico Papirus-Team-Papirus-Apps-Weather.ico -manifest co2monitor.exe.manifest -o rsrc.syso
go get -v -d
go build -ldflags="-H windowsgui -s -w"