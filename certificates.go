package iot

type ListDeviceCertificatesRequest struct {
	AppId string `json:"app_id,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Marker string `json:"marker,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

type ListDeviceCertificatesResponse struct {
	Certificates []CertificatesRspDTO `json:"certificates"`
	Page         Page                 `json:"page"`
}

type CertificatesRspDTO struct {
	CertificateID string `json:"certificate_id"`
	CnName        string `json:"cn_name"`
	Owner         string `json:"owner"`
	Status        bool   `json:"status"`
	VerifyCode    string `json:"verify_code"`
	CreateDate    string `json:"create_date"`
	EffectiveDate string `json:"effective_date"`
	ExpiryDate    string `json:"expiry_date"`
}

type UploadDeviceCertificatesRequest struct {
	Content string `json:"content"`
	AppId string `json:"app_id,omitempty"`
}

type UploadDeviceCertificatesResponse struct {
	CertificateID string `json:"certificate_id"`
	CnName string `json:"cn_name"`
	Owner string `json:"owner"`
	Status bool `json:"status"`
	VerifyCode string `json:"verify_code"`
	CreateDate string `json:"create_date"`
	EffectiveDate string `json:"effective_date"`
	ExpiryDate string `json:"expiry_date"`
}
