//go:build linux

package main

import (
	"github.com/godbus/dbus/v5"
)

const (
	networkManagerService              = "org.freedesktop.NetworkManager"
	networkManagerPath                 = dbus.ObjectPath("/org/freedesktop/NetworkManager")
	networkManagerInterface            = "org.freedesktop.NetworkManager"
	networkManagerGetDevicesMethod     = networkManagerInterface + ".GetDevices"
	dbusPropertiesInterface            = "org.freedesktop.DBus.Properties"
	dbusPropertiesGetMethod            = dbusPropertiesInterface + ".Get"
	networkManagerDeviceInterface      = "org.freedesktop.NetworkManager.Device"
	networkManagerWirelessInterface    = "org.freedesktop.NetworkManager.Device.Wireless"
	networkManagerAccessPointInterface = "org.freedesktop.NetworkManager.AccessPoint"
	networkManagerDeviceTypeProperty   = networkManagerDeviceInterface + ".DeviceType"
	networkManagerDeviceStateProperty  = networkManagerDeviceInterface + ".State"
	networkManagerActiveAccessPoint    = networkManagerWirelessInterface + ".ActiveAccessPoint"
	networkManagerBitrateProperty      = networkManagerWirelessInterface + ".Bitrate"
	networkManagerSSIDProperty         = networkManagerAccessPointInterface + ".Ssid"
	networkManagerStrengthProperty     = networkManagerAccessPointInterface + ".Strength"
	networkManagerNoObjectPath         = dbus.ObjectPath("/")
	networkManagerDeviceTypeEthernet   = uint32(1)
	networkManagerDeviceTypeWiFi       = uint32(2)
	networkManagerDeviceStateActivated = uint32(100)
	networkManagerKilobitsPerMegabit   = float64(1000)
)

func readNetworkLink() networkLink {
	connection, err := dbus.SystemBus()
	if err != nil {
		return networkLink{connectionType: networkConnectionUnknown}
	}

	manager := connection.Object(networkManagerService, networkManagerPath)
	var devicePaths []dbus.ObjectPath
	if err := manager.Call(networkManagerGetDevicesMethod, dbus.FlagNoAutoStart).Store(&devicePaths); err != nil {
		return networkLink{connectionType: networkConnectionUnknown}
	}

	connectionType := networkConnectionNone
	for _, devicePath := range devicePaths {
		device := connection.Object(networkManagerService, devicePath)
		deviceType, ok := networkManagerUint32Property(device, networkManagerDeviceTypeProperty)
		if !ok {
			continue
		}
		deviceState, ok := networkManagerUint32Property(device, networkManagerDeviceStateProperty)
		if !ok || deviceState != networkManagerDeviceStateActivated {
			continue
		}

		if deviceType == networkManagerDeviceTypeWiFi {
			return readNetworkManagerWiFiLink(connection, device)
		}
		if deviceType == networkManagerDeviceTypeEthernet {
			connectionType = networkConnectionEthernet
		}
	}

	return networkLink{available: true, connectionType: connectionType}
}

func readNetworkManagerWiFiLink(connection *dbus.Conn, device dbus.BusObject) networkLink {
	link := networkLink{available: true, connectionType: networkConnectionWiFi}
	activeAccessPoint, ok := networkManagerObjectPathProperty(device, networkManagerActiveAccessPoint)
	if !ok || activeAccessPoint == networkManagerNoObjectPath {
		return link
	}

	accessPoint := connection.Object(networkManagerService, activeAccessPoint)
	if ssid, ok := networkManagerBytesProperty(accessPoint, networkManagerSSIDProperty); ok {
		networkName := string(ssid)
		link.networkName = &networkName
	}
	if strength, ok := networkManagerByteProperty(accessPoint, networkManagerStrengthProperty); ok {
		signalPercent := int(strength)
		link.signalPercent = &signalPercent
	}
	if bitrate, ok := networkManagerUint32Property(device, networkManagerBitrateProperty); ok {
		linkBitrate := float64(bitrate) / networkManagerKilobitsPerMegabit
		link.linkBitrateMbps = &linkBitrate
	}
	return link
}

func networkManagerUint32Property(object dbus.BusObject, property string) (uint32, bool) {
	value, ok := networkManagerPropertyByQualifiedName(object, property)
	if !ok {
		return 0, false
	}
	return networkManagerUint32Value(value)
}

func networkManagerByteProperty(object dbus.BusObject, property string) (byte, bool) {
	value, ok := networkManagerPropertyByQualifiedName(object, property)
	if !ok {
		return 0, false
	}
	return networkManagerByteValue(value)
}

func networkManagerBytesProperty(object dbus.BusObject, property string) ([]byte, bool) {
	value, ok := networkManagerPropertyByQualifiedName(object, property)
	if !ok {
		return nil, false
	}
	return networkManagerBytesValue(value)
}

func networkManagerObjectPathProperty(object dbus.BusObject, property string) (dbus.ObjectPath, bool) {
	value, ok := networkManagerPropertyByQualifiedName(object, property)
	if !ok {
		return networkManagerNoObjectPath, false
	}
	return networkManagerObjectPathValue(value)
}

func networkManagerPropertyByQualifiedName(object dbus.BusObject, qualifiedName string) (dbus.Variant, bool) {
	separatorIndex := len(qualifiedName)
	for separatorIndex > 0 && qualifiedName[separatorIndex-1] != '.' {
		separatorIndex--
	}
	if separatorIndex == 0 || separatorIndex == len(qualifiedName) {
		return dbus.Variant{}, false
	}
	interfaceName := qualifiedName[:separatorIndex-1]
	propertyName := qualifiedName[separatorIndex:]
	var value dbus.Variant
	err := object.Call(dbusPropertiesGetMethod, dbus.FlagNoAutoStart, interfaceName, propertyName).Store(&value)
	return value, err == nil
}
