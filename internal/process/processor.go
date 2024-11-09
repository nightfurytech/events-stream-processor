package process

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nightfurytech/events-stream-processor/internal/models"
	"log"
)

type Processor struct {
	db        *sql.DB
	messageCh chan []byte
}

func NewProcessor(db *sql.DB, ch chan []byte) *Processor {
	return &Processor{
		db:        db,
		messageCh: ch,
	}
}

func (p *Processor) EventProcessor() {
	for {
		select {
		case msg := <-p.messageCh:
			eventMsg := &models.Event{}
			if err := json.Unmarshal(msg, eventMsg); err != nil {
				fmt.Println("Error reading message:", err)
			}
			p.incrementCountInDb(eventMsg)
		}
	}
}

func (p *Processor) incrementCountInDb(eventMsg *models.Event) {
	// Start a transaction
	tx, err := p.db.Begin()
	if err != nil {
		log.Fatal("Failed to begin transaction:", err)
	}

	// SELECT query within the transaction
	var event models.DbEvent
	selectQuery := "SELECT event_type, count FROM events WHERE event_type = $1"
	err = tx.QueryRow(selectQuery, eventMsg.Type).Scan(&event.Type, &event.Count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Failed to execute SELECT query:", err)
		tx.Rollback() // Rollback if select fails
		return
	}

	event.Count += 1
	event.Type = eventMsg.Type
	// UPSERT query within the transaction
	upsertQuery := `
		INSERT INTO events (event_type, count)
		VALUES ($1, $2)
		ON CONFLICT (event_type) 
		DO UPDATE SET count = EXCLUDED.count
		RETURNING event_type
	`
	err = tx.QueryRow(upsertQuery, string(event.Type), event.Count).Scan(&event.Type)
	if err != nil {
		log.Println("Failed to execute UPSERT query:", err)
		tx.Rollback() // Rollback if upsert fails
		return
	}

	fmt.Printf("Upserted event with type: %s\n", event.Type)

	// Commit the transaction if all operations succeed
	if err := tx.Commit(); err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}

	fmt.Println("Transaction committed successfully")
}
