package repository

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"paws/internal/types"
	"time"
)

type ConversationRepository interface {
	Create(c *types.Conversation) error
	GetOrCreate(identifier uuid.UUID, secondaryParticipantID string) (*types.Conversation, error)
	List(participantID string) ([]types.Conversation, error)
	ListHistoricalMessages(conversationID int64, toDate time.Time, lookbackDays int) ([]types.Message, error)
	CreateMessage(m *types.Message) error
	MarkMessageRead(messageId int64, participantID string) error
	RecentMessages(conversationID int64, limit int) ([]types.Message, error)
}

type postgresConversationRepository struct {
	db *sqlx.DB
}

func NewConversationsRepository(db *sqlx.DB) ConversationRepository {
	return &postgresConversationRepository{db: db}
}

func (r *postgresConversationRepository) List(participantID string) ([]types.Conversation, error) {
	stmt := `
		select * 
		from conversations 
		where primary_participant_id = $1 or secondary_participant_id = $1;`

	var cc []types.Conversation
	if err := r.db.Select(&cc, stmt, participantID); err != nil {
		return nil, err
	}
	return cc, nil
}

func (r *postgresConversationRepository) Create(c *types.Conversation) error {
	stmt := `
		insert into conversations (identifier, primary_participant_id, secondary_participant_id)
		values ($1, $2, $3)
		returning id, last_message_at, created_at;`

	if err := r.db.Get(c, stmt, c.Identifier, c.PrimaryParticipantID, c.SecondaryParticipantID); err != nil {
		return err
	}
	return nil
}

// GetOrCreate finds an existing or creates a new conversation.
// The participantID is either the primary or secondary participant for an exising conversation, but can only
// be the secondary participant when creating a conversation as conversations must be initialised by them.
func (r *postgresConversationRepository) GetOrCreate(identifier uuid.UUID, participantID string) (*types.Conversation, error) {
	stmt := `select * from conversations where identifier = $1 and (primary_participant_id = $2 or secondary_participant_id = $2);`
	var conversation types.Conversation
	if err := r.db.Get(&conversation, stmt, identifier, participantID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if conversation.ID != 0 {
		return &conversation, nil
	}

	stmt = `select user_id from pets where id = $1;`
	var primaryParticipantID string
	if err := r.db.Get(&primaryParticipantID, stmt, identifier); err != nil {
		return nil, err
	}

	conversation = types.Conversation{
		Identifier:             identifier,
		PrimaryParticipantID:   primaryParticipantID,
		SecondaryParticipantID: participantID,
	}

	if err := r.Create(&conversation); err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *postgresConversationRepository) CreateMessage(m *types.Message) error {
	stmt := `
		insert into messages (conversation_id, sender_id, text)
		values ($1, $2, $3)
		returning id, created_at, read_at;`

	if err := r.db.Get(m, stmt, m.ConversationID, m.SenderID, m.Text); err != nil {
		return err
	}
	return nil
}

func (r *postgresConversationRepository) ListHistoricalMessages(conversationID int64, toDate time.Time, lookbackDays int) ([]types.Message, error) {
	q := `
		select *
		from messages
		where conversation_id = $1
		  and created_at between $2 and $3
		order by created_at;`

	var mm []types.Message
	fromDate := toDate.AddDate(0, 0, -lookbackDays)
	if err := r.db.Select(&mm, q, conversationID, fromDate, toDate); err != nil {
		return nil, err
	}
	return mm, nil
}

func (r *postgresConversationRepository) RecentMessages(conversationID int64, limit int) ([]types.Message, error) {
	q := `
		with recent_messages as (
			select *
			from messages
			where conversation_id = $1
			order by created_at desc
			limit $2
		)
		select *
		from recent_messages
		order by created_at;`

	var messages []types.Message
	if err := r.db.Select(&messages, q, conversationID, limit); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *postgresConversationRepository) MarkMessageRead(messageID int64, participantID string) error {
	authorizationStmt := `
		select primary_participant_id, secondary_participant_id
		from conversations c
		join messages m on c.id = m.conversation_id
		where m.id = $1;`

	var p1, p2 string
	if err := r.db.QueryRow(authorizationStmt, messageID).Scan(&p1, &p2); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if p1 != participantID && p2 != participantID {
		return ErrNotAuthorized
	}

	stmt := `
		with target_message as (
			select conversation_id, sender_id, created_at
			from messages
			where id = $1
		)
		update messages
		set read_at = now()
		where conversation_id = (select conversation_id from target_message)
		  and sender_id != (select target_message.sender_id from target_message)
		  and created_at <= (select created_at from target_message)
		  and read_at is null;`

	if _, err := r.db.Exec(stmt, messageID); err != nil {
		return err
	}
	return nil
}
