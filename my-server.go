package main

import (
    "database/sql"
    "encoding/json"
	"fmt"
	"time"
    "log"
    "net/http"
    // "errors"
    "os"
	
    _ "github.com/mattn/go-sqlite3"
    "github.com/dgrijalva/jwt-go"
)

var mySignKey = []byte("secureSecretText")

type User struct {
    Username string `json: "username"`
    Password string `json: "password"`
}

func getTodoHandler(w http.ResponseWriter, r *http.Request, user_id int) {
    type Todo struct {
        ActionName string
        ExpectedDuration string
        IsWork string 
        Date string 
        RealDuration int 
        Emotion string 
        IsDone string
        Id int
    }

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

    fmt.Println("here getTodo from", user_id)
    
    db, err := sql.Open("sqlite3", "todo")
    checkErr(err)

    rows, err := db.Query("SELECT id, actionName, expectedDuration, isWork, date, realDuration, emotion, isDone, frontTodoID FROM todos where clientID=?", user_id)
    checkErr(err)

    var todo []*Todo
    for rows.Next() {
        p := new(Todo)
        rows.Scan(&p.Id, &p.ActionName, &p.ExpectedDuration, &p.IsWork, &p.Date, &p.RealDuration, &p.Emotion, &p.IsDone, &p.Id)
        todo = append(todo, p)
    } 
    if err := json.NewEncoder(w).Encode(todo); err != nil {
            log.Println(err)
    }
    
    rows.Close()
}

func getUsernameHandler (w http.ResponseWriter, r *http.Request, user_id int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

    fmt.Println("here getUserName from", user_id)
    

    db, err := sql.Open("sqlite3", "todo")
    checkErr(err)

    userName, err := db.Query("SELECT userName FROM users WHERE id =?", user_id)
    checkErr(err)

    var dbUserName string

    for userName.Next() {
        userName.Scan(&dbUserName)
    }
     
    if err := json.NewEncoder(w).Encode(dbUserName); err != nil {
        log.Println(err)
    }
    db.Close()
}

func createHandler(w http.ResponseWriter, r *http.Request, user_id int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
    if r.Method == http.MethodPost {
 
        fmt.Println("here POST")
        fmt.Println(user_id)
        body := map[string]interface{}{}
        json.NewDecoder(r.Body).Decode(&body)
 
        db, err := sql.Open("sqlite3", "todo")
        checkErr(err)
 
        stmt, err := db.Prepare("INSERT INTO todos (actionName, expectedDuration, isWork, date, realDuration, emotion, isDone, clientID, frontTodoID) values(?,?,?,?,?,?,?,?,?)")
        checkErr(err)
        res, err := stmt.Exec(body["actionName"], body["expectedDuration"], body["isWork"], body["date"], body["realDuration"], body["emotion"], body["isDone"], user_id, body["id"])
 
        checkErr(err)
        fmt.Println(res.RowsAffected())
        
        id, err := db.Query("SELECT id FROM todos WHERE frontTodoID =?", body["id"])
        checkErr(err)
    
        var dbid int
    
        for id.Next() {
            id.Scan(&dbid)
        }

        if err := json.NewEncoder(w).Encode(dbid); err != nil {
            log.Println(err)
        }

        db.Close()
    }
}

func updateHandler(w http.ResponseWriter, r *http.Request, user_id int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
    if r.Method == http.MethodPost {
        fmt.Println("here UPDATE")
        body := map[string]interface{}{}
        json.NewDecoder(r.Body).Decode(&body)
        fmt.Println(body)
        db, err := sql.Open("sqlite3", "todo")
        checkErr(err)

        stmt, err := db.Prepare("UPDATE todos set actionName=?, expectedDuration=?, isWork=?, date=?, realDuration=?, emotion=?, isDone=? WHERE clientID=? AND id=?")
        checkErr(err)

        res, err := stmt.Exec(body["actionName"], body["expectedDuration"], body["isWork"], body["date"], body["realDuration"], body["emotion"], body["isDone"], user_id, body["id"])
        checkErr(err)
    
        fmt.Println(res.RowsAffected())
       
        db.Close()
    }
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    if r.Method == http.MethodPost {
    fmt.Println("Registration request...")
    
    var u User
    
    json.NewDecoder(r.Body).Decode(&u)
    fmt.Println(u)
  
    if err := json.NewEncoder(w).Encode(addRegistrationCard(u)); err != nil {
        log.Println(err)
    }
    }   
}

