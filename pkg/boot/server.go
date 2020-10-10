package boot

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func IPXEScriptHandler(c echo.Context) error {
	return c.String(http.StatusOK, `#!ipxe
set ubuntu http://tcboot:8080/boot/dists/ubuntu/20.04/
initrd ${ubuntu}/initrd
kernel ${ubuntu}/vmlinuz
imgargs vmlinuz initrd=initrd boot=casper ip=dhcp url=${ubuntu}/20.04.1-live-server-amd64.iso debian-installer/language=en
boot
`)
}
