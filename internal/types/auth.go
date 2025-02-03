package types

import validator "github.com/rezakhademix/govalidator/v2"

type RegisterTeamValues struct {
	FirstName       string `json:"first_name" form:"first_name"`
	LastName        string `json:"last_name" form:"last_name"`
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm"`
}

func (values RegisterTeamValues) Validate() map[string]string {

	// if len(values.Password) < 8 {
	// 	errs.Password = "Password must be at least 8 characters."
	// }

	v := validator.New()

	v.RequiredString(values.FirstName, "first_name", "first name is required").
		MinString(values.FirstName, 1, "first_name", "first name must be 1 character")

	v.RequiredString(values.LastName, "last_name", "last name is required").
		MinString(values.LastName, 1, "last_name", "last name must be 1 character")

	v.Email(values.Email, "email", "email address is not valid")
	v.RequiredString(values.Password, "password", "password is required").
		MinString(values.Password, 8, "password", "password must be 8 characters")
	v.RequiredString(values.PasswordConfirm, "password_confirm", "password confirm is required").
		MinString(values.PasswordConfirm, 8, "password_confirm", "password must be 8 characters")

	v.CustomRule(!(values.Password != values.PasswordConfirm), "password_confirm", "passwords must match")

	if v.IsFailed() {
		return v.Errors()
	}

	return map[string]string{}
}

type LoginUserValues struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func (values LoginUserValues) Validate() map[string]string {
	v := validator.New()

	v.Email(values.Email, "email", "email address is not valid")
	v.RequiredString(values.Password, "password", "password is required").
		MinString(values.Password, 8, "password", "password must be 8 characters")

	if v.IsFailed() {
		return v.Errors()
	}

	return map[string]string{}
}

type ValidateUserParams struct {
	ID           string
	Email        string
	PasswordHash string
}

type RegisterFormType struct {
	Type string `json:"type" form:"type" query:"type"`
}
