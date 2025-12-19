package chat

import (
	"net/http"
	"strconv"

	"gin-real-time-talk/internal/entity/interfaces"
	"gin-real-time-talk/pkg/pagination"
	"gin-real-time-talk/pkg/websocket"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

type ChatController struct {
	chatUsecase interfaces.ChatUsecase
	hub         *websocket.Hub
}

func NewChatController(chatUsecase interfaces.ChatUsecase, hub *websocket.Hub) *ChatController {
	return &ChatController{
		chatUsecase: chatUsecase,
		hub:         hub,
	}
}

// GetUserChats godoc
// @Summary Get user chats
// @Description Returns paginated list of chats for the authenticated user
// @Tags chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of items per page (default: 20, max: 100)"
// @Param currentPage query int false "Current page number (default: 1)"
// @Param search query string false "Search query for user name or last message text"
// @Success 200 {object} map[string]interface{} "List of chats with pagination info"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /chats [get]
func (cc *ChatController) GetUserChats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "user not found"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "invalid user ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", strconv.Itoa(pagination.DefaultLimit))
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = pagination.DefaultLimit
	}

	pageStr := c.DefaultQuery("currentPage", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	limit = pagination.NormalizeLimit(limit)
	page = pagination.NormalizePage(page)
	search := c.Query("search")

	chats, totalPages, total, err := cc.chatUsecase.GetUserChats(userIDUint, limit, page, search)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	response := gin.H{
		"success": true,
		"data": gin.H{
			"items":       chats,
			"currentPage": page,
			"totalPages":  totalPages,
			"total":       total,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetChatMessages godoc
// @Summary Get chat messages
// @Description Returns paginated list of messages for a specific chat
// @Tags chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Chat ID"
// @Param limit query int false "Number of items per page (default: 20, max: 100)"
// @Param currentPage query int false "Current page number (default: 1)"
// @Success 200 {object} map[string]interface{} "List of messages with pagination info"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /chats/{id}/messages [get]
func (cc *ChatController) GetChatMessages(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "user not found"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "invalid user ID"})
		return
	}

	chatIDStr := c.Param("id")
	chatID, err := strconv.ParseUint(chatIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid chat ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", strconv.Itoa(pagination.DefaultLimit))
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = pagination.DefaultLimit
	}

	pageStr := c.DefaultQuery("currentPage", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	limit = pagination.NormalizeLimit(limit)
	page = pagination.NormalizePage(page)

	messages, totalPages, total, err := cc.chatUsecase.GetChatMessages(uint(chatID), userIDUint, limit, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	response := gin.H{
		"success": true,
		"data": gin.H{
			"items":       messages,
			"currentPage": page,
			"totalPages":  totalPages,
			"total":       total,
		},
	}

	c.JSON(http.StatusOK, response)
}

type CreateMessageRequest struct {
	RecipientID uint   `json:"recipientId" binding:"required"`
	Text        string `json:"text" binding:"required,min=1"`
}

// CreateMessage godoc
// @Summary Create message
// @Description Creates a new message in a chat. If chat doesn't exist between users, creates a new chat
// @Tags chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateMessageRequest true "Message creation request"
// @Success 200 {object} map[string]interface{} "Created message"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /chat/message [post]
func (cc *ChatController) CreateMessage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "user not found"})
		return
	}

	senderID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "invalid user ID"})
		return
	}

	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if senderID == req.RecipientID {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "cannot send message to yourself"})
		return
	}

	message, err := cc.chatUsecase.CreateMessage(senderID, req.RecipientID, req.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if cc.hub != nil {
		wsMessage := &websocket.Message{
			Type:    "new_message",
			Message: message,
		}
		cc.hub.BroadcastToUser(req.RecipientID, wsMessage)
		cc.hub.BroadcastToUser(senderID, wsMessage)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    message,
	})
}

func (cc *ChatController) HandleWebSocket(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "user not found"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "invalid user ID"})
		return
	}

	upgrader := ws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := websocket.NewClient(cc.hub, conn, userIDUint)
	cc.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
