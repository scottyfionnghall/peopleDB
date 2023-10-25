package controllers

import (
	"net/http"
	"peopledb/models"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

var OrderByPattern = `((^\bage\b.*)|(^\bid\b.*)|(^\bname\b.*)|(^\bsurname\b.*)|(^\bnationality\b.*)|(^\bgender\b.*))((\basc\b)|(\bdesc\b))`

type CreateUserInput struct {
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic,"`
}

type ReturnUserOutout struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type UpdateUserInput struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Nationality string `json:"nationality"`
	Gender      string `json:"gender"`
	Age         string `json:"age"`
}

func AddUser(c *gin.Context) {
	var input CreateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data or missing fields"})
		return
	}

	age, nationality, gender, err := populateInfo(input.Name)
	if err != nil {
		if err.Error() == "third party api not available" {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Third party API not available, try again later"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data or missing fields"})
		return
	}

	user := models.User{
		Name:        input.Name,
		Surname:     input.Surname,
		Patronymic:  input.Patronymic,
		Age:         age,
		Nationality: nationality,
		Gender:      gender,
	}

	models.DB.Create(&user)

	return_user := ReturnUserOutout{
		ID:          user.ID,
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		Age:         user.Age,
		Gender:      user.Gender,
		Nationality: user.Nationality,
	}

	c.JSON(http.StatusOK, gin.H{"user": return_user})
}

func GetUser(c *gin.Context) {
	var user models.User

	name := c.DefaultQuery("name", "")
	surname := c.DefaultQuery("surname", "")
	id := c.DefaultQuery("id", "")

	switch {
	case name != "":
		if err := models.DB.Where("name=?", name).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find specific user"})
			return
		}
	case surname != "":
		if err := models.DB.Where("surname=?", surname).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find specific user"})
			return
		}
	case id != "":
		if err := models.DB.Where("id=?", id).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find specific user"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad query params"})
		return
	}

	return_user := ReturnUserOutout{
		ID:          user.ID,
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		Age:         user.Age,
		Gender:      user.Gender,
		Nationality: user.Nationality,
	}

	c.JSON(http.StatusOK, gin.H{"data": return_user})
}

func GetAllUsers(c *gin.Context) {
	var users []models.User

	name := c.DefaultQuery("name", "")
	surname := c.DefaultQuery("surname", "")
	gender := c.DefaultQuery("gender", "")
	age := c.DefaultQuery("age", "")
	orderby := c.DefaultQuery("order_by", "id asc")
	nationality := c.DefaultQuery("nationality", "")

	check, _ := regexp.MatchString(OrderByPattern, orderby)
	if !check {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad query params"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad query params"})
		return
	}
	if page <= 0 {
		page = 1
	}

	page_size, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad query params"})
		return
	}
	switch {
	case page_size > 100:
		page_size = 100
	case page_size <= 0:
		page_size = 10
	}

	offset := (page - 1) * page_size

	switch {
	case nationality != "":
		if err := models.DB.Order(orderby).Offset(offset).Limit(page_size).Where("nationality", nationality).Find(&users).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find users by specified params"})
			return
		}
	case gender != "":
		if err := models.DB.Order(orderby).Offset(offset).Limit(page_size).Where("gender", gender).Find(&users).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find users by specified params"})
			return
		}
	case age != "":
		if err := models.DB.Order(orderby).Offset(offset).Limit(page_size).Where("age", age).Find(&users).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find users by specified params"})
			return
		}
	case name != "":
		if err := models.DB.Order(orderby).Offset(offset).Limit(page_size).Where("name", name).Find(&users).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find users by specified params"})
			return
		}
	case surname != "":
		if err := models.DB.Order(orderby).Offset(offset).Limit(page_size).Where("surname", surname).Find(&users).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find users by specified params"})
			return
		}
	case surname == "" && name == "":
		if err := models.DB.Order(orderby).Offset(offset).Limit(page_size).Find(&users).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find users by specified params"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad query params"})
		return
	}

	var return_users []ReturnUserOutout

	for _, user := range users {

		return_user := ReturnUserOutout{
			ID:          user.ID,
			Name:        user.Name,
			Surname:     user.Surname,
			Patronymic:  user.Patronymic,
			Age:         user.Age,
			Gender:      user.Gender,
			Nationality: user.Nationality,
		}
		return_users = append(return_users, return_user)
	}

	c.JSON(http.StatusOK, gin.H{"data": return_users})
}

func UpdateUser(c *gin.Context) {
	var user models.User

	id, err := strconv.Atoi(c.DefaultQuery("id", ""))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad query params"})
		return
	}

	if err := models.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad query params"})
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error in request body"})
		return
	}

	models.DB.Model(&user).Updates(input)

	return_user := ReturnUserOutout{
		ID:          user.ID,
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		Age:         user.Age,
		Gender:      user.Gender,
		Nationality: user.Nationality,
	}

	c.JSON(http.StatusOK, gin.H{"data": return_user})

}

func DeleteUser(c *gin.Context) {
	var user models.User

	id, err := strconv.Atoi(c.DefaultQuery("id", ""))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad query params"})
		return
	}

	if err := models.DB.Where("id=?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad query params"})
		return
	}

	models.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
