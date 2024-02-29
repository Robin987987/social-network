package handler

import (
	"backend/pkg/model"
	"backend/pkg/repository"
	"backend/util"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CommentHandler struct {
	commentRepo *repository.CommentRepository
	sessionRepo *repository.SessionRepository
	notificationHandler *NotificationHandler
	postRepo *repository.PostRepository
	userRepo *repository.UserRepository
	VoteHandler *VoteHandler
}

func NewCommentHandler(commentRepo *repository.CommentRepository, sessionRepo *repository.SessionRepository, notificationHandler *NotificationHandler, postRepo *repository.PostRepository, userRepo *repository.UserRepository, voteHandler *VoteHandler) *CommentHandler {
	return &CommentHandler{commentRepo: commentRepo, sessionRepo: sessionRepo, notificationHandler: notificationHandler, postRepo: postRepo, userRepo: userRepo, VoteHandler: voteHandler}
}

func (h *CommentHandler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: id may not come from request and will cause error - added omitempty
	var newComment model.Comment
	err := json.NewDecoder(r.Body).Decode(&newComment)
	if err != nil {
		http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := h.sessionRepo.GetUserIDFromSessionToken(util.GetSessionToken(r))
	if err != nil {
		http.Error(w, "User not authenticated: "+err.Error(), http.StatusUnauthorized)
		return
	}
	newComment.UserID = userID

	// Insert the comment into the database
	// TODO: should it return the new comment?
	createdCommentId, err := h.commentRepo.CreateComment(newComment)
	if err != nil {
		http.Error(w, "Failed to create comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// notify post creator of new comment
	postOwnerId, err := h.postRepo.GetPostOwnerIDByPostID(newComment.PostID)
	if err != nil {
		http.Error(w, "Failed to get post owner id: "+err.Error(), http.StatusInternalServerError)
		return
	}
	post, err := h.postRepo.GetPostByID(newComment.PostID)
	if err != nil {
		http.Error(w, "Failed to get post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// authenticated user username
	username, err := h.userRepo.GetUsernameByID(userID)
	if err != nil {
		http.Error(w, "Failed to get username: "+err.Error(), http.StatusInternalServerError)
		return
	}

	message := username + "commented on your post: " + post.Title
	err = h.notificationHandler.CreateNotification(int(postOwnerId), userID, "post", message)
	if err != nil {
		http.Error(w, "Failed to create notification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Successful response
	response := map[string]interface{}{
		"message": "Comment created successfully",
		"data":    createdCommentId,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *CommentHandler) GetCommentsByUserIDorPostID(w http.ResponseWriter, r *http.Request) {
	var id string

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&id)
	if err != nil {
		http.Error(w, "Error decoding id for comment request: "+err.Error(), http.StatusBadRequest)
		return
	}
	intid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Error converting id to int: "+err.Error(), http.StatusInternalServerError)
		return
	}
	comments, err := h.commentRepo.GetCommentsByID(intid)
	if err != nil {
		http.Error(w, "Error retrieving comments: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// get votes for each comment and append to CommentResponse
	commentsWithVotes, err := h.VoteHandler.AppendVotesToCommentsResponse(comments)
	if err != nil {
		http.Error(w, "Error appending votes to comments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commentsWithVotes)
}

func (h *CommentHandler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the post ID from the URL
	vars := mux.Vars(r)
	commentID, ok := vars["id"]
	intcommentID, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "Failed to parse comment ID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Comment ID is missing in parameters", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Error authenticating user: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Confirm user auth and get userid
	userID, err := h.sessionRepo.GetUserIDFromSessionToken(cookie.Value)
	if err != nil {
		http.Error(w, "Error confirming user authentication: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Delete the comment from the database
	err = h.commentRepo.DeleteComment(intcommentID, userID)
	if err != nil {
		http.Error(w, "Failed to delete the comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Successful response
	response := map[string]string{
		"message": "Comment deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *CommentHandler) EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the post ID from the URL
	vars := mux.Vars(r)
	commentID, ok := vars["id"]
	intcommentID, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "Failed to parse comment ID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Comment ID is missing in parameters", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Error authenticating user: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Confirm user auth and get userid
	userID, err := h.sessionRepo.GetUserIDFromSessionToken(cookie.Value)
	if err != nil {
		http.Error(w, "Error confirming user authentication: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Parse the comment data from the request body
	var commentData model.UpdateCommentRequest
	err = json.NewDecoder(r.Body).Decode(&commentData)
	if err != nil {
		http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Update the comment in the database
	err = h.commentRepo.UpdateComment(intcommentID, userID, commentData)
	if err != nil {
		http.Error(w, "Failed to update the comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Successful response
	response := map[string]string{
		"message": "Comment updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
