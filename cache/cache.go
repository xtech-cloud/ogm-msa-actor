package cache

func Setup() {
	deviceUUID_device_map = make(map[string]*Device)
	guardUUID_guard_map = make(map[string]*Guard)
	domainUUID_guardUUIDS_map = make(map[string]map[string]string)
	domainUUID_domain_map = make(map[string]*Domain)
	domainUUID_applicationManifest_map = make(map[string]string)
	domainUUID_applicationMD5_map = make(map[string]string)
}

func Cancel() {
}
