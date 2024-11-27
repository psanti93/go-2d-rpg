package main

import (
	"encoding/json"
	"os"
)

type TilemapLayersJSON struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

type TilemapJSON struct {
	Layers []TilemapLayersJSON `json:"layers"`
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {

	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tileMap TilemapJSON

	err = json.Unmarshal(contents, &tileMap)

	return &tileMap, err

}
