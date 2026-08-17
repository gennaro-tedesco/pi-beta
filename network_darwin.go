//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework CoreLocation -framework CoreWLAN -framework Foundation
#import <CoreLocation/CoreLocation.h>
#import <CoreWLAN/CoreWLAN.h>
#include <stdlib.h>

typedef struct {
	int available;
	int connected;
	char *ssid;
	int rssi;
	double bitrate;
} NativeNetworkLink;

static NativeNetworkLink readNativeNetworkLink(void) {
	static CLLocationManager *locationManager;
	NativeNetworkLink result = {0, 0, NULL, 0, 0};
	@autoreleasepool {
		if (locationManager == nil) {
			locationManager = [[CLLocationManager alloc] init];
			if (@available(macOS 10.15, *)) {
				[locationManager requestWhenInUseAuthorization];
			}
		}

		CWInterface *interface = [[CWWiFiClient sharedWiFiClient] interface];
		if (interface == nil) {
			return result;
		}

		result.available = 1;
		if (![interface powerOn] || [interface wlanChannel] == nil) {
			return result;
		}

		result.connected = 1;
		NSString *ssid = [interface ssid];
		if (ssid != nil) {
			result.ssid = strdup([ssid UTF8String]);
		}
		result.rssi = (int)[interface rssiValue];
		result.bitrate = [interface transmitRate];
	}
	return result;
}
*/
import "C"

import "unsafe"

const nativeNetworkEnabled = 1

func readNetworkLink() networkLink {
	nativeLink := C.readNativeNetworkLink()
	if nativeLink.ssid != nil {
		defer C.free(unsafe.Pointer(nativeLink.ssid))
	}

	if nativeLink.available != nativeNetworkEnabled {
		return networkLink{connectionType: networkConnectionUnknown}
	}
	if nativeLink.connected != nativeNetworkEnabled {
		return networkLink{available: true, connectionType: networkConnectionNone}
	}

	signalPercent := normalizeRSSI(int(nativeLink.rssi))
	signalDbm := int(nativeLink.rssi)
	link := networkLink{
		available:      true,
		connectionType: networkConnectionWiFi,
		signalPercent:  &signalPercent,
		signalDbm:      &signalDbm,
	}
	if nativeLink.ssid != nil {
		networkName := C.GoString(nativeLink.ssid)
		link.networkName = &networkName
	}
	if nativeLink.bitrate > 0 {
		bitrate := float64(nativeLink.bitrate)
		link.linkBitrateMbps = &bitrate
	}
	return link
}
