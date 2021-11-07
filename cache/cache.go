package cache

func Setup() {
	deviceUUID_device_map = make(map[string]*Device)
	guardUUID_guard_map = make(map[string]*Guard)
    domainUUID_guardUUIDS_map = make(map[string][]string)
    domainUUID_domain_map = make(map[string]*Domain)
}

func Cancel() {
}
