package models

import "time"

// User represents a staff user (admin/librarian) in the users table.
type User struct {
	ID           string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador único del usuario del personal"`
	FirstName    string    `json:"first_name" gorm:"size:100;not null;comment:Nombre de pila"`
	LastName     string    `json:"last_name" gorm:"size:100;not null;comment:Apellidos"`
	Email        string    `json:"email" gorm:"size:255;not null;uniqueIndex;comment:Correo único para inicio de sesión"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;type:text;not null;comment:Hash bcrypt de la contraseña nunca expuesto por API"`
	Role         UserRole  `json:"role" gorm:"type:user_role;not null;default:admin;comment:Rol del personal admin o bibliotecario"`
	AvatarURL    *string   `json:"avatar_url,omitempty" gorm:"type:text;comment:URL opcional de imagen de perfil"`
	IsActive     bool      `json:"is_active" gorm:"not null;default:true;comment:Cuenta activa false suele indicar usuario pendiente o deshabilitado"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime;comment:Marca de tiempo de alta"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:Marca de tiempo de última modificación"`
}
