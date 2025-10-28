package db

import (
    "database/sql"
    "errors"
    "time"
)

type Guest struct {
    ID                  int        `json:"id"`
    Name                string     `json:"name"`
    Email               string     `json:"email"`
    Phone               string     `json:"phone"`
    Status              string     `json:"status"`
    Notes               *string    `json:"notes"`                   // NEW
    EventID             int        `json:"event_id"`                // NEW
    RsvpDate            *time.Time `json:"rsvp_date"`               // NEW
    PlusOnes            int        `json:"plus_ones"`               // NEW
    DietaryRestrictions *string    `json:"dietary_restrictions"`    // NEW
    CreatedAt           time.Time  `json:"created_at"`
}

type Event struct {
    ID          int        `json:"id"`
    Title       string     `json:"title"`
    Description *string    `json:"description"`
    EventDate   *time.Time `json:"event_date"`
    Location    *string    `json:"location"`
    CreatedAt   time.Time  `json:"created_at"`
}
type GuestStats struct {
    Total     int `json:"total"`
    Attending int `json:"attending"`
    Pending   int `json:"pending"`
    Declined  int `json:"declined"`
}


// CheckEmailExists checks if email already exists
func CheckEmailExists(db *sql.DB, email string) (bool, error) {
    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM guests WHERE email = $1)", email).Scan(&exists)
    return exists, err
}


// GetAllGuests retrieves all guests from the database
// GetAllGuests retrieves all guests from the database
func GetAllGuests(db *sql.DB, status string) ([]Guest, error) {
    var rows *sql.Rows
    var err error

    if status != "" {
        rows, err = db.Query(`
            SELECT id, name, email, phone, status, notes, event_id, 
                   rsvp_date, plus_ones, dietary_restrictions, created_at 
            FROM guests 
            WHERE status = $1 
            ORDER BY created_at ASC
        `, status)
    } else {
        rows, err = db.Query(`
            SELECT id, name, email, phone, status, notes, event_id, 
                   rsvp_date, plus_ones, dietary_restrictions, created_at 
            FROM guests 
            ORDER BY created_at ASC
        `)
    }

    if err != nil {
        // ADD THIS LOG
        println("Query error:", err.Error())
        return nil, err
    }
    defer rows.Close()

    var guests []Guest
    for rows.Next() {
        var g Guest
        err := rows.Scan(
            &g.ID, &g.Name, &g.Email, &g.Phone, &g.Status, 
            &g.Notes, &g.EventID, &g.RsvpDate, &g.PlusOnes, 
            &g.DietaryRestrictions, &g.CreatedAt,
        )
        if err != nil {
            // ADD THIS LOG
            println("Scan error:", err.Error())
            return nil, err
        }
        guests = append(guests, g)
    }

    return guests, rows.Err()
}
//count
func GetGuestStats(db *sql.DB) (*GuestStats, error) {
    var stats GuestStats
    err := db.QueryRow(`
        SELECT 
            COALESCE(COUNT(*), 0) as total,
            COALESCE(SUM(CASE WHEN status = 'attending' THEN 1 ELSE 0 END), 0) as attending,
            COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0) as pending,
            COALESCE(SUM(CASE WHEN status = 'declined' THEN 1 ELSE 0 END), 0) as declined
        FROM guests
    `).Scan(&stats.Total, &stats.Attending, &stats.Pending, &stats.Declined)
    
    return &stats, err
}

// GetGuestByID retrieves a single guest by ID
func GetGuestByID(db *sql.DB, id int) (*Guest, error) {
    var g Guest
    err := db.QueryRow(`
        SELECT id, name, email, phone, status, notes, event_id, 
               rsvp_date, plus_ones, dietary_restrictions, created_at 
        FROM guests 
        WHERE id = $1
    `, id).Scan(
        &g.ID, &g.Name, &g.Email, &g.Phone, &g.Status,
        &g.Notes, &g.EventID, &g.RsvpDate, &g.PlusOnes,
        &g.DietaryRestrictions, &g.CreatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, errors.New("guest not found")
    }
    if err != nil {
        return nil, err
    }

    return &g, nil
}

// CreateGuest inserts a new guest into the database
func CreateGuest(db *sql.DB, name, email, phone, status string, notes *string, eventID int, 
                 rsvpDate *time.Time, plusOnes int, dietaryRestrictions *string) (*Guest, error) {
    var g Guest
    err := db.QueryRow(`
        INSERT INTO guests (name, email, phone, status, notes, event_id, 
                           rsvp_date, plus_ones, dietary_restrictions, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
        RETURNING id, name, email, phone, status, notes, event_id, 
                  rsvp_date, plus_ones, dietary_restrictions, created_at
    `, name, email, phone, status, notes, eventID, rsvpDate, plusOnes, dietaryRestrictions).Scan(
        &g.ID, &g.Name, &g.Email, &g.Phone, &g.Status,
        &g.Notes, &g.EventID, &g.RsvpDate, &g.PlusOnes,
        &g.DietaryRestrictions, &g.CreatedAt,
    )

    if err != nil {
        return nil, err
    }

    return &g, nil
}

// DeleteGuest removes a guest from the database
func DeleteGuest(db *sql.DB, id int) error {
    result, err := db.Exec("DELETE FROM guests WHERE id = $1", id)
    if err != nil {
        return err
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rows == 0 {
        return errors.New("guest not found")
    }

    return nil
}
