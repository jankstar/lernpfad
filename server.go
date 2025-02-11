package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

const (
	cTime string = "02.01.2006 15:04:05"
	// to get a string like this run:
	// openssl rand -hex 32
	SECRET_KEY string = "8a07c8e7d1b330c50a8bb9cfa532ee69fa1166feb0e6bbbb17d9298980c66354"
	//ALGORITHM                   = "HS256"
	ACCESS_TOKEN_EXPIRE_MINUTES time.Duration = time.Hour
	TENANT                      string        = "data"
)

var (
	goServer      *gin.Engine
	goEnforcer    *casbin.Enforcer
	identityKey   = "id" //Feldname aus dem Payload des Token zur Identifzierung
	LOCAL_TEMPDIR = "lernpfad"
	PORT          = "8081"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Tenant   string `json:"tenant"`
	password string //hash
	Role     string `json:"role"` // Role of the user; cannot be defined except for admin
	Kurse    []Kurs `json:"kurse"`
}

func GetHash(iPassword string) string {

	h := fnv.New32a()
	h.Write([]byte(iPassword))
	return strconv.FormatUint(uint64(h.Sum32()), 10)

}

func (me *User) PasswordIsValide(iPassword string) bool {

	type DbUser struct {
		ID        string   `json:"id"`             // unique id ( username )
		Name      string   `json:"name"`           // username given be the user
		Password  string   `json:"password"`       // hashed password specified by the user
		Role      string   `json:"role,omitempty"` // Role of the user; cannot be defined except for admin
		Avatar    string   `json:"avatar"`         // Avatar choosen by the user; default: cat1
		LastLogin string   `json:"lastLogin"`      // Datetimestring of lastlogin
		Playtime  int      `json:"playtime"`       // time in Minutes the user can play at same day
		Tips      []string `json:"tips"`           // list of Tips owned by user
		Points    int      `json:"points"`         // Count of Points earned
	}

	if me.Name == "Gast" && iPassword == "Gast" {
		return true
	}

	adminUser := os.Getenv("ADMIN_USER")
	if me.Name == adminUser && iPassword != "" {
		me.Role = "admin"
		return true
	}

	endpoint := os.Getenv("COSMOSDBURI")
	if endpoint == "" {
		return false
	}

	key := os.Getenv("COSMOSDBKEY")
	if key == "" {
		return false
	}

	databaseName := os.Getenv("COSMOSDB")
	if databaseName == "" {
		return false
	}

	cred, _ := azcosmos.NewKeyCredential(key)
	client, _ := azcosmos.NewClientWithKey(endpoint, cred, nil)

	containerClient, _ := client.NewContainer(databaseName, "user")

	pk := azcosmos.NewPartitionKeyString(me.Name)

	ctx := context.TODO()
	itemResponse, _ := containerClient.ReadItem(ctx, pk, me.Name, nil)
	var dbUser DbUser
	json.Unmarshal(itemResponse.Value, &dbUser)

	if dbUser.Password == GetHash(iPassword) {
		me.Role = dbUser.Role
		return true
	}
	return false
}

type Answer struct {
	Answer string `json:"answer"`
	Valid  bool   `json:"valid"`
	Choice bool   `json:"choice"`
}

type Question struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Answers  []Answer `json:"answers"`
	Choice   string   `json:"choice"`
}

type Test struct {
	TestDate  string     `json:"test_date"`
	Points    int        `json:"points"`
	Questions []Question `json:"questions"`
}

type Document struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Token       string `json:"token"`
	ViewDate    string `json:"view_date"`
}

type Task struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	SubTitle    string     `json:"sub_title"`
	Description string     `json:"description"`
	Icon        string     `json:"icon"`
	Movie       string     `json:"movie"`
	Questions   []Question `json:"questions"`
	Documents   []Document `json:"documents"`
	Start       string     `json:"start"`
	Change      string     `json:"change"`
	Check       string     `json:"check"`
	Tests       []Test     `json:"tests"`
}

