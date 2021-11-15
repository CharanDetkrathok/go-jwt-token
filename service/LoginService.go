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
		email:    "EMAIL@ru.ac.th",
		password: "PASS",
	}
}

func (info *loginInfomation) LoginUser(email string, password string) bool {
	return info.email == email && info.password == password
}