func addRegistrationCard(u User) string {
    
 
    fmt.Println("here POST")

    db, err := sql.Open("sqlite3", "todo")
    checkErr(err)
     
    stmt, err := db.Prepare("INSERT INTO users (userName, password) values(?,?)")
    checkErr(err)
    res, err := stmt.Exec(u.Username, u.Password)
    if err != nil {
        return "failed";
    }
    // checkErr(err)

    fmt.Println(res.RowsAffected())

    id, err := db.Query("SELECT id FROM users WHERE userName = '" + u.Username + "'")
    checkErr(err)

    var dbid int

    for id.Next() {
        id.Scan(&dbid)
    }
     
    db.Close()
        
    validToken, err := GenerateGWT(dbid)
    fmt.Println(validToken)

    checkErr(err)

    return validToken
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    if r.Method == http.MethodPost {
    fmt.Println("Login request...")

    var u User
    
    json.NewDecoder(r.Body).Decode(&u)
  
    db, err := sql.Open("sqlite3", "todo")
    checkErr(err)
    fmt.Println("Connected to DB...")

    res, err := db.Query("SELECT id, userName, password FROM users WHERE userName = '" + u.Username + "'")
    if err != nil {
        json.NewEncoder(w).Encode("failed");
    }

    dbusercard := User{}
    var dbid int

    for res.Next() {
        res.Scan(&dbid, &dbusercard.Username, &dbusercard.Password)
    }

    fmt.Println(&dbusercard.Username, &dbusercard.Password, "ghhh")
    res.Close()

    if dbusercard.Username == u.Username && dbusercard.Password == u.Password {
        fmt.Println("working")
        validToken, err := GenerateGWT(dbid)
        fmt.Println(validToken)
        checkErr(err)
        if err := json.NewEncoder(w).Encode(validToken); err != nil {
            log.Println(err)
        }
    } else {
        fmt.Println("not correct")
        // log.Fatal(err)
        if err := json.NewEncoder(w).Encode("failed"); err != nil {
            log.Println(err)
        }       
    }   
}
}


// func checkLogin(u User) (string, error) {

//     db, err := sql.Open("sqlite3", "todo")
//     checkErr(err)

//     res, err := db.Query("SELECT id, userName, password FROM users WHERE userName = '" + u.Username + "'")
//     checkErr(err)

//     dbusercard := User{}
//     var id int

//     for res.Next() {
//         res.Scan(&id, &dbusercard.Username, &dbusercard.Password)
//     }

//     res.Close()

//     if dbusercard.Username == u.Username || dbusercard.Password == u.Password {
//         validToken, err := GenerateGWT(id)
//         fmt.Println(validToken)
//         checkErr(err)
//         return validToken, nil
//     } else {
//         fmt.Println("Not correct!")
//         return "", errors.New("not valid login and (or) password")
//     }
// }

func GenerateGWT(id int) (string, error){
    var err error
   
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["exp"] = time.Now().Add(time.Hour * 1000).Unix()
    claims["id"] = id
    fmt.Println(claims)
    
    tokenString, err := token.SignedString(mySignKey)
    if err != nil {
        return "Ошибка", err
    }

    return tokenString, nil
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request, int)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			         
            claims := jwt.MapClaims{}
            token, err := jwt.ParseWithClaims(r.Header["Token"][0], &claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySignKey, nil
			})


			if err != nil {
				w.WriteHeader(http.StatusForbidden)
                fmt.Println("Ошибка!")
			}

            fmt.Println(claims["id"])
            var userId int
            userId = int(claims["id"].(float64))

			if token.Valid {
                fmt.Println("Доступ есть!")
                fmt.Println(userId)
				endpoint(w, r, userId)
			}

		} else {
			fmt.Fprintf(w, "No Authorization Token provided")
		}
	})
}

func removeTodoHandler(w http.ResponseWriter, r *http.Request, user_id int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
    if r.Method == http.MethodPost {
        fmt.Println("here REMOVE todo from", user_id)
        body := map[string]interface{}{}
        json.NewDecoder(r.Body).Decode(&body)
        fmt.Println(body["id"])
        db, err := sql.Open("sqlite3", "todo")
        checkErr(err)
        stmt, err := db.Prepare("DELETE FROM todos where clientID=? AND id=?")
        checkErr(err)
        res, err := stmt.Exec(user_id, body["id"])
        checkErr(err)
        fmt.Println(res.RowsAffected())
        db.Close()
    }
}

func deleteDatedTodosHandler(w http.ResponseWriter, r *http.Request, user_id int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
    if r.Method == http.MethodPost {
        fmt.Println("here DELETE dated todos from", user_id)
        body := map[string]interface{}{}
        json.NewDecoder(r.Body).Decode(&body)
        fmt.Println(body["date"])
        db, err := sql.Open("sqlite3", "todo")
        checkErr(err)
        stmt, err := db.Prepare("DELETE FROM todos where clientID=? AND date=?")
        checkErr(err)
        res, err := stmt.Exec(user_id, body["date"])
        checkErr(err)
        fmt.Println(res.RowsAffected())
        db.Close()
    }
}

func deleteAllTodosHandler(w http.ResponseWriter, r *http.Request, user_id int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
    if r.Method == http.MethodPost {
        fmt.Println("here DELETE all todos from", user_id)
        db, err := sql.Open("sqlite3", "todo")
        checkErr(err)
        stmt, err := db.Prepare("DELETE FROM todos where clientID=?")
        checkErr(err)
        res, err := stmt.Exec(user_id)
        checkErr(err)
        fmt.Println(res.RowsAffected())
        db.Close()
    }
}

func main() {
    fmt.Println("Работает")
    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/login", loginHandler)
    http.Handle("/get_todo", isAuthorized(getTodoHandler))
    http.Handle("/get_username", isAuthorized(getUsernameHandler))
    http.Handle("/create_todo", isAuthorized(createHandler))
    http.Handle("/update_todo", isAuthorized(updateHandler))
    http.Handle("/remove_todo", isAuthorized(removeTodoHandler))
    http.Handle("/delete_all_todos", isAuthorized(deleteAllTodosHandler))
    http.Handle("/delete_dated_todos", isAuthorized(deleteDatedTodosHandler))
    http.Handle("/", http.FileServer(http.Dir("./build")))
    port := os.Getenv("PORT")
    if (port == "") {
        port = "8080"
    }
    log.Fatal(http.ListenAndServe(":" + port, nil))
      
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}