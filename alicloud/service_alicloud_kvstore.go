package alicloud

import (
	"time"

	"github.com/alibaba/terraform-provider/alicloud/connectivity"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/denverdino/aliyungo/common"
)

type KvstoreService struct {
	client *connectivity.AliyunClient
}

func (s *KvstoreService) DescribeRKVInstanceById(id string) (instance *r_kvstore.DBInstanceAttribute, err error) {
	request := r_kvstore.CreateDescribeInstanceAttributeRequest()
	request.InstanceId = id
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeInstanceAttribute(request)
	})
	if err != nil {
		if IsExceptedError(err, InvalidKVStoreInstanceIdNotFound) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("KVStore instance", id))
		}
		return nil, err
	}
	resp, _ := raw.(*r_kvstore.DescribeInstanceAttributeResponse)
	if resp == nil || len(resp.Instances.DBInstanceAttribute) <= 0 {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("KVStore instance", id))
	}

	return &resp.Instances.DBInstanceAttribute[0], nil
}

func (s *KvstoreService) DescribeRKVInstancebackupPolicy(id string) (policy *r_kvstore.DescribeBackupPolicyResponse, err error) {
	request := r_kvstore.CreateDescribeBackupPolicyRequest()
	request.InstanceId = id
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeBackupPolicy(request)
	})
	if err != nil {
		if IsExceptedError(err, InvalidKVStoreInstanceIdNotFound) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("KVStore Instance Policy", id))
		}
		return nil, err
	}
	policy, _ = raw.(*r_kvstore.DescribeBackupPolicyResponse)

	if policy == nil {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("KVStore Instance Policy", id))
	}

	return
}

func (s *KvstoreService) WaitForRKVInstance(instanceId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := s.DescribeRKVInstanceById(instanceId)
		if err != nil && !NotFoundError(err) {
			return err
		}

		if instance != nil && instance.InstanceStatus == string(status) {
			break
		}

		if timeout <= 0 {
			return common.GetClientErrorFromString("Timeout")
		}

		timeout = timeout - DefaultIntervalMedium
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
	return nil
}
