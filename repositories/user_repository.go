package repositories

import (
	"go-project-testing/config"
	"go-project-testing/models"
)

// UserRepository - like UserRepository extends JpaRepository<User, Long> in Spring Boot
type UserRepository struct{}

func (r *UserRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	return user, result.Error
}

func (r *UserRepository) Create(user *models.User) error {
	result := config.DB.Create(user)
	return result.Error
}

func (r *UserRepository) ExistsByEmail(email string) bool {
	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}
