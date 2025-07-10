package database

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/joaqu1m/goexpert-labs/modules/13/graph/model"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

func (r *CourseRepository) CreateCourse(title string, description *string, categoryID string) (*model.Course, error) {
	id := uuid.New().String()
	query := "INSERT INTO courses (id, title, description, category_id) VALUES ($1, $2, $3, $4) RETURNING id, title, description, category_id"
	var course model.Course
	course.Category = &model.Category{}
	err := r.db.QueryRow(query, id, title, description, categoryID).Scan(&course.ID, &course.Title, &course.Description, &course.Category.ID)
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *CourseRepository) GetAllCourses() ([]*model.Course, error) {
	rows, err := r.db.Query("SELECT id, title, description, category_id FROM courses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*model.Course
	for rows.Next() {
		var course model.Course
		course.Category = &model.Category{}
		if err := rows.Scan(&course.ID, &course.Title, &course.Description, &course.Category.ID); err != nil {
			return nil, err
		}
		courses = append(courses, &course)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return courses, nil
}
