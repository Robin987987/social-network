package repository

import (
	"backend/pkg/model"
	"database/sql"
	"errors"
)

type FriendsRepository struct {
	db *sql.DB
}

func NewFriendsRepository(db *sql.DB) *FriendsRepository {
	return &FriendsRepository{db: db}
}

func (r *FriendsRepository) AddFriend(userID, friendID int) error {
	exists, err := r.FriendRequestExists(userID, friendID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("a friend request already exists between these users")
	}
	query := `
        INSERT INTO friends (user_id1, user_id2, status, action_user_id)
        VALUES (?, ?, 'pending', ?)
    `
	_, err = r.db.Exec(query, userID, friendID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *FriendsRepository) GetFriends(userID int) ([]model.FriendList, error) {
	query := `
        SELECT users.id, users.first_name, users.last_name, users.avatar_url, users.username
        FROM friends
        JOIN users ON (friends.user_id2 = users.id AND friends.user_id1 = ?) OR (friends.user_id1 = users.id AND friends.user_id2 = ?)
        WHERE friends.status = 'accepted'
    `
	rows, err := r.db.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []model.FriendList
	for rows.Next() {
		var friend model.FriendList
		if err := rows.Scan(&friend.UserID, &friend.FirstName, &friend.LastName, &friend.AvatarURL, &friend.Username); err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return friends, nil
}

func (r *FriendsRepository) UpdateFriendStatus(userID, friendID int, status string) error {
	query := `
        UPDATE friends
        SET status = ?, action_user_id = ?
        WHERE (user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)
    `
	_, err := r.db.Exec(query, status, userID, userID, friendID, friendID, userID)
	return err
}

func (r *FriendsRepository) RemoveFriend(userID, friendID int) error {
	query := `
        DELETE FROM friends
        WHERE (user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)
    `
	_, err := r.db.Exec(query, userID, friendID, friendID, userID)
	return err
}

func (r *FriendsRepository) GetFriendStatus(userID, friendID int) (string, error) {
	query := `
        SELECT status, user_id1 FROM friends
        WHERE (user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)
    `
	var (
		status  string
		user1ID int
	)
	err := r.db.QueryRow(query, userID, friendID, friendID, userID).Scan(&status, &user1ID)
	if err != nil {
		return "", err
	}

	// If the status is 'pending', further clarify based on who initiated the request
	if status == "pending" {
		if user1ID != userID {
			// The current user sent the friend request
			status = "pending_confirmation"
		} else {
			// The current user received the friend request
			status = "pending"
		}
	}

	return status, nil
}

func (r *FriendsRepository) FriendRequestExists(userID, friendID int) (bool, error) {
	status, err := r.GetFriendStatus(userID, friendID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No friend request exists
			return false, nil
		}
		// An error occurred
		return false, err
	}
	// A friend request exists if the status is not an empty string
	return status != "", nil
}
