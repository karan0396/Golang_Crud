package controller

import "bootcamp/model"

func ValidationReq(u *model.ReqUser)bool{
	if len(u.FName)+len(u.Lname)>30{
		return false
	}
	if len(u.Password)<8 || len(u.Password)>20{
		return false
	}
	if len(u.Email)>30{
		return false
	}
	return true
}

func ValidationRes(u *model.ResUser)bool{
	if len(u.FName)+len(u.Lname)>30{
		return false
	}
	if len(u.Email)>30{
		return false
	}
	return true
}