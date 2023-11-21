package controllers

func CreateSkill(user User, skill *Skill) error {
	_, err := Db.Exec(`
		INSERT INTO skills (name, logo, textcolor, fillcolor, bgcolor, author)
		VALUES (?, ?, ?, ?, ?, ?)`,
		skill.Name, skill.Logo, skill.TextColor, skill.FillColor, skill.BgColor, user.ID)
	return err
}

func UpdateSkill(skill *Skill) error {
	_, err := Db.Exec(`
		UPDATE skills
		SET name = ?, logo = ?, textcolor = ?, fillcolor = ?, bgcolor = ?
		WHERE id = ?`,
		skill.Name, skill.Logo, skill.TextColor, skill.FillColor, skill.BgColor, skill.ID)
	return err
}

func DeleteSkill(skillID int) error {
	_, err := Db.Exec("DELETE FROM skills WHERE id = ?", skillID)
	return err
}

func CreateProject(user User, project *Project) error {
	_, err := Db.Exec(`
		INSERT INTO projects (title, description, url, image, classes, author)
		VALUES (?, ?, ?, ?, ?, ?)`,
		project.Title, project.Description, project.Url, project.Image, project.Classes, user.ID)
	return err
}

func UpdateProject(project *Project) error {
	_, err := Db.Exec(`
		UPDATE projects
		SET title = ?, description = ?, url = ?, image = ?, classes = ?
		WHERE id = ?`,
		project.Title, project.Description, project.Url, project.Image, project.Classes, project.ID)
	return err
}

func DeleteProject(projectID int) error {
	_, err := Db.Exec("DELETE FROM projects WHERE id = ?", projectID)
	return err
}
