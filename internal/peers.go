package internal

import (
	"database/sql"
	"fmt"
	"log"
)

type PeerState struct {
	DB *sql.DB
}

// This is equivalent to the rows in the CronPeers table
type Peer struct {
	ID       int
	Endpoint string
	Secret   string
}

// Adds a job
func (p PeerState) Add(endpoint, secret string) error {
	_, err := p.DB.Exec("INSERT INTO Peers (endpoint, secret) VALUES(?, ?);", endpoint, secret)
	if err != nil {
		return fmt.Errorf("Could not insert: %s", err)
	}

	return nil
}

// Returns the peers from the database
func (p PeerState) Get() []Peer {
	rows, err := p.DB.Query("SELECT id, endpoint, secret FROM Peers ORDER BY id;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var peer_buf []Peer
	for rows.Next() {
		var peer Peer
		if err := rows.Scan(&peer.ID, &peer.Endpoint, &peer.Secret); err != nil {
			log.Fatal(err)
		}
		if peer.Endpoint != "THIS" {
			peer_buf = append(peer_buf, peer)
		}
	}

	return peer_buf
}

// Removes a job
func (p PeerState) Remove(id int) error {
	_, err := p.DB.Exec("DELETE FROM Peers WHERE id = ?;", id)
	return err
}

func (p PeerState) Secret() (string, error) {
	rows, err := p.DB.Query("SELECT secret FROM Peers WHERE endpoint='THIS';")
	if err != nil {
		return "", err
	}

	var peer Peer
	for rows.Next() {
		if err := rows.Scan(&peer.Secret); err != nil {
			return "", err
		}
		break
	}
	rows.Close()

	return peer.Secret, nil
}

func (p PeerState) Check(secret string) (bool, error) {
	instanceSecret, err := p.Secret()
	return instanceSecret == secret, err
}
