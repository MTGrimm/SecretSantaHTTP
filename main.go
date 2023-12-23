package main

import (
    "fmt"
    "encoding/json"
    "io"
    "math/rand"
    "os"
    "net/smtp"
)

type Person struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Want  string `json:"want"`
}

func assignGroups(people []Person, res [][]Person) [][]Person {
    if len(people) <= 4 {
        return append(res, people)
    }

    temp := make([]Person, 0)
    next := make([]Person, 0)

    usedIndexes := make(map[int]int)
    count := 0
    for count < 3 {
        randIndex := rand.Intn(len(people)-count)
        if usedIndexes[randIndex] == 0 {
            temp = append(temp, people[randIndex])
            usedIndexes[randIndex] = 1
            count++
        }
    }

    for i, person := range people {
        if usedIndexes[i] == 0 {
            next = append(next, person)
        }
    }

    return assignGroups(next, append(res, temp))
}

type Assignment struct {
    Giver    *Person
    Reciever *Person
}

func assignPeople(groups [][]Person) []Assignment {
    assigned := make(map[*Person]bool)
    assignments := make([]Assignment, 0)

    for groupIndex := range groups { 
        for giverIndex := range groups[groupIndex] {
            failtedAttempts := 0
            for failtedAttempts<100 {
                randGroup := rand.Intn(len(groups)) 
                randMember := rand.Intn(len(groups[randGroup])) 
                reciever := &(groups[randGroup][randMember])
                giver := &(groups[groupIndex][giverIndex])
                if groupIndex != randGroup && reciever != giver && !assigned[reciever] {
                    newAssignment := Assignment{giver, reciever}
                    assignments = append(assignments, newAssignment)
                    assigned[reciever] = true 
                    break
                }
                failtedAttempts++
            }
            if failtedAttempts == 100 {
                return assignPeople(groups)
            }
        }
    }
    return assignments
}

func sendMessage(assignment Assignment, groupIndex int, groups [][]Person) {
    auth := smtp.PlainAuth(
        "",
        "aryan.timilsina195@gmail.com",
        "sracqawdqrrwnooc",
        "smtp.gmail.com",
    )

    msg := "Hi " + assignment.Giver.Name + 
           "\nYour person is " + assignment.Reciever.Name +
           "\nAnd they want " + assignment.Reciever.Want + 
           "\nAnd finally your group is: "

   for _, member := range groups[groupIndex] {
       msg += "\n" + member.Name
   }

    err := smtp.SendMail(
        "smtp.gmail.com:587",
        auth,
        "aryan.timilsina195@gmail.com",
        []string{"aryan.timilsina195@gmail.com"},
        []byte(msg),
    )

    if err != nil {
        fmt.Println(err)
    }
}

func main() {
    info, err := os.Open("assests/info.json")
    if err != nil {
        fmt.Println(err)
    }

    defer info.Close()
    byteValue, _ := io.ReadAll(info)
    
    var people []Person 
    json.Unmarshal(byteValue, &people)
    
    groups := assignGroups(people, make([][]Person, 0))
    groupMap := make(map[*Person]int)
    for i:=0; i<len(groups); i++ {
        fmt.Printf("\nGroup%d\n", i+1)
        for j:=0; j<len(groups[i]); j++ {
            fmt.Println(groups[i][j])
            groupMap[&(groups[i][j])] = i
        }
    }

    assignments := assignPeople(groups)
    fmt.Println("\nAssignments:")
    for i:=0; i<len(assignments); i++ {
        fmt.Println(assignments[i].Giver.Name, assignments[i].Reciever.Name)
    }

    // for _, assignment := range assignments {
    //    sendMessage(assignment, groupMap[assignment.Giver], groups)
    // }
}

