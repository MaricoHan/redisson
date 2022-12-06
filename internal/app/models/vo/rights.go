package vo

type RegisterRequest struct {
	RegisterType uint64                 `json:"register_type"`
	OperationID  string                 `json:"operation_id"`
	UserID       string                 `json:"user_id"`
	ProductInfo  ProductInfo            `json:"product_info"`
	RightsInfo   RightsInfo             `json:"rights_info"`
	Authors      Authors                `json:"authors"`
	Copyrights   Copyrights             `json:"copyrights"`
	ContactNum   string                 `json:"contact_num"`
	Email        string                 `json:"email"`
	UrgentTime   uint32                 `json:"urgent_time"`
	CallbackURL  string                 `json:"callback_url"`
	AuthFile     string                 `json:"auth_file"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type EditRegisterRequest struct {
	RegisterType uint64                 `json:"register_type"`
	OperationID  string                 `json:"operation_id"`
	UserID       string                 `json:"user_id"`
	ProductInfo  ProductInfo            `json:"product_info"`
	RightsInfo   RightsInfo             `json:"rights_info"`
	Authors      Authors                `json:"authors"`
	Copyrights   Copyrights             `json:"copyrights"`
	ContactNum   string                 `json:"contact_num"`
	Email        string                 `json:"email"`
	UrgentTime   uint32                 `json:"urgent_time"`
	CallbackURL  string                 `json:"callback_url"`
	AuthFile     string                 `json:"auth_file"`
	Metadata     map[string]interface{} `json:"metadata"`
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
	Individuals []Individual `json:"authors_individual"`
	Corporates  []Corporate  `json:"authors_corporate"`
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

type UserAuthRequest struct {
	RegisterType       uint64             `json:"register_type"`
	OperationID        string             `json:"operation_id"`
	AuthType           uint32             `json:"auth_type"`
	AuthInfoIndividual AuthInfoIndividual `json:"auth_info_individual"`
	AuthInfoCorporate  AuthInfoCorporate  `json:"auth_info_corporate"`
	CallbackUrl        string             `json:"callback_url"`
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
	RegisterType       uint64             `json:"register_type"`
	OperationID        string             `json:"operation_id"`
	AuthType           uint32             `json:"auth_type"`
	AuthInfoIndividual AuthInfoIndividual `json:"auth_info_individual"`
	AuthInfoCorporate  AuthInfoCorporate  `json:"auth_info_corporate"`
	CallbackUrl        string             `json:"callback_url"`
}

type DeliveryRequest struct {
	RegisterType   uint64 `json:"register_type"`
	OperationID    string `json:"operation_id"`
	ProductID      string `json:"product_id"`
	CertificateNum string `json:"certificate_num"`
	Addr           string `json:"addr"`
	Postcode       string `json:"postcode"`
	Recipient      string `json:"recipient"`
	PhoneNum       string `json:"phone_num"`
}

type EditDeliveryRequest struct {
	RegisterType uint64 `json:"register_type"`
	Addr         string `json:"addr"`
	Postcode     string `json:"postcode"`
	Recipient    string `json:"recipient"`
	PhoneNum     string `json:"phone_num"`
}

type ChangeRequest struct {
}

type EditChangeRequest struct {
}

type TransferRequest struct {
}

type EditTransferRequest struct {
}

type RevokeRequest struct {
}

type EditRevokeRequest struct {
}
