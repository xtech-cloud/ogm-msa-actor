package cache

import "time"

func Setup() {
	deviceUUID_device_map = make(map[string]*Device)
	guardUUID_guard_map = make(map[string]*Guard)
	domainUUID_guardUUIDS_map = make(map[string]map[string]string)
	domainUUID_domain_map = make(map[string]*Domain)
	domainUUID_applicationManifest_map = make(map[string]string)
	domainUUID_applicationMD5_map = make(map[string]string)

	//启动健康检测
	go checkHealthy()
}

func Cancel() {
}

func checkHealthy() {
	// 每秒健康值减1
	for {
		select {
		case <-time.After(time.Second):
			for _, v := range deviceUUID_device_map {
				if nil == v {
					continue
				}
				v.Healthy = v.Healthy - 1
				if v.Healthy < 0 {
					v.Healthy = 0
				}
			}
		}
	}
}
