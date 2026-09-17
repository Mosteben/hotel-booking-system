// This seeder is for a fresh, empty development database only.
//
// It is NOT idempotent: User.Email and User.Phone are unique at the DB
// level, so running this twice against a database that already has this
// seed's data fails loudly on the second run (a duplicate-key error) rather
// than silently duplicating or corrupting anything. That matches this
// script's existing behavior for hotels/rooms/bookings before this fix -
// nothing here resets or deletes data, and this change doesn't add that
// either. Point DB_NAME at an empty database before running it.
package main

import (
	"fmt"
	"time"

	"github.com/Mosteben/hotel-booking-system/configs"
	"github.com/Mosteben/hotel-booking-system/internal/booking/model"
	hotelModel "github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	roomModel "github.com/Mosteben/hotel-booking-system/internal/room/model"
	userModel "github.com/Mosteben/hotel-booking-system/internal/user/model"
	"github.com/Mosteben/hotel-booking-system/pkg/database"
	"github.com/Mosteben/hotel-booking-system/pkg/hash"
)

// SeedPassword is the shared login password for every seeded user
// (admin/manager/guests). It's printed at the end of the run so a
// developer can actually log in with it - never stored or logged as
// plaintext anywhere in the database itself.
const SeedPassword = "SeedPass123!"

