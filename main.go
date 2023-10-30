package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type blogPost struct {
	BlogID  int    `json:"BlogID"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var db *sql.DB

func init() {
	// Open a connection to the MySQL database
	var err error
	db, err = sql.Open("mysql", "root:Xdrcftviji007..@tcp(localhost:3306)/blogpost")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Hello")

	// Create a new router using Gorilla Mux
	router := mux.NewRouter()

	//APi routes
	router.HandleFunc("/blog-posts", createBlogPost).Methods("POST")
	router.HandleFunc("/blog-posts", getBlogPosts).Methods("GET")
	router.HandleFunc("/blog-posts/{id:[0-9]+}", getBlogPostsByID).Methods("GET")
	router.HandleFunc("/blog-posts/{id:[0-9]+}", updateBlogPost).Methods("PUT")
	router.HandleFunc("/blog-posts/{id:[0-9]+}", deleteBlogPost).Methods("DELETE")

	//server

	log.Fatal(http.ListenAndServe(":8080", router))

}

// create
func createBlogPost(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Request to create blogpost---->")
	// Parse the request body to extract the post data
	var post blogPost
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the new post into the database
	_, err = db.Exec("INSERT INTO blogs (Title, Content) VALUES (?, ?)", post.Title, post.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Blog post created successfully",
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(post)
	json.NewEncoder(w).Encode(response)

}

// get
func getBlogPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("enter inside getBlogPosts")

	rows, err := db.Query("SELECT BlogID, Title, Content FROM blogs")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []blogPost
	for rows.Next() {
		var post blogPost
		if err := rows.Scan(&post.BlogID, &post.Title, &post.Content); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// get by ID
func getBlogPostsByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("enter inside getBlogPostsByID-->")
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post blogPost
	err = db.QueryRow("SELECT BlogID, Title, Content FROM blogs WHERE BlogID = ? ", postID).Scan(&post.BlogID, &post.Title, &post.Content)
	if err == sql.ErrNoRows {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)

}

// update
func updateBlogPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("enter inside UpdateBlogPost-->")
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var updatedPost blogPost
	err = json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var existingPost blogPost
	err = db.QueryRow("SELECT  Title, Content FROM blogs WHERE BlogID = ? ", postID).Scan(&existingPost.Title, &existingPost.Content)
	if err == sql.ErrNoRows {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE blogs SET Title = ?, content = ? WHERE BlogID = ?", updatedPost.Title, updatedPost.Content, postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Blog post updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// delete
func deleteBlogPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("enter inside deleteBlogPost--->")
	vars := mux.Vars(r)

	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var existingPost blogPost
	err = db.QueryRow("SELECT Title, Content FROM blogs WHERE BlogID = ? ", postID).Scan(&existingPost.Title, &existingPost.Content)
	if err == sql.ErrNoRows {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM blogs WHERE  BlogID = ? ", postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	respose := map[string]interface{}{
		"message": "Blog post deleted successfully",
	}
	json.NewEncoder(w).Encode(respose)

}
