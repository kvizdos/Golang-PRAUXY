IP=$(ipconfig getifaddr en0)
/usr/X11/bin/xhost + $IP
DISPLAY=$IP:0