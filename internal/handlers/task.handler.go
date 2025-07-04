package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/praction-networks/acs-proxy/internal/services"
	"github.com/praction-networks/common/appError"
	"github.com/praction-networks/common/helpers"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/response"
)

type TaskHandler struct {
	TaskService services.TaskService
}

func NewTaskHandler(svc services.TaskService) *TaskHandler {
	return &TaskHandler{
		TaskService: svc,
	}
}

// @Summary      Retry Task
// @Description  Retry a failed or pending ACS task
// @Tags         Tasks
// @Produce      json
// @Param        task_id  path  string  true  "Task ID"
// @Success      200  {object}  models.DeviceResponseModel
// @Failure      400  {object}  models.BaseError
// @Failure      500  {object}  models.BaseError
// @Router       /acs-proxy/tasks/{task_id}/retry [post]
// @Security     ApiKeyAuth
func (h *TaskHandler) RetryTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "task_id")
	if taskID == "" {
		helpers.HandleAppError(w, appError.New(appError.InvalidInputError, "task_id is required", http.StatusBadRequest, nil))
		return
	}

	logger.Info("Retrying task", "taskID", taskID)
	err := h.TaskService.RetryTask(r.Context(), taskID)
	if err != nil {
		logger.Error("Retry task failed", err)
		helpers.HandleAppError(w, err)
		return
	}
	response.Send200OK(w, "Task retried successfully", nil)
}

// @Summary      Delete Task
// @Description  Delete a scheduled or completed ACS task
// @Tags         Tasks
// @Produce      json
// @Param        task_id  path  string  true  "Task ID"
// @Success      200  {object}  models.DeviceResponseModel
// @Failure      400  {object}  models.BaseError
// @Failure      500  {object}  models.BaseError
// @Router       /acs-proxy/tasks/{task_id} [delete]
// @Security     ApiKeyAuth
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "task_id")
	if taskID == "" {
		helpers.HandleAppError(w, appError.New(appError.InvalidInputError, "task_id is required", http.StatusBadRequest, nil))
		return
	}

	logger.Info("Deleting task", "taskID", taskID)
	err := h.TaskService.DeleteTask(r.Context(), taskID)
	if err != nil {
		logger.Error("Delete task failed", err)
		helpers.HandleAppError(w, err)
		return
	}
	response.Send200OK(w, "Task deleted successfully", nil)
}
