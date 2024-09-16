<!-- eslint-disable vue/no-v-text-v-html-on-component -->
<template>
    <q-page class="fit row">
        <q-card style="width: 100%; height: 100%">
            <q-card-section>
                <div class="q-pa-md">
                    <q-form v-if="kurs" style="width: 100%">
                        <div class="col">
                            <div class="row">
                                <div class="col" style="width: 49%">
                                    <div>
                                        <img :src="kurs.image" style="object-fit: cover; width: 250px; height: 180px" flat />
                                        <span v-if="kurs.start" class="text-h4" style="position: absolute; left: 2%; top: 2%">
                                            <span class="material-icons" style="color: black"> verified </span>
                                        </span>
                                    </div>
                                </div>
                                <div class="col" style="width: 49%">
                                    <div class="text-h7">
                                        {{ kurs.id }} &nbsp; &nbsp; Klasse: {{ kurs.grade }}&nbsp; &nbsp; Fach: {{ kurs.subject }}
                                    </div>
                                    <div class="text-h4">{{ kurs.main_title }}</div>

                                    <div class="text-subtitle2">{{ kurs.sub_title }}</div>
                                    <div v-if="kurs.duration" class="q-pt-none">
                                        Dauer: {{ kurs.duration }} &nbsp; Minuten &nbsp;/&nbsp; {{ schwierigkeit[kurs.difficulty] }}
                                    </div>
                                    <div class="q-pt-none">
                                        Teilnehmer: {{ kurs.enrollments }}
                                        &nbsp; &nbsp; Score:
                                        <q-rating
                                            v-model="kurs.rating"
                                            size="01.0em"
                                            color="yellow"
                                            icon="star_border"
                                            icon-selected="star"
                                        />
                                    </div>
                                    <div v-if="kurs.start != ''" class="q-pt-none">
                                        Start: {{ kurs.start }} <br />
                                        Stand: &nbsp;{{ kurs.actual_task_index + 1 }} von {{ kurs.tasks.length }} &nbsp;
                                        <q-circular-progress
                                            :value="kurs.actual_task_index + 1.01"
                                            min="0"
                                            :max="kurs.tasks.length"
                                            size="20px"
                                            :thickness="1"
                                            color="green"
                                            track-color="grey-3"
                                        />
                                        &nbsp; Punkte: {{ kurs.getPoints() }} von {{ kurs.getMaxPoints() }}&nbsp;
                                        <q-circular-progress
                                            :value="kurs.getPoints() + 0.01"
                                            min="0"
                                            :max="kurs.getMaxPoints()"
                                            size="20px"
                                            :thickness="1"
                                            color="blue"
                                            track-color="grey-3"
                                        />
                                    </div>
                                    <div><br /></div>
                                    <q-btn
                                        v-if="!checkStartJet() && kurs.terminated == ''"
                                        style="background: goldenrod; color: white"
                                        label="Teilnehmen"
                                        icon="rocket"
                                        @click="onTeilnehmen()"
                                    />
                                    <q-btn
                                        v-if="checkStartJet() && kurs.terminated == ''"
                                        style="background: goldenrod; color: white"
                                        label="zum Kurs"
                                        icon="rocket_launch"
                                        @click="onTeilnehmen()"
                                    />
                                    <q-btn
                                        v-if="checkEnd(kurs)"
                                        style="background: goldenrod; color: white"
                                        label="beenden"
                                        icon="sports_score"
                                        @click="onEnde()"
                                    />
                                    <q-btn
                                        style="background: goldenrod; color: white"
                                        label="Startseite"
                                        icon="menu_book"
                                        @click="this.$router.push('/')"
                                    />
                                    <div v-if="kurs.terminated != ''">
                                        Kurs beendet am {{ kurs.terminated }} <br />
                                        <q-btn
                                            style="background: goldenrod; color: white"
                                            label="Wieder Öffnen"
                                            icon="undo"
                                            @click="kurs.terminated = ''"
                                        />
                                    </div>
                                </div>
                            </div>
                            <div class="q-pt-none html">
                                <div class="text-h6">Beschreibung:<br /></div>
                                <q-card-section v-html="kurs.description" />
                            </div>
                            <div class="q-pa-md">
                                <div class="text-h6">Video:<br /></div>
                                <q-card-section>
                                    <div style="height: auto; width: 800px">
                                        <q-video :ratio="16 / 9" :src="kurs.description_movie" />
                                    </div>
                                </q-card-section>
                            </div>
                            <div class="q-pt-none html">
                                <div class="text-h6">Voraussetzungen:<br /></div>
                                <q-card-section v-if="kurs.requirements" v-html="kurs.requirements" />
                                <q-card-section v-if="!kurs.requirements">Keine</q-card-section>
                            </div>
                        </div>
                    </q-form>
                </div>
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
    name: "KursOverviewPage",
    props: ["user", "server"],
    data: () => {
        return {
            kurs: {},
            start_jet: false,
            schwierigkeit: ["leicht", "mittel", "schwer"],
        };
    },

    async created() {
        this.kurs = undefined;
        this.start_jet = false;
        if (this.user.token != "" && this.$route.params.kurs_id) {
            if (this.user.kurse) {
                for (let lKurs of this.user.kurse) {
                    if (lKurs.id == this.$route.params.kurs_id) {
                        this.kurs = lKurs;
                        break;
                    }
                }

                if (this.kurs.start) {
                    this.start_jet = true;
                }
            }
        }
        if (this.kurs == undefined) {
            this.$q.notify({
                message: `Kurs ${this.$route.params.kurs_id} nicht gefunden.`,
                color: "negative",
                icon: "warning",
            });
            this.$router.push("/");
        }
    },

    methods: {
        checkEnd(iKurs) {
            return iKurs.terminated == "" && iKurs.actual_task_index + 1 == iKurs.tasks.length && this.start_jet == true;
        },
        async onEnde() {
            this.kurs.endMe();
            this.kurs.saveMe(this.user).then((res) => {
                if (res.status == 200) {
                    this.$router.push("/");
                } else {
                    this.$q.notify({
                        message: `Fehler beim Sichern des Kurses ${this.kurs.id}.`,
                        color: "negative",
                        icon: "warning",
                    });
                }
            });
        },
        async onTeilnehmen() {
            this.kurs.startMe();
            this.start_jet = true;
            this.kurs.saveMe(this.user).then((res) => {
                if (res.status == 200) {
                    this.$router.push("/do_kurs/" + this.user.name + "/" + this.kurs.id);
                } else {
                    this.$q.notify({
                        message: `Keine Berechtigung für Kurs ${this.kurs.id}.`,
                        color: "negative",
                        icon: "warning",
                    });
                }
            });
        },

        checkStartJet() {
            return this.start_jet;
        },
    },
});
</script>
