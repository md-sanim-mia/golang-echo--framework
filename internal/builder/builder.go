package builder

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type QueryBuilder struct {
	db         *gorm.DB
	query      *gorm.DB
	ctx        echo.Context
	countQuery *gorm.DB // Store query before pagination for counting
}

type PaginationMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

func NewQueryBuilder(db *gorm.DB, ctx echo.Context) *QueryBuilder {
	return &QueryBuilder{
		db:    db,
		query: db,
		ctx:   ctx,
	}
}

func (qb *QueryBuilder) Model(model interface{}) *QueryBuilder {
	qb.query = qb.query.Model(model)
	return qb
}

func (qb *QueryBuilder) Search(searchableFields []string) *QueryBuilder {
	searchTerm := qb.ctx.QueryParam("searchTerm")
	if searchTerm == "" {
		searchTerm = qb.ctx.QueryParam("search")
	}

	if searchTerm != "" && len(searchableFields) > 0 {
		orConditions := qb.db.Session(&gorm.Session{})
		for i, field := range searchableFields {
			if i == 0 {
				orConditions = orConditions.Where(field+" ILIKE ?", "%"+searchTerm+"%")
			} else {
				orConditions = orConditions.Or(field+" ILIKE ?", "%"+searchTerm+"%")
			}
		}
		qb.query = qb.query.Where(orConditions)
	}
	return qb
}

func (qb *QueryBuilder) Filter() *QueryBuilder {
	minPrice := qb.ctx.QueryParam("minPrice")
	maxPrice := qb.ctx.QueryParam("maxPrice")
	if minPrice != "" {
		if val, err := strconv.ParseFloat(minPrice, 64); err == nil {
			qb.query = qb.query.Where("price >= ?", val)
		}
	}
	if maxPrice != "" {
		if val, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			qb.query = qb.query.Where("price <= ?", val)
		}
	}

	startDate := qb.ctx.QueryParam("startDate")
	endDate := qb.ctx.QueryParam("endDate")
	if startDate != "" {
		qb.query = qb.query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		qb.query = qb.query.Where("created_at <= ?", endDate)
	}

	excludeFields := map[string]bool{
		"searchTerm": true,
		"search":     true,
		"sort":       true,
		"limit":      true,
		"page":       true,
		"fields":     true,
		"minPrice":   true,
		"maxPrice":   true,
		"startDate":  true,
		"endDate":    true,
	}

	queryParams := qb.ctx.QueryParams()

	for key, values := range queryParams {
		if excludeFields[key] || len(values) == 0 {
			continue
		}

		if strings.Contains(key, "[") && strings.Contains(key, "]") {
			startIdx := strings.Index(key, "[")
			endIdx := strings.Index(key, "]")
			field := key[:startIdx]
			operator := key[startIdx+1 : endIdx]
			value := values[0]

			qb.applyOperator(field, operator, value)
		} else if len(values) > 1 {
			qb.query = qb.query.Where(key+" IN ?", values)
		} else {
			qb.query = qb.query.Where(key+" = ?", values[0])
		}
	}

	return qb
}

func (qb *QueryBuilder) applyOperator(field, operator, value string) {
	switch operator {
	case "gt", "gte", "lt", "lte":
		numValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			switch operator {
			case "gt":
				qb.query = qb.query.Where(field+" > ?", numValue)
			case "gte":
				qb.query = qb.query.Where(field+" >= ?", numValue)
			case "lt":
				qb.query = qb.query.Where(field+" < ?", numValue)
			case "lte":
				qb.query = qb.query.Where(field+" <= ?", numValue)
			}
		}
	case "in":
		values := strings.Split(value, ",")
		qb.query = qb.query.Where(field+" IN ?", values)
	case "nin", "notin":
		values := strings.Split(value, ",")
		qb.query = qb.query.Where(field+" NOT IN ?", values)
	case "like":
		qb.query = qb.query.Where(field+" LIKE ?", "%"+value+"%")
	case "ilike":
		qb.query = qb.query.Where(field+" ILIKE ?", "%"+value+"%")
	case "ne", "neq":
		qb.query = qb.query.Where(field+" != ?", value)
	case "between":
		parts := strings.Split(value, ",")
		if len(parts) == 2 {
			qb.query = qb.query.Where(field+" BETWEEN ? AND ?", parts[0], parts[1])
		}
	default:
		qb.query = qb.query.Where(field+" = ?", value)
	}
}

func (qb *QueryBuilder) RawFilter(filters map[string]interface{}) *QueryBuilder {
	for key, value := range filters {
		qb.query = qb.query.Where(key, value)
	}
	return qb
}

