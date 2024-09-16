<!-- eslint-disable vue/no-v-text-v-html-on-component -->
<template>
    <q-page class="fit row">
        <q-card style="width: 100%; height: 100%">
            <q-card-section>
                <q-card-section class="row">
                    <div class="col">
                        <div class="text-h7">
                            {{ kurs.id }} &nbsp; &nbsp; Klasse: {{ kurs.grade }}&nbsp; &nbsp; Fach: {{ kurs.subject }}
                        </div>
                        <div class="text-h4">{{ kurs.main_title }}</div>
                    </div>
                </q-card-section>
                <q-stepper v-model="step" ref="stepper" color="primary" animated header-nav>
                    <template v-for="task in kurs.tasks" v-bind:key="task.id">
                        <q-step
                            :name="task.id"
                            :title="task.name"
                            :icon="task.icon"
                            :done="step > task.id"
                            :caption="task.sub_title"
                            done-color="green"
                        >
                            <div class="text-body1">{{ task.description }}<br /></div>
                            <!-- Film ------------------------------------------------- -->
                            <div v-if="task.icon == 'theaters'" class="q-pa-md row justify-center">
                                <div style="height: auto; width: 800px">
                                    <q-video :ratio="16 / 9" :src="task.movie" />
                                </div>
                            </div>

                            <!-- Quiz ------------------------------------------------- -->
                            <div v-if="task.icon == 'quiz'" class="text-h6">
                                <template v-for="question in task.questions" v-bind:key="question.id">
                                    <q-card class="col">
                                        <q-card-section>
                                            <div class="text-h5">Frage {{ question.id }}<br /></div>
                                            <div class="text-h6">{{ question.question }}<br /></div>

                                            <div v-if="getMaxTrueAnswer(question) == 1" class="text-caption row">
                                                (Es gibt genau eine richtige Antwort.)
                                                <div v-if="task.check != ''">
                                                    &nbsp;Punkte {{ question.getPoints() }} von {{ question.getMaxPoints() }}
                                                </div>
                                                <br />
                                            </div>
                                            <div v-if="getMaxTrueAnswer(question) > 1" class="text-caption row">
                                                (Es gibt {{ getMaxTrueAnswer(question) }} richtige Antworten.)
                                                <div v-if="task.check != ''">
                                                    &nbsp;Punkte {{ question.getPoints() }} von {{ question.getMaxPoints() }}
                                                </div>
                                                <br />
                                            </div>
                                            <template v-for="answer in question.answers" v-bind:key="answer.id">
                                                <div class="col text-subtitle2">
                                                    <div v-if="getMaxTrueAnswer(question) > 1" class="row">
                                                        <!-- Check-Box -->
                                                        <q-checkbox
                                                            v-model="answer.choice"
                                                            :label="answer.answer"
                                                            :disable="task.check != ''"
                                                        />
                                                        <div v-if="task.check != ''" class="row">
                                                            &nbsp;&nbsp;
                                                            <span v-if="answer.valid == true" class="material-icons" style="color: gray">
                                                                check_box
                                                            </span>
                                                            <span v-if="answer.valid == false" class="material-icons" style="color: gray">
                                                                check_box_outline_blank </span
                                                            >&nbsp;
                                                            <div
                                                                v-if="
                                                                    onCheckResult(answer.valid, answer.choice, true) &&
                                                                    (answer.valid == true || answer.choice == true)
                                                                "
                                                                style="margin-top: 10px; color: green"
                                                            >
                                                                Super - richtig
                                                            </div>
                                                            <div
                                                                v-if="
                                                                    !onCheckResult(answer.valid, answer.choice, true) &&
                                                                    (answer.valid == true || answer.choice == true)
                                                                "
                                                                style="margin-top: 10px; color: red"
                                                                class="text-negativ"
                                                            >
                                                                uups - leider falsch
                                                            </div>
                                                        </div>
                                                    </div>
                                                    <div v-if="getMaxTrueAnswer(question) == 1" class="row">
                                                        <!-- Radio Button-->
                                                        <q-radio
                                                            v-model="question.choice"
                                                            :val="answer.answer"
                                                            :label="answer.answer"
                                                            :disable="task.check != ''"
                                                        />
                                                        <div v-if="task.check != ''" class="row">
                                                            &nbsp;&nbsp;
                                                            <span v-if="answer.valid == true" class="material-icons" style="color: gray">
                                                                radio_button_checked
                                                            </span>
                                                            <span v-if="answer.valid == false" class="material-icons" style="color: gray">
                                                                radio_button_unchecked
                                                            </span>
                                                            &nbsp;
                                                            <div
                                                                v-if="
                                                                    onCheckResult(answer.answer, question.choice, answer.valid) &&
                                                                    (answer.valid == true || answer.answer == question.choice)
                                                                "
                                                                style="margin-top: 10px; color: green"
                                                            >
                                                                Super - richtig
                                                            </div>
                                                            <div
                                                                v-if="
                                                                    !onCheckResult(answer.answer, question.choice, answer.valid) &&
                                                                    (answer.valid == true || answer.answer == question.choice)
                                                                "
                                                                style="margin-top: 10px; color: red"
                                                                class="text-negativ"
                                                            >
                                                                uups - leider falsch
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            </template>
                                        </q-card-section>
                                    </q-card>
                                </template>
                                <q-card-action>
                                    <q-btn
                                        v-if="task.check == ''"
                                        @click="onCheck(user, task, kurs)"
                                        label="Prüfen"
                                        style="background: goldenrod; color: white; margin-top: 5px"
                                    />
                                    <q-btn
                                        v-if="task.check != ''"
                                        @click="onReopen(user, task, kurs)"
                                        label="Noch mal"
                                        style="background: goldenrod; color: white; margin-top: 5px"
                                    />
                                </q-card-action>
                            </div>

                            <!-- Doku ------------------------------------------------- -->
                            <div v-if="task.icon == 'newspaper'" class="text-h6">
                                <div class="row">
                                    <div class="col" style="width: 49%">
                                        <div class="text-h6">Unterlagen:<br /></div>

                                        <template v-for="document in task.documents" v-bind:key="document.id">
                                            <div class="col">
                                                <q-btn
                                                    outline
                                                    style="width: 300px"
                                                    @click="
                                                        document_id = document.id;
                                                        document.startView();
                                                        kurs.saveMe(user);
                                                    "
                                                    :color="isDocViewed(document.view_date)"
                                                >
                                                    <div class="col">
                                                        <div class="text-subtitle2">{{ document.name }}<br /></div>
                                                        <div class="text-caption">{{ document.description }}<br /></div>
                                                    </div>
                                                </q-btn>
                                            </div>
                                        </template>
                                    </div>
                                    <div class="col" style="width: 49%">
                                        <div v-if="getDocument(task, document_id) != ''">
                                            <iframe :src="getDocument(task, document_id)" style="height: 778px; width: 100%"></iframe>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </q-step>
                    </template>

                    <template v-slot:navigation>
                        <q-stepper-navigation>
                            <div class="row justify-between">
                                <q-btn
                                    v-if="step > 1"
                                    style="background: goldenrod; color: white"
                                    @click="$refs.stepper.previous()"
                                    label="Zurück"
                                    class="q-ml-sm"
                                /><q-btn
                                    v-if="kurs.tasks && step != kurs.tasks.length"
                                    @click="
                                        onNextButton();
                                        document_id = -1;
                                        $refs.stepper.next();
                                        kurs.setTask(step - 1);
                                    "
                                    style="background: goldenrod; color: white"
                                    label="Weiter"
                                />
                                <q-btn on-right icon="signpost" style="background: goldenrod; color: white" @click="onOverview(user, kurs)"
                                    >Übersicht</q-btn
                                >
                            </div>
                        </q-stepper-navigation>
                    </template>
                </q-stepper>
            </q-card-section>
        </q-card>
    </q-page>
