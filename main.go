package main

import (
		"context"
    "fmt"
    "log"
		//"strings"
		"regexp"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

		"net/http"
		"encoding/json"
)

//users details collection is u_collection
//posts details collection is p_collection
var u_collection *mongo.Collection
var p_collection *mongo.Collection

//User structure
type User struct {
    Id string
    Name  string
    Email string
	Password string
}

//Post structure
type Post struct{
	Id string
	Caption string
	Image_url string
	Timestamp string
	User_id string
}


//function to create new User via JSON POST Method
func createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method,r.URL.Path)
	if r.Method == "POST" {
			var new_user User
			err := json.NewDecoder(r.Body).Decode(&new_user)
			 if err != nil {
					 http.Error(w, err.Error(), http.StatusBadRequest)
					 return
			 }
				saveUser(*u_collection,&new_user)
				fmt.Fprintf(w, "created user")
	}
}
//helping function to save new user details to MongoDB
func saveUser(u mongo.Collection,user *User) {
	insertResult, err := u.InsertOne(context.TODO(), user)
	if err != nil {
	    log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}


//get user details using user id in GET Method
func getUser(w http.ResponseWriter, r *http.Request){
		fmt.Println(r.Method,r.URL.Path)
		if r.Method == "GET" {

				re := regexp.MustCompile("/users/([0-9a-zA-Z]+)")
				user_id:=re.FindStringSubmatch(r.URL.Path)[0][7:]
				fmt.Println(user_id)
				var result User
				filter := bson.D{{"id", user_id}}

				err := u_collection.FindOne(context.TODO(), filter).Decode(&result)
				if err != nil {
    			log.Fatal(err)
				}

    		fmt.Printf("Found a single document: %+v\n", result)
				b, err := json.Marshal(result)
				if err!=nil {
					fmt.Println("error converting struct to json")
				}else {
					fmt.Fprintf(w, string(b))
				}
		}
}


//create new post using JSON POST Method
func createPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method,r.URL.Path)
	if r.Method == "POST" {
			var new_post Post
			err := json.NewDecoder(r.Body).Decode(&new_post)
			 if err != nil {
					 http.Error(w, err.Error(), http.StatusBadRequest)
					 return
			 }
				savePost(*p_collection,&new_post)
				fmt.Fprintf(w, "created post")
	}
}

//helping function to save new post to MongoDB
func savePost(p mongo.Collection,post *Post) {
	insertResult, err := p.InsertOne(context.TODO(), post)
	if err != nil {
	    log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}


//get details of particular post using post id via GET Method
func getPost(w http.ResponseWriter, r *http.Request){
		fmt.Println(r.Method,r.URL.Path)
		if r.Method == "GET" {

				re := regexp.MustCompile("/posts/([0-9a-zA-Z]+)")
				post_id:=re.FindStringSubmatch(r.URL.Path)[0][7:]
				fmt.Println(post_id)
				var result Post
				filter := bson.D{{"id", post_id}}

				err := p_collection.FindOne(context.TODO(), filter).Decode(&result)
				if err != nil {
    			log.Fatal(err)
				}

    		fmt.Printf("Found a single document: %+v\n", result)
				b, err := json.Marshal(result)
				if err!=nil {
					fmt.Println("error converting struct to json")
				}else {
					fmt.Fprintf(w, string(b))
				}
		}
}


//function to return all the posts posted by particular user
func listOfpostOfUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method,r.URL.Path)
	if r.Method == "GET" {

			re := regexp.MustCompile("/posts/users/([0-9a-zA-Z]+)")
			user_id:=re.FindStringSubmatch(r.URL.Path)[0][13:]
			fmt.Println(user_id)

			var results []*Post

			filter := bson.D{{"user_id", user_id}}

			findOptions := options.Find()
			//findOptions.SetLimit(2)

			cur,err := p_collection.Find(context.TODO(),filter,findOptions)
			if err != nil {
				log.Fatal(err)
			}

			for cur.Next(context.TODO()) {

		    // create a value into which the single document can be decoded
		    var elem Post
		    err := cur.Decode(&elem)
		    if err != nil {
		        log.Fatal(err)
		    }
				
		    results = append(results, &elem)
			}

			if err := cur.Err(); err != nil {
			    log.Fatal(err)
			}

			// Close the cursor once finished
			cur.Close(context.TODO())

			fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
			b, err := json.Marshal(results)
			if err!=nil {
				fmt.Println("error converting struct to json")
			}else {
				fmt.Fprintf(w, string(b))
			}
	}
}

//routes handling function
func handleRequests() {
	 http.HandleFunc("/users",createUser)
	 http.HandleFunc("/users/",getUser)
	 http.HandleFunc("/posts",createPost)
	 http.HandleFunc("/posts/",getPost)
	 http.HandleFunc("/posts/users/",listOfpostOfUser)

	 //server is listening at PORT 4000
	 log.Fatal(http.ListenAndServe(":4000", nil))
 }



func main() {
	  //*************Connecting to MongoDB***********//
		clientOptions:= options.Client().ApplyURI("mongodb+srv://authuser:authuser@cluster0.wrelr.mongodb.net/insta-API?retryWrites=true&w=majority")
		client,err:= mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
				log.Fatal(err)
		}

		// Check the connection
		err = client.Ping(context.TODO(), nil)

		if err != nil {
				log.Fatal(err)
		}

		fmt.Println("Connected to MongoDB!")
		//***********************************************//

		//users details collection
		u_collection= client.Database("test").Collection("users")

		//posts details collection
		p_collection=client.Database("test").Collection("posts")

		//Requests
    handleRequests()
}
