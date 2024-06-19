package Db

import (
	"birthdays/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Storage struct {
	DatabaseURL string
	db          *sql.DB
	//Ua          *UserActions
}

func NewStorage(url string) *Storage {
	return &Storage{
		DatabaseURL: url,
	}

}

func (s *Storage) Open() error {

	db, err := sql.Open("postgres", s.DatabaseURL)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return err
	}

	s.db = db

	return nil
}
func (s *Storage) Close() {
	s.db.Close()
}
func (s *Storage) CreateUser(u *models.User) error {
	query := `INSERT INTO users (email, passwordHash,nme, birthday)VALUES ($1, $2,$3,$4)RETURNING id`

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	u.PasswordHash = string(hashedPassword)
	args := []interface{}{u.Email, u.PasswordHash, u.Name, u.Date}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&u.Id)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:

			return err
		default:
			return err
		}
	}

	return nil
}
func (s *Storage) FindById(id int) (*models.User, error) {
	if id < 1 {
		return nil, errors.New("negative Id")
	}
	query := `SELECT id,email,nme,birthday::date FROM users WHERE id = $1`
	u := models.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, id).Scan(&u.Id, &u.Email, &u.Name, &u.Date)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("No Record")
		default:
			return nil, err

		}
	}
	return &u, nil
}
func (s *Storage) FindByEmail(email string) (*models.User, error) {

	query := `SELECT id,email,passwordhash FROM users WHERE email = $1`
	u := models.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, email).Scan(&u.Id, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (s *Storage) SubscribeTo(u *models.User, email string) error {

	checkQuery := "SELECT COUNT(*) FROM users WHERE email = $1"
	var count int
	err := s.db.QueryRow(checkQuery, email).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("User does Not Exists")
	}
	query := `INSERT INTO subscriptions (subscriber_id, subscribed_to_id)VALUES($1, (SELECT id FROM users WHERE email = $2 ));`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{u.Id, email}
	err = s.db.QueryRowContext(ctx, query, args...).Err()
	if err != nil {

		return err
	}
	return nil
}
func (s *Storage) GetUserSubscriptions(u *models.User) (error, *[]string) {

	query := `SELECT email FROM subscriptions JOIN users ON subscriptions.subscribed_to_id = users.id WHERE subscriptions.subscriber_id = $1;`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, u.Id)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			panic(err)
		}
		emails = append(emails, email)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	//data := make(map[string][]int)
	//data["subscriptions"] = ids

	return nil, &emails
}

// calendar///////////////
func (s *Storage) GetLasCalRecord() time.Time {

	query := `SELECT event_timestamp FROM events  WHERE id = $1;`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var ts time.Time
	row := s.db.QueryRowContext(ctx, query, 1)

	err := row.Scan(&ts)
	if err != nil {
		ts = time.Time{}
		return ts
	}
	return ts

}
func (s *Storage) SetLasCalRecord() error {

	query := ` UPDATE events SET event_timestamp=current_timestamp WHERE id = $1;`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.QueryContext(ctx, query, 1)

	if err != nil {
		return err
	}
	return nil

}
func (s *Storage) SubscriptionsRows() []*models.CalendarJob {

	query := ` SELECT  
     u1.id AS subscriber_id, 
        u1.email AS subscriber_email, 
        u2.nme AS subscribed_to_mane,
        u2.birthday AS subscribed_to_date FROM subscriptions s 
            JOIN users u1 ON s.subscriber_id = u1.id 
            JOIN users u2 ON s.subscribed_to_id = u2.id;`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {

	}
	defer rows.Close()

	var jobs []*models.CalendarJob

	for rows.Next() {

		var job models.CalendarJob
		err := rows.Scan(&job.Id, &job.SubEmail, &job.SubName, &job.Date)
		job.Date = job.Date.UTC()
		job.Text = fmt.Sprintf("your friend %s has birthday today", job.SubName)
		job.Bth = false
		jobs = append(jobs, &job)

		if err != nil {
			panic(err)
		}
	}
	return jobs
}
func (s *Storage) UpdateIds(uid int, ids []int) error {
	query := `UPDATE users SET calendar_ids= $2 WHERE id = $1 `
	args := []any{uid, pq.Array(ids)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s.db.QueryRowContext(ctx, query, args...)

	return nil
}
func (s *Storage) FindJobsById(id int) ([]string, error) {
	if id < 1 {
		return nil, errors.New("negative Id")
	}
	query := `SELECT calendar_ids FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var ids []string
	err := s.db.QueryRowContext(ctx, query, id).Scan(pq.Array(&ids))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("No Record")
		default:
			return nil, err

		}
	}

	if err != nil {
		panic(err)
	}
	return ids, nil
}

/////////////////misc/////////

type testErr struct {
}

func (e *testErr) Error() string {
	return "boom"

}

/*



func (d *DBinUse) User() *UserActions {
	if d.Ua != nil {
		return d.Ua
	}
	d.Ua = &UserActions{
		database: d,
	}
	return d.Ua
}*/
