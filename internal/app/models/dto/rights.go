package dto

type RegisterRequest struct {
	Code         string      `json:"code"`
	Module       string      `json:"module"`
	ProjectID    uint64      `json:"project_id"`
	RegisterType uint64      `json:"register_type"`
	OperationID  string      `json:"operation_id"`
	UserID       string      `json:"user_id"`
	ProductInfo  ProductInfo `json:"product_info"`
	RightsInfo   RightsInfo  `json:"rights_info"`
	Authors      Authors     `json:"authors"`
	Copyrights   Copyrights  `json:"copyrights"`
	ContactNum   string      `json:"contact_num"`
	Email        string      `json:"email"`
	UrgentTime   uint32      `json:"urgent_time"`
	CallbackURL  string      `json:"callback_url"`
	AuthFile     string      `json:"auth_file"`
	Metadata     []byte      `json:"metadata"`
}

type EditRegisterRequest struct {
	Code         string      `json:"code"`
	Module       string      `json:"module"`
	ProjectID    uint64      `json:"project_id"`
	RegisterType uint64      `json:"register_type"`
	OperationID  string      `json:"operation_id"`
	UserID       string      `json:"user_id"`
	ProductInfo  ProductInfo `json:"product_info"`
	RightsInfo   RightsInfo  `json:"rights_info"`
	Authors      Authors     `json:"authors"`
	Copyrights   Copyrights  `json:"copyrights"`
	ContactNum   string      `json:"contact_num"`
	Email        string      `json:"email"`
	UrgentTime   uint32      `json:"urgent_time"`
	CallbackURL  string      `json:"callback_url"`
	AuthFile     string      `json:"auth_file"`
	Metadata     []byte      `json:"metadata"`
}

type ProductInfo struct {
	Name          string `json:"name"`
	CatName       string `json:"cat_name"`
	CoverImg      string `json:"cover_img"`
	File          string `json:"file"`
	Description   string `json:"description"`
	CreateNatName string `json:"create_nat_name"`
	CreateTime    string `json:"create_time"`
	CreateAddr    string `json:"create_addr"`
	IsPublished   uint32 `json:"is_published"`
	PubAddr       string `json:"pub_addr"`
	PubTime       string `json:"pub_time"`
	PubChannel    uint32 `json:"pub_channel"`
	PubAnnex      string `json:"pub_annex"`
	Hash          string `json:"hash"`
}

type RightsInfo struct {
	Hold          uint32 `json:"hold"`
	HoldName      string `json:"hold_name"`
	HoldExp       string `json:"hold_exp"`
	RightDocument string `json:"right_document"`
}

type Authors struct {
	Individuals []Individual `json:"copyrights_individual"`
	Corporates  []Corporate  `json:"copyrights_corporate"`
}

type Copyrights struct {
	Individuals []Individual `json:"copyrights_individual"`
	Corporates  []Corporate  `json:"copyrights_corporate"`
}

type Individual struct {
	IsApplicant uint32 `json:"is_applicant"`
	RealName    string `json:"real_name"`
	AuthNum     string `json:"auth_num"`
}

type Corporate struct {
	IsApplicant uint32 `json:"is_applicant"`
	CardType    string `json:"card_type"`
	CompanyName string `json:"company_name"`
	AuthNum     string `json:"auth_num"`
}

type RegisterResponse struct {
	OperationID string `json:"operation_id"`
}

type EditRegisterResponse struct {
	OperationID string `json:"operation_id"`
}

type QueryRegisterRequest struct {
	Code         string `json:"code"`
	Module       string `json:"module"`
	ProjectID    uint64 `json:"project_id"`
	RegisterType uint64 `json:"register_type"`
	OperationID  string `json:"operation_id"`
}

type QueryRegisterResponse struct {
	OperationID       string   `json:"operation_id"`
	AuditStatus       uint32   `json:"audit_status"`
	AuditFile         []string `json:"audit_file"`
	AuditOpinion      string   `json:"audit_opinion"`
	CertificateStatus uint32   `json:"certificate_status"`
	CertificateNum    string   `json:"certificate_num"`
	ProductID         string   `json:"product_id"`
	CertificateURL    []string `json:"certificate_url"`
	BackTag           string   `json:"back_tag"`
}

type DictRequest struct {
	Code         string `json:"code"`
	Module       string `json:"module"`
	ProjectID    uint64 `json:"project_id"`
	RegisterType uint64 `json:"register_type"`
}

