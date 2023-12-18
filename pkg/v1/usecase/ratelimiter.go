// Package usecase acts as a bridge between the repository and the business logic
package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/thegeekywanderer/fluxy/models"
	interfaces "github.com/thegeekywanderer/fluxy/pkg/v1"
	algorithm "github.com/thegeekywanderer/fluxy/pkg/v1/algorithms"
	"gorm.io/gorm"
)

// UseCase struct
type UseCase struct {
  repo      interfaces.RepoInterface
  algorithm string
  cache     *redis.Client 
}

// New function instantiates a new usecase
func New(repo interfaces.RepoInterface, algorithm string, cache *redis.Client) interfaces.UseCaseInterface {
  return &UseCase{repo, algorithm, cache}
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
func (uc *UseCase) VerifyLimit(name string) (*interfaces.Result, error) {
  strategy, err := algorithm.New(uc.algorithm, uc.cache, time.Now)
  if err != nil {
    log.Fatal(err)
  }
  dataKey := fmt.Sprintf("%s-data", name)
  val, err := uc.cache.Get(dataKey).Result()
  var client models.Client
  request := interfaces.Request{}
  if err != nil {
    client, err = uc.repo.GetClient(name)
    if err != nil {
      return nil, err
    }
    request.Key = client.Name
    request.Limit = client.Limit
    request.Duration = time.Duration(client.Duration) * time.Second
    json, err := json.Marshal(client)
    dataKey := fmt.Sprintf("%s-data", client.Name)
    err = uc.cache.Set(dataKey, json, 0).Err()
    if err != nil {
      return nil, err
    }
	}

  err = json.Unmarshal([]byte(val), &client)
  if err != nil {
    return nil, err
  }
  request.Key = client.Name
  request.Limit = client.Limit
  request.Duration = time.Duration(client.Duration) * time.Second

  res, err := strategy.Run(&request)
  if err != nil {
    return nil, err
  }
  return res, nil
}