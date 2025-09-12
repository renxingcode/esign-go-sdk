package types

// ESignCommonResponse e签宝通用响应结构体
type ESignCommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// GetESignTokenRequest 获取e签宝token的请求体
type GetESignTokenRequest struct {
	AppId     string `json:"appId"`
	Secret    string `json:"secret"`
	GrantType string `json:"grantType,default:client_credentials"`
}

// GetESignTokenResponse 获取e签宝token的响应体
type GetESignTokenResponse struct {
	ExpiresIn    string `json:"expiresIn"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

// GetESignTemplateDetailResponse 获取流程模版详情的响应体
type GetESignTemplateDetailResponse struct {
	OrgId              string              `json:"orgId"`              //机构账号ID
	SignTemplateId     string              `json:"signTemplateId"`     //模板ID
	SignTemplateName   string              `json:"signTemplateName"`   //模板名称
	SignTemplateStatus int                 `json:"signTemplateStatus"` //模板状态
	DedicatedCloudId   string              `json:"dedicatedCloudId"`   //专有云ID
	Docs               []DetailDoc         `json:"docs"`               //底稿文件信息
	Attachments        []DetailAttachment  `json:"attachments"`        //附件信息
	Copiers            []DetailCopier      `json:"copiers"`            //设置抄送方信息
	Participants       []DetailParticipant `json:"participants"`       //参与方信息
}
type DetailDoc struct {
	FileId          string `json:"fileId"`
	FileName        string `json:"fileName"`
	FileDownloadUrl string `json:"fileDownloadUrl"`
}
type DetailAttachment struct {
	FileId      string `json:"fileId"`
	FileName    string `json:"fileName"`
	DownloadUrl string `json:"downloadUrl"`
}
type DetailCopier struct {
	CopierPsnInfo struct {
		PsnId      string `json:"psnId"`
		PsnName    string `json:"psnName"`
		PsnAccount string `json:"psnAccount"`
	} `json:"copierPsnInfo"`
	CopierOrgInfo struct {
		OrgId   string `json:"orgId"`
		OrgName string `json:"orgName"`
	} `json:"copierOrgInfo"`
}
type DetailParticipant struct {
	ParticipantId        string             `json:"participantId"`
	ParticipantFlag      string             `json:"participantFlag"`
	ParticipantType      int                `json:"participantType"`
	ParticipateBizType   string             `json:"participateBizType"`
	ParticipantSetMode   int                `json:"participantSetMode"`
	DraftOrder           int                `json:"draftOrder"`
	SignOrder            int                `json:"signOrder"`
	SealTypes            string             `json:"sealTypes"`
	WillingnessAuthModes string             `json:"willingnessAuthModes"`
	OrgParticipant       OrgParticipantData `json:"orgParticipant"`
	PsnParticipant       PsnParticipantData `json:"psnParticipant"`
	Components           []ComponentsData   `json:"components"`
}
type OrgParticipantData struct {
	OrgId      string `json:"orgId"`
	OrgName    string `json:"orgName"`
	Transactor struct {
		TransactorPsnID      string `json:"transactorPsnId"`
		TransactorPsnAccount string `json:"transactorPsnAccount"`
		TransactorName       string `json:"transactorName"`
	} `json:"transactor"`
}
type PsnParticipantData struct {
	PsnId      string `json:"psnId"`
	PsnName    string `json:"psnName"`
	PsnAccount string `json:"psnAccount"`
}
type ComponentsData struct {
	ComponentId             string `json:"componentId"`
	ComponentKey            string `json:"componentKey"`
	ComponentName           string `json:"componentName"`
	Required                bool   `json:"required"`
	ComponentType           int    `json:"componentType"`
	ComponentDefaultValue   string `json:"componentDefaultValue"`
	OriginCustomComponentID string `json:"originCustomComponentId"`
	FileID                  string `json:"fileId"`
	ComponentPosition       struct {
		ComponentPositionX float64 `json:"componentPositionX"`
		ComponentPositionY float64 `json:"componentPositionY"`
		ComponentPageNum   int     `json:"componentPageNum"`
	} `json:"componentPosition"`
	ComponentSpecialAttribute struct {
		DateFormat         string      `json:"dateFormat"`
		ImageType          string      `json:"imageType"`
		Options            interface{} `json:"options"`
		TableContent       interface{} `json:"tableContent"`
		NumberFormat       string      `json:"numberFormat"`
		ComponentMaxLength string      `json:"componentMaxLength"`
	} `json:"componentSpecialAttribute"`
	ComponentTextFormat struct {
		Font                int     `json:"font"`
		FontSize            float64 `json:"fontSize"`
		TextColor           string  `json:"textColor"`
		Bold                bool    `json:"bold"`
		Italic              bool    `json:"italic"`
		HorizontalAlignment string  `json:"horizontalAlignment"`
		VerticalAlignment   string  `json:"verticalAlignment"`
		TextLineSpacing     float64 `json:"textLineSpacing"`
	} `json:"componentTextFormat"`
	ComponentSize struct {
		ComponentWidth  int `json:"componentWidth"`
		ComponentHeight int `json:"componentHeight"`
	} `json:"componentSize"`
	RemarkSignField struct {
		InputType      int    `json:"inputType"`
		AiCheck        int    `json:"aiCheck"`
		RemarkContent  string `json:"remarkContent"`
		RemarkFontSize int    `json:"remarkFontSize"`
	} `json:"remarkSignField"`
	NormalSignField struct {
		ShowSignDate   int    `json:"showSignDate"`
		DateFormat     string `json:"dateFormat"`
		SignFieldStyle int    `json:"signFieldStyle"`
		SealSpecs      int    `json:"sealSpecs"`
		SealType       int    `json:"sealType"`
		MustSign       bool   `json:"mustSign"`
	} `json:"normalSignField"`
}

// 通过模板创建文件的请求体 https://open.esign.cn/doc/opendoc/saas_api/cz9d65_sh823i?searchText=
type CreateESignFileByTemplateRequest struct {
	TemplateDocFileId   string            `json:"templateId"`       //模板文件id,e签宝文档这里的说明不准确,文档中的字面意思是"模板id",实际上是"模板文件id"
	TemplateDocFileName string            `json:"name"`             //模板文件名称
	SimpleFormFields    map[string]string `json:"simpleFormFields"` //todo 可以根据情况改为结构体,不影响整体流程
}
type CreateESignFileByTemplateResponse struct {
	FileId         string `json:"fileId"`         //文件id
	FileName       string `json:"fileName"`       //文件名称
	DownloadUrl    string `json:"downloadUrl"`    //文件下载url
	DownloadOssUrl string `json:"downloadOssUrl"` //将下载的url上传到自己服务器之后的URL
}

// 创建合同签署数据请求体
type ESignCreateFlowRequestData struct {
	SignerName    string                 `json:"signer_name"`     // 签署人姓名
	SignerPhone   string                 `json:"signer_phone"`    // 签署人手机号
	CompanySealID string                 `json:"company_seal_id"` // 公司印章ID
	ContractFiles []ESignCreateFlowFiles `json:"contract_files"`  // 合同文件列表
}
type ESignCreateFlowFiles struct {
	TemplateId string `json:"template_id"` // 模板ID
	EFileId    string `json:"e_fileid"`    // e签宝文件ID
}

// e签宝发起签署的主结构体
type ESignCreateFlowRequest struct {
	Docs     []ESignCreateFlowDocs   `json:"docs"`
	FlowInfo ESignCreateFlowFlowInfo `json:"flowInfo"`
	Signers  []ESignCreateFlowSigner `json:"signers"`
}

// 文档结构体
type ESignCreateFlowDocs struct {
	FileID string `json:"fileId"`
}

// 流程信息结构体
type ESignCreateFlowFlowInfo struct {
	AutoArchive    bool                `json:"autoArchive"`
	AutoInitiate   bool                `json:"autoInitiate"`
	BusinessScene  string              `json:"businessScene"`
	FlowConfigInfo EsignFlowConfigInfo `json:"flowConfigInfo"`
}

// 流程配置信息结构体
type EsignFlowConfigInfo struct {
	NoticeDeveloperUrl string `json:"noticeDeveloperUrl"`
}

// 签署人结构体
type ESignCreateFlowSigner struct {
	PlatformSign  bool             `json:"platformSign,omitempty"`
	SignerAccount SignerAccount    `json:"signerAccount"`
	SignFields    []EsignSignField `json:"signfields"`
	ThirdOrderNo  string           `json:"thirdOrderNo,omitempty"`
}

// 签署人账户结构体
type SignerAccount struct {
	SignerAccountID string `json:"signerAccountId"`
}

// 签署字段结构体
type EsignSignField struct {
	AutoExecute        bool         `json:"autoExecute"`
	ActorIndentityType int          `json:"actorIndentityType,omitempty"`
	FileID             string       `json:"fileId"`
	SealID             string       `json:"sealId,omitempty"`
	PosBean            EsignPosBean `json:"posBean"`
	SignDateBeanType   int          `json:"signDateBeanType"`
	//SignDateBean       EsignSignDateBean `json:"signDateBean,omitempty"`
}

// 位置信息结构体
type EsignPosBean struct {
	PosPage int     `json:"posPage"`
	PosX    float64 `json:"posX"`
	PosY    float64 `json:"posY"`
}

// 签署日期信息结构体
type EsignSignDateBean struct {
	FontSize int     `json:"fontSize"`
	Format   string  `json:"format"`
	PosPage  int     `json:"posPage"`
	PosX     float64 `json:"posX"`
	PosY     float64 `json:"posY"`
}

// 创建合同签署数据的返回结构
type ESignCreateFlowResponseData struct {
	ContractSignDataId int64  `json:"contract_sign_data_id"`
	BusinessTypeID     int64  `json:"business_type_id"`
	BusinessDataID     string `json:"business_data_id"`
	ESignFlowId        string `json:"e_sign_flow_id"`
	ESignUrl           string `json:"e_sign_url"`
	ESignShortUrl      string `json:"e_sign_short_url"`
}

// 查询个人认证信息的返回结构
type PersonsIdentityData struct {
	AuthorizeUserInfo bool            `json:"authorizeUserInfo"`
	RealnameStatus    int             `json:"realnameStatus"`
	PsnId             string          `json:"psnId"`
	PsnAccount        ESignPsnAccount `json:"psnAccount"`
	PsnInfo           ESignPsnInfo    `json:"psnInfo"`
}
type ESignPsnAccount struct {
	AccountMobile string      `json:"accountMobile"`
	AccountEmail  interface{} `json:"accountEmail"`
}

type ESignPsnInfo struct {
	PsnName        string      `json:"psnName"`
	PsnNationality interface{} `json:"psnNationality"`
	PsnIDCardNum   string      `json:"psnIDCardNum"`
	PsnIDCardType  string      `json:"psnIDCardType"`
	BankCardNum    interface{} `json:"bankCardNum"`
	PsnMobile      string      `json:"psnMobile"`
}

// 创建个人认证信息的请求体
type CreateESignPersonsIdentityRequest struct {
	Name             string `json:"name"`
	Mobile           string `json:"mobile"`
	ThirdPartyUserId string `json:"thirdPartyUserId,optional"`
	Email            string `json:"email,optional"`
	IdType           string `json:"idType,optional"`
	IdNumber         string `json:"idNumber,optional"`
}

// 创建个人认证信息的返回结构
type CreateESignPersonsIdentityResponse struct {
	AccountId string `json:"accountId"`
}
