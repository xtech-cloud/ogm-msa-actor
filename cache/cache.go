package cache

func Setup() {
	deviceUUID_device_map = make(map[string]*Device)
	profileUUID_profile_map = make(map[string]*Profile)
    domainUUID_profileUUIDS_map = make(map[string][]string)
    domainUUID_domain_map = make(map[string]*Domain)
}

func Cancel() {
}
