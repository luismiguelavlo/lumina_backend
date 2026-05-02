package models

import "time"

// Student is a library member (student profile).
type Student struct {
	ID                     string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador interno del estudiante"`
	StudentIDCode          string    `json:"student_id_code" gorm:"size:50;not null;uniqueIndex;comment:Matrícula o código institucional visible ej LUM-2024-001"`
	FirstName              string    `json:"first_name" gorm:"size:100;not null;comment:Nombre de pila"`
	LastName               string    `json:"last_name" gorm:"size:100;not null;comment:Apellidos"`
	Email                  *string   `json:"email,omitempty" gorm:"size:255;uniqueIndex;comment:Correo de contacto opcional único si se informa no es login"`
	AvatarURL              *string   `json:"avatar_url,omitempty" gorm:"type:text;comment:URL opcional de foto del estudiante"`
	DepartmentID           *string   `json:"department_id,omitempty" gorm:"type:uuid;index;comment:Departamento al que pertenece si aplica"`
	DegreeLevel            *string   `json:"degree_level,omitempty" gorm:"size:100;comment:Nivel académico ej estudiante de grado"`
	Major                  *string   `json:"major,omitempty" gorm:"size:150;comment:Carrera o especialidad"`
	ExpectedGraduationYear *int      `json:"expected_graduation_year,omitempty" gorm:"comment:Año previsto de graduación"`
	IsActive               bool      `json:"is_active" gorm:"not null;default:true;comment:Perfil activo false indica baja lógica"`
	MemberSince            time.Time `json:"member_since" gorm:"autoCreateTime;comment:Fecha desde la que es miembro de la biblioteca"`
	RegisteredBy           *string   `json:"registered_by,omitempty" gorm:"type:uuid;comment:Usuario staff que dio de alta el perfil"`
	CreatedAt              time.Time `json:"created_at" gorm:"autoCreateTime;comment:Marca de tiempo de creación del registro"`
	UpdatedAt              time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:Marca de tiempo de última actualización"`

	// Dept is loaded via Preload for API responses (not a column on students).
	Dept *Department `json:"-" gorm:"foreignKey:DepartmentID"`
}
