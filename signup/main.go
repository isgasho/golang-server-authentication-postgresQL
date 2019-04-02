package signup

import "golang.org/x/crypto/bcrypt"

//Signup is for users
type Signup struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//HashPassword is a method of Signup
func (s *Signup) HashPassword() (string, error) {
	by, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(by), nil
}

//ComparePassword is a method of Singup
func (s *Signup) ComparePassword(hp, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hp), []byte(p))
	if err != nil {
		return false, err
	}
	return true, nil
}
