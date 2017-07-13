package statistics

import "time"
import "fmt"

var AllTimeGenerator = &Generator{
	Duration: time.Second,
	Generator: func(s *Statistics) error {
		fmt.Println("called all time gen")
		_, err := s.DB.Exec("SELECT count(*) FROM songs")
		if err != nil {
			return err
		}

		return nil
	},
}

var AllTimeGenerator2 = &Generator{
	Duration: time.Second,
	Generator: func(s *Statistics) error {
		fmt.Println("called all time gen 2")
		_, err := s.DB.Exec("SELECT count(*) FROM songs")
		if err != nil {
			return err
		}

		return nil
	},
}
