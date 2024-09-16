const TENANT = "data"

function LocalDateTime() {
  let lNow = new Date()
  return lNow.toLocaleDateString() + ' ' + lNow.toLocaleTimeString()
}

class User {
  id
  name
  password
  token
  server
  role
  kurse
  onLogin
  onLoadKurse
  setting
  grade_cat
  subject_cat
  admin

  constructor(iUser) {
    let lSet = { klasse: '', fach: '', sel_status: { value: "5", label: "offene und laufende Kurse" }, sortierung: '', suche_text: '' }
    iUser = iUser || {}
    this.id = iUser.name || iUser.id || "Gast"
    this.name = iUser.name || "Gast"
    this.password = iUser.password || "Gast"
    this.token = iUser.token || ""
    this.server = iUser.server || ""
    this.role = iUser.role || ""
    this.kurse = iUser.kurse || []
    this.onLogin = false
    this.onLoadKurse = false
    this.setting = iUser.setting || lSet
    this.grade_cat = []
    this.subject_cat = []
    this.admin = iUser.admin || {}
  }

  async logout() {
    await fetch(this.server + "/logout/" + TENANT, {
      //mode: "no-cors",
      method: "POST",
      headers: {
        Accept: "*/*",
        "Content-Type": "application/json",
      },
    });
    this.id = 'Gast'
    this.name = 'Gast'
    this.password = ''
    this.token = ''
    this.kurse = []
    this.onLogin = false
    this.onLoadKurse = false
    this.grade_cat = []
    this.subject_cat = []
  }

  async login(iServer = "", iName = "Gast", iPassword = "Gast") {
    if (this.onLogin == true) { return }
    this.onLogin = true
    if (!iName || !iPassword) {
      return { status: 400, message: "Name oder Passwort fehlt." }// Fehler
    }

    const res = await fetch(iServer + "/login/" + TENANT, {
      //mode: "no-cors",
      method: "POST",
      headers: {
        Accept: "*/*",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: iName,
        password: iPassword,
      }),
    });

