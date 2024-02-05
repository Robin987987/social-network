package handler

import (
	"backend/pkg/db/sqlite"
	"backend/pkg/model"
	"backend/pkg/repository"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type GroupHandler struct {
    repo *repository.GroupRepository
}

type InvitationHandler struct {
    repo *repository.InvitationRepository
}

func NewGroupHandler(repo *repository.GroupRepository) *GroupHandler {
    return &GroupHandler{repo: repo}
}

func NewInvitationHandler(repo *repository.InvitationRepository) *InvitationHandler {
    return &InvitationHandler{repo: repo}
}

// Group Handlers
func (h *GroupHandler) GetAllGroupsHandler(w http.ResponseWriter, r *http.Request) {

    // DEPRECATED: auth checked as middleware
    // check auth and get userid from cookie
    // cookie, err := r.Cookie("session_token")
    // if err != nil {
    //     if err == http.ErrNoCookie {
    //         // If the session cookie doesn't exist, set isAuthenticated to false
    //         http.Error(w, "User not authenticated", http.StatusUnauthorized)
    //         return
    //     } else {
    //         http.Error(w, "Error checking session token: " + err.Error(), http.StatusInternalServerError)
    //         return
    //     }
    // }
    // _, err = middleware.ConfirmAuthentication(cookie) //(userid, err :=) when later would highlight groups user is a part of
    // if err != nil {
    //     http.Error(w, "Error confirming authentication: " + err.Error(), http.StatusUnauthorized)
    //     return
    // }


    // TODO: Implement logic for getting all groups
    groups, err := h.repo.GetAllGroups()
    if err != nil {
        http.Error(w, "Failed to get groups: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(groups)
}

func (h *GroupHandler) CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: Implement logic for creating a group
    var newGroup model.Group
    err := json.NewDecoder(r.Body).Decode(&newGroup)
    if err != nil {
        http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
        return
    }

    // DEPRECATED: auth checked as middleware
    // // check auth and get userid from cookie
    // cookie, err := r.Cookie("session_token")
    // if err != nil {
    //     if err == http.ErrNoCookie {
    //         // If the session cookie doesn't exist, set isAuthenticated to false
    //         http.Error(w, "User not authenticated", http.StatusUnauthorized)
    //         return
    //     } else {
    //         http.Error(w, "Error checking session token: " + err.Error(), http.StatusInternalServerError)
    //         return
    //     }
    // }
    // userid, err := middleware.ConfirmAuthentication(cookie)
    // if err != nil {
    //     http.Error(w, "Error confirming authentication: " + err.Error(), http.StatusUnauthorized)
    //     return
    // }

    // TODO: check if group with title already exists
    userID := r.Context().Value("AuthUserID").(int)
    newGroup.CreatorId = userID
    groupRepo := repository.NewGroupRepository(sqlite.Dbase)
    createdGroup, err := groupRepo.CreateGroup(newGroup)
    if err != nil {
        http.Error(w, "Failed to create group: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(createdGroup)

}

func (h *GroupHandler) GetGroupByIDHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: Implement logic for getting a group by ID
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }
    
    group, err := h.repo.GetGroupByID(id)
    if err != nil {
        http.Error(w, "Failed to get group: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(group)
}

func (h *GroupHandler) EditGroupHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: Implement logic for editing a group
    var updatedGroup model.Group
    err := json.NewDecoder(r.Body).Decode(&updatedGroup)
    if err != nil {
        http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
        return
    }
    
    err = h.repo.UpdateGroup(updatedGroup)
    if err != nil {
        http.Error(w, "Failed to update group: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedGroup)
}

// DeleteGroupHandler handles the HTTP request for deleting a group.
// It checks the user's authentication, verifies their authorization to delete the group,
// and deletes the group from the repository if all conditions are met.
// If any errors occur during the process, appropriate HTTP error responses are returned.
// TODO: implement notification to all group members that the group has been deleted, and remove all group members;
// implement logging of the deletion or add bool field "deleted"
func (h *GroupHandler) DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {

    // // DEPRECATED: auth checked as middleware
    // // check auth and get userid from cookie
    // cookie, err := r.Cookie("session_token")
    // if err != nil {
    //     if err == http.ErrNoCookie {
    //         // If the session cookie doesn't exist, set isAuthenticated to false
    //         http.Error(w, "User not authenticated", http.StatusUnauthorized)
    //         return
    //     } else {
    //         http.Error(w, "Error checking session token: " + err.Error(), http.StatusInternalServerError)
    //         return
    //     }
    // }
    // userId, err := middleware.ConfirmAuthentication(cookie)
    // if err != nil {
    //     http.Error(w, "Error confirming authentication: " + err.Error(), http.StatusUnauthorized)
    //     return
    // }

    userID := r.Context().Value("AuthUserID").(int)

    // TODO: Implement logic for deleting a group
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }
    
    group, err := h.repo.GetGroupByID(id)
    if err != nil {
        http.Error(w, "Failed to get group: "+err.Error(), http.StatusInternalServerError)
        return
    }
    if group.CreatorId != userID {
        http.Error(w, "User not authorized to delete this group", http.StatusUnauthorized)
        return
    }
    err = h.repo.DeleteGroup(id)
    if err != nil {
        http.Error(w, "Failed to delete group: "+err.Error(), http.StatusInternalServerError)
        return
    }
    // Successful response
	response := map[string]string{
		"message": "Post deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DEPRECATED: moved to groupMemberHandler.go
// // AddMemberToGroup adds a user to a group. It takes two parameters: the ID of the group
// // and the ID of the user. It inserts a new row into the group_members table in the database,
// // which represents the user being a member of the group. If the operation is successful,
// // it returns nil. If there is an error, it returns the error.
// func (h *GroupHandler) AddMemberHandler(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     groupId, ok := vars["groupId"]
//     if !ok {
//         http.Error(w, "Missing group ID", http.StatusBadRequest)
//         return
//     }
//     intGroupId, err := strconv.Atoi(groupId); if err != nil {
//         http.Error(w, "Failed to convert groupid string to int: " + err.Error(), http.StatusBadRequest)
//         return
//     }

//     userId, ok := vars["userId"]
//     if !ok {
//         http.Error(w, "Missing user ID", http.StatusBadRequest)
//         return
//     }
//     intUserId, err := strconv.Atoi(userId); if err != nil {
//         http.Error(w, "Failed to convert userid string to int: " + err.Error(), http.StatusBadRequest)
//         return
//     }

//     err = h.repo.AddMemberToGroup(intGroupId, intUserId); if err != nil {
//         http.Error(w, "Failed to add member to group: "+err.Error(), http.StatusInternalServerError)
//         return
//     }
//     w.WriteHeader(http.StatusOK)
// }

// // RemoveMemberFromGroup removes a user from a group. It takes two parameters: the ID of the group
// // and the ID of the user. It deletes the row from the group_members table in the database that
// // represents the user being a member of the group. If the operation is successful, it returns nil.
// // If there is an error, it returns the error.
// func (h *GroupHandler) RemoveMemberHandler(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     groupId, ok := vars["groupId"]
//     if !ok {
//         http.Error(w, "Missing group ID", http.StatusBadRequest)
//         return
//     }
//     intGroupId, err := strconv.Atoi(groupId); if err != nil {
//         http.Error(w, "Failed to convert groupid string to int: " + err.Error(), http.StatusBadRequest)
//         return
//     }

//     userId, ok := vars["userId"]
//     if !ok {
//         http.Error(w, "Missing user ID", http.StatusBadRequest)
//         return
//     }
//     intUserId, err := strconv.Atoi(userId); if err != nil {
//         http.Error(w, "Failed to convert userid string to int: " + err.Error(), http.StatusBadRequest)
//         return
//     }

//     err = h.repo.RemoveMemberFromGroup(intGroupId, intUserId)
//     if err != nil {
//         http.Error(w, "Failed to remove member from group: "+err.Error(), http.StatusInternalServerError)
//         return
//     }
//     w.WriteHeader(http.StatusOK)
// }

// // ----------------- Group Membership/Invitation/Request Handlers ----------------- 

// // RequestGroupMembershipHandler allows a user to request membership in a group.
// // It creates a membership request in the database that can be approved or denied by the group's admin.
// func (h *InvitationHandler) RequestGroupMembershipHandler(w http.ResponseWriter, r *http.Request) {
//     // Parse the request body to get the group ID and user ID
//     var request model.GroupInvitation
//     err := json.NewDecoder(r.Body).Decode(&request)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     // Create the membership request in the database
//     err = h.repo.CreateGroupInvitation(request)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.WriteHeader(http.StatusCreated)
// }

// // TODO: implement logic for setting the player as group member
// func (h *InvitationHandler) ApproveGroupMembershipHandler(w http.ResponseWriter, r *http.Request) {
//     // Parse the request URL to get the invitation ID
//     vars := mux.Vars(r)
//     id := vars["id"]

//     // Update the status of the membership request to "approved"
//     err := h.repo.AcceptGroupInvitation(id)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.WriteHeader(http.StatusOK)
// }

// func (h *InvitationHandler) DeclineGroupMembershipHandler(w http.ResponseWriter, r *http.Request) {
//     // Parse the request URL to get the invitation ID
//     vars := mux.Vars(r)
//     id := vars["id"]

//     // Update the status of the membership request to "declined"
//     err := h.repo.DeclineGroupInvitation(id)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.WriteHeader(http.StatusOK)
// }

// // InviteMemberHandler sends an invitation to a user to join a group.
// // It creates an invitation in the database that can be accepted or declined by the user.
// func (h *InvitationHandler) InviteGroupMemberHandler(w http.ResponseWriter, r *http.Request) {
//     // Parse the request body to get the group ID and user ID
//     var invitation model.GroupInvitation
//     err := json.NewDecoder(r.Body).Decode(&invitation)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     // Create the invitation in the database
//     err = h.repo.CreateGroupInvitation(invitation)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.WriteHeader(http.StatusCreated)
// }

// // GetAllInvitationsHandler gets all invitations
// func (h *InvitationHandler) GetAllGroupInvitationsHandler(w http.ResponseWriter, r *http.Request) {
//     invitations, err := h.repo.GetAllGroupInvitations()
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     json.NewEncoder(w).Encode(invitations)
// }

// // CreateInvitationHandler creates a new invitation
// // TODO: merge with InviteGroupMemberHandler - same functionality
// func (h *InvitationHandler) CreateGroupInvitationHandler(w http.ResponseWriter, r *http.Request) {
//     var newInvitation model.GroupInvitation
//     err := json.NewDecoder(r.Body).Decode(&newInvitation)
//     if err != nil {
//         http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
//         return
//     }
//     err = h.repo.CreateGroupInvitation(newInvitation)
//     if err != nil {
//         http.Error(w, "Failed to create invitation: "+err.Error(), http.StatusInternalServerError)
//         return
//     }
//     w.WriteHeader(http.StatusCreated)
//     json.NewEncoder(w).Encode(newInvitation)
// }

// // GetInvitationByIDHandler gets an invitation by ID
// // TODO: refactor to use userid from cookie to get all group invitations for the user (pending)
// func (h *InvitationHandler) GetGroupInvitationByIDHandler(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     id := vars["id"]
//     invitation, err := h.repo.GetGroupInvitationByID(id)
//     if err != nil {
//         http.Error(w, "Failed to get invitation: "+err.Error(), http.StatusInternalServerError)
//         return
//     }
//     json.NewEncoder(w).Encode(invitation)
// }

// // AcceptGroupInvitationHandler allows a user to accept an invitation to join a group.
// func (h *InvitationHandler) AcceptGroupInvitationHandler(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     id := vars["id"]
//     err := h.repo.AcceptGroupInvitation(id)
//     if err != nil {
//         http.Error(w, "Failed to accept invitation: "+err.Error(), http.StatusInternalServerError)
//         return
//     }
//     w.WriteHeader(http.StatusOK)
// }

// // DeclineGroupInvitationHandler allows a user to decline an invitation to join a group.
// func (h *InvitationHandler) DeclineGroupInvitationHandler(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     id := vars["id"]
//     err := h.repo.DeclineGroupInvitation(id)
//     if err != nil {
//         http.Error(w, "Failed to decline invitation: "+err.Error(), http.StatusInternalServerError)
//         return
//     }
//     w.WriteHeader(http.StatusOK)
// }