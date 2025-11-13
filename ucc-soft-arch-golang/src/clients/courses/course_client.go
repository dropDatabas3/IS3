package courses

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	customError "github.com/Guidotss/ucc-soft-arch-golang.git/src/domain/errors"
	"github.com/Guidotss/ucc-soft-arch-golang.git/src/model"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type CourseClient struct {
	Db *gorm.DB
}

func NewCourseClient(db *gorm.DB) *CourseClient {
	return &CourseClient{Db: db}
}

func (c *CourseClient) Create(course model.Course) (model.Course, error) {
	result := c.Db.Create(&course)
	if result.Error != nil {
		var err error
		switch {
		case strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint"):
			err = customError.NewError(
				"DUPLICATE_IDENTIFIER",
				"A course with the same identifier already exists. Please use a different identifier.",
				http.StatusConflict)
		case strings.Contains(result.Error.Error(), "connection"):
			err = customError.NewError(
				"DB_CONNECTION_ERROR",
				"Database connection error. Please try again later.",
				http.StatusInternalServerError)
		default:
			err = customError.NewError(
				"UNEXPECTED_ERROR",
				"An unexpected error occurred. Please try again later.",
				http.StatusInternalServerError)
		}
		return model.Course{}, err
	}
	return course, nil
}

func (c *CourseClient) GetAll(filter string) (model.Courses, error) {
	var courses model.Courses
	var rawResults []map[string]interface{}
	if filter == "" {
		err := c.Db.Raw(
			`SELECT 
    		courses.*, 
    		categories.category_name, 
    		COALESCE(r.ratingavg, 0) as ratingavg
		FROM 
    		courses
		LEFT JOIN 
    		(SELECT course_id, AVG(rating) as ratingavg 
		FROM ratings
    	GROUP BY course_id) as r ON 
    		courses.id = r.course_id
		JOIN 
    		categories 
		ON 
    		courses.category_id = categories.id
		WHERE 				
			courses.deleted_at IS NULL;`).Scan(&rawResults).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, customError.NewError("NOT_FOUND", "There is no courses", http.StatusNotFound)
			}
			return nil, customError.NewError("DB_ERROR", "Error retrieving course from database", http.StatusInternalServerError)
		}
	} else {
		err := c.Db.Raw(
			`SELECT
				courses.*,
				categories.category_name,
				COALESCE(r.ratingavg, 0) as ratingavg
			FROM
				courses
			LEFT JOIN
				(SELECT course_id, AVG(rating) as ratingavg
				FROM ratings
				GROUP BY course_id) as r ON
				courses.id = r.course_id
			JOIN
				categories
			ON
				courses.category_id = categories.id
			WHERE
				courses.deleted_at IS NULL AND (
					LOWER(courses.course_name) LIKE LOWER(?) OR
					LOWER(courses.course_description) LIKE LOWER(?) OR
					LOWER(categories.category_name) LIKE LOWER(?)
				)`, "%"+filter+"%", "%"+filter+"%", "%"+filter+"%").Scan(&rawResults).Error
		if err != nil {
			fmt.Println("error: ", err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, customError.NewError("NOT_FOUND", "There is no courses", http.StatusNotFound)
			}
			return nil, customError.NewError("DB_ERROR", "Error retrieving course from database", http.StatusInternalServerError)
		}
	}
	for _, data := range rawResults {
		course := model.Course{
			Id:                parseUUID(data["id"]),
			CourseName:        data["course_name"].(string),
			CourseDescription: data["course_description"].(string),
			CoursePrice:       toFloat64(data["course_price"]),
			CourseDuration:    toInt(data["course_duration"]),
			CourseInitDate:    data["course_init_date"].(string),
			CourseState:       toBool(data["course_state"]),
			CourseCapacity:    toInt(data["course_capacity"]),
			CourseImage:       data["course_image"].(string),
			CategoryID:        parseUUID(data["category_id"]),
			Category: model.Category{
				CategoryName: data["category_name"].(string),
			},
			RatingAvg: toFloat64(data["ratingavg"]),
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (c *CourseClient) GetById(id uuid.UUID) (model.Course, error) {
	var rawResult map[string]interface{}
	err := c.Db.Raw(
		`SELECT courses.*, categories.category_name ,r.ratingavg
			FROM courses, 
				(SELECT course_id , AVG(rating) as ratingavg 
				 FROM ratings GROUP BY course_id) as r, 
				categories
			WHERE 
				courses.id = r.course_id AND 
				courses.deleted_at IS NULL AND
				courses.category_id = categories.id AND
				courses.id = ?`, id).Scan(&rawResult).Error
	fmt.Println("rawresult: ", rawResult)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Course{}, customError.NewError("NOT_FOUND", "Course not found", http.StatusNotFound)
		}
		return model.Course{}, customError.NewError("DB_ERROR", "Error retrieving course from database", http.StatusInternalServerError)
	}
	// When using Raw+Scan into a map, no rows results in a nil/empty map and no error; handle explicitly
	if rawResult == nil || len(rawResult) == 0 {
		return model.Course{}, customError.NewError("NOT_FOUND", "Course not found", http.StatusNotFound)
	}
	course := model.Course{
		Id:                parseUUID(rawResult["id"]),
		CourseName:        rawResult["course_name"].(string),
		CourseDescription: rawResult["course_description"].(string),
		CoursePrice:       toFloat64(rawResult["course_price"]),
		CourseDuration:    toInt(rawResult["course_duration"]),
		CourseInitDate:    rawResult["course_init_date"].(string),
		CourseState:       toBool(rawResult["course_state"]),
		CourseCapacity:    toInt(rawResult["course_capacity"]),
		CourseImage:       rawResult["course_image"].(string),
		CategoryID:        parseUUID(rawResult["category_id"]),
		Category: model.Category{
			CategoryName: rawResult["category_name"].(string),
		},
		RatingAvg: toFloat64(rawResult["ratingavg"]),
	}
	return course, nil
}
func (c *CourseClient) UpdateCourse(course model.Course) (model.Course, error) {

	result := c.Db.Table("courses").Where("id = ?", course.Id).Updates(&course)

	if result.Error != nil {
		var err error
		switch {
		case strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint"):
			err = customError.NewError(
				"DUPLICATE_IDENTIFIER",
				"A course with the same identifier or name already exists. Please use a different identifier or name.",
				http.StatusConflict)
		case strings.Contains(result.Error.Error(), "connection"):
			err = customError.NewError(
				"DB_CONNECTION_ERROR",
				"Database connection error. Please try again later.",
				http.StatusInternalServerError)
		default:
			err = customError.NewError(
				"UNEXPECTED_ERROR",
				"An unexpected error occurred. Please try again later.",
				http.StatusInternalServerError)
		}
		return model.Course{}, err
	}
	return course, nil
}