type DictResponse struct {
	ProCat       []KeyValueDetail `json:"pro_cat"`
	ProCreateNat []KeyValueDetail `json:"pro_create_nat"`
	IndustryCode []KeyValue       `json:"industry_code"`
	AutHold      []KeyValue       `json:"aut_hold"`
}

type RegionRequest struct {
	Code         string `json:"code"`
	Module       string `json:"module"`
	ProjectID    uint64 `json:"project_id"`
	ParentID     uint64 `json:"parent_id"`
	RegisterType uint64 `json:"register_type"`
}

type RegionResponse struct {
	Data []Region `json:"data"`
}

type KeyValueDetail struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Detail string `json:"detail"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Region struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	ParentID   uint64 `json:"parent_id"`
	ShortName  string `json:"short_name"`
	MergerName string `json:"merger_name"`
	PinYin     string `json:"pin_yin"`
}

type UserAuthRequest struct {
	Code               string             `json:"code"`
	Module             string             `json:"module"`
	ProjectID          uint64             `json:"project_id"`
	RegisterType       uint64             `json:"register_type"`
	OperationID        string             `json:"operation_id"`
	AuthType           uint32             `json:"auth_type"`
	AuthInfoIndividual AuthInfoIndividual `json:"auth_info_individual"`
	AuthInfoCorporate  AuthInfoCorporate  `json:"auth_info_corporate"`
	CallbackUrl        string             `json:"callback_url"`
}

type UserAuthResponse struct {
	OperationID      string `json:"operation_id"`
	UserID           string `json:"user_id"`
	AuditStatus      uint32 `json:"audit_status"`
	AuditInstruction string `json:"audit_instruction"`
}

type AuthInfoIndividual struct {
	RealName        string `json:"real_name"`
	IDCardNum       string `json:"idcard_num"`
	IDCardFimg      string `json:"idcard_fimg"`
	IDCardBimg      string `json:"idcard_bimg"`
	IDCardHimg      string `json:"idcard_himg"`
	IDCardStartDate string `json:"idcard_start_date"`
	IDCardEndDate   string `json:"idcard_end_date"`
	IDCardProvince  string `json:"idcard_province"`
	IDCardCity      string `json:"idcard_city"`
	IDCardArea      string `json:"idcard_area"`
	ContactNum      string `json:"contact_num"`
	ContactAddr     string `json:"contact_addr"`
	Postcode        string `json:"postcode"`
	Contact         string `json:"contact"`
	Email           string `json:"email"`
	IndustryCode    string `json:"industry_code"`
	IndustryName    string `json:"industry_name"`
}

type AuthInfoCorporate struct {
	CardType        string `json:"card_type"`
	CompanyName     string `json:"company_name"`
	BusLicNum       string `json:"bus_lic_num"`
	CompanyAddr     string `json:"company_addr"`
	BusLicImg       string `json:"bus_lic_img"`
	BusLicStartDate string `json:"bus_lic_start_date"`
	BusLicEndDate   string `json:"bus_lic_end_date"`
	BusLicProvince  string `json:"bus_lic_province"`
	BusLicCity      string `json:"bus_lic_city"`
	BusLicArea      string `json:"bus_lic_area"`
	Postcode        string `json:"postcode"`
	Contact         string `json:"contact"`
	ContactNum      string `json:"contact_num"`
	Email           string `json:"email"`
	IndustryCode    string `json:"industry_code"`
	IndustryName    string `json:"industry_name"`
}

type EditUserAuthRequest struct {
	Code               string             `json:"code"`
	Module             string             `json:"module"`
	ProjectID          uint64             `json:"project_id"`
	RegisterType       uint64             `json:"register_type"`
	OperationID        string             `json:"operation_id"`
	AuthType           uint32             `json:"auth_type"`
	AuthInfoIndividual AuthInfoIndividual `json:"auth_info_individual"`
	AuthInfoCorporate  AuthInfoCorporate  `json:"auth_info_corporate"`
	CallbackUrl        string             `json:"callback_url"`
}

type EditUserAuthResponse struct {
	Data string `json:"data"`
}

type QueryUserAuthRequest struct {
	Code         string `json:"code"`
	Module       string `json:"module"`
	ProjectID    uint64 `json:"project_id"`
	RegisterType uint64 `json:"register_type"`
	AuthType     uint32 `json:"auth_type"`
	AuthNum      string `json:"auth_num"`
}

type QueryUserAuthResponse struct {
	UserID           string `json:"user_id"`
	AuditStatus      uint32 `json:"audit_status"`
	AuditInstruction string `json:"audit_instruction"`
}

type DeliveryRequest struct {
	Code           string `json:"code"`
	Module         string `json:"module"`
	ProjectID      uint64 `json:"project_id"`
	RegisterType   uint64 `json:"register_type"`
	OperationID    string `json:"operation_id"`
	ProductID      string `json:"product_id"`
	CertificateNum string `json:"certificate_num"`
	Addr           string `json:"addr"`
	Postcode       string `json:"postcode"`
	Recipient      string `json:"recipient"`
	PhoneNum       string `json:"phone_num"`
}

type DeliveryResponse struct {
	OperationID string `json:"operation_id"`
}

type EditDeliveryRequest struct {
	Code         string `json:"code"`
	Module       string `json:"module"`
	ProjectID    uint64 `json:"project_id"`
	RegisterType uint64 `json:"register_type"`
	OperationID  string `json:"operation_id"`
	Addr         string `json:"addr"`
	Postcode     string `json:"postcode"`
	Recipient    string `json:"recipient"`
	PhoneNum     string `json:"phone_num"`
}

type EditDeliveryResponse struct {
}

type DeliveryInfoRequest struct {
	Code           string `json:"code"`
	Module         string `json:"module"`
	ProjectID      uint64 `json:"project_id"`
	RegisterType   uint64 `json:"register_type"`
	ProductID      string `json:"product_id"`
	CertificateNum string `json:"certificate_num"`
}

type DeliveryInfoResponse struct {
	ProductID      string `json:"product_id"`
	CertificateNum string `json:"certificate_num"`
	Addr           string `json:"addr"`
	Postcode       string `json:"postcode"`
	Recipient      string `json:"recipient"`
	PhoneNum       string `json:"phone_num"`
	ExpressNum     string `json:"express_num"`
	Status         string `json:"status"`
}

type ChangeRequest struct {
	Code                  string                `json:"code"`
	Module                string                `json:"module"`
	ProjectID             uint64                `json:"project_id"`
	RegisterType          uint64                `json:"register_type"`
	OperationID           string                `json:"operation_id"`
	ProductID             string                `json:"product_id"`
	CertificateNum        string                `json:"certificate_num"`
	Name                  string                `json:"name"`
	CatName               string                `json:"cat_name"`
	CopyrighterCorporate  CopyrighterCorporate  `json:"copyrighter_info_corporate"`
	CopyrighterIndividual CopyrighterIndividual `json:"copyrighter_info_individual"`
	ProofFiles            string                `json:"proof_files"`
	UrgentTime            uint32                `json:"urgent_time"`
}

type CopyrighterCorporate struct {
	CopyrighterType uint32 `json:"copyrighter_type"`
	CompanyName     string `json:"company_name"`
	BusLicImg       string `json:"bus_lic_img"`
}

type CopyrighterIndividual struct {
	RealName   string `json:"real_name"`
	IDCardFimg string `json:"idcard_fimg"`
	IDCardBimg string `json:"idcard_bimg"`
	IDCardHimg string `json:"idcard_himg"`
}

type ChangeResponse struct {
	OperationID string `json:"operation_id"`
}

type EditChangeRequest struct {
	Code                  string                `json:"code"`
	Module                string                `json:"module"`
	ProjectID             uint64                `json:"project_id"`
	RegisterType          uint64                `json:"register_type"`
	OperationID           string                `json:"operation_id"`
	Name                  string                `json:"name"`
	CatName               string                `json:"cat_name"`
	CopyrighterCorporate  CopyrighterCorporate  `json:"copyrighter_info_corporate"`
	CopyrighterIndividual CopyrighterIndividual `json:"copyrighter_info_individual"`
	ProofFiles            string                `json:"proof_files"`
}

type EditChangeResponse struct {
}

type ChangeInfoRequest struct {
	Code         string `json:"code"`
	Module       string `json:"module"`
	ProjectID    uint64 `json:"project_id"`
	RegisterType uint64 `json:"register_type"`
	OperationID  string `json:"operation_id"`
}

type ChangeInfoResponse struct {
	ProductId            string `json:"product_id"`             // 作品编号
	CertificateNum       string `json:"certificate_num"`        // 版权登记号
	Status               uint32 `json:"status"`                 // 版权变更状态，0：待审核; 1：登记成功; 2：登记失败; 3：复审核中
	ChangeCertificateUrl string `json:"change_certificate_url"` // 变更登记证书
	ErrorMessage         string `json:"error_message"`          // 登记失败时显示失败原因
	ChangeCertificateNum string `json:"change_certificate_num"` // 变更成功后返回的登记号

}

type TransferRequest struct {
	Code             string `json:"code"`
	Module           string `json:"module"`
	ProjectID        uint64 `json:"project_id"`
	RegisterType     uint64 `json:"register_type"`
	OperationID      string `json:"operation_id"`
	CertificateNum   string `json:"certificate_num"`
	ProductID        string `json:"product_id"`
	AuthorityName    string `json:"authority_name"`
	AuthorityIDType  uint32 `json:"authority_id_type"`
	AuthorityIDNum   string `json:"authority_id_num"`
	AuthoritedName   string `json:"authorited_name"`
	AuthoritedIDType uint32 `json:"authorited_id_type"`
	AuthoritedIDNum  string `json:"authorited_id_num"`
	AuthInstructions string `json:"auth_instructions"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	Scope            string `json:"scope"`
	ContractAmount   string `json:"contract_amount"`
	ContractFiles    string `json:"contract_files"`
	UrgentTime       uint32 `json:"urgent_time"`
}