</template>
<style lang="sass" scoped>
.my-card
  width: 100%
  max-width: 250px
</style>
<script>
import { defineComponent } from "vue";

export default defineComponent({
    name: "DoKursPage",
    props: ["user", "server"],
    data: () => {
        return {
            document_id: -1,
            step: 1,
            kurs: {},
        };
    },

    async created() {
        var lParams = this.$route.params;
        if (!this.user.name || !lParams || this.user.name != lParams.user_name) {
            this.$q.notify({
                message: "Fehler beim Usernamen",
                color: "negative",
                icon: "warning",
            });
            this.$router.push("/");
        } else {
            this.onLoadKurs();
        }
    },

    methods: {
        onCheck(iUser, iTask, iKurs) {
            iTask.startCheck();
            iKurs.kurs.setTask(this.step - 1);
            let that = this;
            iKurs.saveMe(iUser).then((res) => {
                if (res.status != 200) {
                    that.$q.notify({
                        message: `Fehler beim Sichern des Kurs ${that.kurs.id}.`,
                        color: "negative",
                        icon: "warning",
                    });
                }
            });
        },
        onReopen(iUser, iTask, iKurs) {
            this.$q.notify({
                message: `Die Fragen wurden neu gemischt.`,
                type: "positive",
                position: "center",
                timeout: 1000,
            });
            iTask.clearCheck();
            iKurs.kurs.setTask(this.step - 1);
            let that = this;
            iKurs.saveMe(iUser).then((res) => {
                if (res.status != 200) {
                    that.$q.notify({
                        message: `Fehler beim Sichern des Kurs ${that.kurs.id}.`,
                        color: "negative",
                        icon: "warning",
                    });
                }
            });
        },
        onOverview(iUser, iKurs) {
            this.kurs.setTask(this.step - 1);
            let that = this;
            iKurs.saveMe(iUser).then((res) => {
                if (res.status != 200) {
                    that.$q.notify({
                        message: `Fehler beim Sichern des Kurs ${that.kurs.id}.`,
                        color: "negative",
                        icon: "warning",
                    });
                }
            });
            this.$router.push("/kurs/" + iKurs.id);
        },
        getDateTime() {
            let lNow = new Date();
            return lNow.toLocaleDateString() + " " + lNow.toLocaleTimeString();
        },
        isDocViewed(iViewDate) {
            if (iViewDate == "") {
                return "black";
            } else {
                return "green";
            }
        },
        getMaxTrueAnswer(iQuestion) {
            let eMaxTrue = 0;
            if (iQuestion.answers) {
                for (let element of iQuestion.answers) {
                    if (element.valid) {
                        eMaxTrue += 1;
                    }
                }
            }
            return eMaxTrue;
        },
        getDocument(iTask, iDocument_id) {
            for (let element of iTask.documents) {
                if (element.id === iDocument_id) {
                    if (element.link.toLowerCase().startsWith("http")) {
                        return element.link;
                    }
                    return this.user.server + element.link + "?token=" + this.user.token;
                }
            }
            return "";
        },
        onNextButton() {
            if (this.step === this.kurs.tasks.length) {
                //ende Button betätigt
                this.$router.push("/kurs/" + this.kurs.id);
            }
        },
        onCheckResult(iValid, iCoice, iVote) {
            if (iCoice == undefined) {
                iCoice = false;
            }
            if (iValid == undefined) {
                iValid = false;
            }
            if (iValid == iCoice) {
                return iVote;
            } else {
                return !iVote;
            }
        },
        async onLoadKurs() {
            var lParams = this.$route.params;
            var lError = false;
            if (!this.user.name || !lParams || !lParams.kurs_id) {
                lError = true;
            }
            for (this.kurs of this.user.kurse) {
                if (this.kurs.id == lParams.kurs_id) {
                    break;
                }
                this.kurs = {};
            }
            if (lError || !this.kurs.id) {
                this.$q.notify({
                    message: "Fehler beim Kurs",
                    color: "negative",
                    icon: "warning",
                });
                this.$router.push("/");
                return;
            }
            //Werte zum Kurs / Task initialisieren
            this.step = this.kurs.actual_task_index + 1;
            this.kurs.setTask(this.step - 1);
        },
    },
});
</script>
