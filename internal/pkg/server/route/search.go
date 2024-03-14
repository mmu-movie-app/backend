package route

import (
	"github.com/gin-gonic/gin"
	"movie-website-backend-core/internal/pkg/database"
	movie_db "movie-website-backend-core/internal/pkg/movie-db"
	"net/http"
	"strconv"
)

// filterType
const (
	Title     = 0
	Genres    = 1
	Year      = 2
	Budget    = 3
	Actors    = 4
	Revenue   = 5
	Directors = 6
)

// state
const (
	exact = 0
	high  = 1
	low   = 2
)

type Filter struct {
	FilterType int    `json:"filter_type"`
	Value      string `json:"value"`
	State      int    `json:"state"`
}

func Search(c *gin.Context) {
	var filters []Filter
	if err := c.BindJSON(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build the query based on the filters
	query := database.DB.Model(&movie_db.MovieDB{})

	// Keep track of applied filters
	appliedFilters := make(map[int]bool)

	for _, filter := range filters {
		switch filter.FilterType {
		case Title:
			if appliedFilters[Title] {
				query = query.Or("title LIKE ?", "%"+filter.Value+"%")
			} else {
				query = query.Where("title LIKE ?", "%"+filter.Value+"%")
				appliedFilters[Title] = true
			}
		case Genres:
			if appliedFilters[Genres] {
				query = query.Or("genre_dbs.name LIKE ?", "%"+filter.Value+"%")
			} else {
				query = query.Joins("JOIN movie_genres ON movie_dbs.id = movie_genres.movie_db_id").
					Joins("JOIN genre_dbs ON movie_genres.genre_db_id = genre_dbs.id").
					Where("genre_dbs.name LIKE ?", "%"+filter.Value+"%")
				appliedFilters[Genres] = true
			}
		case Year:
			year, err := strconv.Atoi(filter.Value)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year value"})
				return
			}
			if filter.State == exact {
				query = query.Where("release_year = ?", year)
			} else if filter.State == high {
				query = query.Where("release_year >= ?", year)
			} else if filter.State == low {
				query = query.Where("release_year <= ?", year)
			}
		case Budget:
			budget, err := strconv.Atoi(filter.Value)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget value"})
				return
			}
			if filter.State == exact {
				query = query.Where("budget = ?", budget)
			} else if filter.State == high {
				if appliedFilters[Budget] {
					query = query.Or("budget >= ?", budget)
				} else {
					query = query.Where("budget >= ?", budget)
					appliedFilters[Budget] = true
				}
			} else if filter.State == low {
				query = query.Where("budget <= ?", budget)
			}
		case Revenue:
			revenue, err := strconv.Atoi(filter.Value)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid revenue value"})
				return
			}
			if filter.State == exact {
				query = query.Where("revenue = ?", revenue)
			} else if filter.State == high {
				query = query.Where("revenue >= ?", revenue)
			} else if filter.State == low {
				query = query.Where("revenue <= ?", revenue)
			}
		case Actors:
			if appliedFilters[Actors] {
				query = query.Or("cast_dbs.name LIKE ?", "%"+filter.Value+"%")
			} else {
				query = query.Joins("JOIN movie_cast ON movie_dbs.id = movie_cast.movie_db_id").
					Joins("JOIN cast_dbs ON movie_cast.cast_db_id = cast_dbs.id").
					Where("cast_dbs.name LIKE ?", "%"+filter.Value+"%")
				appliedFilters[Actors] = true
			}
		case Directors:
			if appliedFilters[Directors] {
				query = query.Or("crew_dbs.name LIKE ?", "%"+filter.Value+"%")
			} else {
				query = query.Joins("JOIN movie_crew ON movie_dbs.id = movie_crew.movie_db_id").
					Joins("JOIN crew_dbs ON movie_crew.crew_db_id = crew_dbs.id").
					Where("crew_dbs.name LIKE ? AND crew_dbs.job = 'Director'", "%"+filter.Value+"%")
				appliedFilters[Directors] = true
			}
		}
	}

	var movies []movie_db.MovieDB
	if err := query.Preload("Genres").
		Preload("Cast").
		Preload("Crew").
		Preload("Backdrops").
		Preload("Posters").
		Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve movies"})
		return
	}

	c.JSON(http.StatusOK, movies)
}
