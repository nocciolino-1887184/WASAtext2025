package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// user endpoint

	rt.router.GET("/all-users", rt.getAllUsers)
	rt.router.GET("/users", rt.searchUser)
	rt.router.PUT("/users/:id/nickname", rt.changeNick)
	rt.router.PUT("/users/:id/propic", rt.changePhoto)
	rt.router.GET("/users/:id/conversations", rt.showMyConvos)
	rt.router.GET("/users/:id", rt.showUserInfo)

	// messages endpoint
	rt.router.DELETE("/conversations/:convoId/messages/:messageId", rt.delMess)
	rt.router.POST("/conversations/:convoId/forwarded", rt.forwardMess)
	rt.router.POST("/conversations/:convoId/messages", rt.sendMessage)
	rt.router.PUT("/conversations/:convoId/letti", rt.putRead)
	rt.router.GET("/conversations/:convoId", rt.showAConv) // chat o gruppo
	rt.router.POST("/conversations", rt.addConvo)          // oneOf

	// reaction endpoint
	rt.router.POST("/messages/:messageId/comments", rt.addReaction)
	rt.router.GET("/messages/:messageId/comments", rt.getReactions)
	rt.router.DELETE("/messages/:messageId/comments/:reactId", rt.delReaction)

	// groups endpoint
	// rt.router.GET("/users/:id/groups", rt.showGroups) NO, creo direttamente le conversazioni
	// rt.router.DELETE("/groups/:groupId", rt.delGroup)
	rt.router.POST("/groups/:groupId/members", rt.addMember)
	rt.router.GET("/groups/:groupId/members", rt.getMembers) // per visualizzare i membri di un gruppo
	rt.router.DELETE("/groups/:groupId/members/:id", rt.delMember)
	rt.router.PUT("/groups/:groupId/name", rt.changeName)
	rt.router.PUT("/groups/:groupId/photo", rt.changeImage)

	// login
	rt.router.POST("/session", rt.login)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
