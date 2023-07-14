package main

import (
   "fmt"
   "database/sql"
   "os"
   "bufio"
   "strings"
   "errors"
   "strconv"
   _ "github.com/mattn/go-sqlite3"
)

type post struct {
   id int
   post_name string
   chat_id int
   message_1 int
   message_2 int 
}

func String(p post) string {
   res := strconv.Itoa(p.id) + "\t" + p.post_name + "\t" + strconv.Itoa(p.chat_id) + "\t" + strconv.Itoa(p.message_1) + "\t" + strconv.Itoa(p.message_2);
   return res;
}

func create_db() {
   if _, err := os.Stat("data.db"); errors.Is(err, os.ErrNotExist) {
      os.Create("data.db")
   }
   db, err := sql.Open("sqlite3", "data.db")
   defer db.Close()
   if err != nil {
      fmt.Println(err)
      os.Exit(1)
   }

   // ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
   // defer cancelfunc()
   // next - ExecContext (if db not local)
   _, err = db.Exec("Create TABLE IF NOT EXISTS data ( id INTEGER PRIMARY KEY AUTOINCREMENT, post_name TEXT, chst_id INTEGER, message_1 INTEGER, message_2 INTEGER )")
   if err != nil {
      fmt.Println(err)
      os.Exit(1)
   }

}

func main() {
   fmt.Println("Program started")
   fmt.Println("Use SHOW to see all db")
   fmt.Println("use ADD to insert")
   fmt.Println("use DEL to delete")
   fmt.Println("Use EXIT to exit the program")
   create_db()
   db, err := sql.Open("sqlite3", "data.db")
   defer db.Close()
   if err != nil {
      fmt.Println(err)
      os.Exit(1)
   }
   reader := bufio.NewReader(os.Stdin)
   for {
      command, _ := reader.ReadString('\n')
      command = strings.TrimSuffix(command, "\n")
      command = strings.ToLower(command)
      switch command {
      case "exit":
         fmt.Println("Goodbye")
         os.Exit(0)
      case "add":
         fmt.Println("type post_name, chat_id, message_1, message_2")
         var v1, v2, v3 int
         var s string
         fmt.Scan(&s, &v1, &v2, &v3)
         _, err := db.Exec("INSERT INTO data (post_name, chst_id, message_1, message_2) VALUES (?, ?, ?, ?)", s, v1, v2, v3)
         if (err != nil) {
            fmt.Println(err)
            os.Exit(1)
         }
         fmt.Println("Done!")
      case "del":
         fmt.Println("type record id")
         var val int
         fmt.Scan(&val)
         _, err := db.Exec(`DELETE FROM data WHERE id = ?`, val)
         if (err != nil) {
            fmt.Println(err)
            os.Exit(1)
         }
         fmt.Println("Done!")
      case "show":
         res, err := db.Query("SELECT * FROM data")
         if err != nil {
            fmt.Println(err)
            os.Exit(1)
         }
         var info []post
         for res.Next() {
            var wr post
            _ = res.Scan(&wr.id, &wr.post_name, &wr.chat_id, &wr.message_1, &wr.message_2)
            info = append(info, wr)
         }
         fmt.Println("There are", len(info), "records in the table")
         for i:=0;i<len(info);i++ {
            fmt.Println(String(info[i]))
         }
      default:
         fmt.Println("Unknown command")
      }
   }
}