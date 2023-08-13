package transport

func errorResponce(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}
