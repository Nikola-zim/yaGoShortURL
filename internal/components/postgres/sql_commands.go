package postgres

const (
	getAllURLs = `
		SELECT full_url, user_id 
		FROM urls;`

	writeURL = `
		INSERT INTO urls (shorted_url, full_url, user_id )
		VALUES ($1, $2, $3);
		`
)
