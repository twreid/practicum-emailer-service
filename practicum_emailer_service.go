package main

import (
    "fmt"
    "net/http"
    "html/template"
    "encoding/csv"
    "strings"
    "io"
    "os"
)

// Column indexes for the csv file.
const(
    NameColumn = 4
    EmailColumn = 16
    MajorColumn = 15
    MNumberColumn = 5
    TbExpirationColumn = 8
    LiabExpirationColumn = 9
    FcsrExpirationColumn = 10
    CourseColumn = 1
    FbiStartColumn = 11
    FbiEndColumn = 13
)

type student struct {
    Name string
    Email string
    Major string
    MNumber string
    TbExpiration string
    LiabExpiration string
    FcsrExpiration string
    FbiExpiration string
    Courses map[string]bool
}

func (s student) String() string {
    return fmt.Sprintf("[%s, %s, %s]", s.Name, s.Email, s.Major)
}

func fromRecord(r []string, s map[string]*student) {
    major := r[MajorColumn]
    if strings.Contains(major, "Exercise & Mov") {
        return
    }

    course := strings.Split(r[CourseColumn], " ")[0]
    _, ok := s[r[MNumberColumn]]
    if ok {
        student := s[r[MNumberColumn]]
        student.Courses[course] = true
    } else {
        student := &student{
            r[NameColumn],
            r[EmailColumn],
            major,
            r[MNumberColumn],
            r[TbExpirationColumn],
            r[LiabExpirationColumn],
            r[FcsrExpirationColumn],
            strings.Join(r[FbiStartColumn:FbiEndColumn], ","),
            make(map[string]bool),
        }

        student.Courses[course] = true

        s[student.MNumber] = student
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "index", nil)
}

func csvHandler(w http.ResponseWriter, r *http.Request) {
    file, err := os.Open("C:\\Users\\reidt\\Desktop\\12-04-2014.csv")
    if err != nil {
        fmt.Println(err)
    }
    defer file.Close()

    rr := csv.NewReader(file)
    students := make(map[string]*student)

    for {
       record, err := rr.Read()
       if err == io.EOF {
           break
       }
       if err != nil {
           fmt.Println(err)
       }

       fromRecord(record, students)
   }

   renderTemplate(w, "csv", students)
}

func renderTemplate(w http.ResponseWriter, tmpl string, r interface{}) {
    t, _ := template.ParseFiles("templates\\" + tmpl + ".html")
    t.Execute(w, r)
}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/csv", csvHandler)
    http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
    http.ListenAndServe(":8080", nil)
}
