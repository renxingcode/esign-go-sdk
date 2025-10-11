package tests

import (
	"github.com/renxingcode/esign-go-sdk"
	"github.com/renxingcode/esign-go-sdk/initialize"
	"github.com/renxingcode/esign-go-sdk/utils"
	"testing"
)

// TestGetESignPersonsIdentityInfo 测试查询个人认证信息 | go test tests/account_test.go -v -run TestGetESignPersonsIdentityInfo
func TestGetESignPersonsIdentityInfo(t *testing.T) {
	testClient, err := initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)

	psnAccount := "13945618971" //手机号或邮箱
	accountDetail, err := client.Account.GetESignPersonsIdentityInfo(psnAccount, true)
	if err != nil {
		t.Errorf("Failed to get account detail: %v", err)
	}
	t.Logf("accountDetail: %v", utils.JsonMarshalNoEscape(accountDetail))
}

// TestCreateESignPersonsIdentity 测试创建个人认证信息 | go test tests/account_test.go -v -run TestCreateESignPersonsIdentity
func TestCreateESignPersonsIdentity(t *testing.T) {
	testClient, err := initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)

	name := "张小雨"           //姓名
	mobile := "13945618971"    //手机号
	thirdPartyUserId := mobile //第三方用户ID,可以和手机号相同
	createResp, err := client.Account.CreateESignPersonsIdentity(name, mobile, thirdPartyUserId, true)
	if err != nil {
		t.Errorf("Failed to get account detail: %v", err)
	}
	t.Logf("createResp: %v", utils.JsonMarshalNoEscape(createResp))
}

/*
[问题备注]:
沙箱环境，我通过 https://smlopenapi.esign.cn/v1/accounts/createByThirdPartyUserId
请求参数：{"name":"张小雨","mobile":"13945618971","thirdPartyUserId":"13945618971","email":"","idType":"","idNumber":""}
创建完用户了，返回成功 {"accountId":"48adc0cfe8e94d13abd0552de0554172"}，

然后通过 https://smlopenapi.esign.cn/v3/persons/identity-info?psnAccount=13945618971 查询用户信息，仍然返回 {"code":1435203,"message":"账号不存在或已注销 :13945618971","data":null}

然后，我再次调用 https://smlopenapi.esign.cn/v1/accounts/createByThirdPartyUserId ，又告诉我：{"code":53000000,"message":"账号已存在","data":{"accountId":"48adc0cfe8e94d13abd0552de0554172"}}

[e签宝回复]
/v1/accounts/createByThirdPartyUserId  和  /v3/persons/identity-info 这两个接口  不是同一个账号体系，
/v1/accounts/createByThirdPartyUserId这个接口创建的账号只存在于你们的appid下面，而/v3/persons/identity-info这个接口查询的是我们SaaS下的数据，两个数据是不通的;
再次调用/v1/accounts/createByThirdPartyUserId提示”账号已存在“是因为，使用同一个ThirdPartyUserId创建过帐号了，再使用这个相同的值去创建账号的时候就会返回"账号已存在;

[结论]
也就是说，根据手机号获取accountId 直接通过 /v1/accounts/createByThirdPartyUserId 就可以了，
如果不存在自动创建，如果存在则忽略code：53000000，直接拿返回的 data.accountId 用就可以了。
不需要调用  /v3/persons/identity-info 这个接口去获取accountId.
*/