    const content = await res.json();
    if (res.status == 200 && content.token) {
      this.id = iName
      this.name = iName
      this.password = iPassword
      this.server = iServer
      this.token = content.token || "";
      this.kurse = [];
      this.getRole() //load the role asynchronously
      this.onLogin = false
      let message = content.message || "ok"
      return { status: 200, message: message }

    } else {
      this.onLogin = false
      let message = content.message || "Fehler"
      return { status: 400, message: message }

    }

  }

  getAdminData() {
    let that = this
    fetch(this.server + "/api/admin/user", {
      //mode: "no-cors",
      method: "GET",
      headers: {
        Accept: "*/*",
        "Content-Type": "application/json",
        "Authorization": "Bearer " + this.token,
      },
    }).then((res) => {
      if (res.status != 200) { return }
      res.json().then((content) => {
        if (content && content["allSPolicy"]) {
          that.admin["allSPolicy"] = []
          for (let p of content["allSPolicy"]) {
            that.admin["allSPolicy"].push({ "id": that.admin["allSPolicy"].length, "Gruppe": p[0], "Kurs": p[3] })
          }
        }
        if (content && content["allSubjects"]) {
          that.admin["allSubjects"] = []
          for (let s of content["allSubjects"]) {
            if (!that.admin["allSubjects"].includes(s)) {
              that.admin["allSubjects"].push(s)
            }
          }
        }
        if (content && content["allGroups"]) {
          that.admin["allGroups"] = []
          for (let g of content["allGroups"]) {
            that.admin["allGroups"].push({ "id": that.admin["allGroups"].length, "User": g[0], "Gruppe": g[1] })
          }
        }
      })
    });
  }

  getRole() {
    let that = this
    fetch(this.server + "/api/me", {
      //mode: "no-cors",
      method: "GET",
      headers: {
        Accept: "*/*",
        "Content-Type": "application/json",
        "Authorization": "Bearer " + this.token,
      },
    }).then((res) => {
      if (res.status != 200) { return }
      res.json().then((content) => {
        if (content && content.user && content.user.role) {
          that.role = content.user.role
          if (that.role.includes("admin")) {
            that.getAdminData()
          }
        }
      })
    });
  }


  async getItem(iKey, iUser = "") {
    try {
      let lUser = iUser || this.name
      if (lUser == "Gast") { lUser = "all" }
      const res = await fetch(this.server + "/api/getitem/" + TENANT + "/" + lUser + "/" + iKey, {
        //mode: "no-cors",
        method: "GET",
        headers: {
          Accept: "*/*",
          "Content-Type": "application/json",
          "Authorization": "Bearer " + this.token,
        },
      });

      const content = await res.json();
      if (res.status == 200 && content.data) {

        return { status: res.status, message: 'iO', data: content.data }

      } else {
        let message = content.message || "Fehler"
        console.log("Status: ", res.status, "Message: ", message)
        return { status: res.status, message: message, data: undefined }

      }
    } catch (err) {
      let message = err.message || "Fehler"
      console.log("Status: ", '500', "Message: ", message)
      return { status: '500', message: message, data: undefined }
    }
  }


  async setItem(iKey, iUser = "", iObject) {
    try {
      //Storage.set({ key: iKey, value: JSON.stringify(iObject) })
      let lUser = iUser || this.name
      if (lUser == "Gast") { lUser = "all" }
      const res = await fetch(this.server + "/api/setitem/" + TENANT + "/" + lUser + "/" + iKey, {
        //mode: "no-cors",
        method: "POST",
        headers: {
          Accept: "*/*",
          "Content-Type": "application/json",
          "Authorization": "Bearer " + this.token,
        },
        body: JSON.stringify(iObject)
      });
      const content = await res.json();
      return { status: res.status, data: content }
    } catch (err) {
      let message = err.message || "Fehler"
      console.log("Status: ", '500', "Message: ", message)
      return { status: '500', message: message, data: undefined }
    }
  }

  async getKurseByUser(iUserID) {
    let lKurse = []
    //eigene Kurse lesen
    const { data: lKuseList } = await this.getItem("Kurse", iUserID)
    if (lKuseList) {
      for (let kursName of lKuseList) {
        const { data: lKurs } = await this.getItem("Kurs_" + kursName, iUserID)
        if (lKurs) {
          let lObjKurs = new Kurs(lKurs)
          lKurse.push(lObjKurs)
        }
      }
    }
    return lKurse
  }

  async loadKurse() {
    if (this.onLoadKurse == true) { return }
    this.onLoadKurse = true
    let lKurse = []

    if (this.token) {
      //eigene Kurse lesen
      const { data: lKuseList } = await this.getItem("Kurse")
      if (lKuseList) {
        for (let kursName of lKuseList) {
          const { data: lKurs } = await this.getItem("Kurs_" + kursName)
          if (lKurs) {
            let lObjKurs = new Kurs(lKurs)
            lKurse.push(lObjKurs)
          }
        }
      }

      //die Kurse von all dazulesen
      const { data: lMyKuseList } = await this.getItem("Kurse", "all")
      if (lMyKuseList) {
        for (let kursName of lMyKuseList) {
          const { data: lKurs } = await this.getItem("Kurs_" + kursName, "all")
          if (lKurs) {
            let lObjKurs = new Kurs(lKurs)
            let lChange = false
            for (let i = 0; i < lKurse.length; i++) {
              if (lKurse[i].id == lObjKurs.id) {
                //es hibt den Kurs schon in der Liste
                // die Teilnehmer aus der Vorlage übernehmen
                lKurse[i].enrollments = lObjKurs.enrollments
                lChange = true
                break
              }
            }

            if (lKurse.length == 0 || !lChange) {
              //es gibt ihn noch nicht - anhängen
              lKurse.push(lObjKurs)
            }
          }
        }
      }
    }

    this.kurse = lKurse
    this.grade_cat = [""]
    this.subject_cat = [""]
    for (let kurs of this.kurse) {
      let ifFound = false
      for (let klasse of this.grade_cat) {
        if (klasse == kurs.grade) {
          ifFound = true
          break
        }
      }
      if (ifFound == false) {
        this.grade_cat.push(kurs.grade)
      }
      ifFound = false
      for (let fach of this.subject_cat) {
        if (fach == kurs.subject) {
          ifFound = true
          break
        }
      }
      if (ifFound == false) {
        this.subject_cat.push(kurs.subject)
      }
    }
    this.onLoadKurse = false

    this.sortKurseByFeature("id")
  }

  sortKurseByFeature(iFeature = "") {
    if (iFeature == 'Klasse') {
      iFeature = 'grade'
    }
    if (iFeature == 'Fach') {
      iFeature = 'subject'
    }
    if (!iFeature || iFeature == "") {
      iFeature = 'id'
    }
    this.kurse.sort((a, b) => {
      if (a[iFeature] < b[iFeature]) {
        return -1
      }
      if (a.id > b.id) {
        return 1
      }
      return 0
    })
  }

  changeKurs(iIndex, iKurs) {
    if (this.kurse && this.kurse[iIndex] && this.kurse[iIndex].id == iKurs.id) {
      this.kurse[iIndex] = iKurs
    }
  }

  setSetting(iSetting) {
    let lSet = { klasse: '', fach: '', sel_status: { value: "5", label: "offene und laufende Kurse" }, sortierung: '', suche_text: '' }
    this.setting = iSetting || lSet;
  }

  async addGruppe(iGruppe, iKurs) {
    try {
      const res = await fetch(this.server + "/api/admin/addGruppe", {
        //mode: "no-cors",
        method: "POST",
        headers: {
          Accept: "*/*",
          "Content-Type": "application/json",
          "Authorization": "Bearer " + this.token,
        },
        body: JSON.stringify({
          "gruppe": iGruppe,
          "kurs": iKurs
        })
      });
      const content = await res.json();
      return { status: res.status, data: content }
    } catch (err) {
      let message = err.message || "Fehler"
      console.log("Status: ", '500', "Message: ", message)
      return { status: '500', message: message, data: undefined }
    }
  }

  async removeGruppe(iGruppe, iKurs) {
    try {
      const res = await fetch(this.server + "/api/admin/removeGruppe", {
        //mode: "no-cors",
        method: "POST",
        headers: {
          Accept: "*/*",
          "Content-Type": "application/json",
          "Authorization": "Bearer " + this.token,
        },
        body: JSON.stringify({
          "gruppe": iGruppe,
          "kurs": iKurs
        })
      });
      const content = await res.json();
      return { status: res.status, data: content }
    } catch (err) {
      let message = err.message || "Fehler"
      console.log("Status: ", '500', "Message: ", message)
      return { status: '500', message: message, data: undefined }
    }
  }

  async addUser(iUser, iGruppe) {
    try {
      const res = await fetch(this.server + "/api/admin/addUser", {
        //mode: "no-cors",
        method: "POST",
        headers: {
          Accept: "*/*",
          "Content-Type": "application/json",
          "Authorization": "Bearer " + this.token,
        },
        body: JSON.stringify({
          "user": iUser,
          "gruppe": iGruppe
        })
      });
      const content = await res.json();
      return { status: res.status, data: content }
    } catch (err) {
      let message = err.message || "Fehler"
      console.log("Status: ", '500', "Message: ", message)
      return { status: '500', message: message, data: undefined }
    }
  }

  async removeUser(iUser, iGruppe) {
    try {
      const res = await fetch(this.server + "/api/admin/removeUser", {
        //mode: "no-cors",
        method: "POST",
        headers: {
          Accept: "*/*",
          "Content-Type": "application/json",
          "Authorization": "Bearer " + this.token,
        },
        body: JSON.stringify({
          "user": iUser,
          "gruppe": iGruppe
        })
      });
      const content = await res.json();
      return { status: res.status, data: content }
    } catch (err) {
      let message = err.message || "Fehler"
      console.log("Status: ", '500', "Message: ", message)
      return { status: '500', message: message, data: undefined }
    }
  }

}

