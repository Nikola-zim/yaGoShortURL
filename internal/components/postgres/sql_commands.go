package postgres

const (
	getAllURLs = `
		SELECT full_url, user_id 
		FROM urls;`

	writeURL = `
		INSERT INTO urls (shorted_url, full_url, user_id, deleted )
		VALUES ($1, $2, $3, false);
		`

	deleteURL = `
		UPDATE urls
		SET deleted = true
		WHERE shorted_url = $2 AND user_id = $1;
		`
)
