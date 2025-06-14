package service

import (
	"context"
	"errors"
	"github.com/spinmozgJr/note-service/internal/mapper"
	"github.com/spinmozgJr/note-service/internal/models"
	"github.com/spinmozgJr/note-service/internal/storage"
	"github.com/spinmozgJr/note-service/internal/storage/postgres"
)

type NoteService struct {
	NoteRepository storage.NoteStorage
}

func NewNoteService(noteRepository storage.NoteStorage) NoteService {
	return NoteService{
		NoteRepository: noteRepository,
	}
}

func (s *NoteService) CreateNote(ctx context.Context, note *NoteInput) (*models.CreateNoteResponse, error) {
	noteID, err := s.NoteRepository.CreateNote(ctx, note.Title, note.Content, note.UserID)
	if err != nil {
		return nil, err
	}

	response := &models.CreateNoteResponse{
		Data: &models.CreateNoteData{NoteID: noteID},
	}
	return response, nil
}

func (s *NoteService) DeleteNote(ctx context.Context, userID, noteID int) (*models.OperationResultResponse, error) {
	task, err := s.NoteRepository.GetNoteByID(ctx, noteID)
	if err != nil {
		if errors.Is(err, postgres.ErrNoteNotFound) {
			//logger.Get().Info(ctx, "Задача для удаления не найдена", logrus.Fields{
			//	"task_id": taskID,
			//})
			return nil, ErrTaskNotFound
		}
		//logger.Get().Error(ctx, "Ошибка при получении задачи перед удалением", logrus.Fields{
		//	"task_id": taskID,
		//	"error":   err.Error(),
		//})
		return nil, err
	}

	if task.UserID != userID {
		//logger.Get().Info(ctx, "Попытка удалить чужую задачу", logrus.Fields{
		//	"task_id":  taskID,
		//	"owner_id": task.UserID,
		//})
		return nil, ErrTaskForbidden
	}

	err = s.NoteRepository.DeleteNoteByID(ctx, noteID)
	if err != nil {
		//logger.Get().Error(ctx, "Ошибка при удалении задачи", logrus.Fields{
		//	"task_id": taskID,
		//	"error":   err.Error(),
		//})
		return nil, err
	}

	//logger.Get().Info(ctx, "Задача успешно удалена", logrus.Fields{
	//	"task_id": taskID,
	//})

	return &models.OperationResultResponse{Data: &models.OperationResultData{Success: true}}, nil
}

func (s *NoteService) GetAllNotes(ctx context.Context, params *ServiceQueryParams) (*models.NoteListResponse, error) {
	getAllNotesParams := &storage.InputGetAllNotes{
		UserID: params.ID,
		Limit:  params.Limit,
		Offset: params.Offset,
		Sort:   params.Sort,
	}
	notes, err := s.NoteRepository.GetAllNotes(ctx, getAllNotesParams)
	if err != nil {
		return nil, err
	}

	taskDTOs := make([]models.NoteDTO, len(notes))
	for i, e := range notes {
		taskDTOs[i] = mapper.MapNoteDTOFromTaskDb(e)
	}

	return &models.NoteListResponse{Data: &taskDTOs}, nil
}

func (s *NoteService) GetNoteByID(ctx context.Context, userId, noteId int) (*models.NoteResponse, error) {
	note, err := s.NoteRepository.GetNoteByID(ctx, noteId)
	if err != nil {
		//if errors.Is(err, postgres.ErrNoteNotFound) {
		//	logger.Get().Info(ctx, "Задача не найдена", logrus.Fields{
		//		"task_id": taskID,
		//	})
		//	return nil, ErrTaskNotFound
		//}
		//logger.Get().Error(ctx, "Ошибка при получении задачи", logrus.Fields{
		//	"task_id": taskID,
		//	"error":   err.Error(),
		//})
		return nil, err
	}

	if note.UserID != userId {
		//logger.Get().Info(ctx, "Попытка доступа к чужой задаче", logrus.Fields{
		//	"task_id":  taskID,
		//	"owner_id": task.UserID,
		//})
		return nil, ErrTaskForbidden
	}

	noteDTO := mapper.MapNoteDTOFromTaskDb(*note)
	return &models.NoteResponse{Data: &noteDTO}, nil
}

func (s *NoteService) UpdateNote(ctx context.Context, input *UpdateNote) (*models.OperationResultResponse, error) {
	note, err := s.NoteRepository.GetNoteByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrNoteNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	if note.UserID != input.UserID {
		return nil, ErrTaskForbidden
	}

	storageInput := &storage.UpdateNote{
		ID:      input.ID,
		Title:   input.Title,
		Content: input.Content,
	}

	err = s.NoteRepository.UpdateNote(ctx, storageInput)
	if err != nil {
		return nil, err
	}

	return &models.OperationResultResponse{Data: &models.OperationResultData{Success: true}}, nil
}
