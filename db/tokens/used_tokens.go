package tokens

import (
	"med/db"
)

func MarkTokenAsUsed(manager *db.Manager, userGuid string, tokenHash string) error {
	query := `
        INSERT INTO used_tokens (user_guid, token_hash)
        VALUES ($1, $2)
    `
	_, err := manager.Conn.Exec(query, userGuid, tokenHash)
	return err
}

func IsTokenUsed(manager *db.Manager, userGuid string, tokenHash string) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM used_tokens 
            WHERE user_guid = $1 AND token_hash = $2
        )
    `
	var exists bool
	err := manager.Conn.QueryRow(query, userGuid, tokenHash).Scan(&exists)
	return exists, err
}

func CleanupOldTokens(manager *db.Manager) error {
	query := `
        DELETE FROM used_tokens
        WHERE used_at < NOW() - INTERVAL '30 days'
    `
	_, err := manager.Conn.Exec(query)
	return err
}
