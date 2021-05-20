package logm

import "strings"

type Environment int

const (
	// DefaultEnvironment is used if the environment the process runs in is not known.
	DefaultEnvironment Environment = iota

	// SystemdEnvironment indicates that the process is started and managed by systemd.
	SystemdEnvironment

	// ContainerEnvironment indicates that the process is running within a container (docker, k8s, rkt, ...).
	ContainerEnvironment

	// MacOSServiceEnvironment indicates that the process is running as a daemon on macOS (e.g. managed via launchctl).
	MacOSServiceEnvironment

	// WindowsServiceEnvironment indicates the the process is run as a windows service.
	WindowsServiceEnvironment

	// InvalidEnvironment indicates that the environment name given is unknown or invalid.
	InvalidEnvironment
)

// String returns the string representation the configured environment
func (v Environment) String() string {
	switch v {
	case DefaultEnvironment:
		return "default"
	case SystemdEnvironment:
		return "systemd"
	case ContainerEnvironment:
		return "container"
	case MacOSServiceEnvironment:
		return "macOS_service"
	case WindowsServiceEnvironment:
		return "windows_service"
	default:
		return "<invalid>"
	}
}

// ParseEnvironment returns the environment type by name.
// The parse is case insensitive.
// InvalidEnvironment is returned if the environment type is unknown.
func ParseEnvironment(in string) Environment {
	switch strings.ToLower(in) {
	case "default":
		return DefaultEnvironment
	case "systemd":
		return SystemdEnvironment
	case "container":
		return ContainerEnvironment
	case "macos_service":
		return MacOSServiceEnvironment
	case "windows_service":
		return WindowsServiceEnvironment
	default:
		return InvalidEnvironment
	}
}