type Kurs struct {
	ID               string   `json:"id" uri:"id" binding:"required"`
	Grade            string   `json:"grade"`
	Subject          string   `json:"subject"`
	MainTitle        string   `json:"main_title"`
	Publication      string   `json:"publication"`
	SubTitle         string   `json:"sub_title"`
	Image            string   `json:"image"`
	DescriptionShort string   `json:"description_short"`
	Description      string   `json:"description"`
	DescriptionMovie string   `json:"description_movie"`
	Requirements     string   `json:"requirements"`
	Duration         int      `json:"duration"`
	Difficulty       int      `json:"difficulty"`
	Enrollments      int      `json:"enrollments"`
	Rating           float32  `json:"rating"`
	Tasks            []Task   `json:"tasks"`
	Start            string   `json:"start"`
	Change           string   `json:"change"`
	Terminated       string   `json:"terminated"`
	ActualTaskIndex  int      `json:"actual_task_index"`
	PredecessorId    []string `json:"predecessor_id"`
	SuccessorId      []string `json:"successor_id"`
}

var GtKurse = []Kurs{{
	ID:               "G000",
	Grade:            "0-99",
	Subject:          "übergreifend",
	MainTitle:        "Lernpfad - Bedienung",
	SubTitle:         "Im Selbststudium an Lernkursen teilnehmen.",
	Image:            "https://cdn.quasar.dev/img/mountains.jpg",
	DescriptionShort: "Eine kleine Einführung in die Bedienung von Lernpfad.",
	Description: "Hier erhältst Du eine Einführung in Lernpfad<br>" +
		"Der Kurs ist in 3 Videos mit ca. 3 min gegliedert:<br>" +
		" 1. Starten - wie bediene ich Lernpfad<br>" +
		" 2. Kleine Übungen - was muss ich beachten?<br>" +
		" 3. Wiederholung und Unterlegen<br><br>" +
		"Viel Spaß beim Lernen!",
	DescriptionMovie: "https://www.youtube.com/embed/EwkOJNgonJk",
	Duration:         10,
	Tasks: []Task{{
		ID:          1,
		Name:        "Information",
		SubTitle:    "",
		Description: "Eine kleine Einführung",
		Icon:        "theaters",
		Movie:       "https://www.youtube.com/embed/EwkOJNgonJk",
		Questions:   []Question{},
		Start:       "",
		Change:      "",
	}},
	Enrollments: 3,
	Rating:      2.5,
}, {
	ID:               "M001",
	Grade:            "2. Klasse",
	Subject:          "Mathematik",
	MainTitle:        "Mathe - Zahlen bis 100",
	SubTitle:         "Rechnen mit Zahlen bis 100",
	Image:            "https://img.fotocommunity.com/seestern-leben-02ad49f8-f09c-44a0-949d-8b162fdcaaf0.jpg?width=1000",
	DescriptionShort: "Übungen zur Addition und Subtraktion mit Zahlen von 0 bis 100",
	Description: "Hier erhältst Du eine Einführung in die Grundrechenarten<br>" +
		"Der Kurs ist in 3 Videos mit ca. 3 min gegliedert:<br>" +
		" 1. Zahlen und Zahlenstrahl<br>" +
		" 2. Addition - alles wird mehr<br>" +
		" 3. Subtraktion - alles wird weniger<br><br>" +
		"Viel Spaß beim Lernen!",
	DescriptionMovie: "https://www.youtube.com/embed/EwkOJNgonJk",
	Requirements:     "Für diesen Kurs sind keine Voraussetzungen notwendig.",
	Duration:         15,
	ActualTaskIndex:  0,
	Tasks: []Task{{
		ID:          1,
		Name:        "Einführung",
		SubTitle:    "",
		Description: "Eine kleine Einführung",
		Icon:        "theaters",
		Movie:       "https://www.youtube.com/embed/EwkOJNgonJk",
		Questions:   []Question{},
		Start:       "",
		Change:      "",
	}, {
		ID:          2,
		Name:        "Quiz",
		SubTitle:    "",
		Description: "Eine kleine Übung",
		Icon:        "quiz",
		Movie:       "",
		Questions: []Question{{
			ID:       1,
			Question: "Welche Aussage zu einem Zahlenstrahl ist richtig?",
			Answers: []Answer{{
				Answer: "Der Zahlenstrahl geht immer von links nach rechts.",
				Valid:  true,
				Choice: false,
			}, {
				Answer: "Der Zahlenstrahl geht immer von oben nach unten.",
				Valid:  false,
				Choice: false,
			}},
		}, {
			ID:       2,
			Question: "Welche Aussage zur Addition ist richtig?",
			Answers: []Answer{{
				Answer: "Summand plus Summand ist gleich Summe",
				Valid:  true,
				Choice: false,
			}, {
				Answer: "Summe plus Summand ist ein Summand",
				Valid:  false,
				Choice: false,
			}, {
				Answer: "Summand minus Summand ist Summe",
				Valid:  false,
				Choice: false,
			}, {
				Answer: "Die Summanden werden addiert",
				Valid:  true,
				Choice: false,
			}},
		}},
		Start:  "",
		Change: "",
	}, {
		ID:          3,
		Name:        "Addition",
		SubTitle:    "Zahlen bis 100",
		Description: "Addition bis 100",
		Icon:        "theaters",
		Movie:       "https://www.youtube.com/embed/EwkOJNgonJk",
		Questions:   []Question{},
		Start:       "",
		Change:      "",
	}, {
		ID:          4,
		Name:        "Quiz",
		SubTitle:    "",
		Description: "Eine kleine Übung",
		Icon:        "quiz",
		Movie:       "",
		Questions: []Question{{
			ID:       1,
			Question: "Welche Aussage zu einem Zahlenstrahl ist richtig?",
			Answers: []Answer{{
				Answer: "Der Zahlenstrahl geht immer von links nach rechts.",
				Valid:  true,
			}, {
				Answer: "Der Zahlenstrahl geht immer von oben nach unten.",
				Valid:  false,
			}},
		}},
		Start:  "",
		Change: "",
	}, {
		ID:          5,
		Name:        "Subtraktion",
		SubTitle:    "Zahlen bis 100",
		Description: "Subtraktion bis 100",
		Icon:        "theaters",
		Movie:       "https://www.youtube.com/embed/EwkOJNgonJk",
		Questions:   []Question{},
		Start:       "",
		Change:      "",
	}, {
		ID:          6,
		Name:        "Quiz",
		Description: "Eine kleine Übung",
		Icon:        "quiz",
		Movie:       "",
		Questions: []Question{{
			ID:       1,
			Question: "Welche Aussage zu einem Zahlenstrahl ist richtig?",
			Answers: []Answer{{
				Answer: "Der Zahlenstrahl geht immer von links nach rechts.",
				Valid:  true,
			}, {
				Answer: "Der Zahlenstrahl geht immer von oben nach unten.",
				Valid:  false,
			}},
		}},
		Start:  "",
		Change: "",
	}, {
		ID:          7,
		Name:        "Unterlagen",
		SubTitle:    "Übungsaufgaben",
		Description: "Beispielaufgaben und Lösungen zum selbständigen Üben.",
		Icon:        "newspaper",
		Movie:       "",
		Questions:   []Question{},
		Documents: []Document{{
			ID:          1,
			Name:        "Aufgabeblatt 1",
			Description: "Beispiele für Addition bis 100",
			Link:        "https://kekula.de/wp-content/uploads/klassenarbeit/Kekula_Klassenarbeit_Mathe_Klasse1_1.pdf",
		}, {
			ID:          2,
			Name:        "Lösungen 1",
			Description: "Lösungen zu den Aufgaben der Addition bis 100",
			Link:        "https://kekula.de/wp-content/uploads/klassenarbeit/Kekula_Klassenarbeit_Mathe_Klasse1_2.pdf",
		}},
		Start:  "",
		Change: "",
	}},
	Enrollments:   3,
	Rating:        3.2,
	Start:         "",
	Change:        "",
	Terminated:    "",
	PredecessorId: []string{"G000"},
	SuccessorId:   []string{},
}}

