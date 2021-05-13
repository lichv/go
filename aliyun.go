package lichv

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

func GetSTSToken(access_id string,access_key string, rolearn string,rolesessionrole string, region string) (response sts.Credentials, err error){
	client, err := sts.NewClientWithAccessKey(region, access_id, access_key)
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = rolearn
	request.RoleSessionName = rolesessionrole
	resp, err :=client.AssumeRole(request)
	if err != nil {
		fmt.Println(err.Error())
		return sts.Credentials{"","","",""},nil
	}
	return (*resp).Credentials, nil
}
