package handlers

import "gitea.kood.tech/ivanandreev/viewer/internal/domain"

type CatalogueProvider interface {
	CarModels(int) ([]domain.Car, error)              // maybe use int to limit the amount of fethed items?
	CarManufacturers() ([]domain.Manufacturer, error) // we need to get all
	CarCategories() ([]domain.Category, error)        // we need to get all
}
