package main
import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"regexp"
	"strconv"
	"math"
	"strings"
)

// This function is for the endpoint '/receipts/process', upon triggering, it will calculate points, store it in memory and return an ID for the given payload.
func processReceipt(c *gin.Context) {
	var receipt Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidTotal(receipt.Total) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid total format"})
		return
	}

	id := uuid.New().String()
	receipts[id] = receipt
	points[id] = calculatePoints(receipt)
	c.JSON(http.StatusOK, ReceiptResponse{ID: id})
}

// This function returns the points when the endpoint '/receipts/:id/points' is triggered. 
func getPoints(c *gin.Context) {
	id := c.Param("id")
	if _, exists := receipts[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"points": points[id]})
}

// This function calculates the points based on the given rules in the problem statement. 
func calculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name.
	re := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(re.FindAllString(receipt.Retailer, -1))

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	if val, err := strconv.ParseFloat(receipt.Total, 64); err == nil {
		if val == float64(int(val)) {
			points += 50
		}
	}

	// Rule 3: 25 points if the total is a multiple of 0.25.
	if val, err := strconv.ParseFloat(receipt.Total, 64); err == nil {
		if int(val*100)%25 == 0 {
			points += 25
		}
	}

	// Rule 4: 5 points for every two items
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		// Trim leading and trailing spaces before checking length
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			val, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				pointsEarned := math.Ceil(val * 0.2) // Multiply price by 0.2 and round up
				points += int(pointsEarned)
			}
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	if isOddDay(receipt.PurchaseDate) {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	if isAfternoon(receipt.PurchaseTime) {
		points += 10
	}
	return points
}

// check if the day is odd
func isOddDay(date string) bool {
	re := regexp.MustCompile(`^\d{4}-(\d{2})-(\d{2})$`)
	matches := re.FindStringSubmatch(date)
	if len(matches) != 3 {
		return false
	}
	day, err := strconv.Atoi(matches[2])
	if err != nil {
		return false
	}
	return day%2 != 0
}

// check if the hour is after 2 pm and before 4 pm
func isAfternoon(time string) bool {
	re := regexp.MustCompile(`^(\d{2}):(\d{2})$`)
	matches := re.FindStringSubmatch(time)
	if len(matches) != 3 {
		return false
	}
	hour, err := strconv.Atoi(matches[1])
	if err != nil {
		return false
	}
	return hour >= 14 && hour < 16
}

//check if the receipt total is valid
func isValidTotal(total string) bool {
	re := regexp.MustCompile(`^\d+\.\d{2}$`)
	return re.MatchString(total)
}

//Checking if the payload is correct and complete. 
type Receipt struct {
	Retailer     string `json:"retailer" binding:"required"`
	PurchaseDate string `json:"purchaseDate" binding:"required"`
	PurchaseTime string `json:"purchaseTime" binding:"required"`
	Items        []Item `json:"items" binding:"required,min=1"`
	Total        string `json:"total" binding:"required"`
}

type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required"`
	Price            string `json:"price" binding:"required"`
}

type ReceiptResponse struct {
	ID string `json:"id"`
}

var receipts = make(map[string]Receipt) // created this map to store the receipt data
var points = make(map[string]int) // created this map to store the points

func main() {
	r := gin.Default()
	r.POST("/receipts/process", processReceipt)
	r.GET("/receipts/:id/points", getPoints)
	// I am using port 8080 to run the server.
	r.Run(":8080")
}
