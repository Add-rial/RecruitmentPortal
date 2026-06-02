package utils

import (
	"database/sql"
	"recruitmentportal/models"
)

func JobsJSONCreator(rows *sql.Rows) ([]models.Job, error) {
	jobsMap := make(map[uint]*models.Job)
	var jobIDs []uint

	for rows.Next() {
		var (
			jobID              uint
			title              string
			description        string
			company            string
			companyDescription string
			companyContactMail string
			createdBy          uint

			skillID   sql.NullInt64
			skillName sql.NullString
		)
		err := rows.Scan(
			&jobID,
			&title,
			&description,
			&company,
			&companyDescription,
			&companyContactMail,
			&createdBy,
			&skillID,
			&skillName,
		)
		if err != nil {
			return nil, err
		}

		job, exists := jobsMap[jobID]
		if !exists {
			job = &models.Job{
				ID:                 jobID,
				Title:              title,
				Description:        description,
				Company:            company,
				CompanyDescription: companyDescription,
				CompanyContactMail: companyContactMail,
				CreatedBy:          createdBy,
				Skills:             []models.Skill{},
			}
			jobsMap[jobID] = job
			jobIDs = append(jobIDs, jobID)
		}

		if skillID.Valid && skillName.Valid {
			job.Skills = append(
				job.Skills,
				models.Skill{
					ID:   uint(skillID.Int64),
					Name: skillName.String,
				},
			)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	jobs := make([]models.Job, 0, len(jobIDs))
	for _, id := range jobIDs {
		jobs = append(jobs, *jobsMap[id])
	}

	return jobs, nil
}