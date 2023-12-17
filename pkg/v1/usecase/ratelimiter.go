package usecase

import (
	"errors"

	"github.com/thegeekywanderer/fluxy/models"
	interfaces "github.com/thegeekywanderer/fluxy/pkg/v1"
	"gorm.io/gorm"
)

// UseCase struct
type UseCase struct {
  repo interfaces.RepoInterface
}

// New function instantiates a new usecase
func New(repo interfaces.RepoInterface) interfaces.UseCaseInterface {
  return &UseCase{repo}
}

// RegisterClient function creates a new client which was supplied as the argument
func (uc *UseCase) RegisterClient(client models.Client) (models.Client, error) {
  if _, err := uc.repo.GetClient(client.Name); !errors.Is(err, gorm.ErrRecordNotFound){
    return models.Client{}, errors.New("the name is already associated with another client")
  }
  return uc.repo.RegisterClient(client)
}

// GetClient function returns an existing client instance
func (uc *UseCase) GetClient(name string) (models.Client, error){
  var client models.Client
  var err error
  if client, err = uc.repo.GetClient(name); err != nil{
    if errors.Is(err, gorm.ErrRecordNotFound){
      return models.Client{}, errors.New("no such client with the name supplied")
    }
    return models.Client{}, err
  }

  return client, nil
}

// UpdateClient function updates the client limits
func (uc *UseCase) UpdateClient(updateClient models.Client) error{
  // check if client exists
  if _, err := uc.GetClient(string(updateClient.Name)); err != nil {
    return err
  }

  err := uc.repo.UpdateClient(updateClient)
  if err != nil {
    return err
  }
  return nil
}

// DeleteClient function deletes an existing client
func (uc *UseCase) DeleteClient(name string) error{
  var err error
  // check if client exists
  if _, err = uc.GetClient(name); err != nil {
    return err
  }

  err = uc.repo.DeleteClient(name)
  if err != nil {
    return err
  }

  return nil
}

// VerifyLimit functions validates the rate limit of the specified client
func (uc *UseCase) VerifyLimit(name string) (bool, error) {
  return true, nil
}