package main

import "github.com/godbus/dbus/v5"

func networkManagerUint32Value(value dbus.Variant) (uint32, bool) {
	result, ok := value.Value().(uint32)
	return result, ok
}

func networkManagerByteValue(value dbus.Variant) (byte, bool) {
	result, ok := value.Value().(byte)
	return result, ok
}

func networkManagerBytesValue(value dbus.Variant) ([]byte, bool) {
	result, ok := value.Value().([]byte)
	return result, ok
}

func networkManagerObjectPathValue(value dbus.Variant) (dbus.ObjectPath, bool) {
	result, ok := value.Value().(dbus.ObjectPath)
	return result, ok
}