type TransferResponse struct {
	OperationID string `json:"operation_id"`
}

type EditTransferRequest struct {
	Code             string `json:"code"`
	Module           string `json:"module"`
	ProjectID        uint64 `json:"project_id"`
	RegisterType     uint64 `json:"register_type"`
	OperationID      string `json:"operation_id"`
	AuthorityName    string `json:"authority_name"`
	AuthorityIDType  uint32 `json:"authority_id_type"`
	AuthorityIDNum   string `json:"authority_id_num"`
	AuthoritedName   string `json:"authorited_name"`
	AuthoritedIDType uint32 `json:"authorited_id_type"`
	AuthoritedIDNum  string `json:"authorited_id_num"`
	AuthInstructions string `json:"auth_instructions"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	Scope            string `json:"scope"`
	ContractAmount   string `json:"contract_amount"`
	ContractFiles    string `json:"contract_files"`
}

type EditTransferResponse struct {
}

type TransferInfoRequest struct {
	Code         string `json:"code"`
	Module       string `json:"module"`
	ProjectID    uint64 `json:"project_id"`
	RegisterType uint64 `json:"register_type"`
	OperationID  string `json:"operation_id"`
}

type TransferInfoResponse struct {
	ProductId              string `json:"product_id"`               // 作品编号
	CertificateNum         string `json:"certificate_num"`          // 版权登记号
	Status                 uint32 `json:"status"`                   // 登记状态，0:待审核、1:转让成功、2:转让失败
	ErrorMessage           string `json:"error_message"`            // 失败时显示失败原因
	TransferCertificateNum string `json:"transfer_certificate_num"` // 转让后证书登记号
	TransferCertificateUrl string `json:"transfer_certificate_url"` // 转让证书地址

}

type RevokeRequest struct {
	Code           string `json:"code"`
	Module         string `json:"module"`
	ProjectID      uint64 `json:"project_id"`
	RegisterType   uint64 `json:"register_type"`
	OperationID    string `json:"operation_id"`
	ProductID      string `json:"product_id"`
	CertificateNum string `json:"certificate_num"`
}

type RevokeResponse struct {
	OperationID string `json:"operation_id"`
}

type EditRevokeRequest struct {
	Code         string `json:"code"`
	Module       string `json:"module"`
	ProjectID    uint64 `json:"project_id"`
	RegisterType uint64 `json:"register_type"`
	OperationID  string `json:"operation_id"`
}

type EditRevokeResponse struct {
}

type RevokeInfoRequest struct {
	Code         string `json:"code"`
	Module       string `json:"module"`
	ProjectID    uint64 `json:"project_id"`
	RegisterType uint64 `json:"register_type"`
	OperationID  string `json:"operation_id"`
}

type RevokeInfoResponse struct {
	ProductId            string `json:"product_id"`
	CertificateNum       string `json:"certificate_num"`
	Status               uint32 `json:"status"`
	ErrMessage           string `json:"err_message"`
	RevokeCertificateNum string `json:"revoke_certificate_num"`
}
