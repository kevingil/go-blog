package models

import (
	"log"
)

// Projects is a model for home page user skills
type Skill struct {
	ID        int
	Name      string
	Logo      string
	TextColor string
	FillColor string
	BgColor   string
}

// Projects is a model for home page projects.
type Project struct {
	ID          int
	Title       string
	Description string
	Url         string
	Image       string
	Classes     string
}

func (user User) FindSkills() []*Skill {
	var skills []*Skill

	rows, err := Db.Query(`SELECT id, name, logo, textcolor, fillcolor, bgcolor FROM skills WHERE author = ?`, user.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        int
			name      string
			logo      string
			textcolor string
			fillcolor string
			bgcolor   string
		)
		err = rows.Scan(&id, &name, &logo, &textcolor, &fillcolor, &bgcolor)
		if err != nil {
			log.Fatal(err)
		}

		skills = append(skills, &Skill{id, name, logo, textcolor, fillcolor, bgcolor})
	}

	return skills
}

func (user User) FindProjects() []*Project {
	var projects []*Project

	rows, err := Db.Query(`SELECT id, title, description, url, image, classes FROM projects WHERE author = ?`, user.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id          int
			title       string
			description string
			url         string
			image       string
			classes     string
		)
		err = rows.Scan(&id, &title, &description, &url, &image, &classes)
		if err != nil {
			log.Fatal(err)
		}

		projects = append(projects, &Project{id, title, description, url, image, classes})
	}

	return projects
}

