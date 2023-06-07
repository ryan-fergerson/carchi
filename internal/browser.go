package internal

import (
	"database/sql"
)

type ConversationBrowser struct {
	DB *sql.DB
}

type RecentConversation struct {
	Id              string `json:"conversation_id"`
	Title           string `json:"conversation_title"`
	CreateTimestamp string `json:"conversation_create_timestamp"`
	ImportTimestamp string `json:"conversation_import_timestamp"`
}

type ConversationPart struct {
	Id           string         `json:"conversation_id"`
	CurrentNode  string         `json:"conversation_current_node"`
	Title        string         `json:"conversation_title"`
	MessageParts []MessageParts `json:"conversation_message_parts"`
}

type MessageParts struct {
	Id              string `json:"message_id"`
	CreateTimestamp string `json:"message_create_timestamp"`
	ImportTimestamp string `json:"message_import_timestamp"`
	AuthorRole      string `json:"message_author_role"`
	Parts           string `json:"message_parts"`
}

type SearchResult struct {
	ConversationId string
	Title          string
	Parts          string
	Rank           float64
	RowNumber      int64
}

func NewConversationBrowserService() (*ConversationBrowser, error) {
	db, err := NewDatabaseService().GetDatabaseHandle()

	if err != nil {
		return nil, err
	}

	return &ConversationBrowser{
		DB: db,
	}, nil
}

func (cb *ConversationBrowser) SearchConversations(query string) ([]SearchResult, error) {
	var searchResults []SearchResult

	rows, err := cb.DB.Query(`
		SELECT
			conversation_id,
			headline_title,
			headline_parts,
			rank,
			ROW_NUMBER() OVER (ORDER BY rank DESC)
		FROM search_conversations($1);
	`, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var sr SearchResult

		err := rows.Scan(
			&sr.ConversationId,
			&sr.Title,
			&sr.Parts,
			&sr.Rank,
			&sr.RowNumber,
		)

		if err != nil {
			return nil, err
		}

		searchResults = append(searchResults, sr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return searchResults, nil
}

func (cb *ConversationBrowser) GetRecentConversations() ([]RecentConversation, error) {
	var recentConversations []RecentConversation

	rows, err := cb.DB.Query(`
		SELECT
			id,
			title,
			to_timestamp(create_time),
			to_timestamp(import_time)
		FROM conversations
		ORDER BY create_time DESC;
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var rc RecentConversation

		err := rows.Scan(
			&rc.Id,
			&rc.Title,
			&rc.CreateTimestamp,
			&rc.ImportTimestamp,
		)

		if err != nil {
			return nil, err
		}

		recentConversations = append(recentConversations, rc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return recentConversations, nil
}

func (cb *ConversationBrowser) GetMessagePartsByConversationId(id string) (*ConversationPart, error) {
	var conversationPart ConversationPart
	var messageParts []MessageParts

	rows, err := cb.DB.Query(`
		SELECT
			c.id conversation_id,
			c.current_node conversation_current_node,
			c.title conversation_title
		FROM conversations c
		WHERE c.id = $1;
	`, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var cp ConversationPart

		err := rows.Scan(
			&cp.Id,
			&cp.CurrentNode,
			&cp.Title,
		)

		if err != nil {
			return nil, err
		}

		conversationPart = cp
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	rows, err = cb.DB.Query(`
		SELECT
			m.id message_id,
			to_timestamp(m.create_time) message_create_timestamp,
			to_timestamp(m.import_time) message_import_timestamp,
			m.author_role message_author_role,
			m.parts message_parts
		FROM messages m
		JOIN nodes n ON m.node_id = n.id
		JOIN conversations c ON n.conversation_id = c.id
		WHERE c.id = $1
		AND content_type = 'text'
		AND author_role != 'system'
		ORDER BY m.create_time;
	`, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var mp MessageParts

		err := rows.Scan(
			&mp.Id,
			&mp.CreateTimestamp,
			&mp.ImportTimestamp,
			&mp.AuthorRole,
			&mp.Parts,
		)

		if err != nil {
			return nil, err
		}

		messageParts = append(messageParts, mp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	conversationPart.MessageParts = messageParts

	return &conversationPart, nil
}
