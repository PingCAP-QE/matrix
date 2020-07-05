package utils

func Keys(_map map[string]interface{}) []string {
	var keys []string
	for k, _ := range _map {
		keys = append(keys, k)
	}
	return keys
}
