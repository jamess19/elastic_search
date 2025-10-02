package business

import (
)

type BusinessRestApiClient struct{}

func GetBusinessRestApiClient() BusinessApiInterface {
	return &BusinessRestApiClient{}
}
