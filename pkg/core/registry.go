package core

var services = make(map[string]interface{})

func RegisterService(name string, svc interface{}) {
	services[name] = svc
}

func GetServices[T interface{}]() []T {
	results := []T{}
	for _, v := range services {
		if svc, ok := v.(T); ok {
			results = append(results, svc)
		}
	}
	return results
}
