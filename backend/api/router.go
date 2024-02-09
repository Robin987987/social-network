package api

import (
	"backend/pkg/handler"
	"backend/pkg/repository"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

// API layer, handlers, and routing
func Router(mux *mux.Router, db *sql.DB) {
	// User registration requires input in the form like RegistrationData struct at /pkg/model/stucts.go
	sessionRepository := repository.NewSessionRepository(db)
    friendsRepository := repository.NewFriendsRepository(db)

	userHandler := handler.NewUserHandler(repository.NewUserRepository(db), sessionRepository)
	mux.HandleFunc("/api/users/register", userHandler.UserRegisterHandler).Methods("POST")
	// User login and logout
	mux.HandleFunc("/api/users/logout", handler.LogoutHandler).Methods("POST")
	mux.HandleFunc("/api/users/login", userHandler.LoginHandler).Methods("POST")
	mux.HandleFunc("/api/users/check-auth", userHandler.CheckAuth)

	// Posts
	postHandler := handler.NewPostHandler(repository.NewPostRepository(db), sessionRepository, friendsRepository)
	mux.HandleFunc("/post", postHandler.GetAllPostsHandler).Methods("GET") // Main feed, all public posts + user groups posts
	mux.HandleFunc("/post", postHandler.CreatePostHandler).Methods("POST")
	//mux.HandleFunc("/post/{id}", handler.GetPostByIDHandler).Methods("GET")
	mux.HandleFunc("/post/{id}", postHandler.EditPostHandler).Methods("PUT")      // Edit a post
	mux.HandleFunc("/post/{id}", postHandler.DeletePostHandler).Methods("DELETE") // Delete a post
	mux.HandleFunc("/groups/posts/{id}", postHandler.GetPostsByGroupIDHandler).Methods("GET")
	// Profile feed, all posts by user
	mux.HandleFunc("/profile/posts/{id}", postHandler.GetAllUserPostsHandler).Methods("GET")

	// Comments
	commentHandler := handler.NewCommentHandler(repository.NewCommentRepository(db), sessionRepository)
	mux.HandleFunc("/post/{id}/comments", commentHandler.GetCommentsByUserIDorPostID).Methods("GET")
	mux.HandleFunc("/post/comment", commentHandler.CreateCommentHandler).Methods("POST")
	mux.HandleFunc("/post/comment/{id}", commentHandler.DeleteCommentHandler).Methods("DELETE")

	// Groups
	groupHandler := handler.NewGroupHandler(repository.NewGroupRepository(db), sessionRepository)
	mux.HandleFunc("/groups", groupHandler.GetAllGroupsHandler).Methods("GET")
	mux.HandleFunc("/groups", groupHandler.CreateGroupHandler).Methods("POST")
	mux.HandleFunc("/groups/{id}", groupHandler.GetGroupByIDHandler).Methods("GET")
	mux.HandleFunc("/groups/{id}", groupHandler.EditGroupHandler).Methods("PUT")
	mux.HandleFunc("/groups/{id}", groupHandler.DeleteGroupHandler).Methods("DELETE")

	// Group invitations & requests

	groupMemberHandler := handler.NewGroupMemberHandler(repository.NewGroupMemberRepository(db), repository.NewInvitationRepository(db), sessionRepository)
	mux.HandleFunc("/invitations", groupMemberHandler.GetAllGroupInvitationsHandler).Methods("GET")
	mux.HandleFunc("/invitations", groupMemberHandler.InviteGroupMemberHandler).Methods("POST")
	mux.HandleFunc("/invitations/{id}", groupMemberHandler.GetGroupInvitationByIDHandler).Methods("GET")
	mux.HandleFunc("/invitations/{id}", groupMemberHandler.DeclineGroupInvitationHandler).Methods("PUT")
	mux.HandleFunc("/invitations/{id}", groupMemberHandler.AcceptGroupInvitationHandler).Methods("PUT")
	mux.HandleFunc("/invitations/request/{id}", groupMemberHandler.RequestGroupMembershipHandler).Methods("POST")
	mux.HandleFunc("/groups/{groupId}/members/{userId}", groupMemberHandler.RemoveMemberHandler).Methods("DELETE")
	mux.HandleFunc("/invitations/approve/{id}", groupMemberHandler.ApproveGroupMembershipHandler).Methods("PUT")

	// TODO: Group events
	mux.HandleFunc("/events", handler.GetAllEventsHandler).Methods("GET")
	mux.HandleFunc("/events", handler.CreateEventHandler).Methods("POST")
	mux.HandleFunc("/events/{id}", handler.GetEventByIDHandler).Methods("GET")
	mux.HandleFunc("/events/{id}", handler.EditEventHandler).Methods("PUT")
	mux.HandleFunc("/events/{id}", handler.DeleteEventHandler).Methods("DELETE")

	// TODO: Notifications
	mux.HandleFunc("/notifications", handler.GetAllNotificationsHandler).Methods("GET")
	mux.HandleFunc("/notifications", handler.CreateNotificationHandler).Methods("POST")
	mux.HandleFunc("/notifications/{id}", handler.GetNotificationByIDHandler).Methods("GET")
	mux.HandleFunc("/notifications/{id}", handler.MarkNotificationAsReadHandler).Methods("PUT")


	friendHandler := handler.NewFriendHandler(friendsRepository, sessionRepository)
	mux.HandleFunc("/friends/request", friendHandler.SendFriendRequestHandler).Methods("POST")
	mux.HandleFunc("/friends/accept", friendHandler.AcceptFriendRequestHandler).Methods("POST")
	mux.HandleFunc("/friends/decline", friendHandler.DeclineFriendRequestHandler).Methods("POST")
	mux.HandleFunc("/friends/block", friendHandler.BlockUserHandler).Methods("POST")
	mux.HandleFunc("/friends/unblock", friendHandler.UnblockUserHandler).Methods("POST")

	mux.HandleFunc("/friends", friendHandler.GetFriendsHandler).Methods("GET")

	// CORS
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},                   // Replace with your frontend's origin
		AllowCredentials: true,                                                // Important for cookies, authorization headers with HTTPS
		AllowedHeaders:   []string{"Authorization", "Content-Type"},           // You can adjust this based on your needs
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Adjust the methods based on your requirements
		// You can include other settings like ExposedHeaders, MaxAge, etc., according to your needs
	})
	mux_cors := corsOptions.Handler(mux)
	http.Handle("/", mux_cors)
}