func getMe(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	type TLocalUser struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Tenant string `json:"tenant"`
		Role   string `json:"role"` // Role of the user; cannot be defined except for admin
	}
	user := &TLocalUser{
		ID:     claims[identityKey].(string),
		Name:   claims["name"].(string),
		Tenant: claims["tenant"].(string),
		Role:   claims["role"].(string),
	}

	c.JSON(200, gin.H{
		"type":    "Info",
		"message": "iO",
		"time":    time.Now().Format(cTime),
		"user":    user,
	})
	return
}

func getAdmin(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	type TLocalUser struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Tenant string `json:"tenant"`
		Role   string `json:"role"` // Role of the user; cannot be defined except for admin
	}
	user := &TLocalUser{
		ID:     claims[identityKey].(string),
		Name:   claims["name"].(string),
		Tenant: claims["tenant"].(string),
		Role:   claims["role"].(string),
	}
	if !strings.Contains(user.Role, "admin") {
		c.JSON(400, gin.H{
			"type":    "Error",
			"message": "Keine Adminberechtigung",
			"time":    time.Now().Format(cTime),
		})
		return
	}

	allSubjects := goEnforcer.GetAllSubjects()
	allSPolicy := goEnforcer.GetPolicy()
	allGroups := goEnforcer.GetNamedGroupingPolicy("g")

	c.JSON(200, gin.H{
		"type":        "Info",
		"message":     "iO",
		"time":        time.Now().Format(cTime),
		"user":        user,
		"allSubjects": allSubjects,
		"allSPolicy":  allSPolicy,
		"allGroups":   allGroups,
	})
	return
}

