package service

type LoginService interface {
	LoginUser(email string, password string) bool
}

type loginInfomation struct {
	email    string
	password string
}

func StaticLoginService() LoginService {
	return &loginInfomation{
		email:    "charan@ru.ac.th",
		password: "charan8",
	}
}

func (info *loginInfomation) LoginUser(email string, password string) bool {
	return info.email == email && info.password == password
}