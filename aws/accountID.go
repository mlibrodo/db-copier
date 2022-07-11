package aws

import "github.com/mlibrodo/db-copier/config"

var AWSAccountId string

func init() {
	AWSAccountId = config.GetConfig().AWS.AccountId
}
