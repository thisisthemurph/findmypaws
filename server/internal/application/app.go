package application

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"paws/internal/chat"
	"paws/internal/database/model"
	"paws/internal/repository"
)

type App struct {
	DB           *sqlx.DB
	ChatManager  *chat.Manager
	Repositories *repository.Repositories
	//ServerMux    *http.ServeMux
	Logger *slog.Logger
	Config AppConfig
}

func NewApp() (*App, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	config := NewAppConfig(os.Getenv)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	clerk.SetKey(config.Clerk.Secret)

	return &App{
		Logger: logger,
		Config: config,
	}, nil
}

func (app *App) Build() error {
	app.Logger.Info("building app")

	if err := app.configureDatabase(); err != nil {
		return err
	}
	app.configureRepositories()
	app.configureChatManager()

	//app.ServerMux = routes.BuildRoutesServerMux(app)
	return nil
}

func (app *App) configureDatabase() error {
	app.Logger.Info("configuring stores")

	db, err := sqlx.Open("postgres", app.Config.Database.ConnectionString)
	if err != nil {
		return fmt.Errorf("could not open to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}

	app.DB = db
	return nil
}

func (app *App) configureRepositories() {
	app.Logger.Info("configuring repositories")
	app.Repositories = repository.NewRepositories(app.DB)
}

func (app *App) configureChatManager() {
	app.Logger.Info("configuring chat manager")
	conversation := app.Repositories.ConversationRepository

	app.ChatManager = chat.NewManager(chat.ManagerConfig{
		Logger: app.Logger,
		Callbacks: chat.ManagerCallbacks{
			HandleRoomCreation: func(identifier uuid.UUID, secondaryParticipantID string) (chat.RoomDetail, error) {
				conv, err := conversation.GetOrCreate(identifier, secondaryParticipantID)
				if err != nil {
					return nil, err
				}
				return ConversationWrapper{
					Conversation: conv,
				}, nil
			},
			HandleNewMessage: func(conversationID int64, messageEvent chat.NewMessageEvent) (int64, error) {
				m := &model.Message{
					ConversationID: conversationID,
					SenderID:       messageEvent.SenderID,
					Text:           messageEvent.Text,
				}

				if err := conversation.CreateMessage(m); err != nil {
					return 0, err
				}
				return m.ID, nil
			},
			HandleEmojiUpdate: func(conversationID, messageID int64, emojiKey *string) error {
				message, err := conversation.GetMessage(conversationID, messageID)
				if err != nil {
					return fmt.Errorf("could not get conversation message: %w", err)
				}
				message.EmojiReaction = emojiKey
				if err := conversation.UpdateMessage(message); err != nil {
					return fmt.Errorf("could not update message: %w", err)
				}
				return nil
			},
			FetchHistoricalMessages: func(conversationID int64) ([]chat.MessageDetail, error) {
				mm, err := conversation.ListHistoricalMessages(conversationID, time.Now(), 10)
				if err != nil {
					return nil, err
				}
				results := make([]chat.MessageDetail, len(mm))
				for i, m := range mm {
					results[i] = MessageWrapper{
						Message: &m,
					}
				}
				return results, nil
			},
		},
	})
}