class Answer {
  answer
  valid
  choice

  constructor(iAnswer) {
    iAnswer = iAnswer || {}
    this.answer = iAnswer.answer || ""
    this.valid = iAnswer.valid || false
    this.choice = iAnswer.choice || false
  }
}

class Question {
  id
  question
  answers
  choice


  constructor(iQuestion) {
    iQuestion = iQuestion || {}
    this.id = iQuestion.id || 0
    this.question = iQuestion.question || ""
    this.answers = iQuestion.answers || []
    for (let i = 0; i < this.answers.length; i += 1) {
      this.answers[i] = new Answer(this.answers[i])
    }
    this.choice = iQuestion.choice || ""

  }

  getPoints() {
    let ePoints = 0
    for (let answer of this.answers) {
      if (answer.valid == true) {
        if (answer.valid == answer.choice) {
          //richtige Markierung check-box
          ePoints += 1
        }
        if (this.choice == answer.answer) {
          //richtige Markierung radio-button
          ePoints += 3
        }
      } else {
        if (answer.choice == true) {
          //falsche Markierung
          ePoints -= 1
        }
      }
    }
    if (ePoints < 0) {
      // keine negativen Summen
      ePoints = 0
    }
    return ePoints
  }

  getMaxPoints() {
    let ePoints = 0
    for (let answer of this.answers) {
      if (answer.valid == true) {
        ePoints += 1
      }
    }
    if (ePoints == 1) {
      ePoints = 3
    }
    return ePoints;
  }
}

class Test {
  test_date
  points
  questions

  constructor(iTest) {
    iTest = iTest || {}
    this.test_date = iTest.test_date || ""
    this.points = iTest.points || 0
    this.questions = iTest.questions || []
    for (let i = 0; i < this.questions.length; i += 1) {
      this.questions[i] = new Question(this.questions[i])
    }
  }
}

class Document {
  id
  name
  description
  link
  token
  view_date