// "/admin/:phrase"
func setAdmin(c *gin.Context) {

	type MyTParam struct {
		Phrase string `json:"tenant" uri:"phrase" binding:"required"`
	}
	var MyParam MyTParam

	if err := c.ShouldBindUri(&MyParam); err != nil {
		c.JSON(400, gin.H{
			"type":    "Error",
			"message": "Parameter falsch",
			"time":    time.Now().Format(cTime),
		})
		return
	}

	claims := jwt.ExtractClaims(c)
	type TLocalUser struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Tenant string `json:"tenant"`
		Role   string `json:"role"` // Role of the user; cannot be defined except for admin
	}
	user := &TLocalUser{
		ID:     claims[identityKey].(string),
		Name:   claims["name"].(string),
		Tenant: claims["tenant"].(string),
		Role:   claims["role"].(string),
	}
	if !strings.Contains(user.Role, "admin") {
		c.JSON(400, gin.H{
			"type":    "Error",
			"message": "Keine Adminberechtigung",
			"time":    time.Now().Format(cTime),
		})
		return
	}

	content, err := io.ReadAll(c.Request.Body)
	if err == nil {

		if MyParam.Phrase == "addUser" {

			type body struct {
				User   string `json:"user"`
				Gruppe string `json:"gruppe"`
			}

			var payload body
			err = json.Unmarshal(content, &payload)
			if err == nil && payload.User != "" && payload.Gruppe != "" {

				added, err2 := goEnforcer.AddNamedGroupingPolicy("g", payload.User, payload.Gruppe, TENANT)
				if added && err2 == nil {

					err = goEnforcer.SavePolicy()
					if err == nil {

						c.JSON(200, gin.H{
							"type":    "Info",
							"message": "iO",
							"time":    time.Now().Format(cTime),
							"user":    user,
						})
						return

					}

				}
			}

		}

		if MyParam.Phrase == "removeUser" {

			type body struct {
				User   string `json:"user"`
				Gruppe string `json:"gruppe"`
			}

			var payload body
			err = json.Unmarshal(content, &payload)
			if err == nil && payload.User != "" && payload.Gruppe != "" {

				added, err2 := goEnforcer.RemoveGroupingPolicy(payload.User, payload.Gruppe, TENANT)
				if added && err2 == nil {

					err = goEnforcer.SavePolicy()
					if err == nil {
						c.JSON(200, gin.H{
							"type":    "Info",
							"message": "iO",
							"time":    time.Now().Format(cTime),
							"user":    user,
						})
						return
					}
				}
			}

		}

		if MyParam.Phrase == "addGruppe" {

			type body struct {
				Gruppe string `json:"gruppe"`
				Kurs   string `json:"kurs"`
			}

			var payload body
			err = json.Unmarshal(content, &payload)
			if err == nil && payload.Gruppe != "" && payload.Kurs != "" {

				added, err2 := goEnforcer.AddNamedPolicy("p", payload.Gruppe, TENANT, "all", payload.Kurs, "get")
				if added && err2 == nil {

					err = goEnforcer.SavePolicy()
					if err == nil {

						c.JSON(200, gin.H{
							"type":    "Info",
							"message": "iO",
							"time":    time.Now().Format(cTime),
							"user":    user,
						})
						return

					}

				}
			}

		}

		if MyParam.Phrase == "removeGruppe" {

			type body struct {
				Gruppe string `json:"gruppe"`
				Kurs   string `json:"kurs"`
			}

			var payload body
			err = json.Unmarshal(content, &payload)
			if err == nil && payload.Gruppe != "" && payload.Kurs != "" {

				added, err2 := goEnforcer.RemovePolicy(payload.Gruppe, TENANT, "all", payload.Kurs, "get")
				if added && err2 == nil {

					err = goEnforcer.SavePolicy()
					if err == nil {
						c.JSON(200, gin.H{
							"type":    "Info",
							"message": "iO",
							"time":    time.Now().Format(cTime),
							"user":    user,
						})
						return
					}
				}
			}

		}

	}

	c.JSON(400, gin.H{
		"type":    "Error",
		"message": "Falsche Phrase oder Operation nicht möglich",
		"time":    time.Now().Format(cTime),
	})
	return
}