func Skills_Test() []*Skill {
	var skills []*Skill
	skill0 := &Skill{
		Name: "Go",
		Logo: `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 640 450">
		<path d="M400.1 194.8C389.2 197.6 380.2 199.1 371 202.4C363.7 204.3 356.3 206.3 347.8 208.5L347.2 208.6C343 209.8 342.6 209.9 338.7 205.4C334 200.1 330.6 196.7 324.1 193.5C304.4 183.9 285.4 186.7 267.7 198.2C246.5 211.9 235.6 232.2 235.9 257.4C236.2 282.4 253.3 302.9 277.1 306.3C299.1 309.1 316.9 301.7 330.9 285.8C333 283.2 334.9 280.5 337 277.5V277.5L337 277.5C337.8 276.5 338.5 275.4 339.3 274.2H279.2C272.7 274.2 271.1 270.2 273.3 264.9C277.3 255.2 284.8 239 289.2 230.9C290.1 229.1 292.3 225.1 296.1 225.1H397.2C401.7 211.7 409 198.2 418.8 185.4C441.5 155.5 468.1 139.9 506 133.4C537.8 127.8 567.7 130.9 594.9 149.3C619.5 166.1 634.7 188.9 638.8 218.8C644.1 260.9 631.9 295.1 602.1 324.4C582.4 345.3 557.2 358.4 528.2 364.3C522.6 365.3 517.1 365.8 511.7 366.3C508.8 366.5 506 366.8 503.2 367.1C474.9 366.5 449 358.4 427.2 339.7C411.9 326.4 401.3 310.1 396.1 291.2C392.4 298.5 388.1 305.6 382.1 312.3C360.5 341.9 331.2 360.3 294.2 365.2C263.6 369.3 235.3 363.4 210.3 344.7C187.3 327.2 174.2 304.2 170.8 275.5C166.7 241.5 176.7 210.1 197.2 184.2C219.4 155.2 248.7 136.8 284.5 130.3C313.8 124.1 341.8 128.4 367.1 145.6C383.6 156.5 395.4 171.4 403.2 189.5C405.1 192.3 403.8 193.9 400.1 194.8zM48.3 200.4C47.05 200.4 46.74 199.8 47.36 198.8L53.91 190.4C54.53 189.5 56.09 188.9 57.34 188.9H168.6C169.8 188.9 170.1 189.8 169.5 190.7L164.2 198.8C163.6 199.8 162 200.7 161.1 200.7L48.3 200.4zM1.246 229.1C0 229.1-.3116 228.4 .3116 227.5L6.855 219.1C7.479 218.2 9.037 217.5 10.28 217.5H152.4C153.6 217.5 154.2 218.5 153.9 219.4L151.4 226.9C151.1 228.1 149.9 228.8 148.6 228.8L1.246 229.1zM75.72 255.9C75.1 256.8 75.41 257.7 76.65 257.7L144.6 258C145.5 258 146.8 257.1 146.8 255.9L147.4 248.4C147.4 247.1 146.8 246.2 145.5 246.2H83.2C81.95 246.2 80.71 247.1 80.08 248.1L75.72 255.9zM577.2 237.9C577 235.3 576.9 233.1 576.5 230.9C570.9 200.1 542.5 182.6 512.9 189.5C483.9 196 465.2 214.4 458.4 243.7C452.8 268 464.6 292.6 487 302.6C504.2 310.1 521.3 309.2 537.8 300.7C562.4 287.1 575.8 268 577.4 241.2C577.3 240 577.3 238.9 577.2 237.9z" />
	</svg>`,
		TextColor: "text-cyan-700",
		FillColor: "fill-cyan-700",
		BgColor:   "bg-cyan-200",
	}
	skill1 := &Skill{
		Name: "Py",
		Logo: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512" width="16" height="16" >
		<path  d="M439.8 200.5c-7.7-30.9-22.3-54.2-53.4-54.2h-40.1v47.4c0 36.8-31.2 67.8-66.8 67.8H172.7c-29.2 0-53.4 25-53.4 54.3v101.8c0 29 25.2 46 53.4 54.3 33.8 9.9 66.3 11.7 106.8 0 26.9-7.8 53.4-23.5 53.4-54.3v-40.7H226.2v-13.6h160.2c31.1 0 42.6-21.7 53.4-54.2 11.2-33.5 10.7-65.7 0-108.6zM286.2 404c11.1 0 20.1 9.1 20.1 20.3 0 11.3-9 20.4-20.1 20.4-11 0-20.1-9.2-20.1-20.4.1-11.3 9.1-20.3 20.1-20.3zM167.8 248.1h106.8c29.7 0 53.4-24.5 53.4-54.3V91.9c0-29-24.4-50.7-53.4-55.6-35.8-5.9-74.7-5.6-106.8.1-45.2 8-53.4 24.7-53.4 55.6v40.7h106.9v13.6h-147c-31.1 0-58.3 18.7-66.8 54.2-9.8 40.7-10.2 66.1 0 108.6 7.6 31.6 25.7 54.2 56.8 54.2H101v-48.8c0-35.3 30.5-66.4 66.8-66.4zm-6.7-142.6c-11.1 0-20.1-9.1-20.1-20.3.1-11.3 9-20.4 20.1-20.4 11 0 20.1 9.2 20.1 20.4s-9 20.3-20.1 20.3z"/>
	</svg>`,
		TextColor: "text-sky-700",
		FillColor: "fill-sky-700",
		BgColor:   "bg-sky-200",
	}
	skill2 := &Skill{
		Name: "JS",
		Logo: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512" width="16" height="16" >
		<path  d="M400 32H48C21.5 32 0 53.5 0 80v352c0 26.5 21.5 48 48 48h352c26.5 0 48-21.5 48-48V80c0-26.5-21.5-48-48-48zM243.8 381.4c0 43.6-25.6 63.5-62.9 63.5-33.7 0-53.2-17.4-63.2-38.5l34.3-20.7c6.6 11.7 12.6 21.6 27.1 21.6 13.8 0 22.6-5.4 22.6-26.5V237.7h42.1v143.7zm99.6 63.5c-39.1 0-64.4-18.6-76.7-43l34.3-19.8c9 14.7 20.8 25.6 41.5 25.6 17.4 0 28.6-8.7 28.6-20.8 0-14.4-11.4-19.5-30.7-28l-10.5-4.5c-30.4-12.9-50.5-29.2-50.5-63.5 0-31.6 24.1-55.6 61.6-55.6 26.8 0 46 9.3 59.8 33.7L368 290c-7.2-12.9-15-18-27.1-18-12.3 0-20.1 7.8-20.1 18 0 12.6 7.8 17.7 25.9 25.6l10.5 4.5c35.8 15.3 55.9 31 55.9 66.2 0 37.8-29.8 58.6-69.7 58.6z"/>
</svg>`,
		TextColor: "text-amber-700",
		FillColor: "fill-amber-700",
		BgColor:   "bg-amber-200",
	}
	skill3 := &Skill{
		Name: "SQL",
		Logo: `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 448 512" >
		<path 
		d="M448 80v48c0 44.2-100.3 80-224 80S0 172.2 0 128V80C0 35.8 100.3 0 224 0S448 35.8 448 80zM393.2 214.7c20.8-7.4 39.9-16.9 54.8-28.6V288c0 44.2-100.3 80-224 80S0 332.2 0 288V186.1c14.9 11.8 34 21.2 54.8 28.6C99.7 230.7 159.5 240 224 240s124.3-9.3 169.2-25.3zM0 346.1c14.9 11.8 34 21.2 54.8 28.6C99.7 390.7 159.5 400 224 400s124.3-9.3 169.2-25.3c20.8-7.4 39.9-16.9 54.8-28.6V432c0 44.2-100.3 80-224 80S0 476.2 0 432V346.1z"/></svg>`,
		TextColor: "text-emerald-700",
		FillColor: "fill-emerald-700",
		BgColor:   "bg-emerald-200",
	}
	skill4 := &Skill{
		Name: "NoSQL",
		Logo: `<svg xmlns="http://www.w3.org/2000/svg"  width="16" height="16" viewBox="0 0 576 512" >
		<path  d="M264.5 5.2c14.9-6.9 32.1-6.9 47 0l218.6 101c8.5 3.9 13.9 12.4 13.9 21.8s-5.4 17.9-13.9 21.8l-218.6 101c-14.9 6.9-32.1 6.9-47 0L45.9 149.8C37.4 145.8 32 137.3 32 128s5.4-17.9 13.9-21.8L264.5 5.2zM476.9 209.6l53.2 24.6c8.5 3.9 13.9 12.4 13.9 21.8s-5.4 17.9-13.9 21.8l-218.6 101c-14.9 6.9-32.1 6.9-47 0L45.9 277.8C37.4 273.8 32 265.3 32 256s5.4-17.9 13.9-21.8l53.2-24.6 152 70.2c23.4 10.8 50.4 10.8 73.8 0l152-70.2zm-152 198.2l152-70.2 53.2 24.6c8.5 3.9 13.9 12.4 13.9 21.8s-5.4 17.9-13.9 21.8l-218.6 101c-14.9 6.9-32.1 6.9-47 0L45.9 405.8C37.4 401.8 32 393.3 32 384s5.4-17.9 13.9-21.8l53.2-24.6 152 70.2c23.4 10.8 50.4 10.8 73.8 0z"/></svg>`,
		TextColor: "text-red-700",
		FillColor: "fill-red-700",
		BgColor:   "bg-red-200",
	}
	skill5 := &Skill{
		Name: "UI/UX",
		Logo: `<svg width="16" height="16" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" xml:space="preserve"
        stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"><path fill-rule="evenodd" clip-rule="evenodd" d="M12 6.036c-2.667 0-4.333 1.325-5 3.976 1-1.325 2.167-1.822 3.5-1.491.761.189 1.305.738 1.906 1.345C13.387 10.855 14.522 12 17 12c2.667 0 4.333-1.325 5-3.976-1 1.325-2.166 1.822-3.5 1.491-.761-.189-1.305-.738-1.907-1.345-.98-.99-2.114-2.134-4.593-2.134zM7 12c-2.667 0-4.333 1.325-5 3.976 1-1.326 2.167-1.822 3.5-1.491.761.189 1.305.738 1.907 1.345.98.989 2.115 2.134 4.594 2.134 2.667 0 4.333-1.325 5-3.976-1 1.325-2.167 1.822-3.5 1.491-.761-.189-1.305-.738-1.906-1.345C10.613 13.145 9.478 12 7 12z"></path></g></svg>
        </svg>`,
		TextColor: "text-fuchsia-700",
		FillColor: "fill-fuchsia-700",
		BgColor:   "bg-fuchsia-200",
	}

	skills = append(skills, skill0, skill1, skill2, skill3, skill4, skill5)
	return skills
}

