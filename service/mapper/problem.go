package mapper

import (
	gormAgent "github.com/OJ-lab/oj-lab-services/core/agent/gorm"
	"github.com/OJ-lab/oj-lab-services/service/model"
	"gorm.io/gorm"
)

func CreateProblem(problem model.Problem) error {
	db := gormAgent.GetDefaultDB()
	return db.Create(&problem).Error
}

func GetProblem(slug string) (*model.Problem, error) {
	db := gormAgent.GetDefaultDB()
	db_problem := model.Problem{}
	err := db.Model(&model.Problem{}).Preload("Tags").Where("Slug = ?", slug).First(&db_problem).Error
	if err != nil {
		return nil, err
	}

	return &db_problem, nil
}

func DeleteProblem(problem model.Problem) error {
	db := gormAgent.GetDefaultDB()
	return db.Delete(&model.Problem{Slug: problem.Slug}).Error
}

func UpdateProblem(problem model.Problem) error {
	db := gormAgent.GetDefaultDB()
	return db.Model(&model.Problem{Slug: problem.Slug}).Updates(problem).Error
}

type GetProblemOptions struct {
	Selection []string
	Slug      *string
	Title     *string
	Tags      []*model.AlgorithmTag
	Offset    *int
	Limit     *int
}

func buildTXByOptions(db *gorm.DB, options GetProblemOptions, isCount bool) *gorm.DB {
	tagsList := []string{}
	for _, tag := range options.Tags {
		tagsList = append(tagsList, tag.Slug)
	}
	tx := db.Model(&model.Problem{})
	if len(options.Selection) > 0 {
		tx = tx.Select(options.Selection)
	}
	if len(tagsList) > 0 {
		tx = tx.Joins("JOIN problem_algorithm_tags ON problem_algorithm_tags.problem_slug = problems.slug").
			Where("problem_algorithm_tags.algorithm_tag_slug in ?", tagsList)
	}
	if options.Slug != nil {
		tx = tx.Where("slug = ?", *options.Slug)
	}
	if options.Title != nil {
		tx = tx.Where("title = ?", *options.Title)
	}
	tx = tx.Distinct().
		Preload("Tags")
	if !isCount {
		if options.Offset != nil {
			tx = tx.Offset(*options.Offset)
		}
		if options.Limit != nil {
			tx = tx.Limit(*options.Limit)
		}
	}

	return tx
}

func CountProblemByOptions(options GetProblemOptions) (int64, error) {
	db := gormAgent.GetDefaultDB()
	var count int64

	tx := buildTXByOptions(db, options, true)
	err := tx.Count(&count).Error

	return count, err
}

func GetProblemListByOptions(options GetProblemOptions) ([]model.Problem, int64, error) {
	total, err := CountProblemByOptions(options)
	if err != nil {
		return nil, 0, err
	}

	db := gormAgent.GetDefaultDB()
	problemList := []model.Problem{}

	tx := buildTXByOptions(db, options, false)
	err = tx.Find(&problemList).Error
	if err != nil {
		return nil, 0, err
	}

	return problemList, total, nil
}

func GetTagsList(problem model.Problem) []string {
	tagsList := []string{}
	for _, tag := range problem.Tags {
		tagsList = append(tagsList, tag.Slug)
	}
	return tagsList
}
