package commands

import (
	"fmt"
	"time"
	"context"
	"strconv"
	"database/sql"
	"github.com/google/uuid"
	"github.com/bitofaphilistine/gator/internal/rss"
	"github.com/bitofaphilistine/gator/internal/config"
	"github.com/bitofaphilistine/gator/internal/database"
)


type State struct {
	Config *config.Config
	Database *sql.DB
	Queries *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct{
	Registered map[string]func(state *State, c Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	if handler, ok := c.Registered[cmd.Name]; ok {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.Name)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Registered[name] = f
}


func HandlerResetDb(s *State, cmd Command) error {
	targetCases := map[string][]func(context.Context) error {
		"all": {s.Queries.ResetUsers, s.Queries.ResetFeeds, s.Queries.ResetFeedFollows, s.Queries.ResetPosts},
		"users": {s.Queries.ResetUsers, s.Queries.ResetFeedFollows},
		"feeds": {s.Queries.ResetFeeds, s.Queries.ResetFeedFollows, s.Queries.ResetPosts},
		"feedfollows": {s.Queries.ResetFeedFollows},
		"posts": {s.Queries.ResetPosts},
	}
	target := "all"

	if cmd.Args != nil && len(cmd.Args) > 0 {
		target = cmd.Args[0]
	}

	if _, ok := targetCases[target]; !ok {
		return fmt.Errorf("invalid reset target (all, users, feeds, feedfollows, posts): ", target)
	}

	for _, reset := range targetCases[target] {
		if err := reset(context.Background()); err != nil {
			return fmt.Errorf("failed to run", reset, err)
		}
	}

	fmt.Println(target, "table(s) reset successfully")
	return nil
}

func HandlerLogin(s *State, cmd Command) error {
	if cmd.Args == nil || len(cmd.Args) < 1 {
		return fmt.Errorf("login command requires a username argument")
	}

	user, err := s.Queries.GetUserByName(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	err = s.Config.SetUser(user.Username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Println("User set to:", user.Username)

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if cmd.Args == nil || len(cmd.Args) < 1 {
		return fmt.Errorf("register command requires a username argument")
	}

	userParams := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		Username: cmd.Args[0],
	}
	user, err := s.Queries.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	fmt.Println("User created:", user.Username)

	err = s.Config.SetUser(user.Username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	return nil
}

func HandlerListUsers(s *State, cmd Command) error {
	users, err := s.Queries.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	fmt.Println("Registered Users:")
	for _, user := range users {
		if user.Username == s.Config.CurrentUserName {
			fmt.Println(" *", user.Username, "(current)")
		} else {
			fmt.Println(" *", user.Username)
		}
	}

	return nil
}

func HandlerAggregate(s *State, cmd Command) error {
	if cmd.Args == nil || len(cmd.Args) < 1 {
		return fmt.Errorf("agg command requires an interval argument")
	}

	interval, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to parse interval: %w", err)
	}

	fmt.Println("Collecting feeds every", cmd.Args[0])

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Println("Error scraping feeds:", err)
		}

		<-ticker.C
	}

	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if cmd.Args == nil || len(cmd.Args) < 2 {
		return fmt.Errorf("addfeed command requires a name and url argument")
	}

	feed, err := s.Queries.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		Name: cmd.Args[0],
		Url: cmd.Args[1],
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}
	
	fmt.Println("Created feed:", feed.Name, "with URL:", feed.Url, "for user:", user.Username)

	_, err = s.Queries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed: %w", err)
	}

	return nil
}

func HandlerListFeeds(s *State, cmd Command) error {
	feeds, err := s.Queries.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list feeds: %w", err)
	}

	fmt.Println("Registered Feeds:")
	usernames := make(map[uuid.UUID]string)
	for _, feed := range feeds {
		username, cached := usernames[feed.UserID]
		if !cached {
			user, err := s.Queries.GetUserById(context.Background(), feed.UserID)
			if err != nil {
				return fmt.Errorf("failed to get user for feed: %w", err)
			}
			username = user.Username
			usernames[feed.UserID] = username
		}

		fmt.Printf(" * %s (URL: %s, User: %s)\n", feed.Name, feed.Url, username)
	}

	return nil
}