// Projects_Test returns test data for projects.
func Projects_Test() []*Project {
	var projects []*Project

	// Create dummy Project objects
	project0 := &Project{
		Title:       "Interior Designer AI",
		Description: "Designer renders with Stable Diffusion XL. Python, Postgress, S3, Nginx, Gunicorn, and SSL backend. React, NextJS, Typescript, frontend",
		Url:         "https://interiordesigner-ai.com/",
		Classes:     "col-span-2",
	}

	project1 := &Project{
		Title:       "Blog",
		Description: "Minimalist Go blog with mysql and htmx frontend",
		Url:         "/post/minimalist-blog-with-go-mysql-htmx-and-tailwind",
	}

	project2 := &Project{
		Title:       "CoffeeGPT",
		Description: "Use OpenAI to dial in your morning specialty coffee.",
		Url:         "/projects/coffeeapp",
	}

	project4 := &Project{
		Title:       "Client Side Moderation",
		Description: "Demo of TensorflowJS toxicity AI model for social media.",
		Url:         "/projects/moderatorjs",
	}

	project3 := &Project{
		Title:       "Document Viewer",
		Description: "Pure JS, drag, zoom, and resize for iframe content.",
		Url:         "/post/document-viewer-for-embedded-html-pages",
	}

	// Append the dummy projects to the list
	projects = append(projects, project0, project1, project2, project3, project4)

	return projects
}
