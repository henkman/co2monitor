# co2monitor

## To allow non-root to open the device
put following
```
SUBSYSTEM=="usb", ATTRS{idVendor}=="04d9", ATTRS{idProduct}=="a052", MODE="0666"
```
into
```
/etc/udev/rules.d/99-co2monitor.rules
```