// -----------------------------------
// URL "/getfile/:tenant/:name/:key/:file"
func getFile(c *gin.Context) {
	type MyTParam struct {
		Tenant string `json:"tenant" uri:"tenant" binding:"required"`
		Name   string `json:"name" uri:"name" binding:"required"`
		Key    string `json:"key" uri:"key" binding:"required"`
		File   string `json:"file" uri:"file" binding:"required"`
	}
	var MyParam MyTParam

	if err := c.ShouldBindUri(&MyParam); err != nil {
		c.JSON(400, gin.H{
			"type":    "Error",
			"message": "Parameter falsch",
			"time":    time.Now().Format(cTime),
		})
		return
	}

	claims := jwt.ExtractClaims(c)
	user := &User{
		ID:     claims[identityKey].(string),
		Name:   claims["name"].(string),
		Tenant: claims["tenant"].(string),
	}

	lDir := LOCAL_TEMPDIR + "/" + MyParam.Tenant + "/all"

	_, err2 := os.Stat(lDir + "/" + MyParam.Key + "/" + MyParam.File)
	if err2 == nil {

		c.File(lDir + "/" + MyParam.Key + "/" + MyParam.File)

		return
	}

	// file ermitteln und zurück liefern
	c.JSON(400, gin.H{
		"type":    "Error",
		"message": "invalide parameter",
		"time":    time.Now().Format(cTime),
		"data": gin.H{
			"user":  user,
			"param": MyParam,
		}})
	return

}

// -----------------------------------
// URL "/getitem/:tenant/:name/:key"
func getItem(c *gin.Context) {
	type MyTParam struct {
		Tenant string `json:"tenant" uri:"tenant" binding:"required"`
		Name   string `json:"name" uri:"name" binding:"required"`
		Key    string `json:"key" uri:"key" binding:"required"`
	}
	var MyParam MyTParam

	if err := c.ShouldBindUri(&MyParam); err != nil {
		c.JSON(400, gin.H{
			"type":    "Error",
			"message": "Parameter falsch",
			"time":    time.Now().Format(cTime),
		})
		return
	}

	claims := jwt.ExtractClaims(c)
	user := &User{
		ID:     claims[identityKey].(string),
		Name:   claims["name"].(string),
		Tenant: claims["tenant"].(string),
	}

	lDir := LOCAL_TEMPDIR + "/" + MyParam.Tenant + "/" + MyParam.Name

	if MyParam.Key == "Kurse" {
		var lKurse = []string{}
		//Directory zusammen bauen -  bei "Kurse" alle Dateien mit "Kurs_*.json" KursID zurück liefern
		_ = filepath.Walk(lDir, func(iPath string, iInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			lName := iInfo.Name()
			if filepath.Ext(lName) == ".json" && strings.HasPrefix(lName, "Kurs_") {
				re := regexp.MustCompile(`Kurs_([a-zA-Z0-9]*).json`)
				lKursName := re.FindAllStringSubmatch(lName, 1)
				if len(lKursName) > 0 && len(lKursName[0]) > 1 && lKursName[0][1] != "" {

					//Berechtigung prüfen
					sub := user.Name
					dom := MyParam.Tenant
					own := MyParam.Name
					obj := "Kurs_" + lKursName[0][1]
					act := strings.ToLower(c.Request.Method)
					valid, _ := goEnforcer.Enforce(sub, dom, own, obj, act)

					if valid == true {

						lKurse = append(lKurse, lKursName[0][1])
					} else {
						fmt.Println("keine Berechtigung:", sub, dom, own, obj, act)
					}
				}
			}
			return nil
		})

		c.JSON(200, gin.H{
			"type":    "Info",
			"message": "iO",
			"time":    time.Now().Format(cTime),
			"data":    lKurse,
		})
		return

	}

	content, err2 := os.ReadFile(lDir + "/" + MyParam.Key + ".json")
	if err2 == nil {
		var payload interface{}
		err2 = json.Unmarshal(content, &payload)
		if err2 == nil {

			// file ermitteln und zurück liefern
			c.JSON(200, gin.H{
				"type":    "Info",
				"message": "iO",
				"time":    time.Now().Format(cTime),
				"data":    payload})
			return

		}
	}

	// file ermitteln und zurück liefern
	c.JSON(400, gin.H{
		"type":    "Error",
		"message": "invalide parameter",
		"time":    time.Now().Format(cTime),
		"data": gin.H{
			"user":  user,
			"param": MyParam,
		}})
	return

}

