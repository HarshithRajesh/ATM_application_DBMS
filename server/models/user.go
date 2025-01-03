package models


type User struct{
  UserID      uint    `json:"user_id" gorm:"primary_key;auto_increment"`
  Username    string  `json:"username" gorm:"unique"`
  FirstName   string  `json:"first_name"`
  LastName    string  `json:"last_name"`
  Email       string  `json:"email" gorm:"unique"`
  PhoneNumber string  `json:"phone_number"`
  Password    string  `json:"password"`

}
func (User) TableName() string {
    return "users" // Replace this with your actual table name if it's different
}
