package repository
package repository




























}	return db.Where("id = ?", id).Take(entity).Errorfunc (r *Repository[T]) FindById(db *gorm.DB, entity *T, id any) error {}	return total, err	err := db.Model(new(T)).Where("id = ?", id).Count(&total).Error	var total int64func (r *Repository[T]) CountById(db *gorm.DB, id any) (int64, error) {}	return db.Delete(entity).Errorfunc (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {}	return db.Save(entity).Errorfunc (r *Repository[T]) Update(db *gorm.DB, entity *T) error {}	return db.Create(entity).Errorfunc (r *Repository[T]) Create(db *gorm.DB, entity *T) error {}	DB *gorm.DBtype Repository[T any] struct {import "gorm.io/gorm"