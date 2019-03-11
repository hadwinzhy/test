package frequent_person

import (
	"net/http"
	"siren/models"
	"siren/pkg/database"
	"siren/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 判定用户是否是回头客
func JudgeFrequentPersonHandler(c *gin.Context) {

	var params JudgeParams
	if err := c.ShouldBindQuery(&params); err != nil {
		return
	}

	var group models.FrequentCustomerGroup
	if dbError := database.POSTGRES.Where("company_id = ? AND shop_id = ?", params.CompanyID, params.ShopID).First(&group).Error; dbError != nil {
		return
	}

	var person models.FrequentCustomerPeople
	captureTime := time.Unix(params.CreatedAt, 0)

	today := utils.CurrentDate(captureTime)
	if dbError := database.POSTGRES.Preload("FrequentCustomerGroup").
		Where("frequent_customer_group_id = ?", group.ID).
		Where("person_id = ?", params.PersonID).
		Where("hour >= ?", today).
		Where("hour < ?", today.AddDate(0, 0, 1)).
		First(&person).Error; dbError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "true"})

}