func (c *CourseClient) DeleteCourse(id uuid.UUID) error {
	result := c.Db.Where("id = ?", id).Delete(&model.Course{})
	if result.Error != nil {
		var err error
		switch {
		case strings.Contains(result.Error.Error(), "connection"):
			err = customError.NewError(
				"DB_CONNECTION_ERROR",
				"Database connection error. Please try again later.",
				http.StatusInternalServerError)
		default:
			err = customError.NewError(
				"UNEXPECTED_ERROR",
				"An unexpected error occurred. Please try again later.",
				http.StatusInternalServerError)
		}
		return err
	}
	return nil
}

// FUNCION PARA PARSEAR UUID
func parseUUID(value interface{}) uuid.UUID {
	if value != nil {
		id, _ := uuid.Parse(value.(string))
		return id
	}
	return uuid.Nil
}

// helpers to normalize sqlite vs postgres scan types
func toBool(v interface{}) bool {
	switch t := v.(type) {
	case bool:
		return t
	case int64:
		return t != 0
	case int:
		return t != 0
	case float64:
		return t != 0
	case []byte:
		s := string(t)
		return s == "1" || strings.EqualFold(s, "true")
	case string:
		return t == "1" || strings.EqualFold(t, "true")
	default:
		return false
	}
}

func toInt(v interface{}) int {
	switch t := v.(type) {
	case int:
		return t
	case int32:
		return int(t)
	case int64:
		return int(t)
	case float32:
		return int(t)
	case float64:
		return int(t)
	case []byte:
		// best effort parse
		s := string(t)
		var n int
		fmt.Sscanf(s, "%d", &n)
		return n
	default:
		return 0
	}
}

func toFloat64(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case []byte:
		s := string(t)
		var f float64
		fmt.Sscanf(s, "%f", &f)
		return f
	default:
		return 0
	}
}