func HandlerFollowFeed(s *State, cmd Command, user database.User) error {
	if cmd.Args == nil || len(cmd.Args) < 1 {
		return fmt.Errorf("follow command requires a feed url argument")
	}

	feed, err := s.Queries.GetFeedByUrlOrName(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to get feed by url: %w", err)
	}

	follow, err := s.Queries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed: %w", err)
	}

	fmt.Printf("User %s followed feed: %s", follow.UserName, follow.FeedName)
	return nil
}

func HandlerUnfollowFeed(s *State, cmd Command, user database.User) error {
	if cmd.Args == nil || len(cmd.Args) < 1 {
		return fmt.Errorf("unfollow command requires a feed url argument")
	}

	feed, err := s.Queries.GetFeedByUrlOrName(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to get feed by url: %w", err)
	}

	err = s.Queries.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow feed: %w", err)
	}
	
	fmt.Printf("User %s unfollowed feed: %s", user.Username, feed.Name)
	return nil
}

func HandlerListFollowedFeeds(s *State, cmd Command, user database.User) error {
	follows, err := s.Queries.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get followed feeds: %w", err)
	}

	fmt.Println("Followed Feeds:")
	usernames := make(map[uuid.UUID]string)
	feedNames := make(map[uuid.UUID]string)
	for _, follow := range follows {
		username, cached := usernames[follow.UserID]
		if !cached {
			user, err := s.Queries.GetUserById(context.Background(), follow.UserID)
			if err != nil {
				return fmt.Errorf("failed to get user for follow: %w", err)
			}
			usernames[follow.UserID] = user.Username
			username = user.Username
		}

		feedName, cached := feedNames[follow.FeedID]
		if !cached {
			feed, err := s.Queries.GetFeedById(context.Background(), follow.FeedID)
			if err != nil {
				return fmt.Errorf("failed to get feed for follow: %w", err)
			}
			feedNames[follow.FeedID] = feed.Name
			feedName = feed.Name
		}

		fmt.Printf(" * User: %s, Feed: %s\n", username, feedName)
	}

	return nil
}

func HandlerListPosts(s *State, cmd Command, user database.User) error {
	limit := int32(2)
	if cmd.Args != nil && len(cmd.Args) > 0 {
		arg, err := strconv.ParseUint(cmd.Args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid limit argument: %w", err)
		}
		limit = int32(arg)
	}
	fmt.Printf("Listing the latest %d posts for user: %s\n", limit, user.Username)

	posts, err := s.Queries.GetPostsByUser(context.Background(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit: limit,
	})
	if err != nil {
		return fmt.Errorf("failed to get posts for user: %w", err)
	}
	
	for _, post := range posts {
		fmt.Printf(" * %s (%s)\n", post.Title, post.Url)
	}

	return nil
}


func scrapeFeeds(s *State) error {
	feed, err := s.Queries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %w", err)
	}

	err = s.Queries.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID: feed.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to mark feed as fetched: %w", err)
	}

	res, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	for _, item := range res.Channel.Item {
		var publishedTime time.Time
		var err error
		
		for _, format := range []string{
			time.Layout,
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			time.RFC1123Z,
			time.RFC1123,
			time.RFC822Z,
			time.RFC822,
			time.RFC850,
			time.RFC3339,
			time.RFC3339Nano,
		} {
			publishedTime, err = time.Parse(format, item.PubDate)
			if err == nil {
				break
			}
		}
		if err != nil {
			fmt.Printf("Warning: failed to parse published date for item '%s': %v\n", item.Title, err)
			publishedTime = time.Time{}
		}

		_, err = s.Queries.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			Title: item.Title,
			Url: item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: publishedTime, Valid: true},
			FeedID: feed.ID,
		})
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("failed to create post: %w", err)
		}
	}
	return nil
}