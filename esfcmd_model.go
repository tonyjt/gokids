package gokids

type EsfCmdResponse struct{
	Result bool
	Data EsfCmdResponseData `json:"data"`
}

type EsfCmdResponseData struct{
	ITotalRecords int `json:"iTotalRecords"`
	AaData []EsfCmdAddr `json:"aaData"`
	ITotalDisplayRecords int `json:"iTotalDisplayRecords"`
}

type EsfCmdAddr struct{
	Address string `json:"address"`
	Port int `json:"port"`
	Appcode string `json:"appcode"`
	Enablestatus int `json:"enablestatus"`
}

type PTModelRes struct{
	Success bool
	Code int
	Msg string
}


