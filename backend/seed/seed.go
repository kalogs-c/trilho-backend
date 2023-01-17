package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/kalogsc/ego/models"
)

var users = []models.User{
	{
		Name:     "Peter",
		LastName: "Parker",
		Email:    "notspiderman@gmail.com",
		Password: "web123",
	},
	{
		Name:     "Alessia",
		LastName: "Cara",
		Email:    "amusic@gmail.com",
		Password: "coxinha123",
	},
}

var transactions = []models.Transaction{
	{
		Name:   "Coffe",
		Amount: 600,
	},
	{
		Name:   "Pão de queijo",
		Amount: -250,
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Transaction{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Transaction{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Transaction{}).AddForeignKey("owner_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range users {
		err = users[i].Save(db)
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		transactions[i].OwnerId = users[i].ID

		err = transactions[i].Save(db)
		if err != nil {
			log.Fatalf("cannot seed transactions table: %v", err)
		}
	}
}