func main() {

	// =========================
	// Configuration
	// =========================

	configs.LoadEnv()

	// =========================
	// Database
	// =========================

	database.Connect()

	fmt.Println("Starting mock data seeder...")

	// =========================
	// Seed Users
	// =========================
	//
	// Every seeded Booking below references one of these users' REAL
	// database-assigned IDs - never a hand-written UUID. GORM's
	// User.BeforeCreate hook assigns each user's ID; Create() (including
	// batch Create on a slice) populates that ID back onto the struct
	// after insert, which is what gets read into userIDs below.

	hashedPassword, err := hash.HashPassword(SeedPassword)
	if err != nil {
		panic(err)
	}

	users := []userModel.User{
		{
			FirstName: "Seed",
			LastName:  "Admin",
			Email:     "admin@seed.nilestay.local",
			Password:  hashedPassword,
			Phone:     "+205000000001",
			Role:      "admin",
			IsActive:  true,
		},
		{
			FirstName: "Seed",
			LastName:  "Manager",
			Email:     "manager@seed.nilestay.local",
			Password:  hashedPassword,
			Phone:     "+205000000002",
			Role:      "manager",
			IsActive:  true,
		},
	}

	const guestCount = 6
	for i := 1; i <= guestCount; i++ {
		users = append(users, userModel.User{
			FirstName: "Guest",
			LastName:  fmt.Sprintf("%d", i),
			Email:     fmt.Sprintf("guest%d@seed.nilestay.local", i),
			Password:  hashedPassword,
			Phone:     fmt.Sprintf("+2050000000%02d", 10+i),
			Role:      "customer",
			IsActive:  true,
		})
	}

	if err := database.DB.Create(&users).Error; err != nil {
		panic(err)
	}

	fmt.Printf("%d users created (1 admin, 1 manager, %d guests).\n", len(users), guestCount)

	// Guest bookings only ever reference guest accounts, never the
	// seeded admin/manager - matches how the real app works (a customer
	// books, an admin manages).
	var guestIDs []string
	for _, u := range users {
		if u.Role == "customer" {
			guestIDs = append(guestIDs, u.ID.String())
		}
	}

	// =========================
	// Seed Hotels
	// =========================

	hotels := []hotelModel.Hotel{
		{
			Name:        "Mock Hotel 01",
			Description: "Luxury hotel in Cairo",
			Address:     "Nile Street 1",
			City:        "Cairo",
			Country:     "Egypt",
			Phone:       "+201000000001",
			Email:       "hotel01@mock.com",
			Stars:       5,
		},
		{
			Name:        "Mock Hotel 02",
			Description: "Modern hotel near downtown",
			Address:     "Tahrir Street 2",
			City:        "Cairo",
			Country:     "Egypt",
			Phone:       "+201000000002",
			Email:       "hotel02@mock.com",
			Stars:       4,
		},
		{
			Name:        "Mock Hotel 03",
			Description: "Comfortable family hotel",
			Address:     "Nasr City Street 3",
			City:        "Cairo",
			Country:     "Egypt",
			Phone:       "+201000000003",
			Email:       "hotel03@mock.com",
			Stars:       4,
		},
		{
			Name:        "Mock Hotel 04",
			Description: "Business hotel",
			Address:     "New Cairo Street 4",
			City:        "Cairo",
			Country:     "Egypt",
			Phone:       "+201000000004",
			Email:       "hotel04@mock.com",
			Stars:       5,
		},
		{
			Name:        "Mock Hotel 05",
			Description: "Affordable city hotel",
			Address:     "Giza Street 5",
			City:        "Giza",
			Country:     "Egypt",
			Phone:       "+201000000005",
			Email:       "hotel05@mock.com",
			Stars:       3,
		},
		{
			Name:        "Mock Hotel 06",
			Description: "Beach style hotel",
			Address:     "Alexandria Corniche 6",
			City:        "Alexandria",
			Country:     "Egypt",
			Phone:       "+201000000006",
			Email:       "hotel06@mock.com",
			Stars:       5,
		},
		{
			Name:        "Mock Hotel 07",
			Description: "Relaxing resort",
			Address:     "Hurghada Road 7",
			City:        "Hurghada",
			Country:     "Egypt",
			Phone:       "+201000000007",
			Email:       "hotel07@mock.com",
			Stars:       5,
		},
		{
			Name:        "Mock Hotel 08",
			Description: "Red Sea resort",
			Address:     "Sharm Road 8",
			City:        "Sharm El Sheikh",
			Country:     "Egypt",
			Phone:       "+201000000008",
			Email:       "hotel08@mock.com",
			Stars:       5,
		},
		{
			Name:        "Mock Hotel 09",
			Description: "Budget friendly hotel",
			Address:     "Luxor Street 9",
			City:        "Luxor",
			Country:     "Egypt",
			Phone:       "+201000000009",
			Email:       "hotel09@mock.com",
			Stars:       3,
		},
		{
			Name:        "Mock Hotel 10",
			Description: "Premium tourist hotel",
			Address:     "Aswan Street 10",
			City:        "Aswan",
			Country:     "Egypt",
			Phone:       "+201000000010",
			Email:       "hotel10@mock.com",
			Stars:       4,
		},
	}

	if err := database.DB.Create(&hotels).Error; err != nil {
		panic(err)
	}

	fmt.Println("10 hotels created.")

	// =========================
	// Seed 50 Rooms
	// =========================
	//
	// Rooms are created after hotels, referencing each hotel's real
	// database-assigned ID (hotels[hotelIndex].ID, populated by the
	// Create() call above) - never a guessed or hardcoded hotel ID.

	var rooms []roomModel.Room

	roomTypes := []string{
		"single",
		"double",
		"twin",
		"deluxe",
		"suite",
	}

	prices := []float64{
		800,
		1200,
		1400,
		2000,
		3000,
	}

	capacities := []int{
		1,
		2,
		2,
		3,
		4,
	}

	for hotelIndex := 0; hotelIndex < 10; hotelIndex++ {

		for roomIndex := 0; roomIndex < 5; roomIndex++ {

			roomNumber := fmt.Sprintf(
				"M-%02d-%02d",
				hotelIndex+1,
				roomIndex+1,
			)

			room := roomModel.Room{
				HotelID:       hotels[hotelIndex].ID,
				RoomNumber:    roomNumber,
				Type:          roomTypes[roomIndex],
				Description:   fmt.Sprintf("Mock room %s", roomNumber),
				PricePerNight: prices[roomIndex],
				Capacity:      capacities[roomIndex],
				Status:        "available",
			}

			rooms = append(rooms, room)
		}
	}

	if err := database.DB.Create(&rooms).Error; err != nil {
		panic(err)
	}

	fmt.Println("50 rooms created.")

	// =========================
	// Seed 50 Bookings
	// =========================
	//
	// Each booking references a real seeded guest's actual database ID
	// (guestIDs, built from users[].ID above) via round-robin, instead of
	// a hand-written UUID that doesn't correspond to any user row.

	var bookings []model.Booking

	baseDate := time.Date(
		2026,
		9,
		10,
		14,
		0,
		0,
		0,
		time.Local,
	)

	statuses := []string{
		"pending",
		"confirmed",
		"cancelled",
		"confirmed",
		"pending",
	}

	for i := 0; i < 50; i++ {

		// Each room gets one mock booking.
		room := rooms[i]

		checkIn := baseDate.AddDate(0, 0, i%10)

		// Different booking lengths.
		nights := (i % 4) + 2

		checkOut := checkIn.AddDate(0, 0, nights)

		userID := guestIDs[i%len(guestIDs)]

		totalPrice := room.PricePerNight * float64(nights)

		booking := model.Booking{
			UserID:     userID,
			RoomID:     room.ID,
			CheckIn:    checkIn,
			CheckOut:   checkOut,
			Guests:     (i % room.Capacity) + 1,
			TotalPrice: totalPrice,
			Status:     statuses[i%len(statuses)],
		}

		bookings = append(bookings, booking)
	}

	if err := database.DB.Create(&bookings).Error; err != nil {
		panic(err)
	}

	fmt.Println("50 bookings created, each referencing a real seeded guest.")

	// =========================
	// Done
	// =========================
	//
	// Payments, reviews, and hotel/room images are not seeded by this
	// script (they weren't before this fix either) - they're created
	// through the real application flows instead. Nothing here invents
	// seed coverage for them.

	fmt.Println()
	fmt.Println("====================================")
	fmt.Println("Mock data seeded successfully!")
	fmt.Println("====================================")
	fmt.Println("Users   :", len(users), "(1 admin, 1 manager,", guestCount, "guests)")
	fmt.Println("Hotels  : 10")
	fmt.Println("Rooms   : 50")
	fmt.Println("Bookings: 50")
	fmt.Println()
	fmt.Println("Seed login password (all seeded users):", SeedPassword)
	fmt.Println("Example: admin@seed.nilestay.local /", SeedPassword)
}
