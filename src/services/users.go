package services

import (
	"go-fiber-template/domain/entities"
	"go-fiber-template/domain/repositories"
	"go-fiber-template/httpclient"
	"go-fiber-template/src/utils"
	"mime/multipart"
	"time"
)

type usersService struct {
	UsersRepository repositories.IUsersRepository
}

type IUsersService interface {
	GetAllUsers() (*[]entities.UserDataModel, error)
	InsertNewUser(data entities.UserDataModel) error
	UpdateImage(userID string, imageFile *multipart.FileHeader) error
}

func NewUsersService(repo0 repositories.IUsersRepository) IUsersService {
	return &usersService{
		UsersRepository: repo0,
	}
}

func (sv *usersService) GetAllUsers() (*[]entities.UserDataModel, error) {
	data, err := sv.UsersRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (sv *usersService) InsertNewUser(data entities.UserDataModel) error {
	data.CreatedAt = time.Now().Add(7 * time.Hour)
	dataIp, err := httpclient.GetIP()
	if err != nil {
		return err
	}
	data.Ip = dataIp

	return sv.UsersRepository.InsertUser(data)
}

func (sv *usersService) UpdateImage(userID string, imageFile *multipart.FileHeader) error {
	keyName := utils.CreateKeyName(userID)

	imageContent, err := imageFile.Open()
	if err != nil {
		return err
	}
	defer imageContent.Close()

	imageBytes := make([]byte, imageFile.Size)
	_, err = imageContent.Read(imageBytes)
	if err != nil {
		return err
	}

	url, err := utils.UploadS3FromString(imageBytes, keyName, "image/png")
	if err != nil {
		return err
	}
	if err := sv.UsersRepository.UpdateImage(userID, url); err != nil {
		return err
	}
	return nil
}
