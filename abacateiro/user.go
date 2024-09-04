package application

import (
	"fmt"
)

// User representa um usuário do sistema
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"user_name"`
    Email    string `json:"user_email"`
    Password string `json:"user_password"`
    Document string `json:"user_document"`
}

// Validate verifica se os campos obrigatórios do usuário estão preenchidos
func (u *User) Validate() error {
    if u.Name == "" {
        return fmt.Errorf("name is required")
    }
    if u.Email == "" {
        return fmt.Errorf("email is required")
    }
    if u.Password == "" {
        return fmt.Errorf("password is required")
    }
    if u.Document == "" {
        return fmt.Errorf("document is required")
    }
    return nil
}

// UserFilterDTO representa os filtros para busca de usuários
type UserFilterDTO struct {
    ID       int    `json:"id"`
    Name     string `json:"user_name"`
    Email    string `json:"user_email"`
    Document string `json:"user_document"`
}

// UserUpdateDTO representa os dados para atualização de um usuário
type UserUpdateDTO struct {
    Name     string `json:"user_name"`
}

// UserService define as operações disponíveis para um usuário
type UserService interface {
    CreateUser(user User) (User, error)
    GetUser(id int) (User, error)
    GetUsers() ([]User, error)
    UpdateUser(user User) (User, error)
    DeleteUser(id int) error
}