// Package repository contains the implementation of data access methods
package repository

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/thegeekywanderer/fluxy/models"
	interfaces "github.com/thegeekywanderer/fluxy/pkg/v1"
	"gorm.io/gorm"
)

// Repo struct
type Repo struct {
  db *gorm.DB
  cache *redis.Client
}

// New function creates a new Repo instance
func New(db *gorm.DB, cache *redis.Client) interfaces.RepoInterface {
  return &Repo{db, cache}
}

// RegisterClient function enters a new client into the fluxy database
func (repo *Repo) RegisterClient(client models.Client) (models.Client, error){
  err := repo.db.Create(&client).Error
  if err != nil {
    return client, err
  }
  json, err := json.Marshal(client)
  dataKey := fmt.Sprintf("%s-data", client.Name)
  err = repo.cache.Set(dataKey, json, 0).Err()
  if err != nil {
    return client, err
  }
  return client, nil
}

// GetClient function returns the client instance searched by name
func (repo *Repo) GetClient(name string) (models.Client, error){
  var client models.Client
  err := repo.db.Where("name = ?", name).First(&client).Error
  return client, err
}

// UpdateClient function updates the client 
func (repo *Repo) UpdateClient(client models.Client) error{
  var dbClient models.Client
  if err := repo.db.Where("name = ?", client.Name).First(&dbClient).Error; err != nil {
    return err
  }
  dbClient.Limit = client.Limit
  dbClient.Duration = client.Duration
  err := repo.db.Save(dbClient).Error
  json, err := json.Marshal(client)
  dataKey := fmt.Sprintf("%s-data", client.Name)
  err = repo.cache.Set(dataKey, json, 0).Err()
  if err != nil {
    return err
  }
  return err
}

// DeleteClient function deletes an existing client 
func (repo *Repo) DeleteClient(name string) error{
  err := repo.db.Unscoped().Where("name = ?", name).Delete(&models.Client{}).Error
  return err
}