// -------------------------------------//
func setItem(c *gin.Context) {
	type MyTParam struct {
		Tenant string `json:"tenant" uri:"tenant" binding:"required"`
		Name   string `json:"name" uri:"name" binding:"required"`
		Key    string `json:"key" uri:"key" binding:"required"`
	}
	var MyParam MyTParam

	if err := c.ShouldBindUri(&MyParam); err != nil {
		c.JSON(400, gin.H{
			"type":    "Error",
			"message": "Parameter falsch",
			"time":    time.Now().Format(cTime),
		})
		return
	}

	claims := jwt.ExtractClaims(c)
	user := &User{
		ID:     claims[identityKey].(string),
		Name:   claims["name"].(string),
		Tenant: claims["tenant"].(string),
	}

	//Berechtigung prüfen
	sub := user.Name
	dom := MyParam.Tenant
	own := MyParam.Name
	obj := MyParam.Key
	act := strings.ToLower(c.Request.Method)
	valid, _ := goEnforcer.Enforce(sub, dom, own, obj, act)

	if valid != true {
		valid, _ = goEnforcer.Enforce(sub, dom, own, '*', act)
	}

	if valid == true {
		//die Berechtigungen liegen vor

		content, err := io.ReadAll(c.Request.Body)
		if err == nil {

			var payload interface{}
			err2 := json.Unmarshal(content, &payload)
			if err2 == nil {

				lDir := LOCAL_TEMPDIR + "/" + MyParam.Tenant + "/" + MyParam.Name

				//gibt es das User-Verzeichnis?
				_, err := os.Stat(lDir)
				if os.IsNotExist(err) {
					err := os.Mkdir(lDir, 0755)
					if err != nil {
						log.Fatal(err)
					}
				}

				//gibt es die Datei schon?
				_, err = os.Stat(lDir + "/" + MyParam.Key + ".json")
				if os.IsNotExist(err) {
					go func() {
						//Kurs wurde neu geöffnet -> Zähler erhöhen
						lDirAll := LOCAL_TEMPDIR + "/" + MyParam.Tenant + "/all"
						contentAll, err2 := os.ReadFile(lDirAll + "/" + MyParam.Key + ".json")
						if err2 == nil {
							var kurs Kurs
							err2 = json.Unmarshal(contentAll, &kurs)
							if err2 == nil {
								kurs.Enrollments += 1
								file, _ := json.MarshalIndent(kurs, "", " ")
								_ = os.WriteFile(lDirAll+"/"+MyParam.Key+".json", file, 0644)

							}
						}
					}()
				}

				err = os.WriteFile(lDir+"/"+MyParam.Key+".json", content, 0644)
				if err == nil {
					// file ermitteln und zurück liefern
					c.JSON(200, gin.H{
						"type":    "Info",
						"message": "iO",
						"time":    time.Now().Format(cTime),
						"data":    payload})
					return
				}
			}
		}

	} else {
		fmt.Println("keine Berechtigung:", sub, dom, own, obj, act)
		c.JSON(400, gin.H{
			"type":    "Error",
			"message": "invalide permission",
			"time":    time.Now().Format(cTime),
			"data": gin.H{
				"sub": sub,
				"dom": dom,
				"own": own,
				"obj": obj,
				"act": act,
			}})
		return
	}

	// file ermitteln und zurück liefern
	c.JSON(400, gin.H{
		"type":    "Error",
		"message": "invalide parameter",
		"time":    time.Now().Format(cTime),
		"data": gin.H{
			"user":  user,
			"param": MyParam,
		}})
	return

}

