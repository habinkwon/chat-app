package graph

import "github.com/habinkwon/chat-app/pkg/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserSvc *service.User
	ChatSvc *service.Chat
}
