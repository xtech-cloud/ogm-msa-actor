package cache

func Setup() {
	deviceMap = make(map[string]*Device)
	profileMap = make(map[string]*Profile)
	propertyMap = make(map[string]string)
}

func Cancel() {
}