func index(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/lernpfad")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// -----------------------------------------//
func initDir() {

	_, err := os.Stat(LOCAL_TEMPDIR)
	if os.IsNotExist(err) {
		err := os.Mkdir(LOCAL_TEMPDIR, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = os.Stat(LOCAL_TEMPDIR + "/data")
	if os.IsNotExist(err) {
		err := os.Mkdir(LOCAL_TEMPDIR+"/data", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = os.Stat(LOCAL_TEMPDIR + "/data/all")
	if os.IsNotExist(err) {
		err := os.Mkdir(LOCAL_TEMPDIR+"/data/all", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	//prüfen, ob casbin model vorhanden
	// _, err = os.Stat(LOCAL_TEMPDIR + "/model.conf")
	// if err != nil {
	// 	content, err2 := os.ReadFile("skibby/model.conf")
	// 	if err2 == nil {
	// 		_ = os.WriteFile(LOCAL_TEMPDIR+"/model.conf", content, 0644)
	// 	}
	// }
	_, err = os.Stat(LOCAL_TEMPDIR + "/policy.csv")
	if err != nil {
		content, err2 := os.ReadFile("lernpfad/policy.csv")
		if err2 == nil {
			_ = os.WriteFile(LOCAL_TEMPDIR+"/policy.csv", content, 0644)
		}
	}
}

// -----------------------------------------//
func main() {
	//Port setzen/merken
	lPort := os.Getenv("PORT")
	if lPort != "" {
		PORT = lPort
	}

	//lokales Verzeichnis setzen/merken/erzeugen
	lDir := os.Getenv("LOCAL_TEMPDIR")
	if lDir != "" {
		LOCAL_TEMPDIR = lDir
	}

	err := godotenv.Load(".env")

	initDir()

	casbin_model, err := model.NewModelFromString(`
  [request_definition]
  r = sub, dom, own, obj, act

  [policy_definition]
  p = sub, dom, own, obj, act

  [role_definition]
  g = _, _, _

  [policy_effect]
  e = some(where (p.eft == allow))

  [matchers]
  m = g(r.sub, p.sub, r.dom) && ((r.dom == p.dom && ( r.own == p.own || p.own == '*' ) && ( r.obj == p.obj || p.obj == '*' ) && ( r.act == p.act || p.act == '*' )) || r.own == r.sub)
  `)
	//m = g(r.sub, p.sub, r.dom) && (r.dom == p.dom && r.own == p.own && r.obj == p.obj && r.act == p.act || r.own == r.sub)
	//`)

	casbin_adapter := fileadapter.NewAdapter(LOCAL_TEMPDIR + "/policy.csv")

	//goEnforcer = casbin.NewEnforcer(LOCAL_TEMPDIR+"/model.conf", LOCAL_TEMPDIR+"/policy.csv")
	goEnforcer, err = casbin.NewEnforcer(casbin_model, casbin_adapter)

	//gin.SetMode(gin.DebugMode)
	gin.SetMode(gin.ReleaseMode)
	goServer = gin.New()

	goServer.StaticFile("favicon.ico", "dist/spa/favicon.ico")
	goServer.Static("/lernpfad", "dist/spa")

	goServer.Use(CORSMiddleware())

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:          "lernpfad",
		Key:            []byte(SECRET_KEY),
		Timeout:        ACCESS_TOKEN_EXPIRE_MINUTES,
		MaxRefresh:     ACCESS_TOKEN_EXPIRE_MINUTES, // Timeout + MaxRefresh for refresh-time
		SendCookie:     true,
		CookieName:     "lernpfad",
		CookieSameSite: http.SameSiteDefaultMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode
		IdentityKey:    identityKey,
		// Token-Claims: Daten zur Einbettung in den Token liefern, z.B. Username
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
					"name":      v.Name,
					"tenant":    v.Tenant,
					"role":      v.Role,
				}
			}
			return jwt.MapClaims{}
		},
		//liefert aus dem PayLoad das Feld zur Indentifizierung: identityKey
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				ID:     claims[identityKey].(string),
				Name:   claims["name"].(string),
				Tenant: claims["tenant"].(string),
				Role:   claims["role"].(string),
			}
		},
		// Authentifizierung des Users
		Authenticator: func(c *gin.Context) (interface{}, error) {

			type Login struct {
				Name     string `json:"username"`
				Password string `json:"password"`
				Tenant   string `json:"tenant" uri:"tenant" binding:"required"`
			}

			var loginVals Login

			//aus uri
			if err := c.ShouldBindUri(&loginVals); err != nil {
				c.JSON(400, gin.H{
					"type":    "Error",
					"message": "Parameter falsch",
					"time":    time.Now().Format(cTime),
				})
				return nil, (errors.New("incorrect Tenant information"))
			}

			//aus body
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Name
			password := loginVals.Password

			user := &User{
				ID:     userID,
				Name:   userID,
				Tenant: loginVals.Tenant,
			}

			if user.PasswordIsValide(password) {
				return user, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		// Berechtigung des Users
		Authorizator: func(data interface{}, c *gin.Context) bool {

			if c.FullPath() == "/api/me" {
				//own user data can always be read
				return true
			}

			type Access struct {
				Name   string `json:"name" uri:"name"`
				Key    string `json:"key" uri:"key"`
				Tenant string `json:"tenant" uri:"tenant" binding:"required"`
			}
			var accessData Access

			if user, ok := data.(*User); ok && user.ID != "" {

				if strings.HasPrefix(c.FullPath(), "/api/admin") && strings.Contains(user.Role, "admin") {
					//Admin may start admin prefix
					return true
				}

				if err := c.ShouldBindUri(&accessData); err == nil {

					sub := user.Name
					dom := accessData.Tenant
					own := accessData.Name
					obj := ""
					if accessData.Key != "" {
						obj = accessData.Key
					} else {
						obj = c.Request.URL.Path
					}
					act := strings.ToLower(c.Request.Method)

					valid, _ := goEnforcer.Enforce(sub, dom, own, obj, act)

					return valid
				}
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: lernpfad",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value.
		//This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	// //globale routen
	goServer.GET("/", index)

	goServer.POST("/login/:tenant", authMiddleware.LoginHandler)
	goServer.POST("/logout/:tenant", authMiddleware.LogoutHandler)
	goServer.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	api := goServer.Group("/api")
	// Refresh time can be longer than token timeout
	api.GET("/refresh_token", authMiddleware.RefreshHandler)
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/getfile/:tenant/:name/:key/:file", getFile)
		api.GET("/getitem/:tenant/:name/:key", getItem)
		api.POST("/setitem/:tenant/:name/:key", setItem)
		api.GET("/me", getMe)
		api.GET("/admin/:phrase", getAdmin)
		api.POST("/admin/:phrase", setAdmin)
	}

	//Beispieldaten ausbauen
	for _, ele := range GtKurse {
		file, _ := json.MarshalIndent(ele, "", " ")

		//gibt es die Datei schon?
		_, err = os.Stat(LOCAL_TEMPDIR + "/data/all/Kurs_" + ele.ID + ".json")
		if os.IsNotExist(err) {

			_ = os.WriteFile(LOCAL_TEMPDIR+"/data/all/Kurs_"+ele.ID+".json", file, 0644)
			if ele.ID == "M001" {
				ele.Start = time.Now().Format(cTime)
				file, _ := json.MarshalIndent(ele, "", " ")

				_ = os.WriteFile(LOCAL_TEMPDIR+"/data/jank/Kurs_"+ele.ID+".json", file, 0644)
			}
		}
	}

	//	if err := goServer.RunTLS(":8080", "./key/server.pem", "./key/server.key"); err != nil {
	if err := goServer.Run(":" + PORT); err != nil {
		log.Fatal("Launch not possible:", err)
	}

}
