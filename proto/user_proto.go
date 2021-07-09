package proto

type CreateUserREQ struct {
	UserName string `json:"user_name,omitempty"`
	PassWord string `json:"pass_word,omitempty"`
}

type UpdatePssswordREQ struct {
	NewPassWord string `json:"new_pass_word,omitempty"`
}

type LoginREQ struct {
	UserName string `json:"user_name,omitempty"`
	PassWord string `json:"pass_word,omitempty"`
}

type LoginRES struct {
	Token string `json:"token,omitempty"`
}
