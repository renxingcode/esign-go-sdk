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
