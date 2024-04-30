package dtos

// SQL request
// @Description Sql query
type SQLRequestDto struct {
	Sql string `json:"sql" validate:"sql" swaggertype:"string" example:"SELECT 1 LIMIT 1;"`
}

// SQL response
// @Description Sql response
type SQLResponseDto struct {
	Data []map[string]interface{} `json:"data"`
}
