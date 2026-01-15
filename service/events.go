package service

import (
	"challenge/models"
	"challenge/repository"
	"challenge/utils"
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server holds the Echo instance and database
type Server struct {
	Echo *echo.Echo
	DB   *repository.Database
}

// NewServer creates a new server instance
func NewServer(db *repository.Database) *Server {
	e := echo.New()

	// Middlewarego
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	server := &Server{
		Echo: e,
		DB:   db,
	}

	// Register routes
	server.registerRoutes()

	return server
}

// createEvent handles POST /events
// Accepts a JSON payload with title, description, start_time, and end_time
// Returns the created event as JSON with HTTP 201 status
func (s *Server) createEvent(c echo.Context) error {
	ctx := context.Background()

	// Parse request body
	var req models.CreateEventRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// Validate request
	if err := models.IsValid(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Parse timestamps
	startTime, _ := utils.ParseTimestamp(req.StartTime)
	endTime, _ := utils.ParseTimestamp(req.EndTime)

	// Create event object
	event := &models.Event{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	// Insert into database (ID and CreatedAt will be generated automatically)
	if err := s.DB.InsertEvent(ctx, event); err != nil {
		log.Printf("Error inserting event: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create event",
		})
	}

	// Return created event with 201 status
	return c.JSON(http.StatusCreated, event)
}

// registerRoutes sets up all the API routes
func (s *Server) registerRoutes() {
	// API v1 routes
	api := s.Echo.Group("/api/v1")
	api.POST("/events", s.createEvent)
	api.GET("/events", s.listEvents)
	api.GET("/events/:id", s.getEventByID)
}

// listEvents handles GET /events
// Returns a JSON array of all events ordered by start_time ascending
func (s *Server) listEvents(c echo.Context) error {
	ctx := context.Background()

	events, err := s.DB.GetAllEvents(ctx)
	if err != nil {
		log.Printf("Error getting events: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve events",
		})
	}

	// Return empty array instead of null if no events
	if events == nil {
		events = []*models.Event{}
	}

	return c.JSON(http.StatusOK, events)
}

// getEventByID handles GET /events/:id
// Returns the event with the specified UUID or 404 if not found
func (s *Server) getEventByID(c echo.Context) error {
	ctx := context.Background()

	// Parse UUID from path parameter
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "Invalid UUID format",
		})
	}

	// Get event from database
	event, err := s.DB.GetEventByID(ctx, id)
	if err != nil {
		if err.Error() == "event not found" {
			return echo.NewHTTPError(http.StatusNotFound, map[string]string{
				"error": "Event not found",
			})
		}
		log.Printf("Error getting event by ID: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve event",
		})
	}

	return c.JSON(http.StatusOK, event)
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	return s.Echo.Start(":" + port)
}
