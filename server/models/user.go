package models


type User struct{
  UserID      uint    `json:"user_id" gorm:"primary_key;auto_increment"`
  FirstName   string  `json:"first_name"`
  LastName    string  `json:"last_name"`
  Email       string  `json:"email" gorm:"unique"`
  PhoneNumber string  `json:"phone_number"`
  Password    string  `json:"password"`

}