  constructor(iDoc) {
    iDoc = iDoc || {}
    this.id = iDoc.id || 0
    this.name = iDoc.name || ""
    this.description = iDoc.description || ""
    this.link = iDoc.link || ""
    this.token = iDoc.token || ""
    this.view_date = iDoc.view_date || ""
  }

  startView() {
    if (this.view_date == "") {
      this.view_date = LocalDateTime()
    }
  }
}

class Task {
  id
  name
  sub_title
  description
  icon
  movie
  questions
  documents
  start
  change
  check
  test

  constructor(iTask) {
    iTask = iTask || {}
    this.id = iTask.id || 0
    this.name = iTask.name || ""
    this.sub_title = iTask.sub_title || ""
    this.description = iTask.description || ""
    this.icon = iTask.icon || ""
    this.movie = iTask.movie || ""
    this.questions = iTask.questions || []
    for (let i = 0; i < this.questions.length; i += 1) {
      this.questions[i] = new Question(this.questions[i])
    }
    this.documents = iTask.documents || []
    for (let i = 0; i < this.documents.length; i += 1) {
      this.documents[i] = new Document(this.documents[i])
    }
    this.start = iTask.start || ""
    this.change = iTask.change || ""
    this.check = iTask.check || ""
    this.tests = iTask.test || []
  }

  startTask() {
    if (this.start == "") {
      this.start = LocalDateTime()
    }
  }
  startCheck() {
    if (this.check == "") {
      this.check = LocalDateTime()
    }
  }
  clearCheck() {
    if (this.icon != "quiz") { return }
    this.check = ""
    for (let question of this.questions) {
      question.choice = ""
      for (let answer of question.answers) {
        answer.choice = false
      }
      for (let i = 0; i < question.answers.length; i += 1) {
        let ii = (Math.floor(Math.random() * (question.answers.length + 1)) - 1);
        if (ii < 0) { ii = 0 }
        if (ii >= question.answers.length) { ii = question.answers.length - 1 }
        let a = question.answers[i]
        question.answers[i] = question.answers[ii]
        question.answers[ii] = a
      }
    }
  }

}

class Kurs {
  id
  grade
  subject
  publication
  main_title
  sub_title
  image
  description_short
  description
  description_movie
  requirements
  duration
  difficulty
  enrollments
  rating
  tasks
  start
  change
  terminated
  actual_task_index
  predecessor_id
  successor_id

  constructor(iKurs) {
    iKurs = iKurs || {}
    this.id = iKurs.id || ""
    this.grade = iKurs.grade || ""
    this.subject = iKurs.subject || ""
    this.publication = iKurs.publication || LocalDateTime()
    this.main_title = iKurs.main_title || ""
    this.sub_title = iKurs.sub_title || ""
    this.image = iKurs.image || ""
    this.description_short = iKurs.description_short || ""
    this.description = iKurs.description || ""
    this.description_movie = iKurs.description_movie || ""
    this.requirements = iKurs.requirements || ""
    this.duration = iKurs.duration || 0
    this.difficulty = iKurs.difficulty || 0
    this.enrollments = iKurs.enrollments || 0
    this.rating = iKurs.rating || 0
    this.tasks = iKurs.tasks || []
    for (let i = 0; i < this.tasks.length; i += 1) {
      this.tasks[i] = new Task(this.tasks[i])
    }
    this.start = iKurs.start || ""
    this.change = iKurs.change || ""
    this.terminated = iKurs.terminated || ""
    this.actual_task_index = iKurs.actual_task_index || 0
    this.predecessor_id = iKurs.predecessor_id || []
    this.successor_id = iKurs.successor_id || []
  }

  startMe() {
    if (this.start == "") {
      this.start = LocalDateTime()
    }
  }

  endMe() {
    if (this.terminated == "") {
      this.terminated = LocalDateTime()
    }
  }

  setTask(iIndex = 0) {
    if (this.tasks.length >= iIndex && (this.actual_task_index < iIndex || iIndex == 0)) {
      this.actual_task_index = iIndex
      this.tasks[iIndex].startTask()
    }

  }

  getPoints() {
    let eSumme = 0
    for (let task of this.tasks) {
      if (task.icon == "quiz" && task.check != "") {
        for (let question of task.questions) {
          eSumme += question.getPoints()
        }
      }
    }
    return eSumme
  }

  getMaxPoints() {
    let eSumme = 0
    for (let task of this.tasks) {
      if (task.icon == "quiz") {
        for (let question of task.questions) {
          eSumme += question.getMaxPoints()
        }
      }
    }
    return eSumme
  }

  saveMe(iUser) {
    if (iUser.name && iUser.name != "Gast") {
      this.change = LocalDateTime()
      return iUser.setItem('Kurs_' + this.id, iUser.name, this)
    }
  }

}

export { User, Kurs, Task, Document, Test, Question, Answer }