func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.query = qb.query.Where(condition, args...)
	return qb
}

func (qb *QueryBuilder) Sort() *QueryBuilder {
	sortParam := qb.ctx.QueryParam("sort")
	if sortParam == "" {
		sortParam = "-created_at"
	}

	sortFields := strings.Split(sortParam, ",")
	for _, field := range sortFields {
		if strings.HasPrefix(field, "-") {
			qb.query = qb.query.Order(field[1:] + " DESC")
		} else {
			qb.query = qb.query.Order(field + " ASC")
		}
	}

	return qb
}

func (qb *QueryBuilder) Paginate() *QueryBuilder {
	// Save query state BEFORE pagination for counting
	qb.countQuery = qb.query.Session(&gorm.Session{})

	page, _ := strconv.Atoi(qb.ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(qb.ctx.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	qb.query = qb.query.Limit(limit).Offset(offset)

	return qb
}

func (qb *QueryBuilder) Fields() *QueryBuilder {
	fieldsParam := qb.ctx.QueryParam("fields")
	if fieldsParam != "" {
		fields := strings.Split(fieldsParam, ",")
		qb.query = qb.query.Select(fields)
	}
	return qb
}

func (qb *QueryBuilder) Select(fields []string) *QueryBuilder {
	qb.query = qb.query.Select(fields)
	return qb
}

func (qb *QueryBuilder) Preload(associations ...string) *QueryBuilder {
	for _, assoc := range associations {
		qb.query = qb.query.Preload(assoc)
	}
	return qb
}

func (qb *QueryBuilder) PreloadWithCondition(association string, conditions ...interface{}) *QueryBuilder {
	qb.query = qb.query.Preload(association, conditions...)
	return qb
}

func (qb *QueryBuilder) PriceRange(minPrice, maxPrice *float64) *QueryBuilder {
	if minPrice != nil {
		qb.query = qb.query.Where("price >= ?", *minPrice)
	}
	if maxPrice != nil {
		qb.query = qb.query.Where("price <= ?", *maxPrice)
	}
	return qb
}

func (qb *QueryBuilder) AutoPriceRange() *QueryBuilder {
	minPrice := qb.ctx.QueryParam("minPrice")
	maxPrice := qb.ctx.QueryParam("maxPrice")

	if minPrice != "" {
		if val, err := strconv.ParseFloat(minPrice, 64); err == nil {
			qb.query = qb.query.Where("price >= ?", val)
		}
	}
	if maxPrice != "" {
		if val, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			qb.query = qb.query.Where("price <= ?", val)
		}
	}
	return qb
}

func (qb *QueryBuilder) Range(field, minParam, maxParam string) *QueryBuilder {
	minVal := qb.ctx.QueryParam(minParam)
	maxVal := qb.ctx.QueryParam(maxParam)

	if minVal != "" {
		if val, err := strconv.ParseFloat(minVal, 64); err == nil {
			qb.query = qb.query.Where(field+" >= ?", val)
		}
	}
	if maxVal != "" {
		if val, err := strconv.ParseFloat(maxVal, 64); err == nil {
			qb.query = qb.query.Where(field+" <= ?", val)
		}
	}
	return qb
}

func (qb *QueryBuilder) DateRange(field, startDate, endDate string) *QueryBuilder {
	if startDate != "" {
		qb.query = qb.query.Where(field+" >= ?", startDate)
	}
	if endDate != "" {
		qb.query = qb.query.Where(field+" <= ?", endDate)
	}
	return qb
}

func (qb *QueryBuilder) Execute(dest interface{}) error {
	return qb.query.Find(dest).Error
}

func (qb *QueryBuilder) First(dest interface{}) error {
	return qb.query.First(dest).Error
}

func (qb *QueryBuilder) CountTotal() (*PaginationMeta, error) {
	var total int64

	// Use saved query from before pagination
	if qb.countQuery != nil {
		if err := qb.countQuery.Count(&total).Error; err != nil {
			return nil, err
		}
	} else {
		// Fallback: remove limit/offset and count
		if err := qb.db.Session(&gorm.Session{}).Table(qb.query.Statement.Table).Count(&total).Error; err != nil {
			return nil, err
		}
	}

	page, _ := strconv.Atoi(qb.ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(qb.ctx.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	totalPage := int(total) / limit
	if int(total)%limit > 0 {
		totalPage++
	}

	return &PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	}, nil
}

func (qb *QueryBuilder) GetQuery() *gorm.DB {
	return qb.query
}
